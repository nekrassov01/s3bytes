package main

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dustin/go-humanize"
	"github.com/nekrassov01/logwrapper/log"
	"github.com/nekrassov01/s3bytes"
	"github.com/urfave/cli/v2"
)

const (
	name  = "s3bytes"
	label = "S3BYTES"
)

var (
	logger          = &log.AppLogger{}
	defaultLogLevel = log.InfoLevel
	defaultLogStyle = log.DefaultStyles()
)

type app struct {
	*cli.App

	completion  *cli.StringFlag
	profile     *cli.StringFlag
	loglevel    *cli.StringFlag
	region      *cli.StringSliceFlag
	metricName  *cli.StringFlag
	storageType *cli.StringFlag
	prefix      *cli.StringFlag
	filter      *cli.StringFlag
	output      *cli.StringFlag
}

func newApp(w, ew io.Writer) *app {
	logger = log.NewAppLogger(ew, defaultLogLevel, defaultLogStyle, label)
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
		EnvVars: []string{label + "_LOG_LEVEL"},
		Value:   log.InfoLevel.String(),
	}
	a.region = &cli.StringSliceFlag{
		Name:        "region",
		Aliases:     []string{"r"},
		Usage:       "set target regions",
		Value:       cli.NewStringSlice(s3bytes.DefaultRegions...),
		DefaultText: "all regions with no opt-in",
	}
	a.metricName = &cli.StringFlag{
		Name:    "metric-name",
		Aliases: []string{"m"},
		Usage:   "set metric name",
		Value:   s3bytes.MetricNameBucketSizeBytes.String(),
	}
	a.storageType = &cli.StringFlag{
		Name:    "storage-type",
		Aliases: []string{"s"},
		Usage:   "set storage type",
		Value:   s3bytes.StorageTypeStandardStorage.String(),
	}
	a.prefix = &cli.StringFlag{
		Name:    "prefix",
		Aliases: []string{"P"},
		Usage:   "set bucket name prefix",
	}
	a.filter = &cli.StringFlag{
		Name:    "filter",
		Aliases: []string{"f"},
		Usage:   "set filter expression for metric values",
	}
	a.output = &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "set output type",
		EnvVars: []string{label + "_OUTPUT_TYPE"},
		Value:   s3bytes.OutputTypeCompressedText.String(),
	}
	a.App = &cli.App{
		Name:                 name,
		Version:              getVersion(),
		Usage:                "S3 size checker",
		Description:          "Check the size of all buckets in S3 in one shot.",
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Writer:               w,
		ErrWriter:            ew,
		Before:               a.before,
		Action:               a.action,
		Flags:                []cli.Flag{a.completion, a.profile, a.loglevel, a.region, a.prefix, a.filter, a.metricName, a.storageType, a.output},
		Metadata:             map[string]any{},
	}
	return &a
}

func (a *app) before(c *cli.Context) error {
	// parse log level passed as string
	level, err := log.ParseLevel(c.String(a.loglevel.Name))
	if err != nil {
		return err
	}

	// set logger for the application
	logger.SetLevel(level)

	// load aws config with the specified profile
	cfg, err := s3bytes.LoadConfig(c.Context, c.String(a.profile.Name))
	if err != nil {
		return err
	}

	// set logger for the AWS SDK
	cfg.Logger = log.NewSDKLogger(a.ErrWriter, level, defaultLogStyle, "SDK")
	cfg.ClientLogMode = aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogSigning | aws.LogDeprecatedUsage

	// set aws config to the metadata
	a.Metadata["config"] = cfg

	return err
}

func (a *app) action(c *cli.Context) error {
	// print completion scripts
	if c.IsSet(a.completion.Name) {
		return a.comp(c)
	}

	// parse metric name passed as string
	metricName, err := s3bytes.ParseMetricName(c.String(a.metricName.Name))
	if err != nil {
		return err
	}

	// parse storage type passed as string
	storageType, err := s3bytes.ParseStorageType(c.String(a.storageType.Name))
	if err != nil {
		return err
	}

	// parse output type passed as string
	outputType, err := s3bytes.ParseOutputType(c.String(a.output.Name))
	if err != nil {
		return err
	}

	// logging at process start
	logger.Info(
		"started",
		"at", time.Now().Format(time.RFC3339),
		"metricName", metricName,
		"storageType", storageType,
		"output", outputType,
	)

	// get aws config from the metadata
	cfg := a.Metadata["config"].(aws.Config)

	// create a new client
	client := s3bytes.NewClient(cfg)

	// initialize the manager
	man := s3bytes.NewManager(c.Context, client)

	// set regions to the manager
	if err := man.SetRegion(c.StringSlice(a.region.Name)); err != nil {
		return err
	}

	// set metric name and strorage type to the manager
	if err := man.SetMetric(metricName, storageType); err != nil {
		return err
	}

	// set prefix to the manager
	if err := man.SetPrefix(c.String(a.prefix.Name)); err != nil {
		return err
	}

	// set filter to the manager
	if err := man.SetFilter(c.String(a.filter.Name)); err != nil {
		return err
	}

	// run list operation
	data, err := man.List()
	if err != nil {
		return err
	}
	debug(man)

	// sort metrics
	s3bytes.SortMetrics(data)

	// render result
	ren := s3bytes.NewRenderer(a.Writer, data, outputType)
	if err := ren.Render(); err != nil {
		return err
	}

	// logging at process stop with total bytes
	logger.Info(
		"stopped",
		"total", humanize.Comma(data.Total),
	)

	return nil
}

func (a *app) comp(c *cli.Context) error {
	n, err := parseShell(c.String(a.completion.Name))
	if err != nil {
		return err
	}
	switch n {
	case bash:
		_, _ = fmt.Fprintln(a.Writer, completionBash)
	case zsh:
		_, _ = fmt.Fprintln(a.Writer, completionZsh)
	case pwsh:
		_, _ = fmt.Fprintln(a.Writer, completionPwsh)
	default:
	}
	return nil
}

func debug(man *s3bytes.Manager) {
	logger.Debug("Manager\n" + man.String() + "\n")
}
