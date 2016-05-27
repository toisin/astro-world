/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var Action = React.createClass({displayName: "Action",

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();

    return  React.createElement("div", {className: "action"}, 
			    React.createElement("h3", null, "Investigating Factor: ", React.createElement("b", null, "Fitness")), 

			    React.createElement("table", null, 
			      React.createElement("tbody", null, 
			        React.createElement("tr", null, 
			          React.createElement("td", null, " "), 
			          React.createElement("td", {colSpan: "3", className: "question"}, "Which record would you like to see?")
			        ), 
			        React.createElement("tr", null, 
                      React.createElement("td", {colSpan: "4"}, 
                      React.createElement(ActionInput, {user: user, prompt: prompt, onComplete: app.changeState}))
			        )
			      )
			    )
            );
  }
});


var ActionInput = React.createClass({displayName: "ActionInput",

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
                    React.createElement("form", {id: "inputForm", onSubmit: this.handleSubmit, onChange: this.handleChange, className: "request"}, 
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
          return React.createElement(ActionPromptOption, {option: option, key: i});
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
    // console.error("Error: Unknown prompt type!");  
    return React.createElement("div", null);  
  },
});

var ActionPromptOption = React.createClass({displayName: "ActionPromptOption",

  render: function() {
    var option = this.props.option;
      return  React.createElement("label", null, 
                React.createElement("input", {type: "radio", name: "input", value: option.Value}), 
                option.Label
              )
  },
});


