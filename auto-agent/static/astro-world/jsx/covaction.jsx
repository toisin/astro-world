/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var CovAction = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();

    switch (prompt.UIActionModeId) {
      case "NO_UIACTION":
        return  <div className="action">
              <h3>Investigating Factor: <b>Fitness</b></h3>

              <table>
                <tbody>
                  <tr>
                    <td>&nbsp;</td>
                    <td colSpan="3" className="question">Which record would you like to see?</td>
                  </tr>
                  <tr>
                          <td colSpan="4">
                          <ActionInput user={user} prompt={prompt} onComplete={app.changeState}/></td>
                  </tr>
                </tbody>
              </table>
                </div>;
      case "RECORD_SELECT_ONE":
        return <div></div>;
      case "RECORD_SELECT_TWO":
        return <TwoRecordSelection user={user} prompt={prompt} onComplete={app.changeState}/>;
      case "ONE_RECORD_PERFORMANCE":
        return <OneRecordPerformance user={user} prompt={prompt} onComplete={app.changeState}/>;
      default:
        return <div></div>;
    }
    //   case "chart":
    //     return <div></div>
    //   case "prediction":
    //     return <div></div>

    //     break;
    // }
  }
});


var ActionInput = React.createClass({

  getInitialState: function() {
    return {enabled: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  handleChange: function(event) {
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault(); // default might be to follow a link, instead, takes control over the event

    var user = this.props.user;
    var onComplete = this.props.onComplete;
    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");
    e = f.elements['covactioninput'];
    var value = e ? e.value : "";
    e.value = "";
    var text, id;

    switch (user.CurrentUIPrompt.Type) {
    case UI_PROMPT_MC:
      var options = user.CurrentUIPrompt.Options;
      for (i = 0; i < options.length; i++) {
        if (options[i].ResponseId == value) {
          text = options[i].Text;
          id = value;
          break;
        }
      }
      break;
    case UI_PROMPT_TEXT:
      text = value;
      id = user.CurrentUIPrompt.ResponseId
      break;
    }

    var response = {};
    response.text = text;
    response.id = id;
    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
    this.setState({mode: 0, enabled:false});
  },

  render: function() {
    var prompt = this.props.prompt;
    var user = this.props.user;

    var promptId = prompt.PromptId;
    var phaseId = user.CurrentPhaseId;
    var type = prompt.Type;
    var human = user.getScreenname() ? user.getScreenname() : user.getUsername();

    if (type == UI_PROMPT_TEXT) {
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange} className="request">
                      <textarea name="covactioninput"></textarea>
                      <br/>
                      <input type="hidden" id="promptId" value={promptId}/>
                      <input type="hidden" id="phaseId" value={phaseId}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    }
    if (type == UI_PROMPT_YES_NO) {
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      <label>
                        <input type="radio" name="covactioninput" value="Yes"/>
                        Yes
                      </label>
                      <label>
                        <input type="radio" name="covactioninput" value="No"/>
                        No
                      </label>
                      <br/>
                      <input type="hidden" id="promptId" value={promptId}/>
                      <input type="hidden" id="phaseId" value={phaseId}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    }
    if (type == UI_PROMPT_MC) {
      if (!prompt.Options) {
        console.error("Error: MC Prompt without options!");    
        return <div></div>;
      }
      var options = prompt.Options.map(
        function(option, i) {
          return <ActionPromptOption option={option} key={i}/>;
        });

      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      {options}
                      <br/>
                      <input type="hidden" id="promptId" value={promptId}/>
                      <input type="hidden" id="phaseId" value={phaseId}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    }
    if (type == UI_PROMPT_NO_RESPONSE) {
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
              </div>;
    }
    // console.error("Error: Unknown prompt type!");  
    return <div></div>;  
  },
});

var ActionPromptOption = React.createClass({

  render: function() {
    var option = this.props.option;
      return  <label>
                <input type="radio" name="covactioninput" value={option.ResponseId}/>
                {option.Text}
              </label>
  },
});




