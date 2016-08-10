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

func MakeCovPrompt(p PromptConfig, uiUserData *UIUserData) *CovPrompt {
	var n *CovPrompt
	n = &CovPrompt{}
	n.GenericPrompt = &GenericPrompt{}
	n.GenericPrompt.currentPrompt = n
	n.init(p, uiUserData)
	return n
}

func (cp *CovPrompt) ProcessResponse(r string, u *db.User, uiUserData *UIUserData, c appengine.Context) {
	if cp.promptConfig.ResponseType == RESPONSE_END {
		// Sequence has ended. Update remaining factors
		uiUserData.State.(*CovPhaseState).updateRemainingFactors()
		cp.nextPrompt = cp.generateFirstPromptInNextSequence(uiUserData)
	} else if r != "" {
		dec := json.NewDecoder(strings.NewReader(r))
		pc := cp.promptConfig
		switch pc.ResponseType {
		case RESPONSE_PRIOR_BELIEF_FACTORS, RESPONSE_PRIOR_BELIEF_LEVELS:
			for {
				var beliefResponse UIPriorBeliefResponse
				if err := dec.Decode(&beliefResponse); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.updatePriorBeliefs(uiUserData, beliefResponse)
				cp.response = &beliefResponse
			}
			break
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
				var recordsResponse UIRecordsSelectResponse
				if err := dec.Decode(&recordsResponse); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.checkRecords(&recordsResponse, uiUserData.CurrentFactorId, c)
				cp.updateStateRecords(uiUserData, recordsResponse)
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
			cp.nextPrompt = cp.expectedResponseHandler.generateNextPrompt(cp.response, uiUserData)
		}
	}
}

func (cp *CovPrompt) initDynamicResponseUIPrompt(uiUserData *UIUserData) {
	pc := cp.promptConfig
	cp.currentUIPrompt = NewUIBasicPrompt()
	cp.currentUIPrompt.setPromptType(pc.PromptType)
	cp.currentPrompt.initUIPromptDynamicText(uiUserData, nil)
	if cp.promptDynamicText != nil {
		cp.currentUIPrompt.setText(cp.promptDynamicText.String())
	}
	cp.currentUIPrompt.setId(pc.Id)

	options := []*UIOption{}
	for i := range pc.ExpectedResponses.Values {
		switch pc.ExpectedResponses.Values[i].Id {
		case EXPECTED_SPECIAL_CONTENT_REF:
			c := uiUserData.State.(*CovPhaseState)
			for _, v := range c.RemainingFactorIds {
				options = append(options, &UIOption{v, GetFactorConfig(v).Name})
			}
		default:
			options = append(options, &UIOption{pc.ExpectedResponses.Values[i].Id, pc.ExpectedResponses.Values[i].Text})
		}
	}
	cp.currentUIPrompt.setOptions(options)
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

func (cp *CovPrompt) checkRecords(rsr *UIRecordsSelectResponse, currentFactorId string, c appengine.Context) {
	rsr.VaryingFactorIds = make([]string, len(appConfig.CovPhase.ContentRef.Factors))
	rsr.VaryingFactorsCount = 0
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
		rsr.dbRecordNoOne = record
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
		rsr.dbRecordNoTwo = record
	}

	// Determine the type of record response
	if rsr.RecordNoOne != nil && rsr.RecordNoTwo != nil {
		// For each factor, check if the two records have different levels
		for i := range rsr.RecordNoOne {
			for j := range rsr.RecordNoTwo {
				if rsr.RecordNoOne[i].FactorId == rsr.RecordNoTwo[j].FactorId {
					if rsr.RecordNoOne[i].SelectedLevelId != rsr.RecordNoTwo[j].SelectedLevelId {
						if rsr.RecordNoOne[i].FactorId == currentFactorId {
							isTargetVarying = true
						}
						rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoOne[i].FactorId
						rsr.VaryingFactorsCount++
					}
				}
			}
		}
	} else if rsr.UseDBRecordNoOne && rsr.RecordNoTwo != nil {
		r := cp.state.(*CovPhaseState).RecordNoOne
		for j := range rsr.RecordNoTwo {
			if r.FactorLevels[rsr.RecordNoTwo[j].FactorId].OppositeLevel != rsr.RecordNoTwo[j].SelectedLevelId {
				if rsr.RecordNoTwo[j].FactorId == currentFactorId {
					isTargetVarying = true
				}
				rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoTwo[j].FactorId
				rsr.VaryingFactorsCount++
			}
		}
	} else if rsr.UseDBRecordNoTwo && rsr.RecordNoOne != nil {
		r := cp.state.(*CovPhaseState).RecordNoTwo
		for j := range rsr.RecordNoOne {
			if r.FactorLevels[rsr.RecordNoOne[j].FactorId].SelectedLevel != rsr.RecordNoOne[j].SelectedLevelId {
				if rsr.RecordNoOne[j].FactorId == currentFactorId {
					isTargetVarying = true
				}
				rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoOne[j].FactorId
				rsr.VaryingFactorsCount++
			}
		}
	} else {
		rsr.Id = COV_RESPONSE_ID_SINGLE_CASE
		return
	}
	if rsr.VaryingFactorsCount == 0 {
		rsr.Id = COV_RESPONSE_ID_NON_VARYING
	} else if !isTargetVarying {
		rsr.Id = COV_RESPONSE_ID_TARGET_NON_VARYING
	} else if rsr.VaryingFactorsCount == 1 {
		if rsr.VaryingFactorIds[0] == currentFactorId {
			rsr.Id = COV_RESPONSE_ID_CONTROLLED
		} else {
			rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
		}
	} else {
		rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
	}

}

