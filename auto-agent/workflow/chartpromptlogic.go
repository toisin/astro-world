package workflow

// import (
// 	"db"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"strings"
// 	"text/template"

// 	"appengine"
// )

// // Prompt logics specific to Chart phase

// type ChartPrompt struct {
// 	// previousPrompt Prompt
// 	response                Response
// 	expectedResponseHandler *ExpectedResponseHandler
// 	currentUIPrompt         UIPrompt
// 	currentUIAction         UIAction
// 	promptConfig            PromptConfig
// 	nextPrompt              Prompt
// 	promptDynamicText       *UICovPromptDynamicText
// 	state                   *CovPhaseState
// }

// func MakeChartPrompt(p PromptConfig) *ChartPrompt {
// 	erh := MakeExpectedResponseHandler(p.ExpectedResponses, PHASE_CHART)

// 	n := &ChartPrompt{}
// 	n.promptConfig = p
// 	n.expectedResponseHandler = erh
// 	return n
// }

// func (cp *ChartPrompt) GetPhaseId() string {
// 	return PHASE_CHART
// }

// func (cp *ChartPrompt) GetPromptId() string {
// 	return cp.promptConfig.Id
// }

// func (cp *ChartPrompt) ProcessResponse(r string, uiUserData *UIUserData, c appengine.Context) {
// 	if r != "" {
// 		dec := json.NewDecoder(strings.NewReader(r))
// 		pc := cp.promptConfig
// 		switch pc.ResponseType {
// 		case RESPONSE_SELECT_TARGET_FACTOR:
// 			for {
// 				var response SimpleResponse
// 				if err := dec.Decode(&response); err == io.EOF {
// 					break
// 				} else if err != nil {
// 					fmt.Fprintf(os.Stderr, "%s", err)
// 					log.Fatal(err)
// 					return
// 				}
// 				// TODO - cleanup double check that db.User.CurrentFactorId does not need
// 				// to be updated now. When does that get updated?
// 				uiUserData.CurrentFactorId = response.Id
// 				cp.response = &response
// 			}
// 			break
// 		case RESPONSE_RECORD:
// 			for {
// 				var recordsResponse RecordsSelectResponse
// 				if err := dec.Decode(&recordsResponse); err == io.EOF {
// 					break
// 				} else if err != nil {
// 					fmt.Fprintf(os.Stderr, "%s", err)
// 					log.Fatal(err)
// 					return
// 				}
// 				recordsResponse.CheckRecords(uiUserData, c)
// 				cp.updateStateRecords(uiUserData, &recordsResponse)
// 				cp.response = &recordsResponse
// 			}
// 			break
// 		default:
// 			for {
// 				var response SimpleResponse
// 				if err := dec.Decode(&response); err == io.EOF {
// 					break
// 				} else if err != nil {
// 					fmt.Fprintf(os.Stderr, "%s", err)
// 					log.Fatal(err)
// 					return
// 				}
// 				cp.response = &response
// 			}
// 		}
// 		if cp.response != nil {
// 			cp.nextPrompt = cp.expectedResponseHandler.getNextPrompt(cp.response.GetResponseId())
// 			cp.nextPrompt.initUIPromptDynamicText(uiUserData, &cp.response)
// 		}
// 	}
// }

// func (cp *ChartPrompt) initUIPromptDynamicText(uiUserData *UIUserData, r *Response) {
// 	if cp.promptDynamicText == nil {
// 		p := &UICovPromptDynamicText{}
// 		p.previousResponse = r
// 		p.promptConfig = cp.promptConfig
// 		cp.updateState(uiUserData)
// 		p.state = cp.state
// 		cp.promptDynamicText = p
// 	}
// }

// func (cp *ChartPrompt) GetNextPrompt() Prompt {
// 	return cp.nextPrompt
// }

// func (cp *ChartPrompt) GetResponseText() string {
// 	return cp.response.GetResponseText()
// }

// func (cp *ChartPrompt) GetResponseId() string {
// 	return cp.response.GetResponseId()
// }

// func (cp *ChartPrompt) GetUIPrompt(uiUserData *UIUserData) UIPrompt {
// 	if cp.currentUIPrompt == nil {
// 		pc := cp.promptConfig
// 		cp.currentUIPrompt = NewUIBasicPrompt()
// 		cp.currentUIPrompt.setPromptType(pc.PromptType)
// 		cp.initUIPromptDynamicText(uiUserData, nil)
// 		if cp.promptDynamicText != nil {
// 			cp.currentUIPrompt.setText(cp.promptDynamicText.String())
// 		}
// 		cp.currentUIPrompt.setId(pc.Id)
// 		options := make([]UIOption, len(pc.ExpectedResponses))
// 		for i := range pc.ExpectedResponses {
// 			options[i] = UIOption{pc.ExpectedResponses[i].Id, pc.ExpectedResponses[i].Text}
// 		}
// 		cp.currentUIPrompt.setOptions(options)
// 	}
// 	return cp.currentUIPrompt
// }

// // Returned UIAction may be nil if not action UI is needed
// func (cp *ChartPrompt) GetUIAction() UIAction {
// 	if cp.currentUIAction == nil {
// 		pc := cp.promptConfig
// 		switch pc.UIActionModeId {
// 		case "RECORD_SELECT_TWO", "RECORD_SELECT_ONE":
// 			p := NewUIRecordAction()
// 			// TODO in progress
// 			// p.SetPromptType(???)
// 			p.Factors = make([]UIFactor, len(appConfig.CovPhase.FactorsOrder))
// 			for i, v := range appConfig.CovPhase.FactorsOrder {
// 				f := GetFactorConfig(v)
// 				p.Factors[i] = UIFactor{
// 					FactorId: f.Id,
// 					Text:     f.Name,
// 				}
// 				p.Factors[i].Levels = make([]UIFactorOption, len(f.Levels))
// 				for j := range f.Levels {
// 					p.Factors[i].Levels[j] = UIFactorOption{
// 						FactorLevelId: f.Levels[j].Id,
// 						Text:          f.Levels[j].Name,
// 						ImgPath:       f.Levels[j].ImgPath,
// 					}
// 				}
// 			}
// 			cp.currentUIAction = p
// 			break
// 		default:
// 			p := NewUIBasicAction()
// 			cp.currentUIAction = p
// 			break
// 		}
// 		if cp.currentUIAction != nil {
// 			cp.currentUIAction.setUIActionModeId(pc.UIActionModeId)
// 		}
// 	}
// 	return cp.currentUIAction
// }

// type SimpleResponse struct {
// 	Text string
// 	Id   string
// }

// func (sr *SimpleResponse) GetResponseText() string {
// 	return sr.Text
// }

// func (sr *SimpleResponse) GetResponseId() string {
// 	return sr.Id
// }
