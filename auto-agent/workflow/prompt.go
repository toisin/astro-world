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
	GetResponseId() string
	GetResponseText() string
	GetNextPrompt() Prompt
	// SetResponse(Response)
	GetPromptId() string
	GetUIPrompt() UIPrompt
	GetUIAction() UIAction
	ProcessResponse(string, *UIUserData)
}

type Response interface {
	GetResponseText() string
	GetResponseId() string
}

type ExpectedResponseHandler struct {
	expectedResponseMap map[string]Prompt
}

func MakeExpectedResponseHandler(ecs []ExpectedResponseConfig, phaseId string) *ExpectedResponseHandler {
	erh := new(ExpectedResponseHandler)
	erh.expectedResponseMap = make(map[string]Prompt)
	for _, v := range ecs {
		// In case if the prompt is already defined
		// erh.expectedResponseMap[strings.ToLower(v.Id)] = MakePromptFromConfig(v.NextPrompt, phaseId)
		erh.expectedResponseMap[strings.ToLower(v.Id)] = MakePrompt(v.NextPrompt.Id, phaseId)
	}
	return erh
}

// Return the next prompt that maps to the expected response
// If there is only one expected response, return that one regardless of the response id
func (erh *ExpectedResponseHandler) GetNextPrompt(rid string) Prompt {
	if len(erh.expectedResponseMap) == 1 {
		for _, v := range erh.expectedResponseMap {
			return v
		}
	}
	return erh.expectedResponseMap[strings.ToLower(rid)]
}

type UIPrompt interface {
	SetText(string)
	SetPromptType(string)
	SetId(string)
	SetOptions([]UIOption)
	Display() string
	GetId() string
}

type UIBasicPrompt struct {
	PromptType string
	Text       string
	PromptId   string
	Options    []UIOption
}

func NewUIBasicPrompt() *UIBasicPrompt {
	return &UIBasicPrompt{}
}

func (ps *UIBasicPrompt) SetText(s string) {
	ps.Text = s
}

func (ps *UIBasicPrompt) SetPromptType(s string) {
	ps.PromptType = s
}

func (ps *UIBasicPrompt) SetId(s string) {
	ps.PromptId = s
}

func (ps *UIBasicPrompt) SetOptions(options []UIOption) {
	ps.Options = options
}

func (ps *UIBasicPrompt) GetId() string {
	return ps.PromptId
}

func (ps *UIBasicPrompt) Display() string {
	return ps.Text
}

type UIOption struct {
	ResponseId string
	Text       string
}

type UIAction interface {
	SetUIActionModeId(string)
}

type UIBasicAction struct {
	UIActionModeId string
}

func NewUIBasicAction() *UIBasicAction {
	return &UIBasicAction{}
}

func (ps *UIBasicAction) SetUIActionModeId(s string) {
	ps.UIActionModeId = s
}

type UIRecordAction struct {
	UIActionModeId string
	Factors        []UIFactor
}

func NewUIRecordAction() *UIRecordAction {
	return &UIRecordAction{}
}

func (ps *UIRecordAction) SetUIActionModeId(s string) {
	ps.UIActionModeId = s
}

type UIFactor struct {
	FactorId string
	Text     string
	Levels   []UIFactorOption
}

type UIFactorOption struct {
	FactorLevelId string
	Text          string
	ImgPath       string
}
