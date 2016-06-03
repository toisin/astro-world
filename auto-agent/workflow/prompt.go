package workflow

import (
	"strings"
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
		erh.expectedResponseMap[strings.ToLower(v.Id)] = MakePromptFromConfig(v.NextPrompt, phaseId)
	}
	return erh
}

func (erh *ExpectedResponseHandler) GetNextPrompt(rid string) Prompt {
	return erh.expectedResponseMap[strings.ToLower(rid)]
}


// type Factor struct {
// 	// TODO
// }

// type PromptGenerator interface {
// 	// GetPromptText() string // actual text ready to be displayed as a prompt
// 	// GetUIActionModeId() string  // the mode of rendering for Action UI
// 	GenerateUIPrompt() UIPrompt
// }



// func (erh *ExpectedResponseHandler) GetNextPrompt(r Response) Prompt {
// 	return erh.expectedResponseMap[r.GetId()]
// }


// func NewTextResponse(t string, id string) *TextResponse {
// 	return &TextResponse{t, id}
// 	// n := new(TextResponse)
// 	// n.text = t
// 	// n.id = id
// 	// return n
// }

// type TextResponse struct {
// 	text string
// 	id string
// }

// func (tr *TextResponse) GetText() string {
// 	return tr.text
// }

// func (tr *TextResponse) GetId() string {
// 	return tr.id
// }