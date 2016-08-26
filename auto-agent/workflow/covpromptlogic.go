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
	"time"
)

// Prompt logics specific to Cov phase

type CovPrompt struct {
	*GenericPrompt
}

func MakeCovPrompt(p PromptConfig, UiUserData *UIUserData) *CovPrompt {
	var n *CovPrompt
	n = &CovPrompt{}
	n.GenericPrompt = &GenericPrompt{}
	n.GenericPrompt.currentPrompt = n
	n.init(p, UiUserData)
	return n
}

func (cp *CovPrompt) ProcessResponse(r string, u *db.User, UiUserData *UIUserData, c appengine.Context) {
	if cp.promptConfig.ResponseType == RESPONSE_END {
		// Sequence has ended. Update remaining factors
		UiUserData.State.(*CovPhaseState).updateRemainingFactors()
		cp.nextPrompt = cp.generateFirstPromptInNextSequence(UiUserData)
	} else if r != "" {
		dec := json.NewDecoder(strings.NewReader(r))
		pc := cp.promptConfig
		switch pc.ResponseType {
		case RESPONSE_MEMO:
			for {
				var memoResponse UIMemoResponse
				if err := dec.Decode(&memoResponse); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				dbmemo := db.Memo{
					FactorId: memoResponse.Id,
					Ask:      memoResponse.Ask,
					Memo:     memoResponse.Memo,
					Evidence: memoResponse.Evidence,
					PhaseId:  cp.GetPhaseId(),
					Date:     time.Now(),
				}
				err := db.PutMemo(c, u.Username, dbmemo)
				if err != nil {
					fmt.Fprint(os.Stderr, "DB Error Adding Memo:"+err.Error()+"!\n\n")
					return
				}
				cp.updateMemo(UiUserData, memoResponse)
				cp.response = &memoResponse
			}
			break
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
				cp.updatePriorBeliefs(UiUserData, beliefResponse)
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
				UiUserData.CurrentFactorId = response.Id
				u.CurrentFactorId = UiUserData.CurrentFactorId
				cp.updateStateCurrentFactor(UiUserData, UiUserData.CurrentFactorId)
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
				cp.checkRecords(&recordsResponse, UiUserData.CurrentFactorId, c)
				cp.updateStateRecords(UiUserData, recordsResponse)
				cp.response = &recordsResponse
			}
			break
		case RESPONSE_CAUSAL_CONCLUSION:
			for {
				var response SimpleResponse
				if err := dec.Decode(&response); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.updateStateCurrentFactorCausal(UiUserData, response.GetResponseId())
				cp.response = &response
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
			cp.nextPrompt = cp.expectedResponseHandler.generateNextPrompt(cp.response, UiUserData)
		}
	}
}

