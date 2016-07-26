package workflow

import (
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"appengine"
)

// Prompt logics specific to Cov phase

type CovPrompt struct {
	*GenericPrompt
}

func MakeCovPrompt(p *PromptConfig) *CovPrompt {
	var n *CovPrompt
	if p != nil {
		erh := MakeExpectedResponseHandler(p)

		n = &CovPrompt{}
		n.GenericPrompt = &GenericPrompt{}
		n.GenericPrompt.currentPrompt = n
		n.promptConfig = p
		n.expectedResponseHandler = erh
	}
	return n
}

func (cp *CovPrompt) ProcessResponse(r string, u *db.User, uiUserData *UIUserData, c appengine.Context) {
	if r != "" {
		dec := json.NewDecoder(strings.NewReader(r))
		pc := cp.promptConfig
		switch pc.ResponseType {
		case RESPONSE_SELECT_TARGET_FACTOR:
			for {
				var response SimpleResponse
				if err := dec.Decode(&response); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				uiUserData.CurrentFactorId = response.Id
				u.CurrentFactorId = uiUserData.CurrentFactorId
				cp.updateStateCurrentFactor(uiUserData, uiUserData.CurrentFactorId)
				cp.response = &response
			}
			break
		case RESPONSE_RECORD:
			for {
				var recordsResponse RecordsSelectResponse
				if err := dec.Decode(&recordsResponse); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				recordsResponse.CheckRecords(uiUserData, c)
				cp.updateStateRecords(uiUserData, &recordsResponse)
				cp.response = &recordsResponse
			}
			break
		default:
			for {
				var response SimpleResponse
				if err := dec.Decode(&response); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.response = &response
			}
		}
		if cp.response != nil {
			cp.nextPrompt = cp.expectedResponseHandler.getNextPrompt(cp.response.GetResponseId())
			cp.nextPrompt.initUIPromptDynamicText(uiUserData, cp.response)
		}
	}
}

func (cp *CovPrompt) initUIPromptDynamicText(uiUserData *UIUserData, r Response) {
	if cp.promptDynamicText == nil {
		p := &UIPromptDynamicText{}
		p.previousResponse = r
		p.promptConfig = cp.promptConfig
		cp.updateState(uiUserData)
		p.state = cp.state
		cp.promptDynamicText = p
	}
}

// Returned UIAction may be nil if not action UI is needed
func (cp *CovPrompt) GetUIAction() UIAction {
	if cp.currentUIAction == nil {
		pc := cp.promptConfig
		switch pc.UIActionModeId {
		case "RECORD_SELECT_TWO", "RECORD_SELECT_ONE":
			p := NewUIRecordAction()
			// TODO in progress
			// p.SetPromptType(???)
			p.Factors = make([]*UIFactor, len(appConfig.CovPhase.ContentRef.Factors))
			for i, v := range appConfig.CovPhase.ContentRef.Factors {
				f := GetFactorConfig(v.Id)
				p.Factors[i] = &UIFactor{
					FactorId: f.Id,
					Text:     f.Name,
				}
				p.Factors[i].Levels = make([]*UIFactorOption, len(f.Levels))
				for j := range f.Levels {
					p.Factors[i].Levels[j] = &UIFactorOption{
						FactorLevelId: f.Levels[j].Id,
						Text:          f.Levels[j].Name,
						ImgPath:       f.Levels[j].ImgPath,
					}
				}
			}
			cp.currentUIAction = p
			break
		default:
			p := NewUIBasicAction()
			cp.currentUIAction = p
			break
		}
		if cp.currentUIAction != nil {
			cp.currentUIAction.setUIActionModeId(pc.UIActionModeId)
		}
	}
	return cp.currentUIAction
}

type RecordsSelectResponse struct {
	RecordNoOne         []*SelectedFactor
	RecordNoTwo         []*SelectedFactor
	Id                  string
	VaryingFactorIds    []string
	VaryingFactorsCount int
	dbRecordNoOne       *db.Record
	dbRecordNoTwo       *db.Record
}

