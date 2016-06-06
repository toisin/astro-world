package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	PHASE_COV        = "Cov"
	PHASE_CHART      = "Chart"
	PHASE_PREDICTION = "Prediction"
	FIRST_PHASE      = "START"
	LAST_PHASE       = "END"

	UI_PROMPT_NO_RESPONSE = "NO_RESPONSE"
	UI_PROMPT_TEXT        = "TEXT"
	UI_PROMPT_YES_NO      = "YES_NO"
	UI_PROMPT_MC          = "MC"
	UI_PROMPT_RECORD      = "RECORD"
	UI_PROMPT_END         = "COMPLETE"

	UIACTION_INACTIVE = "NO_UIACTION"
	// ***TODO MUST FIX!!! server cannot be shut down when json is mulformed
	// PhaseConfig->PromptConfig->ExpectedReponseConfig

	// DOC
	// Configuration Rules:
	// 1. id must be unique and are treated as case insensitive.
	// 2. reference to an already defined prompt by specifying the id only, otherwise, the last definition is used
	// 3. if there is only one expected response, ExpectedResponses. Text might be omitted because it
	//    is likely to be open-ended
	promptTreeJsonStream = `
	{
		"Content":
		{
			"RecordFileName": "cases.csv",
			"RecordSize": 120,
			"CausalFactors":
			[
				{
					"Name": "Fitness",
					"Id": "fitness",
					"Levels":
					[
						{
							"Name": "Excellent",
							"Id": "excellent",
							"ImgPath": "excellent fitness.jpg"
						},
						{
							"Name": "Average",
							"Id": "average",
							"ImgPath": "average fitness.jpg"
						}
					]
				},
				{
					"Name": "Parents' Health",
					"Id": "parentshealth",
					"Levels":
					[
						{
							"Name": "Excellent",
							"Id": "excellent",
							"ImgPath": "excellent parents.jpg"
						},
						{
							"Name": "Fair",
							"Id": "fair",
							"ImgPath": "fair parents.jpg"
						}
					]
				},
				{
					"Name": "Education",
					"Id": "education",
					"Levels":
					[
						{
							"Name": "No College",
							"Id": "no college",
							"ImgPath": "no college.jpg"
						},
						{
							"Name": "Some College",
							"Id": "some college",
							"ImgPath": "some college.jpg"
						},
						{
							"Name": "College",
							"Id": "college",
							"ImgPath": "college.jpg"
						}
					]
				}
			],
			"NonCausalFactors":
			[
				{
					"Name": "Family Size",
					"Id": "familysize",
					"Levels":
					[
						{
							"Name": "Large",
							"Id": "large",
							"ImgPath": "large family.jpg"
						},
						{
							"Name": "Small",
							"Id": "small",
							"ImgPath": "small family.jpg"
						}
					]
				},
				{
					"Name": "Home Climate",
					"Id": "homeclimate",
					"Levels":
					[
						{
							"Name": "Hot",
							"Id": "hot"
						},
						{
							"Name": "Cold",
							"Id": "cold"
						}
					]
				}
			],
			"OutcomeVariable":
			{
				"Name": "Performance",
				"Id": "performance",
				"Levels":
				[
					{
						"Name": "A",
						"Id": "A"
					},
					{
						"Name": "B",
						"Id": "B"
					},
					{
						"Name": "C",
						"Id": "C"
					},
					{
						"Name": "D",
						"Id": "D"
					},
					{
						"Name": "E",
						"Id": "E"
					}
				]
			}
		},
		"Phase":
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
							"Type": "RECORD",
							"UIActionModeId": "RECORD_SELECT_ONE",
							"ExpectedResponses": 
							[
								{
									"Id": "p1r1p1r1",
									"NextPrompt":
									{
										"Id": "p1r1p1r1p1",
										"Text": "Thank you",
										"Type": "COMPLETE",
										"UIActionModeId": "ONE_RECORD_PERFORMANCE"
									}
								}
							]
						}
					},
					{
						"Id": "p1r2",
						"Text": "Two",
						"NextPrompt":
						{
							"Id": "p1r2p1",
							"Text": "Which records would you like to see?",
							"Type": "RECORD",
							"UIActionModeId": "RECORD_SELECT_TWO",
							"ExpectedResponses": 
							[
								{
									"Id": "p1r2p1nonvarying",
									"NextPrompt":
									{
										"Id": "p1r1p1r1p1"
									}
								},
								{
									"Id": "p1r2p1uncontrolled",
									"NextPrompt":
									{
										"Id": "p1r2p1r1p1"
									}
								},
								{
									"Id": "p1r2p1controlled",
									"NextPrompt":
									{
										"Id": "p1r2p1r1p1"
									}
								}
							]
						}
					}
				]
			}
		}
	}`
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
	Type              string
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
	dec := json.NewDecoder(strings.NewReader(promptTreeJsonStream))
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

type UIPrompt interface {
	Display() string
	GetId() string
}

type UITextPrompt struct {
	Type           string
	Text           string
	PromptId       string
	ResponseId     string
	UIActionModeId string
}

func NewUITextPrompt() *UITextPrompt {
	return &UITextPrompt{Type: UI_PROMPT_TEXT}
}

func (ps *UITextPrompt) GetId() string {
	return ps.PromptId
}

func (ps *UITextPrompt) Display() string {
	return ps.Text
}

func NewUIEndPrompt() *UITextPrompt {
	return &UITextPrompt{Type: UI_PROMPT_END}
}

type UIMCPrompt struct {
	Type           string
	Text           string
	Options        []UIOption
	PromptId       string
	UIActionModeId string
}

func NewUIMCPrompt() *UIMCPrompt {
	return &UIMCPrompt{Type: UI_PROMPT_MC}
}

type UIOption struct {
	ResponseId string
	Text       string
}

func (ps *UIMCPrompt) GetId() string {
	return ps.PromptId
	// return ps.WorkflowStateID
}

func (ps *UIMCPrompt) Display() string {
	return ps.Text
}

type UIRecordPrompt struct {
	Type           string
	Text           string
	PromptId       string
	UIActionModeId string
	Factors        []UIFactor
}

func NewUUIRecordPrompt() *UIRecordPrompt {
	return &UIRecordPrompt{Type: UI_PROMPT_RECORD}
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

func (ps *UIRecordPrompt) GetId() string {
	return ps.PromptId
}

func (ps *UIRecordPrompt) Display() string {
	return ps.Text
}
