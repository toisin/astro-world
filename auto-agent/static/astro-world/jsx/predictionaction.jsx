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
      // case "NEW_TARGET_FACTOR":
      //   return  <ChartSelectTargetFactor user={user} onComplete={onComplete} app={app} key={"NEW_TARGET_FACTOR"}/>;
      // case "ALL_RECORDS_ALLOW_TOOLBOX":
      //   return  <Chart user={user} singleColumn allowToolbox app={app} key={"ALL_RECORDS_ALLOW_TOOLBOX"}/>;
      // case "FITNESS_AVERAGE_RECORDS":
      //   return  <Chart user={user} filterFactorName={"Fitness"} filterLevelsLabels={["Average"]} filterRecords={["fitness:average"]} app={app} key={"FITNESS_AVERAGE_RECORDS"}/>;
      // case "FITNESS_AVERAGE_RECORDS_SHOW_TWO_RECORDS":
      //   return  <Chart user={user} recordsToShow={recordsToShow} filterFactorName={"Fitness"} filterLevelsLabels={["Average"]} filterRecords={["fitness:average"]} app={app} key={"FITNESS_AVERAGE_RECORDS"}/>;
      // case "ALL_RECORDS":
      //   return  <Chart user={user} singleColumn app={app} key={"ALL_RECORDS"}/>;
      // case "TARGET_FACTOR_RECORDS":
      //   return  <Chart user={user} showTargetFactorRecords app={app} key={"TARGET_FACTOR_RECORDS"}/>;
      // case "MEMO_FORM":
      //   return <ChartMemoForm user={user} onComplete={onComplete} app={app}/>;
      // case "MEMO":
      //   return <Memo user={user} app={app}/>;
      case "FACTORS_REQUEST_FORM":
        return <FactorsRequestForm user={user} onComplete={onComplete} app={app}/>;
      // case "FACTORS_LEVELS_SUMMARY_FORM":
      //   return <FactorsLevelsSummaryForm user={user} onComplete={onComplete} app={app}/>;
    default:
        return <div></div>;
    }
  }
  return <div></div>;
}





