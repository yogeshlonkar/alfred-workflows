package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yogeshlonkar/alfred-workflows/pkg/app"
)

var (
	start          = time.Now()
	defaultTimeout = 60
	syncFlag       = flag.Bool("sync", false, "force sync items")
	keywordFlag    = flag.String("keyword", "", "type of suggestions. Accepted values "+app.Keywords)
	colorFlag      = flag.Bool("color", false, "print log with colors")
	debugFlag      = flag.Bool("debug", false, "print debug logs")
	timeout        = flag.Int("timeout", timeoutFromEnv(), "timeout in seconds")
)

func timeoutFromEnv() int {
	to, err := strconv.Atoi(os.Getenv("fetch_timeout"))
	if err != nil {
		return defaultTimeout
	}
	return to
}

func main() {
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !*colorFlag})
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	err := app.App(*keywordFlag, *syncFlag, time.Second*time.Duration(*timeout))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Debug().Msgf("completed %s in %s", *keywordFlag, time.Since(start).Round(time.Millisecond))
}