func (rsr *RecordsSelectResponse) CheckRecords(uiUserData *UIUserData, c appengine.Context) {
	rsr.VaryingFactorIds = make([]string, len(appConfig.CovPhase.ContentRef.Factors))
	rsr.VaryingFactorsCount = 0
	var CurrentFactorId = uiUserData.CurrentFactorId
	var isTargetVarying = false

	// Retrieve DB records
	if rsr.RecordNoOne != nil {
		dbOrderedFactorLevels := make([]string, len(rsr.RecordNoOne))
		for _, v := range rsr.RecordNoOne {
			f := GetFactorConfig(v.FactorId)
			j := f.DBIndex
			dbOrderedFactorLevels[j] = v.SelectedLevelId
		}
		record, _, err := db.GetRecord(c, dbOrderedFactorLevels)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting First Record:"+err.Error()+"!\n\n")
			log.Fatal(err)
			return
		}

		rsr.dbRecordNoOne = &record
	}

	if rsr.RecordNoTwo != nil {
		dbOrderedFactorLevels := make([]string, len(rsr.RecordNoTwo))
		for _, v := range rsr.RecordNoTwo {
			f := GetFactorConfig(v.FactorId)
			j := f.DBIndex
			dbOrderedFactorLevels[j] = v.SelectedLevelId
		}
		record, _, err := db.GetRecord(c, dbOrderedFactorLevels)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting Second Record:"+err.Error()+"!\n\n")
			log.Fatal(err)
			return
		}
		rsr.dbRecordNoTwo = &record
	}

	// Determine the type of record response
	if rsr.RecordNoTwo != nil {
		// For each factor, check if the two records have different levels
		for i := range rsr.RecordNoOne {
			for j := range rsr.RecordNoTwo {
				if rsr.RecordNoOne[i].FactorId == rsr.RecordNoTwo[j].FactorId {
					if rsr.RecordNoOne[i].SelectedLevelId != rsr.RecordNoTwo[j].SelectedLevelId {
						if rsr.RecordNoOne[i].FactorId == CurrentFactorId {
							isTargetVarying = true
						}
						rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoOne[i].FactorId
						rsr.VaryingFactorsCount++
					}
				}
			}
		}

		if rsr.VaryingFactorsCount == 0 {
			rsr.Id = COV_RESPONSE_ID_NON_VARYING
		} else if !isTargetVarying {
			rsr.Id = COV_RESPONSE_ID_TARGET_NON_VARYING
		} else if rsr.VaryingFactorsCount == 1 {
			if rsr.VaryingFactorIds[0] == CurrentFactorId {
				rsr.Id = COV_RESPONSE_ID_CONTROLLED
			} else {
				rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
			}
		} else {
			rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
		}
	} else {
		rsr.Id = COV_RESPONSE_ID_SINGLE_CASE
	}
}

func (rsr *RecordsSelectResponse) GetResponseText() string {
	return ""
}

func (rsr *RecordsSelectResponse) GetResponseId() string {
	return rsr.Id
}

func (cp *CovPrompt) updateStateCurrentFactor(uiUserData *UIUserData, fid string) {
	cp.updateState(uiUserData)
	if factorConfigMap[fid] != nil {
		cp.state.setTargetFactor(
			&FactorState{
				FactorName: factorConfigMap[fid].Name,
				FactorId:   fid,
				IsCausal:   factorConfigMap[fid].IsCausal})
	}
	uiUserData.State = cp.state
}

// This method should only update records select
// Unless if no existing state, than create new one, otherwise, only
// update records select
func (cp *CovPrompt) updateStateRecords(uiUserData *UIUserData, r *RecordsSelectResponse) {
	cp.updateState(uiUserData)
	if cp.state != nil {
		s := cp.state.(*CovPhaseState)
		s.RecordNoOne = cp.createRecordStateFromDB(r.dbRecordNoOne, r.RecordNoOne)
		s.RecordNoTwo = cp.createRecordStateFromDB(r.dbRecordNoTwo, r.RecordNoTwo)
		cp.state = s
	}
	uiUserData.State = cp.state
}

func (cp *CovPrompt) updateState(uiUserData *UIUserData) {
	if uiUserData.State != nil {
		// if uiUserData already have a cp state, use that and update it
		if uiUserData.State.GetPhaseId() == appConfig.CovPhase.Id {
			cp.state = uiUserData.State.(*CovPhaseState)
		}
	}
	if cp.state == nil {
		cp.state = &CovPhaseState{}
		cp.state.setUsername(uiUserData.Username)
		cp.state.setScreenname(uiUserData.Screenname)
		fid := uiUserData.CurrentFactorId
		if factorConfigMap[fid] != nil {
			cp.state.setTargetFactor(
				&FactorState{
					FactorName: factorConfigMap[fid].Name,
					FactorId:   fid,
					IsCausal:   factorConfigMap[fid].IsCausal})
		}
	}
	uiUserData.State = cp.state
}

func (cp *CovPrompt) createRecordStateFromDB(r *db.Record, sf []*SelectedFactor) *RecordState {
	rs := &RecordState{}
	if r != nil {
		rs.RecordName = r.Firstname + " " + r.Lastname
		rs.RecordNo = r.RecordNo
		rs.FactorLevels = make(map[string]*FactorState)
		for _, v := range sf {
			rs.FactorLevels[v.FactorId] = CreateCovFactorState(v.FactorId, v.SelectedLevelId)
		}
	} else {
		rs.RecordName = ""
		rs.RecordNo = ""
		rs.FactorLevels = make(map[string]*FactorState)
	}
	return rs
}

type SelectedFactor struct {
	FactorId        string
	SelectedLevelId string
}
