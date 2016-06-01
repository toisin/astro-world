package workflow

type CovPrompt struct {
	state State
	previousPrompt Prompt
	response Response
	expectedResponseHandler *ExpectedResponseHandler
	promptGenerator PromptGenerator
}

func MakeCovPrompt(p PromptConfig) *CovPrompt {
	pg := new(CovPromptGenerator)
	pg.promptID = p.Id
	pg.actionModeId = p.ActionModeId
	pg.promptText = p.Text

	erh := MakeExpectedResponseHandler(p.ExpectedResponses)

	n := new(CovPrompt)
	n.promptGenerator = pg
	n.expectedResponseHandler = erh
	return n
}

func (cp *CovPrompt) GetDisplayText() string {
	return cp.promptGenerator.GetPromptText()
}

func (cp *CovPrompt) GetActionModeId() int {
	return cp.promptGenerator.GetActionModeId()
}

func (cp *CovPrompt) GetState() State {
	return cp.state
}

func (cp *CovPrompt) SetResponse(r Response) {
	cp.response = r
}

func (cp *CovPrompt) GetNextPrompt() Prompt {
	return cp.expectedResponseHandler.GetNextPrompt(cp.response)
}

func (cp *CovPrompt) GetResponseText() string {
	return cp.response.GetText()
}



type CovPromptGenerator struct {
	promptID string
	actionModeId int // the mode of rendering for Action UI
	promptText string // text with place holders for dynamic data
	state State
	previousPrompt Prompt
}

func (cph *CovPromptGenerator) generatePromptText() string {
	// TODO
	// data = state.GetDynamicData() // Get needed dynamic data from state
	// text = generatePromptText(data)
	cph.promptText = ""
	return cph.promptText
}

func (cph *CovPromptGenerator) generateAction() int {
	// TODO
	// data = state.GetDynamicData() // Get needed dynamic data from state
	// text = generatePromptText(data)
	cph.actionModeId = 0
	return cph.actionModeId
}

func (cph *CovPromptGenerator) GetPromptText() string {
	if (cph.promptText == "") {
		cph.generatePromptText()
	}
	return cph.promptText
}

func (cph *CovPromptGenerator) GetActionModeId() int {
	if (cph.actionModeId == 0) {
		cph.generateAction()
	}
	return cph.actionModeId
}
