package main

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dustin/go-humanize"
	"github.com/nekrassov01/logwrapper/log"
	"github.com/nekrassov01/s3bytes"
	"github.com/urfave/cli/v3"
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

func newCmd(w, ew io.Writer) *cli.Command {
	logger = log.NewAppLogger(ew, defaultLogLevel, defaultLogStyle, label)

	profile := &cli.StringFlag{
		Name:    "profile",
		Aliases: []string{"p"},
		Usage:   "set aws profile",
		Sources: cli.EnvVars("AWS_PROFILE"),
	}

	loglevel := &cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Usage:   "set log level",
		Sources: cli.EnvVars(label + "_LOG_LEVEL"),
		Value:   log.InfoLevel.String(),
	}

	region := &cli.StringSliceFlag{
		Name:        "region",
		Aliases:     []string{"r"},
		Usage:       "set target regions",
		Value:       s3bytes.DefaultRegions,
		DefaultText: "all regions with no opt-in",
	}

	metricName := &cli.StringFlag{
		Name:    "metric-name",
		Aliases: []string{"m"},
		Usage:   "set metric name",
		Value:   s3bytes.MetricNameBucketSizeBytes.String(),
	}

	storageType := &cli.StringFlag{
		Name:    "storage-type",
		Aliases: []string{"s"},
		Usage:   "set storage type",
		Value:   s3bytes.StorageTypeStandardStorage.String(),
	}

	prefix := &cli.StringFlag{
		Name:    "prefix",
		Aliases: []string{"P"},
		Usage:   "set bucket name prefix",
	}

	filter := &cli.StringFlag{
		Name:    "filter",
		Aliases: []string{"f"},
		Usage:   "set filter expression for metric values",
	}

	output := &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "set output type",
		Sources: cli.EnvVars(label + "_OUTPUT_TYPE"),
		Value:   s3bytes.OutputTypeCompressedText.String(),
	}

	before := func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// parse log level passed as string
		level, err := log.ParseLevel(cmd.String(loglevel.Name))
		if err != nil {
			return ctx, err
		}

		// set logger for the application
		logger.SetLevel(level)

		// load aws config with the specified profile
		cfg, err := s3bytes.LoadConfig(ctx, cmd.String(profile.Name))
		if err != nil {
			return ctx, err
		}

		// set logger for the AWS SDK
		cfg.Logger = log.NewSDKLogger(ew, level, defaultLogStyle, "SDK")
		cfg.ClientLogMode = aws.LogRequest | aws.LogResponse | aws.LogRetries | aws.LogSigning | aws.LogDeprecatedUsage

		// set aws config to the metadata
		cmd.Metadata["config"] = cfg

		return ctx, nil
	}

	action := func(ctx context.Context, cmd *cli.Command) error {
		// parse metric name passed as string
		metricName, err := s3bytes.ParseMetricName(cmd.String(metricName.Name))
		if err != nil {
			return err
		}

		// parse storage type passed as string
		storageType, err := s3bytes.ParseStorageType(cmd.String(storageType.Name))
		if err != nil {
			return err
		}

		// parse output type passed as string
		outputType, err := s3bytes.ParseOutputType(cmd.String(output.Name))
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
		cfg := cmd.Metadata["config"].(aws.Config)

		// create a new client
		client := s3bytes.NewClient(cfg)

		// initialize the manager
		man := s3bytes.NewManager(client)

		// set regions to the manager
		if err := man.SetRegion(cmd.StringSlice(region.Name)); err != nil {
			return err
		}

		// set metric name and strorage type to the manager
		if err := man.SetMetric(metricName, storageType); err != nil {
			return err
		}

		// set prefix to the manager
		if err := man.SetPrefix(cmd.String(prefix.Name)); err != nil {
			return err
		}

		// set filter to the manager
		if err := man.SetFilter(cmd.String(filter.Name)); err != nil {
			return err
		}

		// run list operation
		data, err := man.List(ctx)
		if err != nil {
			return err
		}
		debug(man)

		// sort metrics
		s3bytes.SortMetrics(data)

		// render result
		ren := s3bytes.NewRenderer(w, data, outputType)
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

	return &cli.Command{
		Name:                  name,
		Version:               getVersion(),
		Usage:                 "S3 size checker",
		Description:           "Check the size of all buckets in S3 in one shot.",
		HideHelpCommand:       true,
		EnableShellCompletion: true,
		Writer:                w,
		ErrWriter:             ew,
		Before:                before,
		Action:                action,
		Flags:                 []cli.Flag{profile, loglevel, region, prefix, filter, metricName, storageType, output},
		Metadata:              map[string]any{},
	}
}

func debug(man *s3bytes.Manager) {
	logger.Debug("Manager\n" + man.String() + "\n")
}
