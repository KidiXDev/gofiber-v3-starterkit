package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	levelStr := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	if levelStr == "" {
		if strings.ToLower(os.Getenv("APP_ENV")) == "production" {
			levelStr = "info"
		} else {
			levelStr = "debug"
		}
	}
	lvl, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}

	output.FormatTimestamp = func(i any) string {
		switch v := i.(type) {
		case time.Time:
			return fmt.Sprintf("[%s]", v.Format(time.DateTime))
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				return fmt.Sprintf("[%s]", t.Format(time.DateTime))
			}
			return fmt.Sprintf("[%s]", v)
		default:
			s := fmt.Sprint(i)
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return fmt.Sprintf("[%s]", t.Format(time.DateTime))
			}
			if s == "<nil>" {
				return ""
			}
			return fmt.Sprintf("[%s]", s)
		}
	}

	output.FormatLevel = func(i any) string {
		const (
			red     = "\x1b[31m"
			yellow  = "\x1b[33m"
			green   = "\x1b[32m"
			blue    = "\x1b[34m"
			magenta = "\x1b[35m"
			reset   = "\x1b[0m"
		)

		lvl := strings.ToUpper(fmt.Sprintf("%-6s", i))
		levelKey := strings.TrimSpace(lvl)

		var color string
		switch levelKey {
		case "DEBUG":
			color = magenta
		case "INFO":
			color = green
		case "WARN", "WARNING":
			color = yellow
		case "ERROR", "FATAL", "PANIC":
			color = red
		default:
			color = blue
		}

		return fmt.Sprintf("%s[%s]%s", color, levelKey, reset)
	}

	writer := output

	log.Logger = log.Output(writer).With().Timestamp().Logger()
}
