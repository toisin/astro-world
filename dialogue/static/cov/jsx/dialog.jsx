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
          return  <div className="dialog" key={i}>
                    <Message message={message} user={user}/>
                  </div>;
        }) : {};
    var prompt = user.getPrompt();

    if ((!prompt) || (Object.keys(prompt).length == 0)) {
        return  <div className="dialog">
                <div className="chat">
                  <Title user={user}/>
                  {messages}
                </div></div>;
    } else {
        return  <div className="dialog">
                <div className="chat">
                  <Title user={user}/>
                  {messages}
                  <Input user={user} prompt={prompt} onComplete={app.changeState}/>
                </div></div>;
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
    var e = document.getElementById("workflowStateID");
    var w = e ? e.value : "";
    var f = document.getElementById("inputForm");
    e = f.elements['input'];
    var i = e ? e.value : "";
    e.value = "";
    user.submitResponse(w, i, onComplete);
    this.setState({mode: 0, enabled:false});
  },

  render: function() {
    var prompt = this.props.prompt;

    var workflowStateID = prompt.WorkflowStateID;
    var type = prompt.Ptype;
    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    if (type == PROMPT_TEXT) {
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="inputForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      <textarea name="input"></textarea>
                      <br/>
                      <input type="hidden" id="workflowStateID" value={workflowStateID}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    }
    if (type == PROMPT_YES_NO) {
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="inputForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      <label>
                        <input type="radio" name="input" value="Yes"/>
                        Yes
                      </label>
                      <label>
                        <input type="radio" name="input" value="No"/>
                        No
                      </label>
                      <br/>
                      <input type="hidden" id="workflowStateID" value={workflowStateID}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    }
    if (type == PROMPT_MC) {
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
                    <form id="inputForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      {options}
                      <br/>
                      <input type="hidden" id="workflowStateID" value={workflowStateID}/>
                      <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;
    }
    if (type == PROMPT_NO_RESPONSE) {
      return  <div>
                <div className="researcher">
                  <div className="name">{DisplayText[MSG_ROBOT]}</div>
                  <div className="message">{prompt.Text}</div>
                </div>
              </div>;
    }
    if (type == PROMPT_END) {
      return  <div></div>;
    }
    console.error("Error: Unknown prompt type!");  
    return <div></div>;  
  },
});

var PromptOption = React.createClass({

  render: function() {
    var option = this.props.option;
      return  <label>
                <input type="radio" name="input" value={option.Value}/>
                {option.Label}
              </label>
  },
});







