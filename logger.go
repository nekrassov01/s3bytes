package s3bytes

import (
	"fmt"
	"io"

	"github.com/aws/smithy-go/logging"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func newAppLogger(w io.Writer) *log.Logger {
	styles := setLoggerStyles()
	logger := log.New(w).WithPrefix(canonicalName)
	logger.SetStyles(styles)
	return logger
}

type sdkLogger struct {
	Logger *log.Logger
}

func (l *sdkLogger) Logf(c logging.Classification, format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	switch c {
	case logging.Warn:
		l.Logger.Warn(s)
	case logging.Debug:
		l.Logger.Debug(s)
	default:
		l.Logger.Info(s)
	}
}

func newSDKLogger(w io.Writer, loglevel log.Level) *sdkLogger {
	styles := setLoggerStyles()
	logger := log.New(w).WithPrefix("SDK")
	logger.SetStyles(styles)
	logger.SetLevel(loglevel)
	return &sdkLogger{
		Logger: logger,
	}
}

func setLoggerStyles() *log.Styles {
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DBG").
		Bold(true).
		MaxWidth(3).
		Foreground(lipgloss.Color("63"))
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INF").
		Bold(true).
		MaxWidth(3).
		Foreground(lipgloss.Color("86"))
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WRN").
		Bold(true).
		MaxWidth(3).
		Foreground(lipgloss.Color("192"))
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERR").
		Bold(true).
		MaxWidth(3).
		Foreground(lipgloss.Color("204"))
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FTL").
		Bold(true).
		MaxWidth(3).
		Foreground(lipgloss.Color("134"))
	return styles
}