func (cp *CovPrompt) checkRecords(rsr *UIRecordsSelectResponse, currentFactorId string, c appengine.Context) {
	rsr.NonVaryingFactorIds = make([]string, len(appConfig.CovPhase.ContentRef.Factors))
	rsr.VaryingFactorIds = make([]string, len(appConfig.CovPhase.ContentRef.Factors))
	rsr.VaryingFactorsCount = 0
	nonVaryingFactorsCount := 0
	var isTargetVarying = false

	// Determine the type of record response
	if rsr.RecordNoOne != nil && rsr.RecordNoTwo != nil {
		// Two records selected from the screen
		// For each factor, check if the two records have different levels

		// Force clean up the state to make sure there are no left over of old state info
		cp.state.(*CovPhaseState).RecordNoOne = nil
		cp.state.(*CovPhaseState).RecordNoTwo = nil
		for i := range rsr.RecordNoOne {
			for j := range rsr.RecordNoTwo {
				if rsr.RecordNoOne[i].FactorId == rsr.RecordNoTwo[j].FactorId {
					if rsr.RecordNoOne[i].SelectedLevelId != rsr.RecordNoTwo[j].SelectedLevelId {
						if rsr.RecordNoOne[i].FactorId == currentFactorId {
							isTargetVarying = true
						} else {
							rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoOne[i].FactorId
							rsr.VaryingFactorsCount++
						}
					} else {
						rsr.NonVaryingFactorIds[nonVaryingFactorsCount] = rsr.RecordNoOne[i].FactorId
						nonVaryingFactorsCount++
					}
				}
			}
		}
	} else if rsr.UseDBRecordNoOne && rsr.RecordNoTwo != nil {
		// One record selected from the screen, compare it with the previously selected one
		r := cp.state.(*CovPhaseState).RecordNoOne

		// Force clean up the state to make sure there are no left over of old state info
		cp.state.(*CovPhaseState).RecordNoTwo = nil
		for j := range rsr.RecordNoTwo {
			if r.FactorLevels[rsr.RecordNoTwo[j].FactorId].SelectedLevelId != rsr.RecordNoTwo[j].SelectedLevelId {
				if rsr.RecordNoTwo[j].FactorId == currentFactorId {
					isTargetVarying = true
				} else {
					rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoTwo[j].FactorId
					rsr.VaryingFactorsCount++
				}
			} else {
				rsr.NonVaryingFactorIds[nonVaryingFactorsCount] = rsr.RecordNoTwo[j].FactorId
				nonVaryingFactorsCount++
			}
		}
	} else if rsr.UseDBRecordNoTwo && rsr.RecordNoOne != nil {
		// One record selected from the screen, compare it with the previously selected one
		r := cp.state.(*CovPhaseState).RecordNoTwo

		// Force clean up the state to make sure there are no left over of old state info
		cp.state.(*CovPhaseState).RecordNoOne = nil
		for j := range rsr.RecordNoOne {
			if r.FactorLevels[rsr.RecordNoOne[j].FactorId].SelectedLevelId != rsr.RecordNoOne[j].SelectedLevelId {
				if rsr.RecordNoOne[j].FactorId == currentFactorId {
					isTargetVarying = true
				} else {
					rsr.VaryingFactorIds[rsr.VaryingFactorsCount] = rsr.RecordNoOne[j].FactorId
					rsr.VaryingFactorsCount++
				}
			} else {
				rsr.NonVaryingFactorIds[nonVaryingFactorsCount] = rsr.RecordNoOne[j].FactorId
				nonVaryingFactorsCount++
			}
		}
	} else {
		// no two records were selected for comparison
		rsr.Id = COV_RESPONSE_ID_SINGLE_CASE
	}
	if rsr.Id != COV_RESPONSE_ID_SINGLE_CASE {
		if !isTargetVarying {
			if rsr.VaryingFactorsCount == 0 {
				// Two records were selected but nothing is varying,
				// practically the same as picking one record
				rsr.Id = COV_RESPONSE_ID_NON_VARYING
				rsr.RecordNoTwo = nil
			} else {
				// Two records were selected but only other factors varying
				// and not the target factor
				rsr.Id = COV_RESPONSE_ID_TARGET_NON_VARYING
			}
		} else {
			if rsr.VaryingFactorsCount == 0 {
				// Target factor and only the target factor is varying
				rsr.Id = COV_RESPONSE_ID_CONTROLLED
			} else {
				// Target factor but other factors are also varying
				rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
			}
		}
	}

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
}

// This method should only update records select
// Unless if no existing state, than create new one, otherwise, only
// update records select
func (cp *CovPrompt) updateStateRecords(UiUserData *UIUserData, r UIRecordsSelectResponse) {
	cp.updateState(UiUserData)
	if cp.state != nil {
		s := cp.state.(*CovPhaseState)
		if r.RecordNoOne != nil {
			s.RecordNoOne = cp.createRecordStateFromDB(r.dbRecordNoOne, r.RecordNoOne)
		}
		if r.RecordNoTwo != nil {
			s.RecordNoTwo = cp.createRecordStateFromDB(r.dbRecordNoTwo, r.RecordNoTwo)
		}
		s.RecordSelectionsTypeId = r.GetResponseId()
		s.VaryingFactorIds = r.VaryingFactorIds
		s.VaryingFactorsCount = r.VaryingFactorsCount
		s.NonVaryingFactorIds = r.NonVaryingFactorIds
		cp.state = s
	}
	UiUserData.State = cp.state
}

func (cp *CovPrompt) updateState(UiUserData *UIUserData) {
	if UiUserData.State != nil {
		// if UiUserData already have a cp state, use that and update it
		if UiUserData.State.GetPhaseId() == appConfig.CovPhase.Id {
			cp.state = UiUserData.State.(*CovPhaseState)
		}
	}
	if cp.state == nil {
		cps := &CovPhaseState{}
		cps.initContents(appConfig.CovPhase.ContentRef.Factors)
		cp.state = cps
		cp.state.setPhaseId(appConfig.CovPhase.Id)
		cp.state.setUsername(UiUserData.Username)
		cp.state.setScreenname(UiUserData.Screenname)
		fid := UiUserData.CurrentFactorId
		if fid != "" {
			cp.state.setTargetFactor(
				FactorState{
					FactorName: factorConfigMap[fid].Name,
					FactorId:   fid,
					IsCausal:   factorConfigMap[fid].IsCausal})
		}
	}
	UiUserData.State = cp.state
}

func (cp *CovPrompt) createRecordStateFromDB(r db.Record, sf []*UISelectedFactor) *RecordState {
	rs := &RecordState{}
	if sf != nil && len(sf) != 0 {
		rs.RecordName = r.Firstname + " " + r.Lastname
		rs.FirstName = r.Firstname
		rs.LastName = r.Lastname
		rs.RecordNo = r.RecordNo
		rs.Performance = r.OutcomeLevel
		rs.PerformanceLevel = GetOutcomeLevelOrder(r.OutcomeLevel)
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
