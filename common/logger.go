package common

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(ctx context.Context) LogEvent
	Info(ctx context.Context) LogEvent
	Warn(ctx context.Context) LogEvent
	Error(ctx context.Context) LogEvent
}

type LogEvent interface {
	Send()
	Msg(msg string)
	Msgf(msg string, v ...any)
	Err(err error) LogEvent
	Strs(strs []struct {
		Key   string
		Value string
	}) LogEvent
}

type logger struct {
	settings Settings
	logger   zerolog.Logger
}

func NewLogger(settings Settings) Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	var zLogger zerolog.Logger

	if settings.IsDev() {
		zLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	} else {
		zLogger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel)
	}

	return &logger{
		settings: settings,
		logger:   zLogger.With().Timestamp().Logger(),
	}
}

func (l *logger) Debug(ctx context.Context) LogEvent {
	return l.newEvent(ctx, l.logger.Debug())
}

func (l *logger) Info(ctx context.Context) LogEvent {
	return l.newEvent(ctx, l.logger.Info())
}

func (l *logger) Warn(ctx context.Context) LogEvent {
	return l.newEvent(ctx, l.logger.Warn())
}

func (l *logger) Error(ctx context.Context) LogEvent {
	return l.newEvent(ctx, l.logger.Error())
}

type logEvent struct {
	event *zerolog.Event
}

// add additional context to the event log
func (l *logger) newEvent(ctx context.Context, event *zerolog.Event) LogEvent {
	event.Str("hostname", l.settings.Hostname())
	event.Str("appname", l.settings.Appname())

	traceID := ctx.Value(ContextKeyTraceID)
	if traceID != nil {
		event.Str("traceid", traceID.(string))
	}

	remoteAddr := ctx.Value(ContextKeyRemoteAddress)
	if remoteAddr != nil && remoteAddr.(string) != "" {
		event.Str("remoteaddr", remoteAddr.(string))
	}

	return &logEvent{
		event,
	}
}

func (e *logEvent) Send() {
	e.event.Send()
}

func (e *logEvent) Msg(msg string) {
	e.event.Msg(msg)
}

func (e *logEvent) Msgf(msg string, v ...any) {
	e.event.Msgf(msg, v...)
}

func (e *logEvent) Err(err error) LogEvent {
	e.event.Stack().Err(err)
	return e
}

func (e *logEvent) Strs(strs []struct {
	Key   string
	Value string
}) LogEvent {
	for _, str := range strs {
		e.event.Str(str.Key, str.Value)
	}
	return e
}
