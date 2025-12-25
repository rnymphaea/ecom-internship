package std

import (
	"log/slog"
)

type StdLogger struct {
	log *slog.Logger
}

func New() *StdLogger {
	logger := &StdLogger{
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
}
