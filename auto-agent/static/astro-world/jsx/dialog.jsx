/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var MSG_ROBOT = 'robot';
var MSG_HUMAN = 'student';
var DisplayText = {};
DisplayText[MSG_ROBOT] = 'Researcher';

var Dialog = React.createClass({
  getInitialState: function() {
    var state = {mode: 0, UIAction:""}
    var user = this.props.user;
    var history = user.getHistory() ? user.getHistory() : {};
    state.isNewUser = history.length == 0;
    state.welcomeText = state.isNewUser ? "Welcome to Astro-world!" : "Welcome back! Let's pick up where we left off.";
    return state;
  },

  changeState: function() {
    var app = this.props.app;
    if (this.props.user.getAction().UIActionModeId != this.state.UIAction) {
      this.state.UIAction = this.props.user.getAction().UIActionModeId;
      app.changeState();
    } else {
      this.setState(this.state);
    }
  },

  showAction: function() {
    var app = this.props.app;
    app.showAction();
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var history = user.getHistory() ? user.getHistory() : {};
    var oldHistoryLength;
    var oldHistory, newHistory;
    if (user.getArchiveHistoryLength() <= MESSAGE_COUNT_LIMIT) {
      oldHistoryLength = user.getArchiveHistoryLength();
      oldHistory = history.slice(0, oldHistoryLength);
      newHistory = history.slice(oldHistoryLength);
    } else {
      oldHistoryLength = user.getArchiveHistoryLength() - history[0].MessageNo + 1;
      oldHistory = user.getHistory() ? user.getHistory().slice(0, oldHistoryLength) : {};
      newHistory = user.getHistory() ? user.getHistory().slice(oldHistoryLength) : {};
    }
    var messages = newHistory.map(
        function(message, i) {
          return <div key={i}>
                  <Message texts={message.Texts} mtype={message.Mtype} app={app} user={user}/>
                </div>
        })
    var prompt = user.getPrompt();
    var welcomeText = this.state.welcomeText;

    if ((!prompt) || (Object.keys(prompt).length == 0)) {
        return  <div>
                  <Title user={user} welcomeText={welcomeText}/>
                  <OldHistory user={user} oldHistory={oldHistory}/>
                  {messages}
                </div>;
    } else {
        return  <div>
                  <Title user={user} welcomeText={welcomeText}/>
                  <OldHistory user={user} oldHistory={oldHistory}/>
                  {messages}
                  <Prompt user={user} prompt={prompt} onShowInput={this.showAction} onComplete={this.changeState} app={app} key={prompt.PromptId}/>
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
                    <Message texts={message.Texts} mtype={message.Mtype} user={user}/>
                  </div>;
        })
    if (messages.length > 0) {
      if (state.showMessages) {
        return <div>
                  {messages}
                  <button type="submit" onClick={this.changeState}>Click to Hide old chat history</button>
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
              <div className="name">{DisplayText[MSG_ROBOT]}</div>
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
    var e = ReactDOM.findDOMNode(this);
    e.scrollIntoView();
  },
    
  componentDidUpdate: function(prevProps, prevState) {
    var e = ReactDOM.findDOMNode(this);
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
        return  <div className="user">
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
      }
      if (this.state.count == 2 || this.state.count == texts.length) {
        this.state.complete = true;
        this.setState(this.state);
        // This should only be necessary if delay is turned on
        // otherwise, everything would have been rendered.
        if (this.props.onComplete) {
          this.props.onComplete();
        }
      }
    }
    // // TODO This should not be needed because everything should have been rendered.
    // } else {
    //   this.state.complete = true;
    //   if (this.props.onComplete) {
    //       this.props.onComplete();
    //   }
    // }
  },

  componentDidMount: function() {
    this.refreshAfterDelay();
  },
    
  componentDidUpdate: function(prevProps, prevState) {
    this.refreshAfterDelay();
  },

  triggerDelay: function() {
    if (this.state.count >= this.props.texts.length) {
      return
    }
    var d = DELAY_PROMPT_TIME_SHORT;
    if (this.props.texts[this.state.count - 1].length > LONG_PROMPT_SIZE) {
      d = DELAY_PROMPT_TIME_LONG;
    } else if (this.props.texts[this.state.count - 1].length < REALLYSHORT_PROMPT_SIZE) {
      d = DELAY_PROMPT_TIME_REALLY_SHORT;
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
    var user = this.props.user;
    var lastCount;

    if (!delay) {
      lastCount = texts.length;
    } else {
      lastCount = this.state.count;
    }

    var messages = texts.slice(0, lastCount).map(
        function(text, i) {
          var message = {};
          message.Mtype = mtype;
          message.Text = text;
          return  <div key={i}>
                    <MessageText message={message} user={user}/>
                  </div>;
        })

    return  <div>{messages}</div>;
  }
});

