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
    state.welcomeText = state.isNewUser ? "Welcome to Astro-world!" : "Welcome back!";
    state.oldHistory = history;
    return state;
  },

  changeState: function() {
    var app = this.props.app;
    app.changeState();
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var newHistory = user.getHistory() ? user.getHistory().slice(state.oldHistory.length) : {};
    var messages = newHistory.map(
        function(message, i) {
          return <div key={i}>
                  <Message texts={message.Texts} mtype={message.Mtype} app={app} user={user}/>
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
                  <Prompt user={user} prompt={prompt} onComplete={this.changeState} app={app}/>
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
    return {count: 1, complete: false};
  },

  refreshAfterDelay: function() {
    var texts = this.props.texts;
    if ((this.props.delay) && !this.state.complete) {
      if (this.state.count < texts.length) {
        this.triggerDelay();
      } else {
        this.state.complete = true;
        this.setState(this.state);
        // This should only be necessary if delay is turned on
        // otherwise, everything would have been rendered.
        if (this.props.onComplete) {
          this.props.onComplete();
        }
      }
    }
  },

  componentDidMount: function() {
    this.refreshAfterDelay();
  },
    
  componentDidUpdate: function(prevProps, prevState) {
    this.refreshAfterDelay();
  },

  triggerDelay: function() {
    var d = DELAY_PROMPT_TIME_SHORT;
    if (this.props.texts[this.state.count].length > 100) {
      d = DELAY_PROMPT_TIME_LONG;
    }
    this.state.interval = window.setInterval(this.unTriggerDelay, d);
  },

  unTriggerDelay: function() {
    window.clearInterval(this.state.interval);
    this.state.count++;
    this.setState(this.state);
  },

  render: function() {
    var texts = this.props.texts;
    var delay = this.props.delay;
    var mtype = this.props.mtype;
    var lastCount;

    if (!delay) {
      lastCount = texts.length;
    } else {
      lastCount = this.state.count;
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

// Renter Prompt
var Prompt = React.createClass({
  getInitialState: function() {
    return {};
  },

  handleChange: function(event) {
    this.setState({});
  },

  render: function() {
    var app = this.props.app;
    var prompt = this.props.prompt;
    var texts = prompt.Texts;
    var user = this.props.user;
    var onComplete = this.props.onComplete;

    var promptId = prompt.PromptId;
    var phaseId = user.CurrentPhaseId;

    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    switch (prompt.PromptType) {
    case UI_PROMPT_TEXT:
    case UI_PROMPT_MC:
    case UI_PROMPT_STRAIGHT_THROUGH:
      return  <div className="chat" key={promptId+user.getHistory().length}>
                <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app} user={user} onComplete={onComplete}/>
                <div className="human">
                  <div className="name">{human}</div>
                  <Input user={user} prompt={prompt} onComplete={onComplete} app={app}/>
                </div>
              </div>;    
    default:
      return  <div className="chat" key={promptId+user.getHistory().length}>
                <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app} user={user} onComplete={onComplete}/>
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


// Renter input window
var Input = React.createClass({
  getInitialState: function() {
    return {enabled: false, passthrough: true};
  },

  componentDidMount: function() {
    if (this.triggerSubmit()) {
      this.handleSubmit()
    }
  },
    
  componentDidUpdate: function(prevProps, prevState) {
    if (this.triggerSubmit()) {
      this.handleSubmit()
    }
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

  triggerSubmit: function() {
    if (this.props.prompt.PromptType == UI_PROMPT_STRAIGHT_THROUGH) {
      if (this.state.passthrough) {
        this.state.passthrough = false;
        return true;
      }
    }
    return false;
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
  },

  render: function() {
    var app = this.props.app;
    var prompt = this.props.prompt;
    var texts = prompt.Texts;
    var user = this.props.user;
    var onComplete = this.props.onComplete;

    var promptId = prompt.PromptId;
    var phaseId = user.CurrentPhaseId;

    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    switch (prompt.PromptType) {
    case UI_PROMPT_TEXT:
      return  <div className="form">
                <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                className="request">
                  <textarea autoFocus name="dialoginput" onKeyDown={this.handleEnter}></textarea>
                  <br/>
                  <input type="hidden" id="promptId" value={promptId}/>
                  <input type="hidden" id="phaseId" value={phaseId}/>
                  <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                </form>
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

      return  <div className="form">
                <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                className="request">
                  {options}
                  <br/>
                  <input type="hidden" id="promptId" value={promptId}/>
                  <input type="hidden" id="phaseId" value={phaseId}/>
                  <button type="submit" disabled={!this.isEnabled()}>Enter</button>
                </form>
              </div>;
    case UI_PROMPT_STRAIGHT_THROUGH:
      return <div className="form">
              <form id="dialogForm" onSubmit={this.handleSubmit} className="request">
                <input type="text" name="dialoginput" disabled/>
                <br/>
                <input type="hidden" id="promptId" value={promptId}/>
                <input type="hidden" id="phaseId" value={phaseId}/>
                <button type="submit" id="submitButton" disabled={!this.isEnabled()}>Enter</button>
              </form>
            </div>;
    default:
      return  <div></div>;
    }
  },
});





