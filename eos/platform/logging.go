package platform

import (
	"log/slog"

	"github.com/mydev/go-eos/eos/internal/cbinding"
)

func initLogging() {
	cbinding.EOS_Logging_SetCallback(func(msg *cbinding.EOS_LogMessage) {
		attrs := []slog.Attr{
			slog.String("category", msg.Category),
		}

		level := slogLevel(msg.Level)
		slog.Default().LogAttrs(nil, level, msg.Message, attrs...)
	})

	cbinding.EOS_Logging_SetLogLevel(cbinding.EOS_LC_AllCategories, cbinding.EOS_LOG_VeryVerbose)
}

func slogLevel(level cbinding.EOS_ELogLevel) slog.Level {
	switch {
	case level <= cbinding.EOS_LOG_Error:
		return slog.LevelError
	case level <= cbinding.EOS_LOG_Warning:
		return slog.LevelWarn
	case level <= cbinding.EOS_LOG_Info:
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}
