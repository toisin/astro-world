/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var Challenge = React.createClass({displayName: 'Challenge',
  getInitialState: function() {
    return {};
  },

  handleChange: function(event) {
    var state = {};
    state[event.target.id] = event.target.value;
    this.setState(state);
  },

  handleSubmit: function(event) {
    event.preventDefault();

    this.post(this.state);
  },

  post: function(data) {
/*    if (!this.isEnabled())
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
*/
  },

  isEnabled: function() {
    if (!('findout' in this.state)) {
      return false;
    }
    return true;
  },

  render: function() {
    var user = this.props.user;
    var variableModels = this.props.variableModels;
    var ivnames = variableModels.iVariables.map(function(iv) {
      return iv.name;
    });

    var variables = this.props.variableModels.iVariables.map(function(variable) {
      return IndependentVariable( {iv:variable});
    });

    switch (this.state.mode) {
      default:
        return React.DOM.form( {onSubmit:this.handleSubmit, onChange:this.handleChange,
                className:"request"}, 
          React.DOM.table(null, React.DOM.tbody(null, 
            React.DOM.tr(null, 
              React.DOM.td(null, 
                " What did you find out about whether the Handle Length makes a difference? "
              ),
              React.DOM.td(null, 
                React.DOM.textarea( {id:"handlelength"})
              )
            ),
            React.DOM.tr(null, 
              React.DOM.td(null, 
                " What results show you are right? "
              ),
              React.DOM.td(null, 
                React.DOM.textarea( {id:"results"})
              )
            )
          )),
          React.DOM.button( {type:"submit", disabled:!this.isEnabled()}, "Enter")
        );
    }      
  
  }
});

function User(name) {
  this.username = username;
  this.oldCart = null;
  this.newCart = null;
}

User.prototype = {
  getUserData: function(username, callback) {
    var xhr = new XMLHttpRequest();
    var self = this;
    xhr.onload = function() {
      self.results = JSON.parse(xhr.responseText);
      callback();
    };
    xhr.open('GET', '/usercart/' + this.username + '/findallcarts');
    xhr.send(null);
  },

  updateCart: function(result) {
    if (this.oldCart == null) {
      this.oldCart = result;
      return;
    }
    var latestCart = this.oldCart;
    if (this.newCart != null) {
      latestCart = this.newCart;
    }
    var ivnames = variableModels.iVariables.map(function(iv) {
      return iv.name;
    });
    for (var i = 0; i < ivnames.length; i++) {
      if (result[ivnames[i]] != latestCart[ivnames[i]]) {
        this.oldCart = latestCart;
        this.newCart = result;
        return;
      }
    }
  },

  addResult: function(result, callback) {
    this.updateCart(result);
    var xhr = new XMLHttpRequest();
    var self = this;
    xhr.onload = function() {
      self.getUserData(username, callback);
    };
    xhr.open('POST', '/usercart/' + this.username + '/addcartdata');
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(result));
  },

  updateResult: function(callback) {
    if (!this.results) {
      this.getUserData(this.username, callback);
    } else {
      callback();
    }
  }
};

