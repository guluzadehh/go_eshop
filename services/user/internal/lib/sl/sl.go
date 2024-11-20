package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

type LoggerSetter interface {
	SetLog(log *slog.Logger)
}

func HandlerJob(log *slog.Logger, op string, request_id string, ls LoggerSetter) *slog.Logger {
	log = log.With(slog.String("request_id", request_id))
	ls.SetLog(log)
	return log.With("op", op)
}
