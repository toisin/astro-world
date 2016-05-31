package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"os"
)

const (
	// // TODO for now use these as respolnse id
	// P1 = `How many records would you like to see? One or Two?`
	// R1_P1 = `one`
	// R2 = `two`
	// P_2 = `Which one would you like to see?`
	// P_3 = `Can you tell me what you would be able to figure out by looking at this record?`
	promptTreeJsonStream = `
	{
		"Id": "p1",
		"Text": "How many records would you like to see? One or Two?",
		"ExpectedResponses": 
		[
			{
				"Id": "p1r1",
				"Text": "one",
				"NextPrompt":
				{
					"Label": "p1r1p1",
					"Text": "Which one would you like to see?",
					"ExpectedResponses": []
				}
			},
			{
				"Id": "p1r2",
				"Text": "two",
				"NextPrompt":
				{
					"Label": "p1r2p1",
					"Text": "Which recoards would you like to see?",
					"ExpectedResponses": []
				}
			}
		]
	}`
)

type ExpectedResponseConfig struct {
	Id string
	Text string
	NextPrompt PromptConfig
}

type PromptConfig struct {
	Id string
	Text string
	ExpectedResponses []ExpectedResponseConfig
}

var promptTree PromptConfig

func InitCovPromptLogic() {
	dec := json.NewDecoder(strings.NewReader(promptTreeJsonStream))
	for {
		if err := dec.Decode(&promptTree); err == io.EOF {
			fmt.Fprintf(os.Stderr, "%s", err)
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, " %s: %s\n", promptTree.Id, (promptTree.ExpectedResponses[0]).Id)
	}
}

type CovPrompt struct {
	state State
	previousPrompt Prompt
	response Response
	expectedResponseHandler ExpectedResponseHandler
	promptGenerator PromptGenerator
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
