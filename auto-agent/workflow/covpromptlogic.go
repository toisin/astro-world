package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"os"
)

type CovPrompt struct {
	// previousPrompt Prompt
	response Response
	expectedResponseHandler *ExpectedResponseHandler
	// promptGenerator PromptGenerator
	currentUIPrompt UIPrompt
	promptConfig PromptConfig
	nextPrompt Prompt
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

func (cp *CovPrompt) GetNextPrompt() Prompt {
	return cp.nextPrompt
}

func (cp *CovPrompt) GetResponse() Response {
	return cp.response
}

func (cp *CovPrompt) GetUIPrompt() UIPrompt {
	if (cp.currentUIPrompt == nil) {
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
    case UI_PROMPT_RECORD,UI_PROMPT_TEXT:
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

