package workflow

import (
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"appengine"
)

// Prompt logics specific to Prediction phase

type PredictionPrompt struct {
	*GenericPrompt
}

func MakePredictionPrompt(p PromptConfig, uiUserData *UIUserData) *PredictionPrompt {
	var n *PredictionPrompt
	n = &PredictionPrompt{}
	n.GenericPrompt = &GenericPrompt{}
	n.GenericPrompt.currentPrompt = n
	n.init(p, uiUserData)
	return n
}

func (cp *PredictionPrompt) ProcessResponse(r string, u *db.User, uiUserData *UIUserData, c appengine.Context) {
	if cp.promptConfig.ResponseType == RESPONSE_END {
		// TODO - how to handle final phase
		uiUserData.State.(*PredictionPhaseState).updateToNextTargetPrediction()
		cp.nextPrompt = cp.generateFirstPromptInNextSequence(uiUserData)
	} else if r != "" {
		dec := json.NewDecoder(strings.NewReader(r))
		pc := cp.promptConfig
		switch pc.ResponseType {
		case RESPONSE_PREDICTION_REQUESTED_FACTORS:
			for {
				var beliefResponse UIMultiFactorsCausalityResponse
				if err := dec.Decode(&beliefResponse); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.updateRequestedFactors(uiUserData, beliefResponse)
				cp.response = &beliefResponse
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
				cp.updateStateCurrentFactorCausal(uiUserData, response.GetResponseId())
				cp.response = &response
			}
			break
		case RESPONSE_PREDICTION_PERFORMANCE:
			for {
				var response SimpleResponse
				if err := dec.Decode(&response); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err)
					log.Fatal(err)
					return
				}
				cp.updateStateCurrentPredictionPerformance(uiUserData, response.GetResponseId())
				cp.response = &response
			}
			break
		case RESPONSE_PREDICTION_NEXT_FACTOR:
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
			cp.updateFirstNextFactor(uiUserData)
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

func (cp *PredictionPrompt) updateRequestedFactors(uiUserData *UIUserData, r UIMultiFactorsCausalityResponse) {
	// Re-use updateMultiFactorsCausalityResponse, assume students requested only causal factors
	cp.GenericPrompt.updateMultiFactorsCausalityResponse(uiUserData, r)
	cp.updateFirstNextFactor(uiUserData)
}

// This is to update the next factor of causal factor not requested
func (cp *PredictionPrompt) updateFirstNextFactor(uiUserData *UIUserData) {
	if cp.state != nil {
		s := cp.state.(*PredictionPhaseState)

		// clear previous state in case if there are left over session data
		// only happens is screen was refreshed before intro_predict has completed
		s.RequestedFactors = make([]UIFactor, 0)
		s.DisplayFactors = make([]UIFactor, 0)
		s.DisplayFactorsReady = false

		wrongFactors := s.Beliefs.IncorrectFactors

		if len(wrongFactors) > 0 {
			// there are more than 1 wrong factors

			beliefNoncausalCount := 0
			for _, v := range wrongFactors {
				if !v.IsBeliefCausal {
					// count the number of wrong factors that were believed to be non-causal,
					// i.e. they are actually causal factors
					beliefNoncausalCount++
				}
			}

			// Enter this if statement to determine if any wrong factors should
			// be presented to walk the student through the chart again
			// Begins with the causal factors that were wrongly not requested,
			// then if there were more than 1 non-causal factors that were wrongly requested,
			// present the chart for that after all of the causal factors were re-examined.
			if beliefNoncausalCount > 0 {
				// there is at least 1 causal factors that were wrongly believed to be non-causal
				factors := make([]UIFactor, beliefNoncausalCount)

				beliefNoncausalCount = 0
				for _, v := range wrongFactors {
					if !v.IsBeliefCausal {
						factors[beliefNoncausalCount] = v
						beliefNoncausalCount++
					}
				}

				s.initMissingCausalFactors(factors)
				// TODO - hard coding the first incorrect factor as target factor
				// maybe too much UI logic. Would be better if it can be triggered
				// by workflow.json
				fid := factors[0].FactorId
				cp.updateStateCurrentFactor(uiUserData, fid)
				return
			} else if (len(wrongFactors) - beliefNoncausalCount) > 0 {
				// there is at least 1 non-causal factors that were wrongly believed to be causal
				factors := make([]UIFactor, len(wrongFactors)-beliefNoncausalCount)

				i := 0
				for _, v := range wrongFactors {
					if v.IsBeliefCausal {
						factors[i] = v
						i++
					}
				}

				s.initRequestedNonCausalFactors(factors)

				if len(factors) > 1 {
					// TODO - hard coding the UI logic here: We present all causal factors plue 1 non-causal
					// to see if students are able to not use the non-causal factor when making prediction,
					// so if only 1 non causal factors was requested, we move on to the prediction.
					// However, if there was more than 1, we set the first one as target factor so
					// that students will be asked about it.

					// TODO - hard coding the first incorrect factor as target factor
					// maybe too much UI logic. Would be better if it can be triggered
					// by workflow.json
					fid := factors[0].FactorId
					cp.updateStateCurrentFactor(uiUserData, fid)
					return
				}
			}
		}

		// There are no more factors to ask the students about.
		// This could mean that students picked only the correct causal factors
		// OR that they picked all correct causal factors plus 1 non-causal factor

		count := 0

		// Check if there was 1 non-causal factor requested
		for _, v := range s.GetContentFactors() {
			if v.IsCausal {
				count++
			}
		}
		tempFactors := make([]UIFactor, count+1)

		count = 0
		for _, v := range s.GetContentFactors() {
			if v.IsCausal {
				// capture all causal factors
				tempFactors[count] = v
				count++
			}
		}

		if len(wrongFactors) > 0 {
			// if there is a non-causal factor requested, include it for display
			tempFactors[count] = wrongFactors[0]
			s.RequestedFactors = tempFactors
		} else {
			// no non-causal factor was requested, add one to display
			s.RequestedFactors = tempFactors
			for _, v := range s.GetContentFactors() {
				if !v.IsCausal {
					tempFactors[count] = v
					break
				}
			}

		}
		s.DisplayFactors = tempFactors
		s.DisplayFactorsReady = true
	}

}

func (cp *PredictionPrompt) updateStateCurrentPredictionPerformance(uiUserData *UIUserData, performanceResponse string) {
	// invoking the initialization methods in the "subclass"
	// in case if they have been overriden
	cp.currentPrompt.updateState(uiUserData)
	s := cp.state.(*PredictionPhaseState)
	s.TargetPrediction.PredictedPerformanceLevel, _ = strconv.Atoi(performanceResponse)
	// if isCausalResponse == "true" {
	// 	targetFactor.IsConcludeCausal = true
	// 	targetFactor.HasConclusion = true
	// } else if isCausalResponse == "false" {
	// 	targetFactor.IsConcludeCausal = false
	// 	targetFactor.HasConclusion = true
	// }
	// // TODO - not sure if it's a good idea
	// // by changing the ContentFactors, we lose track of what the student originally believed
	// tempContentFactors := cp.state.GetContentFactors()
	// tempFactor := tempContentFactors[targetFactor.FactorId]
	// tempFactor.IsBeliefCausal = targetFactor.IsConcludeCausal
	// tempContentFactors[targetFactor.FactorId] = tempFactor

	// allCorrect := true
	// causalFactors := []UIFactor{}
	// incorrectFactors := []UIFactor{}

	// for _, v := range tempContentFactors {
	// 	if v.IsBeliefCausal != v.IsCausal {
	// 		allCorrect = false
	// 		incorrectFactors = append(incorrectFactors, v)
	// 	}
	// 	if v.IsBeliefCausal {
	// 		causalFactors = append(causalFactors, v)
	// 	}
	// }

	// tempBeliefs := cp.state.GetBeliefs()
	// tempBeliefs.IncorrectFactors = incorrectFactors
	// tempBeliefs.CausalFactors = causalFactors
	// tempBeliefs.AllCorrect = allCorrect
	// cp.state.setBeliefs(tempBeliefs)

	// cp.state.setContentFactors(tempContentFactors)
	// cp.state.setTargetFactor(targetFactor)
	uiUserData.State = cp.state
}

// func (cp *PredictionPrompt) updateAppilcant(uiUserData *UIUserData, r UIChartRecordSelectResponse) {
// cp.updateState(uiUserData)
// if cp.state != nil {
// 	s := cp.state.(*ChartPhaseState)
// 	if r.RecordNo != "" {
// 		s.Record = CreateRecordStateFromDB(r.dbRecord)
// 	} else {
// 		s.Record = RecordState{}
// 	}
// 	cp.state = s
// }
// uiUserData.State = cp.state
// }

func (cp *PredictionPrompt) updateState(uiUserData *UIUserData) {
	if uiUserData.State != nil {
		// if uiUserData already have a cp state, use that and update it
		if uiUserData.State.GetPhaseId() == appConfig.PredictionPhase.Id {
			cp.state = uiUserData.State.(*PredictionPhaseState)
		}
	}
	if cp.state == nil {
		cps := &PredictionPhaseState{}
		cps.initContents()
		cp.state = cps
		cp.state.setPhaseId(appConfig.PredictionPhase.Id)
		cp.state.setUsername(uiUserData.Username)
		cp.state.setScreenname(uiUserData.Screenname)

		// TODO - hard coding the first incorrect factor as target factor
		// maybe too much UI logic. Would be better if it can be triggered
		// by workflow.json
		// fid := uiUserData.CurrentFactorId
		// if fid != "" {
		// 	cp.state.setTargetFactor(
		// 		FactorState{
		// 			FactorName: factorConfigMap[fid].Name,
		// 			FactorId:   fid,
		// 			IsCausal:   factorConfigMap[fid].IsCausal})
		// }
	}
	uiUserData.State = cp.state
}
