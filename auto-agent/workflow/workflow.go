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

type AppConfig struct {
	CovPhase        PhaseConfig
	ChartPhase      PhaseConfig `json:"omitempty"`
	PredictionPhase PhaseConfig `json:"omitempty"`
	Content         ContentConfig
}

type ContentConfig struct {
	RecordFileName   string
	RecordSize       int
	CausalFactors    []Factor
	NonCausalFactors []Factor
	OutcomeVariable  Factor
}

type Factor struct {
	Name    string
	Id      string
	ImgPath string
	Levels  []Level
	DBIndex int
}

type Level struct {
	Name    string
	Id      string
	ImgPath string
}

type PhaseConfig struct {
	Id              string
	FirstPrompt     PromptConfig
	PreviousPhaseId string
	NextPhaseId     string
	FactorsOrder    []string // ordered factor ids
}

type PromptConfig struct {
	Id                string
	Text              string
	UIActionModeId    string
	PromptType        string
	ResponseType      string
	ExpectedResponses []ExpectedResponseConfig
}

type ExpectedResponseConfig struct {
	Id         string
	Text       string
	NextPrompt *PromptConfig
}

// Implements workflow.StateEntities
type CovPhaseState struct {
	Username     string
	Screenname   string
	RecordNoOne  *RecordState
	RecordNoTwo  *RecordState
	TargetFactor *CovFactorState
}

func (c *CovPhaseState) GetPhaseId() string {
	return appConfig.CovPhase.Id
}

type RecordState struct {
	RecordName   string
	RecordNo     string
	FactorLevels map[string]*CovFactorState
	// Factor id as keys, such as:
	// "fitness",
	// "parentshealth",
	// "education",
	// "familysize"
}

// This type is used in multiple contexts.
// Not all members may be relevant.
type CovFactorState struct {
	FactorName    string
	FactorId      string
	SelectedLevel string
	OppositeLevel string
}

// var phaseConfigMap = make(map[string]PhaseConfig)
var promptConfigMap = make(map[string]*PromptConfig) //key:PhaseConfig.Id+PromptConfig.Id
var factorConfigMap = make(map[string]*Factor)       //key:PhaseConfig.Id+PromptConfig.Id
var contentConfig ContentConfig
var appConfig AppConfig

func InitWorkflow() {
	f, err := os.Open(promptTreeJsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		log.Fatal(err)
	}
	dec := json.NewDecoder(bufio.NewReader(f))

	for {
		if err := dec.Decode(&appConfig); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			log.Fatal(err)
			return
		}
		covPhaseConfig := appConfig.CovPhase
		populatePromptConfigMap(&covPhaseConfig.FirstPrompt, covPhaseConfig.Id)

		contentConfig = appConfig.Content
	}
	populateFactorConfigMap(&contentConfig)
}

func populatePromptConfigMap(pc *PromptConfig, phaseId string) {
	if promptConfigMap[phaseId+pc.Id] == nil {
		promptConfigMap[phaseId+pc.Id] = pc
		for i := range pc.ExpectedResponses {
			populatePromptConfigMap(pc.ExpectedResponses[i].NextPrompt, phaseId)
		}
	}
}

func populateFactorConfigMap(cf *ContentConfig) {
	for i := range cf.CausalFactors {
		factorConfigMap[cf.CausalFactors[i].Id] = &cf.CausalFactors[i]
	}
	for i := range cf.NonCausalFactors {
		factorConfigMap[cf.NonCausalFactors[i].Id] = &cf.NonCausalFactors[i]
	}
}

func MakeFirstPrompt() Prompt {
	// TODO Hardcoding the first prompt as CovPrompt
	p := MakeCovPrompt(appConfig.CovPhase.FirstPrompt)
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

func GetContentConfig() *ContentConfig {
	return &contentConfig
}

func GetFactorConfig(factorId string) *Factor {
	return factorConfigMap[factorId]
}

func CreateCovFactorState(factorId string, levelId string) *CovFactorState {
	f := GetFactorConfig(factorId)
	allLevels := f.Levels
	var selectedLevel string
	var oppositeLevel string
	// The opposite level of the given level:
	//  - If the given level is at index 0, return the level id of the highest index
	//  - Otherwise, return the level id of index 0
	for i, v := range allLevels {
		if v.Id == levelId {
			selectedLevel = v.Name
			if i == 0 {
				oppositeLevel = allLevels[len(allLevels)-1].Name
			} else {
				oppositeLevel = allLevels[0].Name
			}
		}
	}
	return &CovFactorState{
		FactorName:    f.Name,
		FactorId:      f.Id,
		SelectedLevel: selectedLevel,
		OppositeLevel: oppositeLevel}
}

func UnstringifyState(b []byte, phaseId string) (se StateEntities, err error) {
	if (b != nil) && (len(b)) > 0 {
		switch phaseId {
		case appConfig.CovPhase.Id:
			var cps CovPhaseState
			err = unstringify(b, &cps)
			se = &cps
		}
	}
	return
}

func unstringify(b []byte, v interface{}) (err error) {
	err = json.Unmarshal(b, v)
	return
}
