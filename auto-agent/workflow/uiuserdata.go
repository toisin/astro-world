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
