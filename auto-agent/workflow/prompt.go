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
	GetSequenceOrder() int
	GetPhaseId() string
	GetResponseId() string
	GetResponseText() string
	GetNextPrompt() Prompt
	GetPromptId() string
	GetUIPrompt() UIPrompt
	GetUIAction() UIAction
	ProcessResponse(string, *db.User, *UIUserData, appengine.Context)
	initUIPromptDynamicText(*UIUserData, Response)
	initDynamicResponseUIPrompt(*UIUserData)
	initUIPrompt(uiUserData *UIUserData)
	initUIAction()
	updateState(*UIUserData)
}

type GenericPrompt struct {
	response                Response
	expectedResponseHandler ExpectedResponseHandler
	currentUIPrompt         UIPrompt
	currentUIAction         UIAction
	promptConfig            PromptConfig
	nextPrompt              Prompt
	currentPrompt           Prompt
	promptDynamicText       *UIPromptDynamicText
	state                   StateEntities
}

func (cp *GenericPrompt) GetPhaseId() string {
	return cp.promptConfig.PhaseId
}

func (cp *GenericPrompt) GetResponseId() string {
	if cp.response != nil {
		return cp.response.GetResponseId()
	}
	return ""
}

func (cp *GenericPrompt) GetResponseText() string {
	if cp.response != nil {
		return cp.response.GetResponseText()
	}
	return ""
}

func (cp *GenericPrompt) GetNextPrompt() Prompt {
	return cp.nextPrompt
}

func (cp *GenericPrompt) GetPromptId() string {
	return cp.promptConfig.Id
}

func (cp *GenericPrompt) GetSequenceOrder() int {
	return cp.promptConfig.sequenceOrder
}

// Returned UIAction may be nil if not action UI is needed
func (cp *GenericPrompt) GetUIAction() UIAction {
	return cp.currentUIAction
}

func (cp *GenericPrompt) initUIPrompt(uiUserData *UIUserData) {
	pc := cp.promptConfig
	if !pc.IsDynamicExpectedResponses {
		cp.currentUIPrompt = NewUIBasicPrompt()
		cp.currentUIPrompt.setPromptType(pc.PromptType)
		cp.currentPrompt.initUIPromptDynamicText(uiUserData, nil)
		cp.currentUIPrompt.setText(cp.promptDynamicText.String())
		cp.currentUIPrompt.setId(pc.Id)
		options := make([]*UIOption, len(pc.ExpectedResponses.Values))
		for i := range pc.ExpectedResponses.Values {
			options[i] = &UIOption{pc.ExpectedResponses.Values[i].Id, pc.ExpectedResponses.Values[i].Text}
		}
		cp.currentUIPrompt.setOptions(options)
	} else {
		cp.currentPrompt.initDynamicResponseUIPrompt(uiUserData)
	}
}

func (cp *GenericPrompt) GetUIPrompt() UIPrompt {
	return cp.currentUIPrompt
}

func (cp *GenericPrompt) init(p PromptConfig, uiUserData *UIUserData) {
	cp.promptConfig = p
	cp.expectedResponseHandler = cp.makeExpectedResponseHandler(p)
	// invoking the initialization methods in the "subclass"
	// in case if they have been overriden
	cp.currentPrompt.initUIPrompt(uiUserData)
	cp.currentPrompt.initUIAction()
}

func (cp *GenericPrompt) initUIAction() {
	if cp.currentUIAction == nil {
		cp.currentUIAction = NewUIBasicAction()
		cp.currentUIAction.setUIActionModeId(cp.promptConfig.UIActionModeId)
	}
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
			cp.nextPrompt = cp.expectedResponseHandler.generateNextPrompt(cp.response, uiUserData)
		}
	}
}

func (cp *GenericPrompt) generateFirstPromptInNextSequence(uiUserData *UIUserData) Prompt {
	phaseId := cp.promptConfig.PhaseId
	currentPhase := GetPhase(phaseId)

	var nextPromptId string
	var currentS *Sequence
	var nextS *Sequence
	sequenceOrder := cp.promptConfig.sequenceOrder
	currentS = &currentPhase.OrderedSequences[sequenceOrder]

	if currentS.RepeatOverContent {
		// Check if all content has been through the current sequence
		// if not, go to the next content, otherwise, repeat sequence for the remaining content
		if !cp.state.isContentCompleted() {
			nextS = currentS
		}
	}
	if nextS == nil {
		// Go to the next sequence within the same phase
		// If no next sequence, then go to the first sequence of the next phase
		sequenceOrder++
		if len(currentPhase.OrderedSequences) > sequenceOrder {
			nextS = &currentPhase.OrderedSequences[sequenceOrder]
		} else {
			phaseId = currentPhase.NextPhaseId
			nextS = &GetPhase(phaseId).OrderedSequences[0]
		}
	}
	nextPromptId = nextS.FirstPrompt.Id

	return MakePrompt(nextPromptId, phaseId, uiUserData)

}

