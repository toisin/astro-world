/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var MSG_ROBOT = 'robot';
var MSG_HUMAN = 'student';
var DisplayText = {};
DisplayText[MSG_ROBOT] = 'Researcher';

var Dialog = React.createClass({displayName: "Dialog",
  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var messages = user.getHistory() ? user.getHistory().map(
        function(message, i) {
          return  React.createElement("div", {className: "dialog", key: i}, 
                    React.createElement(Message, {message: message, user: user})
                  );
        }) : {};
    var prompt = user.getPrompt();

    if ((!prompt) || (Object.keys(prompt).length == 0)) {
        return  React.createElement("div", {className: "dialog"}, 
                React.createElement("div", {className: "chat"}, 
                  React.createElement(Title, {user: user}), 
                  messages
                ));
    } else {
        return  React.createElement("div", {className: "dialog"}, 
                React.createElement("div", {className: "chat"}, 
                  React.createElement(Title, {user: user}), 
                  messages, 
                  React.createElement(Input, {user: user, prompt: prompt, onComplete: app.changeState})
                ));
    }

  },
});


// Render the title of the chat window
var Title = React.createClass({displayName: "Title",
  render: function() {
    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    return  React.createElement("div", {className: "researcher"}, 
              React.createElement("div", {className: "name"}, "Researcher"), 
              React.createElement("div", {className: "message"}, 
                "Hello ", human, React.createElement("br", null), 
                "Welcome to the Challenge"
              )
            );
  }
});


// Render each message
var Message = React.createClass({displayName: "Message",
  render: function() {
    var message = this.props.message;
    var human = this.props.user.getScreenname() ? this.props.user.getScreenname() : this.props.user.getUsername();

    if (message.Mtype == MSG_ROBOT) {
      return  React.createElement("div", {className: "researcher"}, 
                React.createElement("div", {className: "name"}, DisplayText[MSG_ROBOT]), 
                React.createElement("div", {className: "message"}, message.Text)
              );
    } else if (message.Mtype == MSG_HUMAN) {
      return  React.createElement("div", {className: "human"}, 
                React.createElement("div", {className: "name"}, human), 
                React.createElement("div", {className: "message"}, message.Text)
              );
    }
    console.error("Unknown sender!", error);
    return  React.createElement("div", {className: "researcher"}, 
              React.createElement("div", {className: "message"}, this.props.message.Text)
            );
  }
});

// Renter input window
var Input = React.createClass({displayName: "Input",

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
      return  React.createElement("div", null, 
                React.createElement("div", {className: "researcher"}, 
                  React.createElement("div", {className: "name"}, DisplayText[MSG_ROBOT]), 
                  React.createElement("div", {className: "message"}, prompt.Text)
                ), 
                React.createElement("div", {className: "human"}, 
                  React.createElement("div", {className: "name"}, human), 
                  React.createElement("div", {className: "form"}, 
                    React.createElement("form", {id: "inputForm", onSubmit: this.handleSubmit, onChange: this.handleChange, 
                    className: "request"}, 
                      React.createElement("textarea", {name: "input"}), 
                      React.createElement("br", null), 
                      React.createElement("input", {type: "hidden", id: "workflowStateID", value: workflowStateID}), 
                      React.createElement("button", {type: "submit", disabled: !this.isEnabled()}, "Enter")
                    )
                  )
                )
              );
    }
    if (type == PROMPT_YES_NO) {
      return  React.createElement("div", null, 
                React.createElement("div", {className: "researcher"}, 
                  React.createElement("div", {className: "name"}, DisplayText[MSG_ROBOT]), 
                  React.createElement("div", {className: "message"}, prompt.Text)
                ), 
                React.createElement("div", {className: "human"}, 
                  React.createElement("div", {className: "name"}, human), 
                  React.createElement("div", {className: "form"}, 
                    React.createElement("form", {id: "inputForm", onSubmit: this.handleSubmit, onChange: this.handleChange, 
                    className: "request"}, 
                      React.createElement("label", null, 
                        React.createElement("input", {type: "radio", name: "input", value: "Yes"}), 
                        "Yes"
                      ), 
                      React.createElement("label", null, 
                        React.createElement("input", {type: "radio", name: "input", value: "No"}), 
                        "No"
                      ), 
                      React.createElement("br", null), 
                      React.createElement("input", {type: "hidden", id: "workflowStateID", value: workflowStateID}), 
                      React.createElement("button", {type: "submit", disabled: !this.isEnabled()}, "Enter")
                    )
                  )
                )
              );
    }
    if (type == PROMPT_MC) {
      if (!prompt.Options) {
        console.error("Error: MC Prompt without options!");    
        return React.createElement("div", null);
      }
      var options = prompt.Options.map(
        function(option, i) {
          return React.createElement(PromptOption, {option: option, key: i});
        });

      return  React.createElement("div", null, 
                React.createElement("div", {className: "researcher"}, 
                  React.createElement("div", {className: "name"}, DisplayText[MSG_ROBOT]), 
                  React.createElement("div", {className: "message"}, prompt.Text)
                ), 
                React.createElement("div", {className: "human"}, 
                  React.createElement("div", {className: "name"}, human), 
                  React.createElement("div", {className: "form"}, 
                    React.createElement("form", {id: "inputForm", onSubmit: this.handleSubmit, onChange: this.handleChange, 
                    className: "request"}, 
                      options, 
                      React.createElement("br", null), 
                      React.createElement("input", {type: "hidden", id: "workflowStateID", value: workflowStateID}), 
                      React.createElement("button", {type: "submit", disabled: !this.isEnabled()}, "Enter")
                    )
                  )
                )
              );
    }
    if (type == PROMPT_NO_RESPONSE) {
      return  React.createElement("div", null, 
                React.createElement("div", {className: "researcher"}, 
                  React.createElement("div", {className: "name"}, DisplayText[MSG_ROBOT]), 
                  React.createElement("div", {className: "message"}, prompt.Text)
                )
              );
    }
    if (type == PROMPT_END) {
      return  React.createElement("div", null);
    }
    console.error("Error: Unknown prompt type!");  
    return React.createElement("div", null);  
  },
});

var PromptOption = React.createClass({displayName: "PromptOption",

  render: function() {
    var option = this.props.option;
      return  React.createElement("label", null, 
                React.createElement("input", {type: "radio", name: "input", value: option.Value}), 
                option.Label
              )
  },
});







