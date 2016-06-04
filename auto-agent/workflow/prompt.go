package workflow

import (
	"strings"
  // "fmt"
  // "os"
)

type Prompt interface {
	// TODO add these bigger structure
	// GetParentPhase() Phase
	// GetParentStrategy() Strategy
	// GetDisplayText() string
	// GetUIActionModeId() string
	GetPhaseId() string
	GetResponse() Response
	GetNextPrompt() Prompt
	// SetResponse(Response)
	GetPromptId() string
	GetUIPrompt() UIPrompt
	ProcessResponse(string)
}

type Response struct {
	Text string
	Id string
}

type ExpectedResponseHandler struct {
	expectedResponseMap map[string]Prompt
}

func MakeExpectedResponseHandler(ecs[]ExpectedResponseConfig, phaseId string) *ExpectedResponseHandler {
	erh := new(ExpectedResponseHandler)
	erh.expectedResponseMap = make(map[string]Prompt)
	for _, v := range ecs {
    // In case if the prompt is already defined
		// erh.expectedResponseMap[strings.ToLower(v.Id)] = MakePromptFromConfig(v.NextPrompt, phaseId)
    erh.expectedResponseMap[strings.ToLower(v.Id)] = MakePrompt(v.NextPrompt.Id, phaseId)
	}
	return erh
}

func (erh *ExpectedResponseHandler) GetNextPrompt(rid string) Prompt {
  // TODO cleanup
  // fmt.Fprintf(os.Stderr, "expectedResponseMap: %s", erh.expectedResponseMap, "\n\n")
  return erh.expectedResponseMap[strings.ToLower(rid)]
}
