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

func (uiUserData *UIUserData) initPhase(pId string) {
	if pId != "" && uiUserData.CurrentPhaseId != pId {
		uiUserData.CurrentPhaseId = pId

		uiUserData.ContentFactors = make([]*UIFactor, len(GetPhase(uiUserData.CurrentPhaseId).ContentRef.Factors))
		for i, v := range GetPhase(uiUserData.CurrentPhaseId).ContentRef.Factors {
			f := GetFactorConfig(v.Id)
			uiUserData.ContentFactors[i] = &UIFactor{
				FactorId: f.Id,
				Text:     f.Name,
				IsCausal: f.IsCausal,
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
	}
}

func (uiUserData *UIUserData) initPrompt(currentPrompt Prompt) {
	uiUserData.initPhase(currentPrompt.GetPhaseId())
	uiUserData.CurrentUIAction = currentPrompt.GetUIAction()
	uiUserData.CurrentUIPrompt = currentPrompt.GetUIPrompt()
}

type StateEntities interface {
	setPhaseId(string)
	setUsername(string)
	setScreenname(string)
	setTargetFactor(FactorState)
	setRemainingFactorIds([]string)
	setBeliefs(BeliefsState)
	setLastMemo(UIMemoResponse)
	GetPhaseId() string
	GetBeliefs() BeliefsState
	isContentCompleted() bool
	GetTargetFactor() FactorState
	GetRemainingFactorIds() []string
	GetLastMemo() UIMemoResponse
}

type GenericState struct {
	PhaseId            string
	Username           string
	Screenname         string
	TargetFactor       FactorState
	RemainingFactorIds []string
	Beliefs            BeliefsState
	LastMemo           UIMemoResponse
}

type BeliefsState struct {
	HasCausalFactors         bool
	CausalFactors            []string
	HasMultipleCausalFactors bool
}

func (c *GenericState) GetPhaseId() string {
	return c.PhaseId
}

func (c *GenericState) GetLastMemo() UIMemoResponse {
	return c.LastMemo
}

func (c *GenericState) GetBeliefs() BeliefsState {
	return c.Beliefs
}

// Not applicable to all phases
func (c *GenericState) GetTargetFactor() FactorState {
	return c.TargetFactor
}

// Not applicable to all phases
func (c *GenericState) GetRemainingFactorIds() []string {
	return c.RemainingFactorIds
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
func (c *GenericState) setRemainingFactorIds(ss []string) {
	c.RemainingFactorIds = ss
}

func (c *GenericState) setBeliefs(s BeliefsState) {
	c.Beliefs = s
}

func (c *GenericState) setLastMemo(s UIMemoResponse) {
	c.LastMemo = s
}

func (cp *GenericState) isContentCompleted() bool {
	if len(cp.RemainingFactorIds) > 0 {
		return false
	}
	return true
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
type ChartPhaseState struct {
	GenericState
}

func (c *ChartPhaseState) GetPhaseId() string {
	return appConfig.ChartPhase.Id
}

func (c *ChartPhaseState) initContents(factors []Factor) {
	c.RemainingFactorIds = make([]string, len(factors))
	for i, v := range factors {
		c.RemainingFactorIds[i] = v.Id
	}
}

// Implements workflow.StateEntities
type CovPhaseState struct {
	GenericState
	RecordNoOne            RecordState
	RecordNoTwo            RecordState
	RecordSelectionsTypeId string
	VaryingFactorIds       []string
	VaryingFactorsCount    int
	NonVaryingFactorIds    []string
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

type RecordState struct {
	RecordName       string
	FirstName        string
	LastName         string
	RecordNo         string
	FactorLevels     map[string]FactorState // factor id as key
	Performance      string
	PerformanceLevel int
	// Factor id as keys, such as:
	// "fitness",
	// "parentshealth",
	// "education",
	// "familysize"
}

// For workflow.json to reference
// as state variable using go template
// (Used by StateEntities.TargetFactor & CovPhaseState.RecordState)
//
// This type is used in multiple contexts.
// Not all members may be relevant.
type FactorState struct {
	FactorName       string
	FactorId         string
	SelectedLevel    string // Level name
	SelectedLevelId  string // Level id
	OppositeLevel    string // Level name
	OppositeLevelId  string // Level id
	IsCausal         bool
	IsConcludeCausal bool
	HasConclusion    bool
}

func (fs *FactorState) String() string {
	s := fs.FactorId + ":" + fs.SelectedLevelId
	return s
}

// For jsx to reference all factors
// configured for the particular phase
// in the workflow.json
// (Used by UIUserData.ContentFactors)
type UIFactor struct {
	FactorId       string
	Text           string
	Levels         []*UIFactorOption
	IsBeliefCausal bool
	BestLevelId    string
	IsCausal       bool
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

	if u.UIState != nil {
		s, err := UnstringifyState(u.UIState, u.CurrentPhaseId)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting json to StateEntities: %s\n\n", err)
		}
		uiUserData.State = s
	}

	uiUserData.initPhase(uiUserData.CurrentPhaseId)

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
	NonVaryingFactorIds []string
	dbRecordNoOne       db.Record
	dbRecordNoTwo       db.Record
	UseDBRecordNoOne    bool
	UseDBRecordNoTwo    bool
}

func (rsr *UIRecordsSelectResponse) GetResponseText() string {
	responseText := ""
	count := 0
	if rsr.dbRecordNoOne.RecordNo != "" {
		responseText = "Record #" + rsr.dbRecordNoOne.RecordNo
		count++
	}
	if rsr.dbRecordNoTwo.RecordNo != "" {
		if count > 0 {
			responseText = responseText + " and " + "Record #" + rsr.dbRecordNoTwo.RecordNo
		} else {
			responseText = "Record #" + rsr.dbRecordNoTwo.RecordNo
		}
	}

	return responseText
}

func (rsr *UIRecordsSelectResponse) GetResponseId() string {
	return rsr.Id
}

// For Prior belief screen UI jsx
// (Used by UIPriorBeliefResponse.BeliefFactors)
type UIPriorBeliefFactor struct {
	FactorId       string
	IsBeliefCausal bool
	BestLevelId    string
}

type UIPriorBeliefResponse struct {
	BeliefFactors []*UIPriorBeliefFactor
	Id            string
}

func (rsr *UIPriorBeliefResponse) GetResponseText() string {
	responseText := ""
	count := 0
	totalcausal := 0

	for _, v := range rsr.BeliefFactors {
		// quickly count number of causal
		if v.IsBeliefCausal {
			totalcausal++
		}
	}

	for _, v := range rsr.BeliefFactors {
		if v.IsBeliefCausal {
			factorName := GetFactorConfig(v.FactorId).Name
			if count == 0 {
				responseText = factorName
			} else if count == (totalcausal - 1) {
				responseText = responseText + " and " + factorName
			} else {
				responseText = responseText + ", " + factorName
			}
			if v.BestLevelId != "" {
				responseText = responseText + ": " + v.BestLevelId
			}
			count++
		}
	}
	return responseText
}

func (rsr *UIPriorBeliefResponse) GetResponseId() string {
	return rsr.Id
}

type UIMemoResponse struct {
	Ask        string
	Memo       string
	Evidence   string
	Id         string // Name of the factor for the memo
	FactorName string
}

func (rsr *UIMemoResponse) GetResponseText() string {
	return rsr.Memo
}

func (rsr *UIMemoResponse) GetResponseId() string {
	return rsr.Id
}

type Performance struct {
	Grade   string
	Records map[string][]RecordState
}

func GetAllPerformanceRecords(records []db.Record) []Performance {
	// e.g. pd[0].Grade = A
	//      pd[0].["all"] = { .. all of records with A performance .. }
	//      pd[0].["fitness:average"] = { .. all of records with
	//  																	fitness = average and have A performance .. }
	pd := make([]Performance, len(appConfig.Content.OutcomeVariable.Levels))
	// e.g. counts[0]["all"] = number of records with A performance
	//      counts[0]["fitness:average"] = number of records with
	//	  																 fitness = average and have A performance
	counts := make([]map[string]int, len(appConfig.Content.OutcomeVariable.Levels))

	for _, v := range records {
		pData := &pd[v.OutcomeLevelOrder]
		if pData.Records == nil {
			pData.Grade = v.OutcomeLevel
			pData.Records = make(map[string][]RecordState)
			counts[v.OutcomeLevelOrder] = make(map[string]int)
		}
		r := CreateRecordStateFromDB(v)
		for _, w := range r.FactorLevels {
			if pData.Records[w.String()] == nil {
				pData.Records[w.String()] = make([]RecordState, 40)
			} else {
				pData.Records[w.String()][counts[v.OutcomeLevelOrder][w.String()]] = r
				counts[v.OutcomeLevelOrder][w.String()]++
			}
		}
		if pData.Records["all"] == nil {
			pData.Records["all"] = make([]RecordState, 40)
		} else {
			pData.Records["all"][counts[v.OutcomeLevelOrder]["all"]] = r
			counts[v.OutcomeLevelOrder]["all"]++
		}
	}
	for i := range pd {
		// In the case when no records had a certain performance level
		// Add the grade for completing. appConfig.Content.OutcomeVariable.Levels[i].Name should be
		// the same as OutcomeLevel from a record
		if pd[i].Grade == "" {
			pd[i].Grade = appConfig.Content.OutcomeVariable.Levels[i].Name
			pd[i].Records = make(map[string][]RecordState)
		}
		for k, _ := range pd[i].Records {
			pd[i].Records[k] = pd[i].Records[k][0 : counts[i][k]+1]
		}
	}
	return pd
}