func (cp *GenericPrompt) makeExpectedResponseHandler(p PromptConfig) ExpectedResponseHandler {
	var erh ExpectedResponseHandler
	if p.IsDynamicExpectedResponses {
		erh = &DynamicExpectedResponseHandler{}
	} else {
		erh = &StaticExpectedResponseHandler{}
	}
	erh.init(p)
	return erh
}

type Response interface {
	GetResponseText() string
	GetResponseId() string
}

type StaticExpectedResponseHandler struct {
	expectedResponseMap map[string]*PromptConfigRef
	currentPromptConfig PromptConfig
}

type DynamicExpectedResponseHandler struct {
	// member StaticExpectedResponseHandler not a pointer so that it is automatically instantiated
	// when DynamicExpectedResponseHandler is instantiated
	StaticExpectedResponseHandler
}

type ExpectedResponseHandler interface {
	generateNextPrompt(Response, *UIUserData) Prompt
	init(PromptConfig)
}

// For now only call "super" init
// May have more to add later
func (derh *DynamicExpectedResponseHandler) init(p PromptConfig) {
	derh.StaticExpectedResponseHandler.init(p)
}

func (erh *StaticExpectedResponseHandler) init(p PromptConfig) {
	erh.expectedResponseMap = make(map[string]*PromptConfigRef)
	erh.currentPromptConfig = p

	ecs := p.ExpectedResponses.Values
	phaseId := p.PhaseId
	var promptId string

	for _, v := range ecs {
		if v.NextPromptRef.Id != "" {
			promptId = v.NextPromptRef.Id
			if v.NextPromptRef.PhaseId != "" {
				phaseId = v.NextPromptRef.PhaseId
			}
		}
		// NextPromptRef and NextPrompt should not co-exist
		// If both were present, NextPrompt takes over
		if v.NextPrompt.Id != "" {
			promptId = v.NextPrompt.Id
			if v.NextPrompt.PhaseId != "" {
				phaseId = v.NextPrompt.PhaseId
			}
		}
		erh.expectedResponseMap[strings.ToLower(v.Id)] = &PromptConfigRef{Id: promptId, PhaseId: phaseId}
	}
}

// Return the next prompt that maps to the expected response
// If there is only one expected response, return that one regardless of the response id
func (erh *StaticExpectedResponseHandler) generateNextPrompt(r Response, uiUserData *UIUserData) Prompt {
	var rid string
	var p *PromptConfigRef
	if len(erh.expectedResponseMap) == 1 {
		// If there is only one expected response, use it regardless of the response
		for _, v := range erh.expectedResponseMap {
			p = v
		}
	} else {
		// If there are more than one expected responses, find the appropriate
		// next prompt based on the current response

		if erh.currentPromptConfig.ExpectedResponses.StateTemplateRef != "" {
			// If StateTemplateRef is provided, evaluate it by applying
			// StateEntities to find the matching expected response
			text := erh.currentPromptConfig.ExpectedResponses.StateTemplateRef
			t := template.Must(template.New("expectedResponses").Parse(text))
			var doc bytes.Buffer
			err := t.Execute(&doc, uiUserData.State)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing expectedResponses template: %s\n\n", err)
				log.Println("executing expectedResponses template:", err)
			}
			rid = doc.String()
		} else {
			// If StateTemplateRef is not provided, use the response id directly
			// to find the matching expected response
			rid = r.GetResponseId()
		}
		p = erh.expectedResponseMap[strings.ToLower(rid)]
	}
	nextPrompt := MakePrompt(p.Id, p.PhaseId, uiUserData)
	nextPrompt.initUIPromptDynamicText(uiUserData, r)

	return nextPrompt
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
	promptConfig     PromptConfig
	state            StateEntities
}

func (ps *UIPromptDynamicText) String() []string {
	display := make([]string, len(ps.promptConfig.Text))
	for i, v := range ps.promptConfig.Text {
		t := template.Must(template.New("display").Parse(v))
		var doc bytes.Buffer
		err := t.Execute(&doc, ps.state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing template: %s\n\n", err)
			log.Println("executing template:", err)
		}
		display[i] = doc.String()
	}
	return display
}

type UIPrompt interface {
	setText([]string)
	setPromptType(string)
	setId(string)
	setOptions([]*UIOption)
	Display() []string
	GetId() string
}

type UIBasicPrompt struct {
	PromptType string
	Texts      []string
	PromptId   string
	Options    []*UIOption
}

func NewUIBasicPrompt() *UIBasicPrompt {
	return &UIBasicPrompt{}
}

func (ps *UIBasicPrompt) setText(s []string) {
	ps.Texts = s
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

func (ps *UIBasicPrompt) Display() []string {
	return ps.Texts
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

// // TODO Remove - not needed, make always available through UIUserData
// type UIRecordAction struct {
// 	UIActionModeId string
// 	Factors        []*UIFactor
// }

// func NewUIRecordAction() *UIRecordAction {
// 	return &UIRecordAction{}
// }

// func (ps *UIRecordAction) setUIActionModeId(s string) {
// 	ps.UIActionModeId = s
// }
