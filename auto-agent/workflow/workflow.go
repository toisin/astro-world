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

	UI_PROMPT_TEXT             = "Text"
	UI_PROMPT_MC               = "MC"
	UI_PROMPT_NO_INPUT         = "NO_INPUT"         // No input from dialog but expect input from action screen
	UI_PROMPT_STRAIGHT_THROUGH = "STRAIGHT_THROUGH" // Differ from NO_INPUT, no input expected, goes to next prompt directly
	// UI_PROMPT_SELECT_FACTOR = "SELECT_TARGET_FACTOR"

	RESPONSE_BASIC                = "Basic"
	RESPONSE_END                  = "COMPLETE"
	RESPONSE_RECORD               = "RECORD"
	RESPONSE_SELECT_TARGET_FACTOR = "SELECT_TARGET_FACTOR"
	RESPONSE_PRIOR_BELIEF_FACTORS = "PRIOR_BELIEF_FACTORS"
	RESPONSE_PRIOR_BELIEF_LEVELS  = "PRIOR_BELIEF_LEVELS"
	RESPONSE_SYSTEM_GENERATED     = "SYSTEM_GENERATED" // For when a submit is triggered by the system

	EXPECTED_SPECIAL_CONTENT_REF = "CONTENT_REF"

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

type GenericState struct {
	PhaseId            string
	Username           string
	Screenname         string
	TargetFactor       FactorState
	RemainingFactorIds []string
	Beliefs            BeliefsState
}

type BeliefsState struct {
	HasCausalFactors         bool
	CausalFactors            []string
	HasMultipleCausalFactors bool
}

func (c *GenericState) setPhaseId(s string) {
	c.PhaseId = s
}

func (c *GenericState) setUsername(s string) {
	c.Username = s
}

func (c *GenericState) setScreenname(s string) {
	c.Screenname = s
}

// Not applicable to all phases
func (c *GenericState) setTargetFactor(t FactorState) {
	c.TargetFactor = t
}

// Not applicable to all phases
func (c *GenericState) updateRemainingFactors() {
	factorId := c.TargetFactor.FactorId
	if c.RemainingFactorIds != nil {
		for i, v := range c.RemainingFactorIds {
			if v == factorId {
				c.RemainingFactorIds = append(c.RemainingFactorIds[:i], c.RemainingFactorIds[i+1:]...)
				break
			}
		}
	}
}

// Not applicable to all phases
func (c *GenericState) getRemainingFactorIds() []string {
	return c.RemainingFactorIds
}

// Implements workflow.StateEntities
type CovPhaseState struct {
	GenericState
	RecordNoOne *RecordState
	RecordNoTwo *RecordState
}

func (c *CovPhaseState) GetPhaseId() string {
	return appConfig.CovPhase.Id
}

func (c *CovPhaseState) initContents(factors []Factor) {
	c.RemainingFactorIds = make([]string, len(factors))
	for i, v := range factors {
		c.RemainingFactorIds[i] = v.Id
	}
}

func (cp *CovPhaseState) isContentCompleted() bool {
	if len(cp.RemainingFactorIds) > 0 {
		return false
	}
	return true
}

type RecordState struct {
	RecordName   string
	RecordNo     string
	FactorLevels map[string]*FactorState
	Performance  string
	// Factor id as keys, such as:
	// "fitness",
	// "parentshealth",
	// "education",
	// "familysize"
}

// This type is used in multiple contexts.
// Not all members may be relevant.
type FactorState struct {
	FactorName    string
	FactorId      string
	SelectedLevel string
	OppositeLevel string
	IsCausal      bool
}

// Implements workflow.StateEntities
type ChartPhaseState struct {
	GenericState
}

func (c *ChartPhaseState) GetPhaseId() string {
	return appConfig.ChartPhase.Id
}

func (cp *ChartPhaseState) isContentCompleted() bool {
	if len(cp.RemainingFactorIds) > 0 {
		return false
	}
	return true
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

func CreateCovFactorState(factorId string, levelId string) *FactorState {
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
	return &FactorState{
		FactorName:    f.Name,
		FactorId:      f.Id,
		SelectedLevel: selectedLevel,
		OppositeLevel: oppositeLevel,
		IsCausal:      f.IsCausal}
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
