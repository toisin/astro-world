/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var VariableImage = React.createClass({displayName: 'VariableImage',
  render: function() {
    return React.DOM.img( {src:'images/' + this.props.name + '.png', height:"100"});
  }
});

var IndependentVariable = React.createClass({displayName: 'IndependentVariable',
  render: function() {
    var iv = this.props.iv;
    var name = iv.name;
    var handleChange = this.handleChange;
    var options = iv.options.map(function(option) {
      return IndependentVariableOption( {name:name, ivOption:option});
    });

    return React.DOM.tr( {className:"iv"}, 
      React.DOM.td(null, VariableImage( {name:iv.name})),
      React.DOM.td(null, iv.label),
      React.DOM.td(null, options)
    );
  }
});

var IndependentVariableOption = React.createClass({displayName: 'IndependentVariableOption',
  render: function() {
    var ivOption = this.props.ivOption;
    return React.DOM.label(null, 
      " // this.props.name is the name of the IndependentVariable this option "+
      "//     is associated with. "+
      "// ivOption.value is the value that gets saved when the option is selected. ",
      React.DOM.input( {type:"radio", name:this.props.name, value:ivOption.value}),
      ivOption.label
    );
  }
});

var Request = React.createClass({displayName: 'Request',
  getInitialState: function() {
    return {};
  },

  // @param {Event} e The event within the variable-level-selection form,
  //     for now, they are only events from an IndependentVariableOption
  // The hashtable is to keep track of the IndependentVariable that has its option selected:
  //     e.target.name is the name of the IndependentVariable
  //     e.target.value is the value of the IndependentVariableOption
  handleChange: function(e) {
    var state = {};
    state[e.target.name] = e.target.value;
    this.setState(state);
  },

  handleSubmit: function(e) {
    e.preventDefault();
    this.post(this.state);
  },

  post: function(data) {
    if (!this.isEnabled())
      return;

    var xhr = new XMLHttpRequest();
    var self = this;
    xhr.onload = function() {
      if (self.props.onComplete) {
        self.props.onComplete(JSON.parse(xhr.responseText));
      }
    };
    xhr.open('POST', '/carts/gettrips');
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(data));
  },

  // Check all IndependentVariable, only allow form submit if all variables
  //     have been selected.
  // (See handleChange for keeping track of the variables with options selected)
  isEnabled: function() {
    var variables = this.props.variableModels.iVariables;
    for (var i = 0; i < variables.length; i++) {
      if (!(variables[i].name in this.state)) {
        return false;
      }
    }
    return true;
  },

  render: function() {
    var variables = this.props.variableModels.iVariables.map(function(variable) {
      return IndependentVariable( {iv:variable});
    });

    return React.DOM.form( {onSubmit:this.handleSubmit, onChange:this.handleChange,
            className:"request"}, 
      React.DOM.table(null, React.DOM.tbody(null, variables)),
      React.DOM.button( {type:"submit", disabled:!this.isEnabled()}, "See Results")
    );
  }
});

var Result = React.createClass({displayName: 'Result',
  render: function() {
    var variableModels = this.props.variableModels;
    var data = this.props.data;
    var dvValues = data[variableModels.dvName].join(', ');

    var variables = variableModels.iVariables.map(function(variable) {
      return ResultSelection( {iv:variable, value:data[variable.name]});
    });

    return React.DOM.table( {className:"result"}, React.DOM.tbody(null, 
      React.DOM.tr(null, 
        React.DOM.td(null),
        React.DOM.td(null, variableModels.dvLabel,":"),
        React.DOM.td(null, dvValues)
      ),
      variables
    ));
  }
});

var ResultSelection = React.createClass({displayName: 'ResultSelection',
  getDisplayValue: function(value) {
    var options = this.props.iv.options;
    for (var i = 0; i < options.length; i++) {
      if (options[i].value == value) {
        return options[i].label;
      }
    }
    return null;
  },

  render: function() {
    var iv = this.props.iv;
    var ivValue = this.getDisplayValue(this.props.value);
    return React.DOM.tr(null, 
      React.DOM.td(null, VariableImage( {name:iv.name})),
      React.DOM.td(null, iv.label,":"),
      React.DOM.td(null, ivValue)
    );
  }
});

