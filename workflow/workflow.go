package workflow

const (
	PROMPT_NO_RESPONSE = "NO_RESPONSE"
	PROMPT_TEXT = "TEXT"
	PROMPT_YES_NO = "YES_NO"
	PROMPT_MC = "MC"

	END_STATE = "COMPLETE"
)

var stateMap = make(map[string]State)
var variableMap = make(map[string]Variable)

func InitWorkflowMaps() {
	variableMap["Y"] = Variable{ Text: "Performance"}
	variableMap["X1"] = Variable{"Health Index", true, false}
	variableMap["X2"] = Variable{"Height", false, false}
	variableMap["X3"] = Variable{"Weight", true, false}
	variableMap["X4"] = Variable{"Laughters", true, false}

	// stateMap["1"] = &MCPromptState{PROMPT_NO_RESPONSE, "1", "Let's get started!", ""}
	stateMap["1"] = &MCPromptState{
						PROMPT_MC,
						"1",
						"Let's get started! What feature have you investigated?",
						"",
					 	[]Option{
					 		Option{variableMap["X1"].Text,"X1"},
					 		Option{variableMap["X2"].Text,"X2"},
					 		Option{variableMap["X3"].Text,"X3"},
					 		Option{variableMap["X4"].Text,"X4"}}}
	stateMap["2"] = &TextPromptState{PROMPT_YES_NO, "2", "Do you think it makes a difference?", "1", "3"}
	//stateMap["2"] = &LogicPromptState{PROMPT_YES_NO, "2", "Do you think it makes a difference?", "1"}
	stateMap["3"] = &TextPromptState{
						PROMPT_TEXT,
						"3",
						"When %X1 goes up, what happens to %Y?",
						"2",
						"4"}
	stateMap["4"] = &TextPromptState{PROMPT_TEXT, "4", "What did you find out about %X1?", "3", "5"}
	stateMap["5"] = &TextPromptState{PROMPT_TEXT, "5", "How do you know?", "4", "6"}
	stateMap["6"] = &TextPromptState{PROMPT_TEXT, "6", "Which records show you are right?", "5", "COMPLETE"}
	// stateMap["8"] = &MCPromptState{"8", "What level is your?", ""}
	// stateMap["9"] = &MCPromptState{"9", "How do you know?", ""}
	// stateMap["10"] = &MCPromptState{"10", "How do you know?", ""}
	// stateMap["11"] = &MCPromptState{"11", "How do you know?", ""}
}

func GetVariableMap() map[string]Variable{
	return variableMap
}

func GetStateMap() map[string]State{
	return stateMap
}

func GetFirstState() State {
	return stateMap["1"]
}

type Variable struct {
	Text string
	IsCausal bool
	IsPostiveCorr bool
}


type State interface {
	Display() string
	GetNextStateId() string
	GetId() string
}

type TextPromptState struct {
	Ptype string
	WorkflowStateID string
	Text string
	LastStateId string
	NextStateId string
}

func (ps *TextPromptState) GetId() string {
	return ps.WorkflowStateID
}

func (ps *TextPromptState) Display() string {
	return ps.Text
}

func (ps *TextPromptState) GetNextStateId() string {
	return ps.NextStateId
}

type MCPromptState struct {
	Ptype string
	WorkflowStateID string
	Text string
	LastStateId string
	Options []Option
}

type Option struct {
	Label string
	Value string
}

func (ps *MCPromptState) GetId() string {
	return ps.WorkflowStateID
}

func (ps *MCPromptState) Display() string {
	return ps.Text
}

func (ps *MCPromptState) GetNextStateId() string {
	if ps.WorkflowStateID == "1" {
		return "2"
	}
	return ""
}

//TODO Not sure if LogicState is needed. Not used

type LogicState struct {
	ptype string
	id string
	text string
	lastStateId string
}

func (ls *LogicState) Id() string {
	return ls.id
}

func (ls *LogicState) Text() string {
	return ls.text
}

func (ls *LogicState) LastStateId() string {
	return ls.lastStateId
}

func (ls *LogicState) Ptype() string {
	return ls.ptype
}



