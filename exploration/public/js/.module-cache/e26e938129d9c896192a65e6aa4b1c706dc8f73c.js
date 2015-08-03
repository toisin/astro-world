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

// this.props.name is the name of the IndependentVariable this option
//     is associated with.
// ivOption.value is the value that gets saved when the option is selected.
var IndependentVariableOption = React.createClass({displayName: 'IndependentVariableOption',
  render: function() {
    var ivOption = this.props.ivOption;
    return React.DOM.label(null, 
      React.DOM.input( {type:"radio", name:this.props.name, value:ivOption.value}),
      ivOption.label
    );
  }
});

// Request renders and submit the Variable-level-selection form where user select a level
//     for each IndependentVariable.
// Request has a list of properties defined as attributes in the Request components.
// (For example, see <Result variableModels={variableModels} data={state.newResult}/> in app.jsx)
// They are stored as properties in this.props
// - variableModels
// - data, 
// - onComplete, call back function provided by callers of Request
var Request = React.createClass({displayName: 'Request',
  getInitialState: function() {
    return {};
  },

  // @param {Event} e The event within the Variable-level-selection form,
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
  //     have their options selected.
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
// End --- Request

// Result renders the outcome screen based on the return results from the backend
//     given the levels of the IndependentVariables.
// It displays the outcome and the list of IndependentVariables and their levels
// dvValues is the outcome variable (aka dependent variable)
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

// ResultSelection renders one independent variable with its selected level
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


