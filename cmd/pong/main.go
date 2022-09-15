package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"

	"github.com/worldline-go/pong/internal/load"
	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/route"
)

var (
	// Populated by goreleaser during build
	version = "v0.0.0"
	commit  = "?"
	date    = ""
)

const helpText = `pong [OPTIONS] <multiple-file.[json|yaml|yml]>
version:[%s] commit:[%s] buildDate:[%s]

Check server up and running

Options:
  -v, --version
    Show version number
  -h, --help
    Show help

Examples:
  pong test.yml test2.json
`

func usage() {
	fmt.Printf(helpText, version, commit, date)
	os.Exit(0)
}

var (
	flagVersion bool
)

func flagParse() []string {
	flag.Usage = usage

	flag.BoolVar(&flagVersion, "v", false, "")
	flag.BoolVar(&flagVersion, "version", false, "")

	flag.Parse()

	// Check Values
	if flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return flag.Args()
}

func main() {
	logz.InitializeLog(nil)

	files := flagParse()

	exitCode := 0
	wg := &sync.WaitGroup{}

	defer func() {
		wg.Wait()
		os.Exit(exitCode)
	}()

	// check length of the arguments
	if len(files) == 0 {
		if err := load.Response(&model.ModuleResponse{
			Msg:    "Missing argument file",
			Failed: true,
		}); err != nil {
			log.Error().Err(err).Msg("error while responding")
		}

		exitCode = 1

		return
	}

	// start operation
	chNotify := make(chan os.Signal, 1)

	signal.Notify(chNotify, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		case <-chNotify:
			log.Info().Msg("shutting down...")
			log.Info().Msg("send signal again to exit force")
			signal.Stop(chNotify)
			close(chNotify)
			cancel()
		}
	}()

	var errRequests []error
	// read config
	for _, file := range files {
		args, err := load.ReadConfig(file)
		if err != nil {
			if err := load.Response(&model.ModuleResponse{
				Msg:    err.Error(),
				Failed: true,
			}); err != nil {
				log.Error().Err(err).Msg("error while responding")
			}

			exitCode = 1

			return
		}

		errs := route.Request(ctx, args)
		if len(errs) > 0 {
			errRequests = append(errRequests, errs...)
		}
	}

	if len(errRequests) == 0 {
		load.ResponseLog(&model.ModuleResponse{
			Msg:    "OK",
			Failed: false,
		})

		return
	}

	var errStrings []string
	for _, err := range errRequests {
		log.Warn().Err(err).Msg("while doing request")
		errStrings = append(errStrings, fmt.Sprintf("%v", err.Error()))
	}

	load.ResponseLog(&model.ModuleResponse{
		Msg:    strings.Join(errStrings, " "),
		Failed: true,
	})
}
