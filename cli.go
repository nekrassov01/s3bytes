package s3bytes

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

type app struct {
	*cli.App
	completion  *cli.StringFlag
	profile     *cli.StringFlag
	loglevel    *cli.StringFlag
	regions     *cli.StringSliceFlag
	metricName  *cli.StringFlag
	storageType *cli.StringFlag
	prefix      *cli.StringFlag
	expression  *cli.StringFlag
	output      *cli.StringFlag
}

func CLI(ctx context.Context) {
	app := newApp(os.Stdout, os.Stderr)
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func newApp(w, ew io.Writer) *app {
	logger := newAppLogger(ew)
	log.SetDefault(logger)
	a := app{}
	a.completion = &cli.StringFlag{
		Name:    "completion",
		Aliases: []string{"c"},
		Usage:   "print completion scripts",
	}
	a.profile = &cli.StringFlag{
		Name:    "profile",
		Aliases: []string{"p"},
		Usage:   "set aws profile",
		EnvVars: []string{"AWS_PROFILE"},
	}
	a.loglevel = &cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Usage:   "set log level",
		EnvVars: []string{canonicalName + "_LOG_LEVEL"},
		Value:   log.InfoLevel.String(),
	}
	a.regions = &cli.StringSliceFlag{
		Name:        "regions",
		Aliases:     []string{"r"},
		Usage:       "set target regions",
		Value:       cli.NewStringSlice(regions...),
		DefaultText: "all regions with no opt-in required",
	}
	a.metricName = &cli.StringFlag{
		Name:    "metric-name",
		Aliases: []string{"m"},
		Usage:   "set metric name",
		Value:   MetricNameBucketSizeBytes.String(),
	}
	a.storageType = &cli.StringFlag{
		Name:    "storage-type",
		Aliases: []string{"s"},
		Usage:   "set storage type",
		Value:   StorageTypeStandardStorage.String(),
	}
	a.prefix = &cli.StringFlag{
		Name:    "prefix",
		Aliases: []string{"P"},
		Usage:   "set bucket name prefix",
	}
	a.expression = &cli.StringFlag{
		Name:    "expression",
		Aliases: []string{"e"},
		Usage:   "set filter expression for metric values",
	}
	a.output = &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "set output type",
		EnvVars: []string{canonicalName + "_OUTPUT_TYPE"},
		Value:   OutputTypeText.String(),
	}
	a.App = &cli.App{
		Name:                 appName,
		Version:              Version,
		Usage:                "S3 size checker CLI",
		Description:          "A tool to get the size of s3 buckets, or number of objects",
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Writer:               w,
		ErrWriter:            ew,
		Before:               a.before,
		Action:               a.action,
		Flags:                []cli.Flag{a.completion, a.profile, a.loglevel, a.regions, a.prefix, a.expression, a.metricName, a.storageType, a.output},
		Metadata:             map[string]any{},
	}
	return &a
}

func (a *app) before(c *cli.Context) error {
	level, err := log.ParseLevel(c.String(a.loglevel.Name))
	if err != nil {
		return err
	}
	log.SetLevel(level)
	cfg, err := LoadAWSConfig(c.Context, a.ErrWriter, c.String(a.profile.Name), level)
	if err != nil {
		return err
	}
	a.Metadata["config"] = cfg
	return err
}

func (a *app) action(c *cli.Context) error {
	if c.IsSet(a.completion.Name) {
		return a.comp(c)
	}
	metricName, storageType, outputType, err := a.parse(c)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("started: metric name: %s, storage type: %s, output: %s", metricName, storageType, outputType))
	metrics, err := a.run(c, metricName, storageType)
	if err != nil {
		return err
	}
	a.sort(metrics)
	if err := a.render(metrics, metricName, outputType); err != nil {
		return err
	}
	log.Info("completed")
	return nil
}

func (a *app) comp(c *cli.Context) error {
	n, err := parseShell(c.String(a.completion.Name))
	if err != nil {
		return err
	}
	switch n {
	case bash:
		fmt.Fprintln(a.Writer, completionBash)
	case zsh:
		fmt.Fprintln(a.Writer, completionZsh)
	case pwsh:
		fmt.Fprintln(a.Writer, completionPwsh)
	default:
	}
	return nil
}

func (a *app) parse(c *cli.Context) (MetricName, StorageType, OutputType, error) {
	metricName, err := ParseMetricName(c.String(a.metricName.Name))
	if err != nil {
		return 0, 0, 0, err
	}
	storageType, err := ParseStorageType(c.String(a.storageType.Name))
	if err != nil {
		return 0, 0, 0, err
	}
	outputType, err := ParseOutputType(c.String(a.output.Name))
	if err != nil {
		return 0, 0, 0, err
	}
	return metricName, storageType, outputType, nil
}

func (a *app) run(c *cli.Context, metricName MetricName, storageType StorageType) ([]Metric, error) {
	var (
		eg, ctx    = errgroup.WithContext(c.Context)
		metricChan = make(chan []Metric, len(regions))
		errorChan  = make(chan error, 1)
		cfg        = a.Metadata["config"].(aws.Config)
		client     = NewClient(cfg)
	)
	for _, region := range c.StringSlice(a.regions.Name) {
		region := region
		eg.Go(func() error {
			man, err := NewManager(ctx, client, region, c.String(a.prefix.Name), c.String(a.expression.Name), metricName, storageType)
			if err != nil {
				return err
			}
			if err := man.SetBuckets(); err != nil {
				return err
			}
			if err := man.SetQueries(); err != nil {
				return err
			}
			if err := man.SetData(); err != nil {
				return err
			}
			man.Debug()
			select {
			case metricChan <- man.Metrics:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}
	go func() {
		if err := eg.Wait(); err != nil {
			select {
			case errorChan <- err:
			default:
			}
		}
		close(metricChan)
	}()
	metrics := make([]Metric, 0, len(regions)*maxQueries)
	for m := range metricChan {
		metrics = append(metrics, m...)
	}
	select {
	case err := <-errorChan:
		return nil, err
	default:
	}
	return metrics, nil
}

func (a *app) sort(metrics []Metric) {
	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].Value() != metrics[j].Value() {
			return metrics[i].Value() > metrics[j].Value()
		}
		return metrics[i].Label() < metrics[j].Label()
	})
}

func (a *app) render(metrics []Metric, metricName MetricName, outputType OutputType) error {
	ren := NewRenderer(a.Writer, metrics, metricName, outputType)
	if err := ren.Render(); err != nil {
		return err
	}
	return nil
}
