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

	UI_PROMPT_ENTER_TO_CONTINUE = "ENTER_TO_CONTINUE"
	UI_PROMPT_TEXT              = "Text"
	UI_PROMPT_MC                = "MC"
	UI_PROMPT_NO_INPUT          = "NO_INPUT"         // No input from dialog but expect input from action screen
	UI_PROMPT_STRAIGHT_THROUGH  = "STRAIGHT_THROUGH" // Differ from NO_INPUT, no input expected, goes to next prompt directly
	// UI_PROMPT_SELECT_FACTOR = "SELECT_TARGET_FACTOR"

	RESPONSE_BASIC                = "Basic"
	RESPONSE_END                  = "COMPLETE"
	RESPONSE_RECORD               = "RECORD"
	RESPONSE_MEMO                 = "MEMO"
	RESPONSE_CAUSAL_CONCLUSION    = "CAUSAL_CONCLUSION"
	RESPONSE_SELECT_TARGET_FACTOR = "SELECT_TARGET_FACTOR"
	RESPONSE_PRIOR_BELIEF_FACTORS = "PRIOR_BELIEF_FACTORS"
	RESPONSE_PRIOR_BELIEF_LEVELS  = "PRIOR_BELIEF_LEVELS"
	RESPONSE_SYSTEM_GENERATED     = "SYSTEM_GENERATED" // For when a submit is triggered by the system

	EXPECTED_SPECIAL_CONTENT_REF = "CONTENT_REF"
	EXPECTED_ANY_RESPONSE        = "ANY_RESPONSE"

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
	ChartPhase      PhaseConfig
	PredictionPhase PhaseConfig
	Content         ContentConfig
}

type ContentConfig struct {
	RecordFileName  string
	RecordSize      int
	Factors         []Factor
	OutcomeVariable Factor
}

type Factor struct {
	Name     string
	Id       string
	ImgPath  string
	Levels   []Level
	DBIndex  int
	IsCausal bool
}

type Level struct {
	Name    string
	Id      string
	ImgPath string
}

type PhaseConfig struct {
	Id              string
	PreviousPhaseId string
	NextPhaseId     string
	// ContentRef: Type of content being referenced is not configurable
	// each phase expects a specific type of content, e.g. CovPhase expects Factors
	ContentRef       ContentConfig
	OrderedSequences []Sequence
}

type Sequence struct {
	RepeatOverContent bool // If true, repeat content over the list of contents sepecify in ContentRef
	FirstPrompt       PromptConfig
}

type PromptConfig struct {
	Id                         string // Id must be unique within the phase
	PhaseId                    string
	Text                       []string
	UIActionModeId             string
	PromptType                 string
	ResponseType               string
	ExpectedResponses          ExpectedResponseConfig
	IsDynamicExpectedResponses bool
	sequenceOrder              int
}

type ExpectedResponseConfig struct {
	StateTemplateRef string // If StateTemplateRef is empty, treat value as is
	Values           []ExpectedResponseValue
}

type ExpectedResponseValue struct {
	Id            string
	Text          string
	NextPrompt    PromptConfig
	NextPromptRef PromptConfigRef
}

type PromptConfigRef struct {
	Id      string
	PhaseId string
}

// var phaseConfigMap = make(map[string]PhaseConfig)
var promptConfigMap = make(map[string]*PromptConfig) //key:PhaseConfig.Id+PromptConfig.Id
var factorConfigMap = make(map[string]Factor)        //key:PhaseConfig.Id+PromptConfig.Id
var contentConfig ContentConfig
var appConfig AppConfig

func GetPhase(phaseId string) *PhaseConfig {
	var currentPhase *PhaseConfig
	switch phaseId {
	case PHASE_COV:
		currentPhase = &appConfig.CovPhase
		break
	case PHASE_CHART:
		currentPhase = &appConfig.ChartPhase
		break
	}
	return currentPhase
}

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
		chartPhaseConfig := appConfig.ChartPhase
		for i, _ := range covPhaseConfig.OrderedSequences {
			populatePromptConfigMap(&covPhaseConfig.OrderedSequences[i].FirstPrompt, covPhaseConfig.Id, i)
		}
		for i, _ := range chartPhaseConfig.OrderedSequences {
			populatePromptConfigMap(&chartPhaseConfig.OrderedSequences[i].FirstPrompt, chartPhaseConfig.Id, i)
		}

		contentConfig = appConfig.Content
	}
	populateFactorConfigMap(&contentConfig)
}

