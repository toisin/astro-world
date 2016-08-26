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

// Prompt logics specific to Chart phase

type ChartPrompt struct {
	*GenericPrompt
}

func MakeChartPrompt(p PromptConfig, UiUserData *UIUserData) *ChartPrompt {
	var n *ChartPrompt
	n = &ChartPrompt{}
	n.GenericPrompt = &GenericPrompt{}
	n.GenericPrompt.currentPrompt = n
	n.init(p, UiUserData)
	return n
}

func (cp *ChartPrompt) ProcessResponse(r string, u *db.User, UiUserData *UIUserData, c appengine.Context) {
	if cp.promptConfig.ResponseType == RESPONSE_END {
		// Sequence has ended. Update remaining factors
		UiUserData.State.(*ChartPhaseState).updateRemainingFactors()
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

func (cp *ChartPrompt) updateState(UiUserData *UIUserData) {
	if UiUserData.State != nil {
		// if UiUserData already have a cp state, use that and update it
		if UiUserData.State.GetPhaseId() == appConfig.ChartPhase.Id {
			cp.state = UiUserData.State.(*ChartPhaseState)
		}
	}
	if cp.state == nil {
		cps := &ChartPhaseState{}
		cps.initContents(appConfig.ChartPhase.ContentRef.Factors)
		cp.state = cps
		cp.state.setPhaseId(appConfig.ChartPhase.Id)
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