// Renter Prompt
var Prompt = React.createClass({
  getInitialState: function() {
    return {completePrompt: false};
  },

  handleChange: function(event) {
    this.setState({});
  },

  showInput: function() {
    this.state.completePrompt = true;
    this.setState(this.state);
    if (this.props.onShowInput) {
      this.props.onShowInput();
    }
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

    if (this.state.completePrompt) {
      switch (prompt.PromptType) {
      case UI_PROMPT_ENTER_TO_CONTINUE:
      case UI_PROMPT_TEXT:
      case UI_PROMPT_MC:
      case UI_PROMPT_STRAIGHT_THROUGH:
        return  <div key={promptId+user.getHistory().length}>
                  <Message texts={texts} mtype={MSG_ROBOT} app={app} user={user}/>
                  <div className="user">
                    <div className="name">{human}</div>
                    <Input user={user} prompt={prompt} onComplete={onComplete} app={app}/>
                  </div>
                </div>;    
      default:
        return  <div key={promptId+user.getHistory().length}>
                  <Message texts={texts} mtype={MSG_ROBOT} app={app} user={user}/>
                </div>;
      }
    }
    return  <div key={promptId+user.getHistory().length}>
              <Message texts={texts} delay={true} mtype={MSG_ROBOT} app={app} user={user} onComplete={this.showInput}/>
            </div>;
  },
});

var PromptOption = React.createClass({

  render: function() {
    var option = this.props.option;
      return  <label>
                <input type="radio" name="dialoginput" value={option.ResponseId}/>
                {option.Text}&nbsp;&nbsp;&nbsp;
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
    var user = this.props.user;
    var f = document.getElementById("dialogForm");
    var e = f.elements['dialoginput'];
    var value = e ? e.value : "";
    switch (user.CurrentUIPrompt.PromptType) {
    case UI_PROMPT_TEXT:
      if (value.trim() != "") {
        this.setState({enabled:true});
      }
      break;
    default:
        this.setState({enabled:true});
    }
  },

  handleEnter: function(event) {
    if (this.state.enabled) {
      if (!event.shiftKey) {
        if (event.which == 13) {  // "Enter" key was pressed.
          this.handleSubmit(event);
        }
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
      for (var i = 0; i < options.length; i++) {
        if (options[i].ResponseId == value) {
          text = options[i].Text;
          id = value;
          break;
        }
      }
      break;
    case UI_PROMPT_TEXT:
      text = value;
      id = value;
      break;
    case UI_PROMPT_ENTER_TO_CONTINUE:
    case UI_PROMPT_STRAIGHT_THROUGH:
      text = RESPONSE_SYSTEM_GENERATED;
      id = RESPONSE_SYSTEM_GENERATED;
      break;
    }

    var response = {};
    response.text = text;
    response.id = id;
    var jsonResponse = JSON.stringify(response);
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
    case UI_PROMPT_ENTER_TO_CONTINUE:
      return  <div className="form">
                <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                className="request">
                  <input type="hidden" id="dialoginput" disabled/>
                  <input type="hidden" id="promptId" value={promptId}/>
                  <input type="hidden" id="phaseId" value={phaseId}/>
                  <button type="submit" autoFocus>Enter</button>
                </form>
              </div>;
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
          return <PromptOption option={option} key={i}/>;
        });

      return  <div className="form">
                <form id="dialogForm" onSubmit={this.handleSubmit} onChange={this.handleChange}
                className="request">
                  {options}
                  <br/>
                  <input type="hidden" id="promptId" value={promptId}/>
                  <input type="hidden" id="phaseId" value={phaseId}/>
                  <button autoFocus type="submit" disabled={!this.isEnabled()}>Enter</button>
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





