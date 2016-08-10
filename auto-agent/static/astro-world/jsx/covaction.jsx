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

    var targetFactor
    if (user.getState().TargetFactor) {
     targetFactor = user.getState().TargetFactor.FactorId != "" ? <h3>Investigating Factor: <b>{user.getState().TargetFactor.FactorName}</b></h3> : null;
    }

    if (action) {
      switch (action.UIActionModeId) {
        case "NEW_TARGET_FACTOR":
          return  <SelectTargetFactor user={user} onComplete={onComplete} app={app}/>;// key={key}/>;
        case "PRIOR_BELIEF_FACTORS":
          return  <PriorBeliefFactors user={user} onComplete={onComplete} app={app}/>;// key={key}/>;
        case "PRIOR_BELIEF_LEVELS":
          return  <PriorBeliefLevels user={user} onComplete={onComplete} app={app}/>;// key={key}/>;
        case "RECORD_SELECT_ONE":
          return <div>
                  {targetFactor}
                  <RecordSelection user={user} onComplete={onComplete} app={app} />
                </div>;
        case "RECORD_SELECT_TWO":
          return <div>
                  {targetFactor}
                  <RecordSelection user={user} onComplete={onComplete} app={app} doubleRecord/>
                </div>;
        case "RECORD_NO_PERFORMANCE":
          return <div>
                  {targetFactor}
                  <RecordPerformance user={user} onComplete={onComplete} app={app} hidePerformance/>
                </div>;
        case "RECORD_SELECT_ONE_AND_SHOW_PERFORMANCE":
          return <div>
                  {targetFactor}
                  <RecordSelection user={user} onComplete={onComplete} app={app} doubleRecord comparePrevious/>
                </div>;
        case "RECORD_PERFORMANCE":
          return <div>
                  {targetFactor}
                  <RecordPerformance user={user} app={app}/>
                </div>
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





