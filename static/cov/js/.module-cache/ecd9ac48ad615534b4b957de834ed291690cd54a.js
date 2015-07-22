/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

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
    var renderCallback = function() {
      switch (self.state.mode) {
        case 1:
        case 3:
        case 4:
          self.setState({mode: 5});
          break;
      }
    };
    this.props.user.enterChallenge(renderCallback);
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

