package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Prompt logics specific to Cov phase

type CovPrompt struct {
	// previousPrompt Prompt
	response                Response
	expectedResponseHandler *ExpectedResponseHandler
	currentUIPrompt         UIPrompt
	currentUIAction         UIAction
	promptConfig            PromptConfig
	nextPrompt              Prompt
}

func MakeCovPrompt(p PromptConfig) *CovPrompt {
	erh := MakeExpectedResponseHandler(p.ExpectedResponses, PHASE_COV)

	n := new(CovPrompt)
	n.promptConfig = p
	n.expectedResponseHandler = erh
	return n
}

func (cp *CovPrompt) GetPhaseId() string {
	return PHASE_COV
}

func (cp *CovPrompt) GetPromptId() string {
	return cp.promptConfig.Id
}

func (cp *CovPrompt) ProcessResponse(r string, uiUserData *UIUserData) {
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
				uiUserData.User.CurrentFactorId = response.Id
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
				//TODO totally hard coding
				recordsResponse.CheckRecords(uiUserData)
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
			cp.nextPrompt = cp.expectedResponseHandler.GetNextPrompt(cp.response.GetResponseId())
		}
	}
}

func (cp *CovPrompt) GetNextPrompt() Prompt {
	return cp.nextPrompt
}

func (cp *CovPrompt) GetResponseText() string {
	return cp.response.GetResponseText()
}

func (cp *CovPrompt) GetResponseId() string {
	return cp.response.GetResponseId()
}

func (cp *CovPrompt) GetUIPrompt() UIPrompt {
	if cp.currentUIPrompt == nil {
		pc := cp.promptConfig
		cp.currentUIPrompt = NewUIBasicPrompt()
		cp.currentUIPrompt.SetPromptType(pc.PromptType)
		cp.currentUIPrompt.SetText(pc.Text) // TODO need to process dynamic data
		cp.currentUIPrompt.SetId(pc.Id)
		options := make([]UIOption, len(pc.ExpectedResponses))
		for i := range pc.ExpectedResponses {
			options[i] = UIOption{pc.ExpectedResponses[i].Id, pc.ExpectedResponses[i].Text}
		}
		cp.currentUIPrompt.SetOptions(options)
	}
	return cp.currentUIPrompt
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
			p.Factors = make([]UIFactor, len(contentConfig.CausalFactors)+len(contentConfig.NonCausalFactors))
			count := 0
			for i := range contentConfig.CausalFactors {
				p.Factors[count] = UIFactor{
					FactorId: contentConfig.CausalFactors[i].Id,
					Text:     contentConfig.CausalFactors[i].Name,
				}
				p.Factors[count].Levels = make([]UIFactorOption, len(contentConfig.CausalFactors[i].Levels))
				for j := range contentConfig.CausalFactors[i].Levels {
					p.Factors[count].Levels[j] = UIFactorOption{
						FactorLevelId: contentConfig.CausalFactors[i].Levels[j].Id,
						Text:          contentConfig.CausalFactors[i].Levels[j].Name,
						ImgPath:       contentConfig.CausalFactors[i].Levels[j].ImgPath,
					}
				}
				count++
			}
			for i := range contentConfig.NonCausalFactors {
				p.Factors[count] = UIFactor{
					FactorId: contentConfig.NonCausalFactors[i].Id,
					Text:     contentConfig.NonCausalFactors[i].Name,
				}
				p.Factors[count].Levels = make([]UIFactorOption, len(contentConfig.NonCausalFactors[i].Levels))
				for j := range contentConfig.NonCausalFactors[i].Levels {
					p.Factors[count].Levels[j] = UIFactorOption{
						FactorLevelId: contentConfig.NonCausalFactors[i].Levels[j].Id,
						Text:          contentConfig.NonCausalFactors[i].Levels[j].Name,
						ImgPath:       contentConfig.NonCausalFactors[i].Levels[j].ImgPath,
					}
				}
				count++
			}
			cp.currentUIAction = p
			break
		default:
			p := NewUIBasicAction()
			cp.currentUIAction = p
			break
		}
		if cp.currentUIAction != nil {
			cp.currentUIAction.SetUIActionModeId(pc.UIActionModeId)
		}
	}
	return cp.currentUIAction
}

type SimpleResponse struct {
	Text string
	Id   string
}

func (sr *SimpleResponse) GetResponseText() string {
	return sr.Text
}

func (sr *SimpleResponse) GetResponseId() string {
	return sr.Id
}

type RecordsSelectResponse struct {
	RecordNoOne         []SelectedFactor
	RecordNoTwo         []SelectedFactor
	Id                  string
	VaryingFactorIds    []string
	CountVaryingFactors int
}

func (rsr *RecordsSelectResponse) CheckRecords(uiUserData *UIUserData) {
	rsr.VaryingFactorIds = make([]string, len(contentConfig.CausalFactors)+len(contentConfig.NonCausalFactors))
	rsr.CountVaryingFactors = 0
	var CurrentFactorId = uiUserData.User.CurrentFactorId
	var isTargetVarying = false

	for i := range rsr.RecordNoOne {
		for j := range rsr.RecordNoTwo {
			if rsr.RecordNoOne[i].FactorId == rsr.RecordNoTwo[j].FactorId {
				if rsr.RecordNoOne[i].SelectedLevelId != rsr.RecordNoTwo[j].SelectedLevelId {
					if rsr.RecordNoOne[i].FactorId == CurrentFactorId {
						isTargetVarying = true
					}
					rsr.VaryingFactorIds[rsr.CountVaryingFactors] = rsr.RecordNoOne[i].FactorId
					rsr.CountVaryingFactors++
				}
			}
		}
	}

	if rsr.CountVaryingFactors == 0 {
		rsr.Id = COV_RESPONSE_ID_NON_VARYING
	} else if !isTargetVarying {
		rsr.Id = COV_RESPONSE_ID_TARGET_NON_VARYING
	} else if rsr.CountVaryingFactors == 1 {
		if rsr.VaryingFactorIds[0] == CurrentFactorId {
			rsr.Id = COV_RESPONSE_ID_CONTROLLED
		} else {
			rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
		}
	} else {
		rsr.Id = COV_RESPONSE_ID_UNCONTROLLED
	}

}

func (rsr *RecordsSelectResponse) GetResponseText() string {
	//TODO
	return fmt.Sprint(rsr)
}

func (rsr *RecordsSelectResponse) GetResponseId() string {
	// TODO hard coded
	return rsr.Id
}

type SelectedFactor struct {
	FactorId        string
	SelectedLevelId string
}
