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
	ActionModeId int
	ExpectedResponses []ExpectedResponseConfig
}

var promptTreeConfig PromptConfig

func InitCovPromptLogic() {
	dec := json.NewDecoder(strings.NewReader(promptTreeJsonStream))
	for {
		if err := dec.Decode(&promptTreeConfig); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			log.Fatal(err)
		}
		//TODO cleanup
		//fmt.Fprintf(os.Stderr, " %s: %s\n", promptTree.Id, (promptTree.ExpectedResponses[0]).Id)
	}
	MakeFirstPrompt()
}

func MakeFirstPrompt()  {
	p:= MakeCovPrompt(promptTreeConfig)
	fmt.Fprintf(os.Stderr, " %s", p.GetDisplayText())
}



type State struct {
	xxx []string
	// GetActivePhase() Phase
	// GetCompletedPhases() []Phase
	// GetFuturePhases() []Phase
	// GetHistory() []Prompt
}




const (
	UI_PROMPT_NO_RESPONSE = "NO_RESPONSE"
	UI_PROMPT_TEXT = "TEXT"
	UI_PROMPT_YES_NO = "YES_NO"
	UI_PROMPT_MC = "MC"

	UI_PROMPT_END = "COMPLETE"
)

var stateMap = make(map[string]UIState)
var variableMap = make(map[string]Variable)

func InitWorkflow() {
	InitCovPromptLogic() // TODO only do cov for now
	variableMap["Y"] = Variable{ Text: "Performance"}
	variableMap["X1"] = Variable{"Health Index", true, false}
	variableMap["X2"] = Variable{"Height", false, false}
	variableMap["X3"] = Variable{"Weight", true, false}
	variableMap["X4"] = Variable{"Laughters", true, false}

	// stateMap["1"] = &UIMCPromptState{UI_PROMPT_NO_RESPONSE, "1", "Let's get started!", ""}
	stateMap["1"] = &UIMCPromptState{
						UI_PROMPT_MC,
						"1",
						"Let's get started! What feature have you investigated?",
						"",
					 	[]UIOption{
					 		UIOption{variableMap["X1"].Text,"X1"},
					 		UIOption{variableMap["X2"].Text,"X2"},
					 		UIOption{variableMap["X3"].Text,"X3"},
					 		UIOption{variableMap["X4"].Text,"X4"}}}
	stateMap["2"] = &UITextPromptState{UI_PROMPT_YES_NO, "2", "Do you think it makes a difference?", "1", "3"}
	//stateMap["2"] = &LogicPromptState{UI_PROMPT_YES_NO, "2", "Do you think it makes a difference?", "1"}
	stateMap["3"] = &UITextPromptState{
						UI_PROMPT_TEXT,
						"3",
						"When %X1 goes up, what happens to %Y?",
						"2",
						"4"}
	stateMap["4"] = &UITextPromptState{UI_PROMPT_TEXT, "4", "What did you find out about %X1?", "3", "5"}
	stateMap["5"] = &UITextPromptState{UI_PROMPT_TEXT, "5", "How do you know?", "4", "6"}
	stateMap["6"] = &UITextPromptState{UI_PROMPT_TEXT, "6", "Which records show you are right?", "5", UI_PROMPT_END}
	stateMap[UI_PROMPT_END] = &UITextPromptState{UI_PROMPT_END, UI_PROMPT_END, "You have done!", "6", UI_PROMPT_END}
	// stateMap["8"] = &UIMCPromptState{"8", "What level is your?", ""}
	// stateMap["9"] = &UIMCPromptState{"9", "How do you know?", ""}
	// stateMap["10"] = &UIMCPromptState{"10", "How do you know?", ""}
	// stateMap["11"] = &UIMCPromptState{"11", "How do you know?", ""}
}

func GetVariableMap() map[string]Variable{
	return variableMap
}

func GetStateMap() map[string]UIState{
	return stateMap
}

func GetFirstState() UIState {
	return stateMap["1"]
}

type Variable struct {
	Text string
	IsCausal bool
	IsPostiveCorr bool
}


type UIState interface {
	Display() string
	GetNextStateId() string
	GetId() string
}

type UITextPromptState struct {
	Ptype string
	WorkflowStateID string
	Text string
	LastStateId string
	NextStateId string
}

func (ps *UITextPromptState) GetId() string {
	return ps.WorkflowStateID
}

func (ps *UITextPromptState) Display() string {
	return ps.Text
}

func (ps *UITextPromptState) GetNextStateId() string {
	return ps.NextStateId
}

type UIMCPromptState struct {
	Ptype string
	WorkflowStateID string
	Text string
	LastStateId string
	Options []UIOption
}

type UIOption struct {
	Label string
	Value string
}

func (ps *UIMCPromptState) GetId() string {
	return ps.WorkflowStateID
}

func (ps *UIMCPromptState) Display() string {
	return ps.Text
}

func (ps *UIMCPromptState) GetNextStateId() string {
	//TODO Totally just hardcoding
	if ps.WorkflowStateID == "1" {
		return "2"
	}
	return ""
}