func populatePromptConfigMap(pc *PromptConfig, phaseId string, sequenceOrder int) {
	if pc.PhaseId != "" {
		phaseId = pc.PhaseId
	} else {
		pc.PhaseId = phaseId
	}
	pc.sequenceOrder = sequenceOrder
	promptConfigMap[phaseId+pc.Id] = pc
	for i := range pc.ExpectedResponses.Values {
		populatePromptConfigMap(&pc.ExpectedResponses.Values[i].NextPrompt, phaseId, sequenceOrder)
	}
}

func GetOutcomeLevelOrder(levelName string) int {
	for i, v := range appConfig.Content.OutcomeVariable.Levels {
		if v.Name == levelName {
			return i
		}
	}
	return -1
}

func populateFactorConfigMap(cf *ContentConfig) {
	for i := range cf.Factors {
		factorConfigMap[cf.Factors[i].Id] = cf.Factors[i]
	}
}

func GetFirstPhase() *PhaseConfig {
	return &appConfig.CovPhase
}

func MakeFirstPrompt(uiUserData *UIUserData) Prompt {
	// Hardcoding the first prompt is the first prompt of CovPrompt
	p := MakePrompt(appConfig.CovPhase.OrderedSequences[0].FirstPrompt.Id, appConfig.CovPhase.Id, uiUserData)
	return p
}

func MakePrompt(promptId string, phaseId string, uiUserData *UIUserData) Prompt {
	pc := GetPromptConfig(promptId, phaseId)
	return MakePromptFromConfig(*pc, uiUserData)
}

func MakePromptFromConfig(pc PromptConfig, uiUserData *UIUserData) Prompt {
	phaseId := pc.PhaseId
	switch phaseId {
	case PHASE_COV:
		return MakeCovPrompt(pc, uiUserData)
	case PHASE_CHART:
		return MakeChartPrompt(pc, uiUserData)
	}
	return nil
}

func GetFirstPromptConfigInSequence(sequenceOrder int, phaseId string) *PromptConfig {
	phase := GetPhase(phaseId)
	// Call GetPromptConfig just in case if some kind of initialization happened when
	// promptConfigMap was populated
	return GetPromptConfig(phase.OrderedSequences[sequenceOrder].FirstPrompt.Id, phaseId)
}

func GetPromptConfig(promptId string, phaseId string) *PromptConfig {
	return promptConfigMap[phaseId+promptId]
}

func GetContentConfig() *ContentConfig {
	return &contentConfig
}

func GetFactorConfig(factorId string) Factor {
	return factorConfigMap[factorId]
}

func CreateCovFactorState(factorId string, selectedLevelId string) FactorState {
	f := GetFactorConfig(factorId)
	allLevels := f.Levels
	var selectedLevel string
	var oppositeLevel string
	var oppositeLevelId string
	// The opposite level of the given level:
	//  - If the given level is at index 0, return the level id of the highest index
	//  - Otherwise, return the level id of index 0
	for i, v := range allLevels {
		if v.Id == selectedLevelId {
			selectedLevel = v.Name
			if i == 0 {
				oppositeLevel = allLevels[len(allLevels)-1].Name
				oppositeLevelId = allLevels[len(allLevels)-1].Id
			} else {
				oppositeLevel = allLevels[0].Name
				oppositeLevelId = allLevels[0].Id
			}
		}
	}
	return FactorState{
		FactorName:      f.Name,
		FactorId:        f.Id,
		SelectedLevel:   selectedLevel,
		SelectedLevelId: selectedLevelId,
		OppositeLevel:   oppositeLevel,
		OppositeLevelId: oppositeLevelId,
		IsCausal:        f.IsCausal}
}

func UnstringifyState(b []byte, phaseId string) (se StateEntities, err error) {
	if (b != nil) && (len(b)) > 0 {
		switch phaseId {
		case appConfig.CovPhase.Id:
			var cps CovPhaseState
			err = unstringify(b, &cps)
			se = &cps
		case appConfig.ChartPhase.Id:
			var cps ChartPhaseState
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

func Stringify(v interface{}) (b []byte, err error) {
	b, err = json.Marshal(v)
	return
}
