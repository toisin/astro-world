/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

function CurrentChallenge(name) {
  this.username = username;
  this.oldCart = null;
  this.newCart = null;
  this.results = null;
  this.currentChallenge = null;
}

User.prototype = {

  loadAllUserData: function(renderCallback) {
    var self = this;
    var cartPromise = self.loadUserResultData(self.username);

    var challengePromise = cartPromise.then(function(username) {
                                              return self.loadUserChallengeData(username);
                                            });
    challengePromise.then(renderCallback, function(error) {
                                            console.error("Failed!", error);
                                          });
  },

};














var Challenge = React.createClass({displayName: 'Challenge',
  getInitialState: function() {
    this.setState({enabled:false});
    debugger;
    return {};
  },

  handleChange: function(event) {
    // var state = {};
    // state[event.target.id] = event.target.value;
    this.setState({enabled:true});
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
    return this.state.enabled;
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
