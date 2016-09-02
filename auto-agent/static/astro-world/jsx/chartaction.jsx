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

    var recordsToShow = [{grade:1, filter:"fitness:average", no:64},{grade:4, filter:"fitness:average", no:47}]

    if (action) {
      switch (action.UIActionModeId) {
        case "NEW_TARGET_FACTOR":
          return  <ChartSelectTargetFactor user={user} onComplete={onComplete} app={app} key={"NEW_TARGET_FACTOR"}/>;
        case "ALL_RECORDS_ALLOW_TOOLBOX":
          return  <Chart user={user} singleColumn allowToolbox onComplete={onComplete} app={app} key={"ALL_RECORDS_ALLOW_TOOLBOX"}/>;
        case "FITNESS_AVERAGE_RECORDS":
          return  <Chart user={user} filterFactorName={"Fitness"} filterLevels={["Average"]} filterRecords={["fitness:average"]} onComplete={onComplete} app={app} key={"FITNESS_AVERAGE_RECORDS"}/>;
        case "FITNESS_AVERAGE_RECORDS_SHOW_TWO_RECORDS":
          return  <Chart user={user} recordsToShow={recordsToShow} filterFactorName={"Fitness"} filterLevels={["Average"]} filterRecords={["fitness:average"]} onComplete={onComplete} app={app} key={"FITNESS_AVERAGE_RECORDS"}/>;
        case "ALL_RECORDS":
          return  <Chart user={user} singleColumn onComplete={onComplete} app={app} key={"ALL_RECORDS"}/>;
        case "TARGET_FACTOR_RECORDS":
          return  <Chart user={user} showTargetFactorRecords onComplete={onComplete} app={app} key={"TARGET_FACTOR_RECORDS"}/>;
        default:
          return <div></div>;
      }
    }
    return <div></div>;
  }
});





