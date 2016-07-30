/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var MSG_ROBOT = 'robot';
var MSG_HUMAN = 'student';
var DisplayText = {};
DisplayText[MSG_ROBOT] = 'Researcher';

var Dialog = React.createClass({
  getInitialState: function() {
    state = {mode: 0}
    var user = this.props.user;
    var history = user.getHistory() ? user.getHistory() : {};
    state.isNewUser = history.length == 0;
    state.welcomeText = state.isNewUser ? "Welcome to the Mission!" : "Welcome back!";
    state.oldHistory = history;
    return state;
  },

  changeState: function() {
    var user = this.props.user;
    action = user.getAction();
    var app = this.props.app;
    this.setState({mode: 0});

  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var newHistory = user.getHistory() ? user.getHistory().slice(state.oldHistory.length) : {};
    var messages = newHistory.map(
        function(message, i) {
          return <div key={i}>
                  <Message texts={message.Texts} mtype={message.Mtype} app={app}/>
                </div>
        })
    var prompt = user.getPrompt();
    var welcomeText = this.state.welcomeText;

    if ((!prompt) || (Object.keys(prompt).length == 0)) {
        return  <div className="dialog">
                  <Title user={user} welcomeText={welcomeText}/>
                  <OldHistory user={user} oldHistory={state.oldHistory}/>
                  {messages}
                </div>;
    } else {
        return  <div className="dialog">
                  <Title user={user} welcomeText={welcomeText}/>
                  <OldHistory user={user} oldHistory={state.oldHistory}/>
                  {messages}
                  <Input user={user} prompt={prompt} onComplete={this.changeState} app={app}/>
                </div>;
    }

  },
});


// Render the title of the chat window
var OldHistory = React.createClass({
  getInitialState: function() {
    return {showMessages: false};
  },

  changeState: function() {
    this.state.showMessages = !this.state.showMessages;
    this.setState(this.state); // This call triggers re-rendering
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var oldHistory = this.props.oldHistory ? this.props.oldHistory : {};
    var messages = oldHistory.map(
        function(message, i) {
          return  <div key={i}>
                    <Message texts={message.Texts} mtype={message.Mtype} user={user} delay={false}/>
                  </div>;
        })
    if (messages.length > 0) {
      if (state.showMessages) {
        return <div>
                  <button type="submit" onClick={this.changeState}>Click to Hide old chat history</button>
                  {messages}
               </div>;
      } else {
        return <div>
                  <button type="submit" onClick={this.changeState}>Click to show old chat history</button>
               </div>;
      }
    }
    return <div></div>;
  }
});


// Render the title of the chat window
var Title = React.createClass({

  render: function() {
    var user = this.props.user;
    var human = user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();
    var welcomeText = this.props.welcomeText;
    return  <div className="researcher">
              <div className="name">Researcher</div>
              <div className="message">
                Hello {human}.<br/>
                {welcomeText}<br/>
              </div>
            </div>;
  }
});


// Render each message
var MessageText = React.createClass({
  componentDidMount: function() {
    var e = React.findDOMNode(this);
    e.scrollIntoView();
  },
    
  componentDidUpdate: function(prevProps, prevState) {
    var e = React.findDOMNode(this);
    e.scrollIntoView();
  },

  render: function() {
    var message = this.props.message;
    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    if (message.Text) {
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
    return <div></div>;
  }
});

// Render each message
var Message = React.createClass({
  getInitialState: function() {
    return {count: 0};
  },

  componentDidMount: function() {
    var e = React.findDOMNode(this);
    e.scrollIntoView();
  },
    
  componentDidUpdate: function(prevProps, prevState) {
    var e = React.findDOMNode(this);
    e.scrollIntoView();
  },

  triggerDelay: function() {
    var d = DELAY_PROMPT_TIME_SHORT;
    if (this.props.texts[this.state.count].length > 100) {
      d = DELAY_PROMPT_TIME_LONG;
    }
    this.state.count++;
    this.state.interval = window.setInterval(this.unTriggerDelay, d);
  },

  unTriggerDelay: function() {
    window.clearInterval(this.state.interval);
    this.setState({})
    if (this.state.count == this.props.texts.length) {
      // TODO - This implicitly only allow action UI to be shown
      // when <Message> has delay = true
      this.props.app.showAction(true);
    }
  },

  render: function() {
    var texts = this.props.texts;
    var delay = this.props.delay;
    var mtype = this.props.mtype;
    var lastCount = 0;
    if (delay) {
      if (this.state.count < texts.length) {
        this.triggerDelay();
      }
      lastCount = this.state.count;
    } else {
      lastCount = texts.length;
    }

    messages = texts.slice(0, lastCount).map(
        function(text, i) {
          var message = {};
          message.Mtype = mtype;
          message.Text = text;
          return  <div className="chat" key={i}>
                    <MessageText message={message} user={user}/>
                  </div>;
        })

    return  <div>{messages}</div>;
  }
});

// Renter input window
var Input = React.createClass({

  getInitialState: function() {
    return {enabled: false, passthrough: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  handleChange: function(event) {
    this.setState({enabled:true});
  },

  handleEnter: function(event) {
    if (!event.shiftKey) {
      if (event.which == 13) {  // "Enter" key was pressed.
        this.handleSubmit(event);
      }
    }
  },

  triggerPassThrough: function() {
    this.state.passthrough = true;
    var d = DELAY_PROMPT_TIME_SHORT;
    if (this.props.prompt.Texts[0].length > 100) {
      d = DELAY_PROMPT_TIME_LONG;
    }
    this.state.interval = window.setInterval(this.unTriggerPassThrough, d);
  },

  unTriggerPassThrough: function() {
    window.clearInterval(this.state.interval);
    this.handleSubmit()
  },

  handleSubmit: function(event) {
    if (event) {
      event.preventDefault();
    }

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
    case UI_PROMPT_STRAIGHT_THROUGH:
      text = RESPONSE_SYSTEM_GENERATED;
      id = RESPONSE_SYSTEM_GENERATED;
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
    var app = this.props.app;
    var prompt = this.props.prompt;
    var texts = prompt.Texts
    var user = this.props.user;

    var promptId = prompt.PromptId;
    var phaseId = user.CurrentPhaseId;

    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    switch (prompt.PromptType) {
    case UI_PROMPT_TEXT:
      return  <div className="chat" key={promptId+user.getHistory().length}>
                <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app}/>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                    className="request">
                      <textarea autoFocus name="dialoginput" onKeyDown={this.handleEnter}></textarea>
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
          return <div key={i}><PromptOption option={option}/></div>;
        });

      return  <div className="chat" key={promptId+user.getHistory().length}>
                <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app}/>
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
    case UI_PROMPT_STRAIGHT_THROUGH:
      if (this.state.passthrough) {
        this.state.passthrough = false
      } else {
        this.triggerPassThrough();
      }
      return  <div className="chat" key={promptId+user.getHistory().length}>
                <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app}/>
                <div className="human">
                  <div className="name">{human}</div>
                  <div className="form">
                    <form id="dialogForm" onSubmit={this.handleSubmit} className="request">
                      <input type="text" name="dialoginput" disabled/>
                      <br/>
                      <input type="hidden" id="promptId" value={promptId}/>
                      <input type="hidden" id="phaseId" value={phaseId}/>
                      <button type="submit" id="submitButton" disabled={!this.isEnabled()}>Enter</button>
                    </form>
                  </div>
                </div>
              </div>;    
    default:
      return  <div className="chat" key={promptId+user.getHistory().length}>
                <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app}/>
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