func (cp *CovPrompt) updatePriorBeliefs(uiUserData *UIUserData, r UIPriorBeliefResponse) {
	causalFactors := []string{}
	var hasCausal bool
	var hasMultipleCausal bool
	for i, v := range uiUserData.ContentFactors {
		uiUserData.ContentFactors[i].IsBeliefCausal = r.CausalFactors[i].IsCausal
		uiUserData.ContentFactors[i].BestLevelId = r.CausalFactors[i].BestLevelId
		if r.CausalFactors[i].IsCausal {
			causalFactors = append(causalFactors, v.Text)
		}
	}
	cp.updateState(uiUserData)
	if len(causalFactors) > 0 {
		hasCausal = true
		if len(causalFactors) > 1 {
			hasMultipleCausal = true
		}
	}

	if cp.state != nil {
		s := cp.state.(*CovPhaseState)
		s.Beliefs = BeliefsState{
			HasCausalFactors:         hasCausal,
			CausalFactors:            causalFactors,
			HasMultipleCausalFactors: hasMultipleCausal}
		cp.state = s
	}
	uiUserData.State = cp.state
}

func (cp *CovPrompt) updateStateCurrentFactor(uiUserData *UIUserData, fid string) {
	cp.updateState(uiUserData)
	if fid != "" {
		cp.state.setTargetFactor(
			FactorState{
				FactorName: factorConfigMap[fid].Name,
				FactorId:   fid,
				IsCausal:   factorConfigMap[fid].IsCausal})
	}
	uiUserData.State = cp.state
}

// This method should only update records select
// Unless if no existing state, than create new one, otherwise, only
// update records select
func (cp *CovPrompt) updateStateRecords(uiUserData *UIUserData, r UIRecordsSelectResponse) {
	cp.updateState(uiUserData)
	if cp.state != nil {
		s := cp.state.(*CovPhaseState)
		if !r.UseDBRecordNoOne {
			s.RecordNoOne = cp.createRecordStateFromDB(r.dbRecordNoOne, r.RecordNoOne)
		}
		if !r.UseDBRecordNoTwo {
			s.RecordNoTwo = cp.createRecordStateFromDB(r.dbRecordNoTwo, r.RecordNoTwo)
		}
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
		cps := &CovPhaseState{}
		cps.initContents(appConfig.CovPhase.ContentRef.Factors)
		cp.state = cps
		cp.state.setPhaseId(appConfig.CovPhase.Id)
		cp.state.setUsername(uiUserData.Username)
		cp.state.setScreenname(uiUserData.Screenname)
		fid := uiUserData.CurrentFactorId
		if fid != "" {
			cp.state.setTargetFactor(
				FactorState{
					FactorName: factorConfigMap[fid].Name,
					FactorId:   fid,
					IsCausal:   factorConfigMap[fid].IsCausal})
		}
	}
	uiUserData.State = cp.state
}

func (cp *CovPrompt) createRecordStateFromDB(r db.Record, sf []*UISelectedFactor) *RecordState {
	rs := &RecordState{}
	if r.RecordNo != "" {
		rs.RecordName = r.Firstname + " " + r.Lastname
		rs.RecordNo = r.RecordNo
		rs.Performance = r.OutcomeLevel
		rs.FactorLevels = make(map[string]FactorState)
		for _, v := range sf {
			rs.FactorLevels[v.FactorId] = CreateCovFactorState(v.FactorId, v.SelectedLevelId)
		}
	} else {
		rs.RecordName = ""
		rs.RecordNo = ""
		rs.FactorLevels = make(map[string]FactorState)
	}
	return rs
}
