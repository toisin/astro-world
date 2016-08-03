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

// Prompt logics specific to Chart phase

type ChartPrompt struct {
	*GenericPrompt
}

func MakeChartPrompt(p PromptConfig, uiUserData *UIUserData) *ChartPrompt {
	var n *ChartPrompt
	n = &ChartPrompt{}
	n.GenericPrompt = &GenericPrompt{}
	n.GenericPrompt.currentPrompt = n
	n.init(p, uiUserData)
	return n
}

func (cp *ChartPrompt) ProcessResponse(r string, u *db.User, uiUserData *UIUserData, c appengine.Context) {
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

//TODO - cleanup copied from CovPrompt
func (cp *ChartPrompt) initDynamicResponseUIPrompt(uiUserData *UIUserData) {
	pc := cp.promptConfig
	cp.currentUIPrompt = NewUIBasicPrompt()
	cp.currentUIPrompt.setPromptType(pc.PromptType)
	cp.currentPrompt.initUIPromptDynamicText(uiUserData, nil)
	cp.currentUIPrompt.setText(cp.promptDynamicText.String())
	cp.currentUIPrompt.setId(pc.Id)

	options := []*UIOption{}
	for i := range pc.ExpectedResponses.Values {
		switch pc.ExpectedResponses.Values[i].Id {
		case EXPECTED_SPECIAL_CONTENT_REF:
			options = append(options, &UIOption{pc.ExpectedResponses.Values[i].Id, pc.ExpectedResponses.Values[i].Text})
		default:
			options = append(options, &UIOption{pc.ExpectedResponses.Values[i].Id, pc.ExpectedResponses.Values[i].Text})
		}
	}
	cp.currentUIPrompt.setOptions(options)
}

func (cp *ChartPrompt) initUIPromptDynamicText(uiUserData *UIUserData, r Response) {
	p := &UIPromptDynamicText{}
	p.previousResponse = r
	p.promptConfig = cp.promptConfig
	cp.updateState(uiUserData)
	p.state = cp.state
	cp.promptDynamicText = p
}

func (cp *ChartPrompt) updateStateCurrentFactor(uiUserData *UIUserData, fid string) {
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

func (cp *ChartPrompt) updateState(uiUserData *UIUserData) {
	if uiUserData.State != nil {
		// if uiUserData already have a cp state, use that and update it
		if uiUserData.State.GetPhaseId() == appConfig.ChartPhase.Id {
			cp.state = uiUserData.State.(*ChartPhaseState)
		}
	}
	if cp.state == nil {
		cp.state = &ChartPhaseState{}
		cp.state.setPhaseId(appConfig.ChartPhase.Id)
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
