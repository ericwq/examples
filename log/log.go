package main

import (
	"fmt"
	"log"
	"net"
	"time"

	slogsyslog "github.com/samber/slog-syslog"
	"golang.org/x/exp/slog"
)

// https://github.com/samber/slog-syslog

func main() {
	// ncat -u -l 9999 -k
	writer, err := net.Dial("udp", "localhost:514")
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slogsyslog.Option{Level: slog.LevelInfo, Writer: writer}.NewSyslogHandler())
	logger = logger.
		With("environment", "dev").
		With("release", "v1.0.0")

	// log error
	logger.
		With("category", "sql").
		With("query.statement", "SELECT COUNT(*) FROM users;").
		With("query.duration", 1*time.Second).
		With("error", fmt.Errorf("could not count users")).
		Error("caramba!")

	// log user signup
	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		Info("user registration")
}
