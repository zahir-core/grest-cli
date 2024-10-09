package app

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/jeffry-luqman/zlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger returns a pointer to the loggerUtil instance (logger).
// If logger is not initialized, it creates a new loggerUtil instance, configures it, and assigns it to logger.
// It ensures that only one instance of loggerUtil is created and reused.
func Logger() *loggerUtil {
	if logger == nil {
		l := &loggerUtil{}
		l.configure()
	}
	return logger
}

// logger is a pointer to a loggerUtil instance.
// It is used to store and access the singleton instance of loggerUtil.
var logger *loggerUtil

// loggerUtil represents a logger utility.
type loggerUtil struct {
	json        *slog.Logger
	textFile    *slog.Logger
	textConsole *slog.Logger
}

// configure sets up the logging framework
//
// Even if you’re shipping your logs to a central platform, we recommend writing them to a file on your local
// machine first. You will want to make sure your logs are always available locally and not lostin the network.
// In addition, writing to a file means that you can decouple the task of writing your logs from the task of
// sending them to a central platform. Your applications themselves will not need to establish connections or
// stream your logs, and you can leave these jobs to specialized software like the Datadog Agent. If you’re
// running your Go applications within a containerized infrastructure that does not already include persistent
// storage (e.g., containers running on AWS Fargate) you may want to configure your log management tool to
// collect logs directly from your containers’ STDOUT and STDERR streams (this is handled differently in Docker
// and Kubernetes).
//
// The output log file will be located at LOG_FILE_FILENAME and will be rolled according to configuration set.
func (l *loggerUtil) configure() {
	logger = &loggerUtil{}
	mapLevel := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"info":    slog.LevelInfo,
		"warning": slog.LevelWarn,
		"error":   slog.LevelError,
	}
	level, ok := mapLevel[LOG_LEVEL]
	if !ok {
		level = slog.LevelInfo
	}
	var jsonWriters []io.Writer
	if LOG_CONSOLE_ENABLED {
		if LOG_CONSOLE_WITH_JSON {
			jsonWriters = append(jsonWriters, os.Stdout)
		} else {
			zlog.HandlerOptions = &slog.HandlerOptions{Level: level}
			zlog.TimeFormat = LOG_CONSOLE_TIME_FORMAT
			logger.textConsole = zlog.New()
		}
	}
	if LOG_FILE_ENABLED && ENV_FILE == "" {
		logFileWriter := &lumberjack.Logger{
			Filename:   LOG_FILE_FILENAME,
			MaxSize:    LOG_FILE_MAX_SIZE,
			MaxAge:     LOG_FILE_MAX_AGE,
			MaxBackups: LOG_FILE_MAX_BACKUPS,
		}
		if LOG_FILE_WITH_JSON {
			jsonWriters = append(jsonWriters, logFileWriter)
		} else {
			logger.textFile = slog.New(slog.NewTextHandler(logFileWriter, &slog.HandlerOptions{Level: level}))
		}
	}
	if len(jsonWriters) > 0 {
		logger.json = slog.New(slog.NewJSONHandler(io.MultiWriter(jsonWriters...), &slog.HandlerOptions{Level: level}))
	}
}

// Attrs return common log attributes from Ctx
func (l *loggerUtil) Attrs(c Ctx, attrss ...[]any) []any {
	var attrs []any
	if len(attrss) > 0 {
		attrs = attrss[0]
	}
	attrs = append(attrs, slog.String("method", c.Action.Method))
	attrs = append(attrs, slog.String("path", c.Action.Path))
	if c.Err != nil {
		attrs = append(attrs, slog.String("err_message", Error().GetError(c.Err).OriginalMessage()+"\n"))
		detailByte, _ := json.MarshalIndent(Error().GetError(c.Err).Body(), "", "  ")
		if detail := string(detailByte); detail != "" && detail != "null" {
			attrs = append(attrs, slog.String("err_detail", detail+"\n"))
		}
		traceByte, _ := json.MarshalIndent(Error().GetError(c.Err).TraceSimple(), "", "  ")
		if trace := string(traceByte); trace != "" && trace != "null" {
			attrs = append(attrs, slog.String("err_trace", trace+"\n"))
		}
	}
	return attrs
}

