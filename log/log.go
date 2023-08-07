package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

// https://github.com/samber/slog-syslog

func main() {
	/*
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
	*/
	// writer, err := net.Dial("udp", "localhost:514")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	writer := os.Stdout

	if writer != nil {
		// defer writer.Close()
		textHandler := slog.NewTextHandler(writer, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
		logger := slog.New(textHandler)

		logger.Debug("from slog!")
		logger.Info("be careful about 70 handcup!")

		logger.
			With("category", "sql").
			With("query.statement", "SELECT COUNT(*) FROM users;").
			With("query.duration", 1*time.Second).
			With("error", fmt.Errorf("could not count users")).
			Error("caramba!")

		logger.
			With(
				slog.Group("user",
					slog.String("id", "user-123"),
					slog.Time("created_at", time.Now()),
				),
			).
			Info("user registration")
	}

	// here is the syslog sample
	// s, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_LOCAL7, "aprilsh")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// s.Warning("from syslog")
	// s.Info("I am 98 handcup.")

}
