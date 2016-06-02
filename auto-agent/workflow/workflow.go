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
	PHASE_COV = "Cov"
	PHASE_CHART = "Chart"
	PHASE_PREDICTION = "Prediction"
	FIRST_PHASE = "START"
	LAST_PHASE = "END"

	UI_PROMPT_NO_RESPONSE = "NO_RESPONSE"
	UI_PROMPT_TEXT = "TEXT"
	UI_PROMPT_YES_NO = "YES_NO"
	UI_PROMPT_MC = "MC"
	UI_PROMPT_END = "COMPLETE"

	UIACTION_INACTIVE = "NO_UIACTION"
	// ***TODO MUST FIX!!! server cannot be shut down when json is mulformed
	// PhaseConfig->PromptConfig->ExpectedReponseConfig
	promptTreeJsonStream = `
	{
		"Id": "Cov",
		"PreviousPhaseId": "START",
		"NextPhaseId": "Chart",
		"FirstPrompt":
		{
			"Id": "p1",
			"Text": "How many records would you like to see? One or Two?",
			"Type": "MC",
			"UIActionModeId": "NO_UIACTION",
			"ExpectedResponses": 
			[
				{
					"Id": "p1r1",
					"Text": "One",
					"NextPrompt":
					{
						"Id": "p1r1p1",
						"Text": "Which one would you like to see?",
						"Type": "TEXT",
						"UIActionModeId": "RECORD_SELECT_ONE",
						"ExpectedResponses": []
					}
				},
				{
					"Id": "p1r2",
					"Text": "Two",
					"NextPrompt":
					{
						"Id": "p1r2p1",
						"Text": "Which records would you like to see?",
						"Type": "TEXT",
						"UIActionModeId": "RECORD_SELECT_TWO",
						"ExpectedResponses": []
					}
				}
			]
		}
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
	UIActionModeId string
	Type string
	ExpectedResponses []ExpectedResponseConfig
}

type PhaseConfig struct {
	Id string
	FirstPrompt PromptConfig
	PreviousPhaseId string
	NextPhaseId string
}

var phaseConfigMap = make(map[string]PhaseConfig)
var promptConfigMap = make(map[string]*PromptConfig) //key:PhaseConfig.Id+PromptConfig.Id

func InitCovPromptLogic() {
	dec := json.NewDecoder(strings.NewReader(promptTreeJsonStream))
	for {
		var phaseConfig PhaseConfig
		if err := dec.Decode(&phaseConfig); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			log.Fatal(err)
			return
		}
		//TODO cleanup
		//fmt.Fprintf(os.Stderr, " %s: %s\n", promptTree.Id, (promptTree.ExpectedResponses[0]).Id)
		phaseConfigMap[phaseConfig.Id]=phaseConfig
		populatePromptConfigMap(&phaseConfig.FirstPrompt, phaseConfig.Id)		
	}
	//TODO cleanup
	// for k,_ := range promptConfigMap {
	// 	fmt.Fprintf(os.Stderr, " %s: %s\n", k, promptConfigMap[k].Text)
	// }
}

func populatePromptConfigMap(pc *PromptConfig, phaseId string) {
	promptConfigMap[phaseId+pc.Id] = pc
	for i := range pc.ExpectedResponses {
		populatePromptConfigMap(&(pc.ExpectedResponses[i].NextPrompt), phaseId)
	}
}

func MakeFirstPrompt() Prompt {
	// TODO Hardcoding the first prompt as CovPrompt
	p:= MakeCovPrompt(phaseConfigMap[PHASE_COV].FirstPrompt)
	// fmt.Fprintf(os.Stderr, " %s", p.GetDisplayText())
	return p
}

func MakePrompt(pId string, phaseId string) Prompt {
	pc := GetPromptConfig(pId, phaseId)
	p:= MakeCovPrompt(*pc)
	// fmt.Fprintf(os.Stderr, " %s", p.GetDisplayText())
	return p
}

func GetPromptConfig(pId string, phaseId string) *PromptConfig {
	return promptConfigMap[phaseId+pId]
}

// type State struct {
// 	xxx []string
// 	// GetActivePhase() Phase
// 	// GetCompletedPhases() []Phase
// 	// GetFuturePhases() []Phase
// 	// GetHistory() []Prompt
// }




// var stateMap = make(map[string]UIPrompt)
// var variableMap = make(map[string]Variable)

func InitWorkflow() {
	InitCovPromptLogic() // TODO only do cov for now
	// variableMap["Y"] = Variable{ Text: "Performance"}
	// variableMap["X1"] = Variable{"Health Index", true, false}
	// variableMap["X2"] = Variable{"Height", false, false}
	// variableMap["X3"] = Variable{"Weight", true, false}
	// variableMap["X4"] = Variable{"Laughters", true, false}

	// stateMap["1"] = &UIMCPrompt{UI_PROMPT_NO_RESPONSE, "1", "Let's get started!", ""}
	// stateMap["1"] = &UIMCPrompt{
	// 					UI_PROMPT_MC,
	// 					"1",
	// 					"Let's get started! What feature have you investigated?",
	// 					"",
	// 				 	[]UIOption{
	// 				 		UIOption{variableMap["X1"].Text,"X1"},
	// 				 		UIOption{variableMap["X2"].Text,"X2"},
	// 				 		UIOption{variableMap["X3"].Text,"X3"},
	// 				 		UIOption{variableMap["X4"].Text,"X4"}},
	// 				 	"p1"}
	// stateMap["2"] = &UITextPrompt{UI_PROMPT_YES_NO, "2", "Do you think it makes a difference?", "1", "3", "p1"}
	// //stateMap["2"] = &LogicPromptState{UI_PROMPT_YES_NO, "2", "Do you think it makes a difference?", "1"}
	// stateMap["3"] = &UITextPrompt{
	// 					UI_PROMPT_TEXT,
	// 					"3",
	// 					"When %X1 goes up, what happens to %Y?",
	// 					"2",
	// 					"4", "p1"}
	// stateMap["4"] = &UITextPrompt{UI_PROMPT_TEXT, "4", "What did you find out about %X1?", "3", "5", "p1"}
	// stateMap["5"] = &UITextPrompt{UI_PROMPT_TEXT, "5", "How do you know?", "4", "6", "p1"}
	// stateMap["6"] = &UITextPrompt{UI_PROMPT_TEXT, "6", "Which records show you are right?", "5", UI_PROMPT_END,"p1"}
	// stateMap[UI_PROMPT_END] = &UITextPrompt{UI_PROMPT_END, UI_PROMPT_END, "You have done!", "6", UI_PROMPT_END,"p1"}
	// stateMap["8"] = &UIMCPrompt{"8", "What level is your?", ""}
	// stateMap["9"] = &UIMCPrompt{"9", "How do you know?", ""}
	// stateMap["10"] = &UIMCPrompt{"10", "How do you know?", ""}
	// stateMap["11"] = &UIMCPrompt{"11", "How do you know?", ""}
}

// func GetVariableMap() map[string]Variable{
// 	return variableMap
// }

// func GetStateMap() map[string]UIPrompt{
// 	return stateMap
// }

// type Variable struct {
// 	Text string
// 	IsCausal bool
// 	IsPostiveCorr bool
// }


type UIPrompt interface {
	Display() string
	// GetNextStateId() string
	GetId() string
}

type UITextPrompt struct {
	Type string
	// WorkflowStateID string
	Text string
	LastStateId string
	NextStateId string
	PromptId string
}

func NewUITextPrompt() *UITextPrompt {
    return &UITextPrompt{Type:UI_PROMPT_TEXT}
}

func (ps *UITextPrompt) GetId() string {
	return ps.PromptId
}

func (ps *UITextPrompt) Display() string {
	return ps.Text
}

// func (ps *UITextPrompt) GetNextStateId() string {
// 	return ps.NextStateId
// }

type UIMCPrompt struct {
	Type string
	// WorkflowStateID string
	Text string
	LastStateId string
	Options []UIOption
	PromptId string
}

func NewUIMCPrompt() *UIMCPrompt {
    return &UIMCPrompt{Type:UI_PROMPT_MC}
}

type UIOption struct {
	Label string
	Value string
}

func (ps *UIMCPrompt) GetId() string {
	return ps.PromptId
	// return ps.WorkflowStateID
}

func (ps *UIMCPrompt) Display() string {
	return ps.Text
}

// func (ps *UIMCPrompt) GetNextStateId() string {
// 	//TODO Totally just hardcoding
// 	if ps.WorkflowStateID == "1" {
// 		return "2"
// 	}
// 	return ""
// }




