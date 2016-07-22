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

func MakeChartPrompt(p *PromptConfig) *ChartPrompt {
	var n *ChartPrompt
	if p != nil {
		erh := MakeExpectedResponseHandler(p)
		n = &ChartPrompt{}
		n.GenericPrompt = &GenericPrompt{}
		n.GenericPrompt.currentPrompt = n
		n.promptConfig = p
		n.expectedResponseHandler = erh
	}
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
			cp.nextPrompt = cp.expectedResponseHandler.getNextPrompt(cp.response.GetResponseId())
			cp.nextPrompt.initUIPromptDynamicText(uiUserData, &cp.response)
		}
	}
}

func (cp *ChartPrompt) initUIPromptDynamicText(uiUserData *UIUserData, r *Response) {
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
func (cp *ChartPrompt) GetUIAction() UIAction {
	if cp.currentUIAction == nil {
		pc := cp.promptConfig
		switch pc.UIActionModeId {
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

func (cp *ChartPrompt) updateStateCurrentFactor(uiUserData *UIUserData, fid string) {
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

func (cp *ChartPrompt) updateState(uiUserData *UIUserData) {
	if uiUserData.State != nil {
		// if uiUserData already have a cp state, use that and update it
		if uiUserData.State.GetPhaseId() == appConfig.ChartPhase.Id {
			cp.state = uiUserData.State.(*ChartPhaseState)
		}
	}
	if cp.state == nil {
		cp.state = &ChartPhaseState{}
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