var App = React.createClass({displayName: 'App',
  getInitialState: function() {
    return {mode: 0};
  },

  continueFrom: function(e) {
    switch (this.state.mode) {
      case 1:
        this.setState({mode: 2});
        break;
      case 3:
        this.setState({mode: 2});
        break;
      case 4:
        this.setState({mode: 2});
        break;
      case 5:
        this.setState({mode: 2});
        break;
    }
  },

  showAllResultsForChallenge: function(e) {
    switch (this.state.mode) {
      case 5:
        this.setState({mode: 6});
        break;
    }
  },

  hideAllResultsForChallenge: function(e) {
    switch (this.state.mode) {
      case 6:
        this.setState({mode: 5});
        break;
    }
  },

  showChallenge: function(e) {
    var self = this;
    this.props.user.updateResult(function() {
      switch (self.state.mode) {
        case 1:
        case 3:
        case 4:
          self.setState({mode: 5});
          break;
      }
    });
  },

  saveResult: function(e) {
    var self = this;
    this.props.user.addResult(this.state.newResult, function() {
      switch (self.state.mode) {
        case 1:
          self.setState({mode: 4});
          break;
        case 3:
          self.setState({mode: 4});
          break;
      }
    });
  },

  handleComplete: function(data) {
    var state = this.state;
    switch (state.mode) {
      case 0:
        this.setState({mode: 1, newResult: data});
        break;
      case 2:
        this.setState({mode: 3, oldResult: state.newResult, newResult: data});
        break;
    }
  },

  render: function() {
    var variableModels = this.props.variableModels;
    var state = this.state;

    switch (this.state.mode) {
      case 0: //Show Initial Request Form
        return React.DOM.div( {className:"app single-column"}, 
          Request( {variableModels:variableModels, onComplete:this.handleComplete,
              style:{width: '100%'}})
        );
    
      case 1: //Show One Case Results
        return React.DOM.div( {className:"app single-column"}, 
          Result( {variableModels:variableModels, data:state.newResult}),
          React.DOM.button( {onClick:this.continueFrom}, "Go to Next Case"),
          React.DOM.button( {onClick:this.saveResult}, "Save Result to Notebook"),
          React.DOM.button( {onClick:this.showChallenge}, "Show Challenge")
        );
    
      case 2: //Show Request Form With Last Result
        return React.DOM.table( {className:"app"}, React.DOM.tbody(null, 
          React.DOM.tr(null, 
            React.DOM.td(null, "New Case:"),
            React.DOM.td(null, "Last Case:")
          ),
          React.DOM.tr(null, 
            React.DOM.td(null, Request( {variableModels:variableModels, onComplete:this.handleComplete})),
            React.DOM.td(null, Result( {variableModels:variableModels, data:state.newResult}))
          )
        ));
    
      case 3: //Show Two Cases Results
        return React.DOM.table( {className:"app"}, React.DOM.tbody(null, 
          React.DOM.tr(null, 
            React.DOM.td(null, "New Case:"),
            React.DOM.td(null, "Last Case:")
          ),
          React.DOM.tr(null, 
            React.DOM.td(null, Result( {variableModels:variableModels, data:state.newResult})),
            React.DOM.td(null, Result( {variableModels:variableModels, data:state.oldResult}))
          ),
          React.DOM.tr(null, 
            React.DOM.td( {colSpan:"2", style:{textAlign: 'center'}}, 
              React.DOM.button( {onClick:this.continueFrom}, "Go to Next Case"),
              React.DOM.button( {onClick:this.saveResult}, "Save Result to Notebook"),
              React.DOM.button( {onClick:this.showChallenge}, "Show Challenge")
            )
          )
        ));

      case 4: //Show Notebook
        return React.DOM.div( {className:"app single-column"}, 
          UserResultData( {variableModels:this.props.variableModels, user:this.props.user, mode:'notebook'}),
          React.DOM.button( {onClick:this.continueFrom}, "Go to Next Case"),
          React.DOM.button( {onClick:this.showChallenge}, "Show Challenge")
        );

      case 5: //Show Challenge
        return React.DOM.div( {className:"app single-column"}, 
          React.DOM.table( {className:"app"}, React.DOM.tbody(null, 
            React.DOM.tr(null, 
              React.DOM.td(null, React.DOM.div( {className:"app single-column"}, 
                Challenge( {variableModels:this.props.variableModels, user:this.props.user}),
                React.DOM.button( {onClick:this.showAllResultsForChallenge}, "Show Notebook")
              ))
            )
        ))
        );

      case 6: //Show Challenge with Notebook
        return React.DOM.div( {className:"app single-column"}, 
          React.DOM.table( {className:"app"}, React.DOM.tbody(null, 
            React.DOM.tr(null, 
              React.DOM.td(null, React.DOM.div( {className:"app single-column"}, 
                Challenge( {variableModels:this.props.variableModels, user:this.props.user}),
                React.DOM.button( {onClick:this.hideAllResultsForChallenge}, "Hide Notebook")
              )),
              React.DOM.td(null, React.DOM.div( {className:"app single-column"}, 
                UserResultData( {variableModels:this.props.variableModels, user:this.props.user, mode:'notebook'})
              ))
            )
        ))
        );

    }

    throw new Error('Unexpected mode');
  }
});


