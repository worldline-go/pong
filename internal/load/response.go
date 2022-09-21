package load

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/worldline-go/pong/internal/model"
)

func Response(response *model.ModuleResponse) error {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}

	fmt.Println(string(responseBytes))

	return nil
}

func ResponseError(err error) {
	response := model.ModuleResponse{
		Msg:    err.Error(),
		Failed: true,
	}

	ResponseLog(&response)
}

func ResponseLog(response *model.ModuleResponse) {
	if err := Response(response); err != nil {
		log.Error().Err(err).Msg("while responding")
	}
}
