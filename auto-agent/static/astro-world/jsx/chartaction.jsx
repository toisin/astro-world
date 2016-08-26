/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var ChartAction = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();
    var action = user.getAction();
    var onComplete = this.props.onComplete;

    if (action) {
      switch (action.UIActionModeId) {
        case "ALL_RECORDS":
          return  <Chart user={user} onComplete={onComplete} app={app}/>;
        default:
          return <div></div>;
      }
    }
    return <div></div>;
  }
});





