package workflow

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	// "strings"
)

const (
	PHASE_COV        = "Cov"
	PHASE_CHART      = "Chart"
	PHASE_PREDICTION = "Prediction"
	FIRST_PHASE      = "START"
	LAST_PHASE       = "END"

	COV_RESPONSE_ID_SINGLE_CASE        = "Single record"
	COV_RESPONSE_ID_NON_VARYING        = "Two records non-varying"
	COV_RESPONSE_ID_TARGET_NON_VARYING = "Target factor non-varying"
	COV_RESPONSE_ID_UNCONTROLLED       = "Two records uncontrolled"
	COV_RESPONSE_ID_CONTROLLED         = "Two records controlled"

	UI_PROMPT_TEXT     = "Text"
	UI_PROMPT_MC       = "MC"
	UI_PROMPT_NO_INPUT = "NO_INPUT"
	// UI_PROMPT_SELECT_FACTOR = "SELECT_TARGET_FACTOR"

	RESPONSE_BASIC                = "Basic"
	RESPONSE_END                  = "COMPLETE"
	RESPONSE_RECORD               = "RECORD"
	RESPONSE_SELECT_TARGET_FACTOR = "SELECT_TARGET_FACTOR"

	UIACTION_INACTIVE = "NO_UIACTION"
	// ***TODO MUST FIX!!! server cannot be shut down when json is mulformed
	// PhaseConfig->PromptConfig->ExpectedReponseConfig

	// DOC
	// Configuration Rules:
	// 1. id must be unique and are treated as case insensitive.
	// 2. reference to an already defined prompt by specifying the id only, otherwise, the last definition is used
	// 3. if there is only one expected response, ExpectedResponses, it becomes the default next prompt
	promptTreeJsonFile = "workflow.json"
)

type ExpectedResponseConfig struct {
	Id         string
	Text       string
	NextPrompt *PromptConfig
}

type PromptConfig struct {
	Id                string
	Text              string
	UIActionModeId    string
	PromptType        string
	ResponseType      string
	ExpectedResponses []ExpectedResponseConfig
}

type PhaseConfig struct {
	Id              string
	FirstPrompt     PromptConfig
	PreviousPhaseId string
	NextPhaseId     string
}

type Level struct {
	Name    string
	Id      string
	ImgPath string
}

type Factor struct {
	Name    string
	Id      string
	ImgPath string
	Levels  []Level
}

type ContentConfig struct {
	RecordFileName   string
	RecordSize       int
	CausalFactors    []Factor
	NonCausalFactors []Factor
	OutcomeVariable  Factor
}

type AppConfig struct {
	Phase   PhaseConfig
	Content ContentConfig
}

var phaseConfigMap = make(map[string]PhaseConfig)
var promptConfigMap = make(map[string]*PromptConfig) //key:PhaseConfig.Id+PromptConfig.Id
var contentConfig ContentConfig

func InitWorkflow() {
	f, err := os.Open(promptTreeJsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		log.Fatal(err)
	}
	dec := json.NewDecoder(bufio.NewReader(f))

	for {
		var appConfig AppConfig
		if err := dec.Decode(&appConfig); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			log.Fatal(err)
			return
		}
		phaseConfig := appConfig.Phase
		//TODO cleanup
		//fmt.Fprintf(os.Stderr, " %s: %s\n", promptTree.Id, (promptTree.ExpectedResponses[0]).Id)
		phaseConfigMap[phaseConfig.Id] = phaseConfig
		populatePromptConfigMap(&phaseConfig.FirstPrompt, phaseConfig.Id)

		contentConfig = appConfig.Content
	}
}

func populatePromptConfigMap(pc *PromptConfig, phaseId string) {
	if promptConfigMap[phaseId+pc.Id] == nil {
		promptConfigMap[phaseId+pc.Id] = pc
		for i := range pc.ExpectedResponses {
			populatePromptConfigMap(pc.ExpectedResponses[i].NextPrompt, phaseId)
		}
	}
}

func MakeFirstPrompt() Prompt {
	// TODO Hardcoding the first prompt as CovPrompt
	p := MakeCovPrompt(phaseConfigMap[PHASE_COV].FirstPrompt)
	return p
}

func MakePrompt(promptId string, phaseId string) Prompt {
	pc := GetPromptConfig(promptId, phaseId)
	return MakePromptFromConfig(pc, phaseId)
}

func MakePromptFromConfig(pc *PromptConfig, phaseId string) Prompt {
	switch phaseId {
	case PHASE_COV:
		return MakeCovPrompt(*pc)
	}
	return nil
}

func GetPromptConfig(promptId string, phaseId string) *PromptConfig {
	return promptConfigMap[phaseId+promptId]
}

func GetContentConfig() ContentConfig {
	return contentConfig
}
