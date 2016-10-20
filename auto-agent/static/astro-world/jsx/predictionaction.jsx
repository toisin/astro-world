/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


function PredictionAction(props) {
  var user = props.user;
  var app = props.app;
  var prompt = user.getPrompt();
  var action = user.getAction();
  var onComplete = props.onComplete;

  var recordsToShow = [{grade:1, filter:"fitness:average", no:64},{grade:4, filter:"fitness:average", no:47}]

  if (action) {
    switch (action.UIActionModeId) {
      case "TARGET_FACTOR_RECORDS":
        return  <Chart user={user} showTargetFactorRecords app={app} key={"TARGET_FACTOR_RECORDS"}/>;
      case "FACTORS_REQUEST_FORM":
        return <FactorsRequestForm user={user} onComplete={onComplete} app={app}/>;
      case "PREDICTION_RECORD":
        return <div>
                <PredictionRecord user={user} onComplete={onComplete} app={app}/>
                <ChartButtons user={user} app={app}/>
                <PredictionRecord user={user} onComplete={onComplete} app={app} predictionHistory showPerformancePrediction/>
              </div>;
      case "PREDICTION_RECORD_SHOW_PREDICTION":
        return <PredictionRecord user={user} onComplete={onComplete} app={app} showPerformancePrediction/>;
      case "CONTRIBUTING_FACTORS_FORM":
        return <ContributingFactorsForm user={user} onComplete={onComplete} app={app}/>; 
      case "SELECT_TEAM_SUMMARY":
        return <SelectTeam user={user} onComplete={onComplete} app={app} isSummary/>; 
      case "SELECT_TEAM":
        return <SelectTeam user={user} onComplete={onComplete} app={app}/>; 
    default:
        return <div></div>;
    }
  }
  return <div></div>;
}





