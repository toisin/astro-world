package workflow

import (
	"bytes"
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"appengine"
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
	promptDynamicText       *UICovPromptDynamicText
	state                   *CovPhaseState
}

func MakeCovPrompt(p PromptConfig) *CovPrompt {
	erh := MakeExpectedResponseHandler(p.ExpectedResponses, PHASE_COV)

	n := &CovPrompt{}
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

func (cp *CovPrompt) ProcessResponse(r string, uiUserData *UIUserData, c appengine.Context) {
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
			cp.nextPrompt.initUIPromptDynamicText(uiUserData, &cp.response)
		}
	}
}

func (cp *CovPrompt) initUIPromptDynamicText(uiUserData *UIUserData, r *Response) {
	if cp.promptDynamicText == nil {
		p := &UICovPromptDynamicText{}
		p.previousResponse = r
		p.promptConfig = cp.promptConfig
		cp.updateState(uiUserData)
		p.state = cp.state
		cp.promptDynamicText = p
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

func (cp *CovPrompt) GetUIPrompt(uiUserData *UIUserData) UIPrompt {
	if cp.currentUIPrompt == nil {
		pc := cp.promptConfig
		cp.currentUIPrompt = NewUIBasicPrompt()
		cp.currentUIPrompt.setPromptType(pc.PromptType)
		cp.initUIPromptDynamicText(uiUserData, nil)
		if cp.promptDynamicText != nil {
			cp.currentUIPrompt.setText(cp.promptDynamicText.String())
		}
		cp.currentUIPrompt.setId(pc.Id)
		options := make([]UIOption, len(pc.ExpectedResponses))
		for i := range pc.ExpectedResponses {
			options[i] = UIOption{pc.ExpectedResponses[i].Id, pc.ExpectedResponses[i].Text}
		}
		cp.currentUIPrompt.setOptions(options)
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
			p.Factors = make([]UIFactor, len(appConfig.CovPhase.FactorsOrder))
			for i, v := range appConfig.CovPhase.FactorsOrder {
				f := GetFactorConfig(v)
				p.Factors[i] = UIFactor{
					FactorId: f.Id,
					Text:     f.Name,
				}
				p.Factors[i].Levels = make([]UIFactorOption, len(f.Levels))
				for j := range f.Levels {
					p.Factors[i].Levels[j] = UIFactorOption{
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
	dbRecordNoOne       *db.Record
	dbRecordNoTwo       *db.Record
}

func (rsr *RecordsSelectResponse) CheckRecords(uiUserData *UIUserData, c appengine.Context) {
	rsr.VaryingFactorIds = make([]string, len(appConfig.CovPhase.FactorsOrder))
	rsr.CountVaryingFactors = 0
	var CurrentFactorId = uiUserData.User.CurrentFactorId
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

type SelectedFactor struct {
	FactorId        string
	SelectedLevelId string
}

type UICovPromptDynamicText struct {
	previousResponse *Response
	promptConfig     PromptConfig
	state            *CovPhaseState
}

func (ps *UICovPromptDynamicText) String() string {
	t := template.Must(template.New("display").Parse(ps.promptConfig.Text))
	var doc bytes.Buffer
	err := t.Execute(&doc, ps.state)
	if err != nil {
		log.Println("executing template:", err)
	}
	display := doc.String()

	return display
}

// This method should only update records select
// Unless if no existing state, than create new one, otherwise, only
// update records select
func (cp *CovPrompt) updateStateRecords(uiUserData *UIUserData, r *RecordsSelectResponse) {
	cp.updateState(uiUserData)
	if cp.state != nil {
		cp.state.RecordNoOne = createRecordStateFromDB(r.dbRecordNoOne, r.RecordNoOne)
		cp.state.RecordNoTwo = createRecordStateFromDB(r.dbRecordNoTwo, r.RecordNoTwo)
		uiUserData.State = cp.state
	}
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
		cp.state.Username = uiUserData.User.Username
		cp.state.Screenname = uiUserData.User.Screenname
		fid := uiUserData.User.CurrentFactorId
		if factorConfigMap[fid] != nil {
			cp.state.TargetFactor = &CovFactorState{FactorName: factorConfigMap[fid].Name, FactorId: fid}
		}
	}
	uiUserData.State = cp.state
}

func createRecordStateFromDB(r *db.Record, sf []SelectedFactor) *RecordState {
	rs := &RecordState{}
	if r != nil {
		rs.RecordName = r.Name
		rs.RecordNo = r.RecordNo
		rs.FactorLevels = make(map[string]*CovFactorState)
		for _, v := range sf {
			f := GetFactorConfig(v.FactorId)
			cfs := &CovFactorState{
				FactorName:    f.Name,
				FactorId:      f.Id,
				SelectedLevel: v.SelectedLevelId,
				OppositeLevel: GetFactorOppositeLevel(f.Id, v.SelectedLevelId)}
			rs.FactorLevels[v.FactorId] = cfs
		}
	} else {
		rs.RecordName = ""
		rs.RecordNo = ""
		rs.FactorLevels = make(map[string]*CovFactorState)
	}
	return rs
}
