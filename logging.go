package DnsCompare

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func ToTimeFormat(s string) string {
	switch strings.ToLower(s) {
	case "ansi":
		fallthrough
	case "ansic":
		return time.ANSIC
	case "unix":
		fallthrough
	case "unixdate":
		return time.UnixDate
	case "ruby":
		fallthrough
	case "rubydate":
		return time.RubyDate
	case "rfc822":
		return time.RFC822
	case "rfc822z":
		return time.RFC822Z
	case "rfc850":
		return time.RFC850
	case "rfc1123":
		return time.RFC1123
	case "rfc1123z":
		return time.RFC1123Z
	case "json":
		fallthrough
	case "rfc3339":
		return time.RFC3339
	case "rfc3339nano":
		return time.RFC3339Nano
	case "kitchen":
		return time.Kitchen
	case "stamp":
		return time.Stamp
	case "stampmili":
		return time.StampMilli
	case "stampmicro":
		return time.StampMicro
	case "stampnano":
		return time.StampNano
	default:
		return s
	}
}

// InitializeLogging - Setup the logger
func InitializeLogger() zerolog.Logger {

	logOutput := viper.GetString("Log.Output")
	if strings.ToLower(logOutput) == "console" {
		switch strings.ToLower(viper.GetString("Log.Format")) {
		case "plain":
			log = zerolog.New(zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: ToTimeFormat(viper.GetString("Log.TimeFormat")),
				NoColor:    !viper.GetBool("Log.Colors"),
			}).With().Timestamp().Logger()
		case "json":
			log = zerolog.New(os.Stdout).With().Timestamp().Logger()
		}

	} else {
		file, err := os.OpenFile(viper.GetString("Log.Output"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log = zerolog.New(file).With().Timestamp().Logger()
		} else {
			log = zerolog.New(zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: ToTimeFormat(viper.GetString("Log.TimeFormat")),
				NoColor:    !viper.GetBool("Log.Colors"),
			}).With().Timestamp().Logger()
			log.Warn().Msg("Failed to log to file, using stderr")
		}
	}

	logLevel, err := zerolog.ParseLevel(viper.GetString("Log.Level"))
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	log = log.Level(logLevel)

	return log.With().Logger()
}
