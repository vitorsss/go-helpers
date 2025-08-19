package logs

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

var goroot string

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func init() {
	goroot = runtime.GOROOT()

	zerolog.ErrorStackMarshaler = ErrorStackMarshaler
	output := os.Getenv("ZEROLOG_OUTPUT")
	if output == "json" {
		Logger = log.Logger
	} else {
		// will probably only use a console decorator during development
		Logger = log.Logger.Output(zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				e := os.Getenv("NO_COLOR")
				if e != "" {
					w.NoColor, _ = strconv.ParseBool(e)
				}
				w.TimeFormat = time.RFC3339
				w.FormatLevel = func(i any) string {
					level := strings.ToUpper(fmt.Sprintf("| %-6s|", i))
					if w.NoColor {
						return level
					}
					color := colorBlack
					switch level {
					case "| INFO  |":
						color = colorCyan
					case "| WARN  |":
						color = colorYellow
					case "| ERROR |":
						color = colorRed
					default:
						return level
					}
					return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, level)
				}
				w.FormatMessage = func(i any) string {
					return fmt.Sprintf("# %s #", i)
				}
			},
		))
	}
	Logger = Logger.Level(zerolog.InfoLevel).With().Stack().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(Logger)
}
