/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var CovAction = React.createClass({

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
        case "RECORD_SELECT_ONE":
          return <RecordSelection user={user} onComplete={app.changeState} app={app} singleRecord={true}/>;
        case "RECORD_SELECT_TWO":
          return <RecordSelection user={user} onComplete={app.changeState} app={app} singleRecord={false}/>;
        case "RECORD_PERFORMANCE":
          return <RecordPerformance user={user} onComplete={app.changeState} app={app}/>;
        default:
          return <div></div>;
      }
    }
    return <div></div>;
    //   case "chart":
    //     return <div></div>
    //   case "prediction":
    //     return <div></div>

    //     break;
    // }
  }
});




