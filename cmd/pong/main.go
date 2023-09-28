package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/initializer"
	"github.com/worldline-go/logz"

	"github.com/worldline-go/pong/internal/load"
	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/route"
	"github.com/worldline-go/pong/pkg/template"
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

var flagVersion bool

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
	files := flagParse()

	initFn := Init{
		files: files,
	}

	initializer.Init(
		initFn.run,
		initializer.WithMsgf("pong [%s]", version),
		initializer.WithOptionsLogz(logz.WithCaller(false)),
	)
}

type Init struct {
	files []string
}

func (s Init) run(ctx context.Context, wg *sync.WaitGroup) error {
	// check length of the arguments
	if len(s.files) == 0 {
		err := fmt.Errorf("missing argument file")
		load.ResponseError(err)

		return err
	}

	var errRequests []error
	// read config
	for _, file := range s.files {
		args, err := load.ReadConfig(file)
		if err != nil {
			load.ResponseError(err)

			return err
		}

		if len(args.Delimeters) >= 2 {
			if err := template.SetDelimeters(args.Delimeters); err != nil {
				load.ResponseError(err)

				return err
			}
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

		return nil
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

	return nil
}
