package workflow

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"db"

	"appengine"
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
	GetUIPrompt(*UIUserData) UIPrompt
	GetUIAction() UIAction
	ProcessResponse(string, *db.User, *UIUserData, appengine.Context)
	initUIPromptDynamicText(*UIUserData, *Response)
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
func (erh *ExpectedResponseHandler) getNextPrompt(rid string) Prompt {
	if len(erh.expectedResponseMap) == 1 {
		for _, v := range erh.expectedResponseMap {
			return v
		}
	}
	return erh.expectedResponseMap[strings.ToLower(rid)]
}

type UIPromptDynamicText interface {
	String() string
}

func generateDynamicText(ttext string, state StateEntities) string {
	t := template.Must(template.New("display").Parse(ttext))
	var doc bytes.Buffer
	err := t.Execute(&doc, state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %s\n\n", err)
		log.Println("executing template:", err)
	}
	display := doc.String()

	return display
}

type UIPrompt interface {
	setText(string)
	setPromptType(string)
	setId(string)
	setOptions([]UIOption)
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

func (ps *UIBasicPrompt) setText(s string) {
	ps.Text = s
}

func (ps *UIBasicPrompt) setPromptType(s string) {
	ps.PromptType = s
}

func (ps *UIBasicPrompt) setId(s string) {
	ps.PromptId = s
}

func (ps *UIBasicPrompt) setOptions(options []UIOption) {
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
	setUIActionModeId(string)
}

type UIBasicAction struct {
	UIActionModeId string
}

func NewUIBasicAction() *UIBasicAction {
	return &UIBasicAction{}
}

func (ps *UIBasicAction) setUIActionModeId(s string) {
	ps.UIActionModeId = s
}

type UIRecordAction struct {
	UIActionModeId string
	Factors        []UIFactor
}

func NewUIRecordAction() *UIRecordAction {
	return &UIRecordAction{}
}

func (ps *UIRecordAction) setUIActionModeId(s string) {
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
