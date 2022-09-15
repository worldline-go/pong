package load

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"
	"github.com/worldline-go/pong/internal/model"
)

func ReadConfig(file string) (*model.ModuleArgs, error) {
	// read file
	v, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// unmarshal with yaml
	args := model.ModuleArgs{
		LogLevel: "info",
	}

	err = yaml.Unmarshal(v, &args)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling yaml: %w", err)
	}

	// set default values
	defaultArgs(&args)

	if err := logz.SetLogLevel(args.LogLevel); err != nil {
		return nil, fmt.Errorf("setting log level: %w", err)
	}

	log.Debug().Msgf("read file %s", file)

	return &args, nil
}

func defaultArgs(args *model.ModuleArgs) {
	for i := range args.Check.Rest {
		if args.Check.Rest[i].Concurrent == 0 {
			args.Check.Rest[i].Concurrent = model.DefaultRestClient.Concurrent
		}
		for j := range args.Check.Rest[i].Check {
			if args.Check.Rest[i].Check[j].Method == "" {
				args.Check.Rest[i].Check[j].Method = model.DefaultRestCheck.Method
			}
			if args.Check.Rest[i].Check[j].Status == 0 {
				args.Check.Rest[i].Check[j].Status = model.DefaultRestCheck.Status
			}
		}
	}
}
