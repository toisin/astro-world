package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CovPrompt struct {
	// previousPrompt Prompt
	response                Response
	recordsResponse         RecordsSelectResponse
	expectedResponseHandler *ExpectedResponseHandler
	// promptGenerator PromptGenerator
	currentUIPrompt UIPrompt
	promptConfig    PromptConfig
	nextPrompt      Prompt
}

type RecordsSelectResponse struct {
	RecordNoOne []SelectedFactor
	RecordNoTwo []SelectedFactor
}

type SelectedFactor struct {
	FactorId        string
	SelectedLevelId string
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

func (cp *CovPrompt) ProcessResponse(r string) {
	if r != "" {
		dec := json.NewDecoder(strings.NewReader(r))
		pc := cp.promptConfig
		switch pc.Type {
		case UI_PROMPT_RECORD:
			for {
				if err := dec.Decode(&cp.recordsResponse); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				//TODO totally hard coding
				cp.response.Text = fmt.Sprint(cp.recordsResponse)
				cp.response.Id = "p1r2p1nonvarying"
				cp.nextPrompt = cp.expectedResponseHandler.GetNextPrompt(cp.response.Id)
			}
		default:
			for {
				if err := dec.Decode(&cp.response); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.nextPrompt = cp.expectedResponseHandler.GetNextPrompt(cp.response.Id)
			}
		}
	}
}

func (cp *CovPrompt) GetNextPrompt() Prompt {
	return cp.nextPrompt
}

func (cp *CovPrompt) GetResponseText() string {
	return cp.response.Text
}

func (cp *CovPrompt) GetResponseId() string {
	return cp.response.Id
}

func (cp *CovPrompt) GetUIPrompt() UIPrompt {
	if cp.currentUIPrompt == nil {
		pc := cp.promptConfig
		switch pc.Type {
		case UI_PROMPT_MC:
			p := NewUIMCPrompt()
			p.Text = cp.promptConfig.Text // TODO need to process dynamic data
			p.Options = make([]UIOption, len(pc.ExpectedResponses))
			for i := range pc.ExpectedResponses {
				p.Options[i] = UIOption{pc.ExpectedResponses[i].Id, pc.ExpectedResponses[i].Text}
			}
			p.PromptId = pc.Id
			p.UIActionModeId = pc.UIActionModeId
			cp.currentUIPrompt = p
			break
		case UI_PROMPT_RECORD:
			p := NewUUIRecordPrompt()
			p.Text = cp.promptConfig.Text
			p.PromptId = pc.Id
			p.UIActionModeId = pc.UIActionModeId
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
			cp.currentUIPrompt = p
			break
		case UI_PROMPT_TEXT:
			p := NewUITextPrompt()
			p.Text = cp.promptConfig.Text // TODO need to process dynamic data
			p.PromptId = pc.Id
			p.UIActionModeId = pc.UIActionModeId
			p.ResponseId = pc.ExpectedResponses[0].Id
			cp.currentUIPrompt = p
			break
		case UI_PROMPT_END:
			p := NewUIEndPrompt()
			p.Text = cp.promptConfig.Text // TODO need to process dynamic data
			p.PromptId = pc.Id
			p.UIActionModeId = pc.UIActionModeId
			cp.currentUIPrompt = p
			break
		}
	}
	return cp.currentUIPrompt
}
