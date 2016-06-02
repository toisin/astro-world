package workflow



type Prompt interface {
	// TODO add these bigger structure
	// GetParentPhase() Phase
	// GetParentStrategy() Strategy
	// GetDisplayText() string
	// GetUIActionModeId() string
	GetPhaseId() string
	// GetResponseText() string
	// GetNextPrompt() Prompt
	// SetResponse(Response)
	GetUIPrompt() UIPrompt
}


// type Factor struct {
// 	// TODO
// }

type PromptGenerator interface {
	// GetPromptText() string // actual text ready to be displayed as a prompt
	// GetUIActionModeId() string  // the mode of rendering for Action UI
	GenerateUIPrompt() UIPrompt
}



type ExpectedResponseHandler struct {
	expectedResponseMap map[string]Prompt
}

func MakeExpectedResponseHandler(ecs[]ExpectedResponseConfig) *ExpectedResponseHandler {
	erh := new(ExpectedResponseHandler)
	erh.expectedResponseMap = make(map[string]Prompt)
	for _, v := range ecs {
		erh.expectedResponseMap[v.Id] = MakeCovPrompt(v.NextPrompt)
	}
	return erh
}

// func (erh *ExpectedResponseHandler) GetNextPrompt(r Response) Prompt {
// 	return erh.expectedResponseMap[r.GetId()]
// }

// type Response interface {
// 	GetText() string
// 	GetId() string
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