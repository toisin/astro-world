/** @jsx React.DOM */

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
    var action = user.getAction();

    if (action) {
      switch (action.UIActionModeId) {
        case "NEW_TARGET_FACTOR":
          return  <SelectTargetFactor user={user} onComplete={app.changeState} app={app}/>;
        default:
          return <div></div>;
      }
    }
    return <div></div>;
  }
});





