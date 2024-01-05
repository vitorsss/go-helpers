package logs

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

var goroot string

func init() {
	goroot = runtime.GOROOT()

	zerolog.ErrorStackMarshaler = ErrorStackMarshaler
	output := os.Getenv("ZEROLOG_OUTPUT")
	if output == "json" {
		Logger = log.Logger
	} else {
		// will probably only use a console decorator during development
		removeGoPath = func(source string) string { return source }
		Logger = log.Logger.Output(zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				w.TimeFormat = time.RFC3339
				w.FormatLevel = func(i interface{}) string {
					return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
				}
				w.FormatMessage = func(i interface{}) string {
					return fmt.Sprintf("# %s #", i)
				}
			},
		))
	}
	Logger = Logger.Level(zerolog.InfoLevel).With().Stack().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(Logger)
}
