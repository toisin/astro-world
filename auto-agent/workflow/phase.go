package workflow

type Phase interface {
	GetName() string
}

type Strategy interface {
	GetName() string
}

type Action interface {
	GetModeId() string
}

// type CovPhase struct {
// 	// Type string
// 	// WorkflowStateID string
// 	// Text string
// 	// LastStateId string
// 	// NextStateId string
// }

// func (p *CovPhase) GetName() string {
// 	return "CovPhase";
// }

// type Strategy struct {
// 	ParentPhase Phase
// 	getName() string
// }

// type Prompt struct {
// 	ParentPhase Phase
// 	ParentStrategy Strategy
// 	Handler PromtHandler
// 	ExpectedResponses []ExpectedResponse
// 	ActualResponse
// 	PreviousPrompt Prompt
	
// }
