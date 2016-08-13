/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var Action = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  changeState: function() {
    // this.setState({mode: 0});
    var app = this.props.app;
    app.changeState();
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();
    // debugger;
    switch (user.getCurrentPhaseId()) {
      case PHASE_COV:
        return  <div><CovAction user={user} onComplete={this.changeState} app={app}/></div>;
      case PHASE_CHART:
        return  <div><ChartAction user={user} onComplete={this.changeState} app={app}/></div>;
      // case "chart":
      //   return <div></div>
      // case "prediction":
      //   return <div></div>

    //     break;
    }
  }
});


