/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var MSG_ROBOT = 'robot';
var MSG_HUMAN = 'student';
var DisplayText = {};
DisplayText[MSG_ROBOT] = 'Researcher';

var Dialog = React.createClass({
  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var messages = user.getHistory() ? user.getHistory().map(
        function(message, i) {
          return  <div className="chat" key={i}>
                    <Message message={message} user={user}/>
                  </div>;
        }) : {};
    var prompt = user.getPrompt();

    if ((!prompt) || (Object.keys(prompt).length == 0)) {
        return  <div className="chat">
                  <Title user={user}/>
                  {messages}
                </div>;
    } else {
        return  <div className="chat">
                  <Title user={user}/>
                  {messages}
                  <Input user={user} prompt={prompt} onComplete={app.changeState}/>
                </div>;
    }

  },
});


// Render the title of the chat window
var Title = React.createClass({
  render: function() {
    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    return  <div className="researcher">
              <div className="name">Researcher</div>
              <div className="message">
                Hello {human}<br/>
                Welcome to the Challenge
              </div>
            </div>;
  }
});


// Render each message
var Message = React.createClass({
  render: function() {
    var message = this.props.message;
    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    if (message.Mtype == MSG_ROBOT) {
      return  <div className="researcher">
                <div className="name">{DisplayText[MSG_ROBOT]}</div>
                <div className="message">{message.Text}</div>
              </div>;
    } else if (message.Mtype == MSG_HUMAN) {
      return  <div className="human">
                <div className="name">{human}</div>
                <div className="message">{message.Text}</div>
              </div>;
    }
    console.error("Unknown sender!", error);
    return  <div className="researcher">
              <div className="message">{this.props.message.Text}</div>
            </div>;
  }
});

// Renter input window
var Input = React.createClass({

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
    event.preventDefault();

    var user = this.props.user;
    var onComplete = this.props.onComplete;
    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("dialogForm");
    e = f.elements['dialoginput'];
    var value = e ? e.value : "";
    e.value = "";
    var text, id;
    var options = user.CurrentUIPrompt.Options;

    switch (user.CurrentUIPrompt.PromptType) {
    case UI_PROMPT_MC:
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
      id = options[0].ResponseId
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

    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    switch (prompt.PromptType) {
    case UI_PROMPT_TEXT:
      return  <div  className="chat">
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      <textarea name="dialoginput"></textarea>
                      <br/>
                      <input type="hidden" id="promptId" value={promptId}/>
                      <input type="hidden" id="phaseId" value={phaseId}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    case UI_PROMPT_MC:
      if (!prompt.Options) {
        console.error("Error: MC Prompt without options!");    
        return <div></div>;
      }
      var options = prompt.Options.map(
        function(option, i) {
          return <PromptOption option={option} key={i}/>;
        });

      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
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
    default:
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
              </div>;
    }
  },
});

var PromptOption = React.createClass({

  render: function() {
    var option = this.props.option;
      return  <label>
                <input type="radio" name="dialoginput" value={option.ResponseId}/>
                {option.Text}
              </label>
  },
});







