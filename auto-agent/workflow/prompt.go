package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"db"

	"appengine"
)

type Prompt interface {
	GetPhaseId() string
	GetResponseId() string
	GetResponseText() string
	GetNextPrompt() Prompt
	GetPromptId() string
	GetUIPrompt(*UIUserData) UIPrompt
	GetUIAction() UIAction
	ProcessResponse(string, *db.User, *UIUserData, appengine.Context)
	initUIPromptDynamicText(*UIUserData, Response)
}

type GenericPrompt struct {
	response                Response
	expectedResponseHandler *ExpectedResponseHandler
	currentUIPrompt         UIPrompt
	currentUIAction         UIAction
	promptConfig            *PromptConfig
	nextPrompt              Prompt
	currentPrompt           Prompt
	promptDynamicText       *UIPromptDynamicText
	state                   StateEntities
}

func (cp *GenericPrompt) GetPhaseId() string {
	return cp.promptConfig.PhaseId
}

func (cp *GenericPrompt) GetResponseId() string {
	return cp.response.GetResponseId()
}

func (cp *GenericPrompt) GetResponseText() string {
	return cp.response.GetResponseText()
}

func (cp *GenericPrompt) GetNextPrompt() Prompt {
	return cp.nextPrompt
}
func (cp *GenericPrompt) GetPromptId() string {
	return cp.promptConfig.Id
}

func (cp *GenericPrompt) GetUIPrompt(uiUserData *UIUserData) UIPrompt {
	if cp.currentUIPrompt == nil {
		pc := cp.promptConfig
		cp.currentUIPrompt = NewUIBasicPrompt()
		cp.currentUIPrompt.setPromptType(pc.PromptType)
		cp.currentPrompt.initUIPromptDynamicText(uiUserData, nil)
		if cp.promptDynamicText != nil {
			cp.currentUIPrompt.setText(cp.promptDynamicText.String())
		}
		cp.currentUIPrompt.setId(pc.Id)
		options := make([]*UIOption, len(pc.ExpectedResponses))
		for i := range pc.ExpectedResponses {
			options[i] = &UIOption{pc.ExpectedResponses[i].Id, pc.ExpectedResponses[i].Text}
		}
		cp.currentUIPrompt.setOptions(options)
	}
	return cp.currentUIPrompt
}

func (cp *GenericPrompt) processSimpleResponse(r string, u *db.User, uiUserData *UIUserData, c appengine.Context) {
	if r != "" {
		dec := json.NewDecoder(strings.NewReader(r))
		for {
			var response SimpleResponse
			if err := dec.Decode(&response); err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				log.Fatal(err)
				return
			}
			cp.response = &response
		}
		if cp.response != nil {
			cp.nextPrompt = cp.expectedResponseHandler.getNextPrompt(cp.response.GetResponseId())
			cp.nextPrompt.initUIPromptDynamicText(uiUserData, cp.response)
		}
	}
}

type Response interface {
	GetResponseText() string
	GetResponseId() string
}

type ExpectedResponseHandler struct {
	expectedResponseMap map[string]*PromptConfigRef
	currentPromptConfig *PromptConfig
}

func MakeExpectedResponseHandler(p *PromptConfig) *ExpectedResponseHandler {
	erh := new(ExpectedResponseHandler)
	erh.expectedResponseMap = make(map[string]*PromptConfigRef)
	erh.currentPromptConfig = p

	ecs := p.ExpectedResponses
	phaseId := p.PhaseId
	var promptId string

	for _, v := range ecs {
		if v.NextPromptRef != nil {
			promptId = v.NextPromptRef.Id
			if v.NextPromptRef.PhaseId != "" {
				phaseId = v.NextPromptRef.PhaseId
			}
		} else {
			promptId = v.NextPrompt.Id
			if v.NextPrompt.PhaseId != "" {
				phaseId = v.NextPrompt.PhaseId
			}
		}
		erh.expectedResponseMap[strings.ToLower(v.Id)] = &PromptConfigRef{Id: promptId, PhaseId: phaseId}
	}
	return erh
}

// Return the next prompt that maps to the expected response
// If there is only one expected response, return that one regardless of the response id
func (erh *ExpectedResponseHandler) getNextPrompt(rid string) Prompt {
	var p *PromptConfigRef
	if erh.currentPromptConfig.ResponseType == RESPONSE_END {
		p = GetFirstPromptInNextSequence(erh.currentPromptConfig)
	} else {
		if len(erh.expectedResponseMap) == 1 {
			for _, v := range erh.expectedResponseMap {
				p = v
			}
		} else {
			p = erh.expectedResponseMap[strings.ToLower(rid)]
		}
	}
	return MakePrompt(p.Id, p.PhaseId)
}

type SimpleResponse struct {
	Text string
	Id   string
}

func (sr *SimpleResponse) GetResponseText() string {
	if sr.Text != RESPONSE_SYSTEM_GENERATED {
		return sr.Text
	}
	return ""
}

func (sr *SimpleResponse) GetResponseId() string {
	return sr.Id
}

type UIPromptDynamicText struct {
	previousResponse Response
	promptConfig     *PromptConfig
	state            StateEntities
}

func (ps *UIPromptDynamicText) String() string {
	t := template.Must(template.New("display").Parse(ps.promptConfig.Text))
	var doc bytes.Buffer
	err := t.Execute(&doc, ps.state)
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
	setOptions([]*UIOption)
	Display() string
	GetId() string
}

type UIBasicPrompt struct {
	PromptType string
	Text       string
	PromptId   string
	Options    []*UIOption
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

func (ps *UIBasicPrompt) setOptions(options []*UIOption) {
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
	Factors        []*UIFactor
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
	Levels   []*UIFactorOption
}

type UIFactorOption struct {
	FactorLevelId string
	Text          string
	ImgPath       string
}
