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
    var prompt = user.getPrompt();
    var action = user.getAction();
    // var key = prompt.PromptId + action.UIActionModeId;
    var onComplete = this.props.onComplete;

    if (action) {
      switch (action.UIActionModeId) {
        case "NEW_TARGET_FACTOR":
          return  <SelectTargetFactor user={user} onComplete={onComplete} app={app}/>;// key={key}/>;
        case "PRIOR_BELIEF_FACTORS":
          return  <PriorBeliefFactors user={user} onComplete={onComplete} app={app}/>;// key={key}/>;
        case "PRIOR_BELIEF_LEVELS":
          return  <PriorBeliefLevels user={user} onComplete={onComplete} app={app}/>;// key={key}/>;
        case "RECORD_SELECT_ONE":
          return <RecordSelection user={user} onComplete={onComplete} app={app} singleRecord={true}/>;// key={key}/>;
        case "RECORD_SELECT_TWO":
          return <RecordSelection user={user} onComplete={onComplete} app={app} singleRecord={false}/>;// key={key}/>;
        case "RECORD_NO_PERFORMANCE":
          return <RecordPerformance user={user} onComplete={onComplete} app={app} showPerformance={false}/>;// key={key}/>;
        case "RECORD_PERFORMANCE":
          return <RecordPerformance user={user} onComplete={onComplete} app={app} showPerformance/>;// key={key}/>;
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