func (l *loggerUtil) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if l.json != nil {
		l.json.Log(ctx, level, msg, args...)
	}
	if l.textFile != nil {
		l.textFile.Log(ctx, level, msg, args...)
	}
	if l.textConsole != nil {
		var inclArgs []any
		exclKeys := strings.Split(LOG_CONSOLE_EXCLUDED_KEYS, ",")
		for _, a := range args {
			if attr, ok := a.(slog.Attr); !ok || !slices.Contains(exclKeys, attr.Key) {
				inclArgs = append(inclArgs, a)
			}
		}
		l.textConsole.Log(ctx, level, msg, inclArgs...)
	}
}

// Debug logs a message at LevelDebug.
func (l *loggerUtil) Debug(msg string, attrs ...any) {
	attrs = l.addBaseAttr(attrs...)
	l.Log(context.Background(), slog.LevelDebug, msg, attrs...)
}

// Info logs a message at LevelInfo.
func (l *loggerUtil) Info(msg string, attrs ...any) {
	attrs = l.addBaseAttr(attrs...)
	l.Log(context.Background(), slog.LevelInfo, msg, attrs...)
}

// Warn logs a message at LevelWarn.
func (l *loggerUtil) Warn(msg string, attrs ...any) {
	attrs = l.addBaseAttr(attrs...)
	l.Log(context.Background(), slog.LevelWarn, msg, attrs...)
}

// Error logs a message at LevelError and send alert to telegram.
func (l *loggerUtil) Error(msg string, attrs ...any) {
	attrs = l.addBaseAttr(attrs...)
	l.sendAlert(msg, attrs)
	l.Log(context.Background(), slog.LevelError, msg, attrs...)
}

// Fatal is equivalent to Error() followed by a call to os.Exit(1).
func (l *loggerUtil) Fatal(msg string, attrs ...any) {
	l.Error(msg, attrs...)
	os.Exit(1)
}

// Panic is equivalent to Error() followed by a call to panic().
func (l *loggerUtil) Panic(msg string, attrs ...any) {
	l.Error(msg, attrs...)
	panic(msg)
}

func (l *loggerUtil) addBaseAttr(args ...any) []any {
	var attrs []any
	hostname, _ := os.Hostname()
	if hostname != "" {
		attrs = append(attrs, slog.String("hostname", hostname))
	}
	attrs = append(attrs, slog.String("env", APP_ENV))
	attrs = append(attrs, slog.String("version", APP_VERSION))
	return append(attrs, args...)
}

// sendAlert send an alert to telegram
func (l *loggerUtil) sendAlert(msg string, args []any) {
	detail, trace := "", ""
	for _, a := range args {
		if attr, ok := a.(slog.Attr); ok {
			if e, ok := attr.Value.Any().(error); ok {
				msg = Error().GetError(e).OriginalMessage()
				detailByte, _ := json.MarshalIndent(Error().GetError(e).Body(), "", "  ")
				detail = string(detailByte)
				traceByte, _ := json.MarshalIndent(Error().GetError(e).TraceSimple(), "", "  ")
				trace = string(traceByte)
			} else if attr.Key == "err_message" {
				msg = attr.Value.String()
			} else if attr.Key == "err_detail" {
				detail = attr.Value.String()
			} else if attr.Key == "err_trace" {
				trace = attr.Value.String()
			}
		}
	}

	writeln := func(s *strings.Builder, args ...string) {
		for _, arg := range args {
			s.WriteString(arg)
		}
		s.WriteByte('\n')
	}

	b := &strings.Builder{}
	writeln(b, "```")
	writeln(b, msg)
	acceptedKey := []string{"env", "version", "method", "url", "base_url", "end_point", "referer", "ip", "hostname", "time", "debug"}
	for _, a := range args {
		if attr, ok := a.(slog.Attr); ok && slices.Contains(acceptedKey, attr.Key) {
			writeln(b, attr.Key, ": ", attr.Value.String())
		}
	}
	if detail != "" && detail != "null" {
		writeln(b, "error: ", detail)
	}
	if trace != "" && trace != "null" {
		writeln(b, "trace: ", trace)
	}
	writeln(b, "```")
	Telegram(b.String()).Send()
}
