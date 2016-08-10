package workflow

import (
	"db"
	"fmt"
	"os"
)

// Includes only the variables that are needed on the client side
type UIUserData struct {
	Username             string
	Screenname           string
	CurrentPhaseId       string
	CurrentFactorId      string
	History              []*db.Message
	CurrentUIPrompt      UIPrompt
	CurrentUIAction      UIAction
	State                StateEntities
	ContentFactors       []*UIFactor
	ArchiveHistoryLength int
}

type StateEntities interface {
	setPhaseId(string)
	setUsername(string)
	setScreenname(string)
	setTargetFactor(FactorState)
	GetPhaseId() string
	isContentCompleted() bool
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
	RecordNoOne            *RecordState
	RecordNoTwo            *RecordState
	RecordSelectionsTypeId string
	VaryingFactorIds       []string
	VaryingFactorsCount    int
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
	FactorLevels map[string]FactorState // factor id as key
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
	FactorName      string
	FactorId        string
	SelectedLevel   string // Level name
	SelectedLevelId string // Level id
	OppositeLevel   string // Level name
	OppositeLevelId string // Level id
	IsCausal        bool
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

type UIFactor struct {
	FactorId       string
	Text           string
	Levels         []*UIFactorOption
	IsBeliefCausal bool
	BestLevelId    string
}

type UIFactorOption struct {
	FactorLevelId string
	Text          string
	ImgPath       string
}

func MakeUIUserData(u db.User) *UIUserData {
	uiUserData := &UIUserData{}

	// update new UserData with everything that is available on db.User
	uiUserData.Username = u.Username
	uiUserData.Screenname = u.Screenname
	uiUserData.CurrentFactorId = u.CurrentFactorId
	uiUserData.CurrentPhaseId = u.CurrentPhaseId

	if u.UIState != nil {
		s, err := UnstringifyState(u.UIState, u.CurrentPhaseId)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting json to StateEntities: %s\n\n", err)
		}
		uiUserData.State = s
	}

	uiUserData.ContentFactors = make([]*UIFactor, len(appConfig.CovPhase.ContentRef.Factors))
	for i, v := range appConfig.CovPhase.ContentRef.Factors {
		f := GetFactorConfig(v.Id)
		uiUserData.ContentFactors[i] = &UIFactor{
			FactorId: f.Id,
			Text:     f.Name,
		}
		uiUserData.ContentFactors[i].Levels = make([]*UIFactorOption, len(f.Levels))
		for j := range f.Levels {
			uiUserData.ContentFactors[i].Levels[j] = &UIFactorOption{
				FactorLevelId: f.Levels[j].Id,
				Text:          f.Levels[j].Name,
				ImgPath:       f.Levels[j].ImgPath,
			}
		}
	}
	return uiUserData
}

type UISelectedFactor struct {
	FactorId        string
	SelectedLevelId string
}

type UIRecordsSelectResponse struct {
	RecordNoOne         []*UISelectedFactor
	RecordNoTwo         []*UISelectedFactor
	Id                  string
	VaryingFactorIds    []string
	VaryingFactorsCount int
	dbRecordNoOne       db.Record
	dbRecordNoTwo       db.Record
	UseDBRecordNoOne    bool
	UseDBRecordNoTwo    bool
}

func (rsr *UIRecordsSelectResponse) GetResponseText() string {
	return ""
}

func (rsr *UIRecordsSelectResponse) GetResponseId() string {
	return rsr.Id
}

type UIPriorBeliefFactor struct {
	FactorId    string
	IsCausal    bool
	BestLevelId string
}

type UIPriorBeliefResponse struct {
	CausalFactors []*UIPriorBeliefFactor
	Id            string
}

func (rsr *UIPriorBeliefResponse) GetResponseText() string {
	return ""
}

func (rsr *UIPriorBeliefResponse) GetResponseId() string {
	return rsr.Id
}