var UserResultData = React.createClass({displayName: 'UserResultData',
  render: function() {
    var user = this.props.user;
    var variableModels = this.props.variableModels;
    var ivnames = variableModels.iVariables.map(function(iv) {
      return iv.name;
    });
    var resultsDisplay = [];
    var allDisplay = [user.results.length];

    switch (this.props.mode) {
/*      case 'all':
        var results = user.results;
        resultsDisplay = results.map(function(result) {
          var index = result[variableModels.dvResultCount];
          return <UserResult variableModels={variableModels} data={result} index={'#' + index}/>;
        });
        break;*/
      case 'notebook':
        var newDisplay = null;
        var oldDisplay = null;
        for (var j=0; j < user.results.length; j++) {
          var result = user.results[j];
          var isNew = true;
          var isOld = true;
          for (var i = 0; i < ivnames.length; i++) {
            if ((!user.newCart) || (result[ivnames[i]] != user.newCart[ivnames[i]])) {
              isNew = false;
            }
            if ((!user.oldCart) || (result[ivnames[i]] != user.oldCart[ivnames[i]])) {
              isOld = false;
            }
          }
          var index = result[variableModels.dvResultCount];
          allDisplay[index-1] = UserResult( {variableModels:variableModels, data:result, index:'#' + index});

          if (isNew) {
            newDisplay = UserResult( {variableModels:variableModels, data:result, index:'#' + index + ' (Newly Saved)'});
          } else if (isOld) {
            oldDisplay = UserResult( {variableModels:variableModels, data:result, index:'#' + index + ' (Last Saved)'});
          }
        }
        if (newDisplay) {
          resultsDisplay.push(newDisplay);
        }
        if (oldDisplay) {
          resultsDisplay.push(oldDisplay);
        }
        break;
    }

    var headers = variableModels.iVariables.map(function(iv) {
      //return <th><VariableImage name={iv.name}/>{iv.label}</th>;
      return React.DOM.th(null, iv.label);
    });

    return React.DOM.table( {className:"result"}, 
      React.DOM.thead(null, 
        React.DOM.tr(null, 
          React.DOM.th(null
          ),
          headers,
          React.DOM.th(null, 
            variableModels.dvLabel
          )
        )
      ),
      React.DOM.tbody(null, 
        resultsDisplay,
        React.DOM.tr(null, React.DOM.td( {co:headers.length}, "All Results")),
        allDisplay
      )
    );
  },
});

var UserResult = React.createClass({displayName: 'UserResult',
  render: function() {
    var variableModels = this.props.variableModels;
    var data = this.props.data;
    var dvValues = data[variableModels.dvName].join(', ');
    var index = this.props.index;

    var variables = variableModels.iVariables.map(function(variable) {
      return UserResultSelection( {iv:variable, value:data[variable.name]});
    });

    return React.DOM.tr(null, 
      React.DOM.td(null, 
        " Cart ", index, " : "
      ),
      variables,
      React.DOM.td(null, 
        dvValues
      )
      );
  }
});

var UserResultSelection = React.createClass({displayName: 'UserResultSelection',
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
    return React.DOM.td(null, ivValue);
  }
});


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
      React.DOM.input( {type:"radio", name:this.props.name, value:ivOption.value}),
      ivOption.label
    );
  }
});

var Request = React.createClass({displayName: 'Request',
  getInitialState: function() {
    return {};
  },

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
        case 3:
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
          React.DOM.button( {onClick:this.saveResult}, "Save Result to Notebook")
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
          React.DOM.button( {onClick:this.continueFrom}, "Go to Next Case")
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

var variableModels = {

  dvLabel: 'Number of trips per hour',
  dvName: 'trips',
  dvResultCount: 'cartNumber',


  iVariables: [
    {
      name: 'handleLength',
      label: 'Handle length',
      options: [
        {value: 'Long', label: 'Long'},
        {value: 'Short', label: 'Short'},
      ]
    },
    {
      name: 'wheelSize',
      label: 'Wheel Size',
      options: [
        {value: 'Large(4)', label: 'Large(4)'},
        {value: 'Small(3)', label: 'Small(3)'}
      ],
    },
    {
      name: 'bucketSize',
      label: 'Bucket Size',
      options: [
        {value: 'Big(13)', label: 'Big(13)'},
        {value: 'Small(10)', label: 'Small(10)'},
      ]
    },
    {
      name: 'bucketPlacement',
      label: 'Bucket Placement',
      options: [
        {value: 'Far', label: 'Far'},
        {value: 'Near', label: 'Near'},
      ]
    }
  ]
};

// Begin side-effect
var username = window.location.hash.substring(1);

if (!username) {
  window.location = "index.html";
} else {

  var user = new User(username);

  React.renderComponent(
    App( {variableModels:variableModels, user:user}),
    document.body);

}

// TODO
window.onbeforeunload = function() {
  return "";
};


