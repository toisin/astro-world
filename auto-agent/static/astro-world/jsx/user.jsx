/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var MESSAGE_COUNT_LIMIT = 100;

var DELAY_PROMPT_TIME_REALLY_SHORT = 500;
var DELAY_PROMPT_TIME_SHORT = 1000;
var DELAY_PROMPT_TIME_LONG = 5000;
// var DELAY_PROMPT_TIME_REALLY_SHORT = 50;
// var DELAY_PROMPT_TIME_SHORT = 100;
// var DELAY_PROMPT_TIME_LONG = 500;
var LONG_PROMPT_SIZE = 150;
var REALLYSHORT_PROMPT_SIZE = 50;

var UI_PROMPT_ENTER_TO_CONTINUE = "ENTER_TO_CONTINUE";
var UI_PROMPT_TEXT = "Text";
var UI_PROMPT_MC = "MC";
var UI_PROMPT_NO_INPUT = "NO_INPUT";
var UI_PROMPT_STRAIGHT_THROUGH = "STRAIGHT_THROUGH";

var RESPONSE_SYSTEM_GENERATED = "SYSTEM_GENERATED";

var PHASE_COV = "Cov";
var PHASE_CHART = "Chart";
var PHASE_PREDICTION = "Prediction";
var FIRST_PHASE = "START";
var LAST_PHASE = "END";

var UIACTION_INACTIVE = "NO_UIACTION";


function User(name) {
  this.Username = name;
  this.Screenname = "";
  this.History = [];
  this.CurrentPhaseId = "";
  this.CurrentUIPrompt = {};
  this.CurrentUIAction = {};
  this.ContentFactors = {};
  this.State = {};
  this.ArchiveHistoryLength = 0;
  // this.AllPerformanceRecords = {};
}

User.prototype = {

  loadAllUserData: function(renderCallback) {
    var historyPromise = this.loadHistory();

    historyPromise.then(renderCallback, function(error) {
                                               console.error("Failed to load history!", error);
                                           });
  },

  getUsername: function() {
    return this.Username;
  },

  getHistory: function() {
    return this.History;
  },

  getCurrentPhaseId: function() {
    return this.CurrentPhaseId;
  },

  getPrompt: function() {
    return this.CurrentUIPrompt;
  },

  getAction: function() {
    return this.CurrentUIAction;
  },

  getContentFactors: function() {
    return this.ContentFactors;
  },

  getScreenname: function() {
    return this.Screenname;
  },

  getState: function() {
    return this.State;
  },

  getArchiveHistoryLength: function() {
    return this.ArchiveHistoryLength;
  },

  getAllPerformanceRecords: function() {
    return this.AllPerformanceRecords;
  },

  updateUser: function(j) {
    this.Screenname = j.Screenname;
    this.History = j.History;
    this.CurrentUIPrompt = j.CurrentUIPrompt;
    this.CurrentUIAction = j.CurrentUIAction;
    this.CurrentPhaseId = j.CurrentPhaseId;
    this.ContentFactors = j.ContentFactors;
    this.State = j.State;
    this.ArchiveHistoryLength = j.ArchiveHistoryLength;
  },

  updateAllPerformanceRecords: function(j) {
    this.AllPerformanceRecords = j;
  },

  loadHistory: function() {
    var self = this;
    var promise = new Promise(function(resolve, reject) {
      var historyReq = new XMLHttpRequest();
      historyReq.onload = function() {
        self.updateUser(JSON.parse(historyReq.responseText));
        resolve();
      };
      historyReq.onerror = function() {
        reject(Error("It broke"));
      };
      historyReq.open('GET', 'history?user='+self.Username);
      historyReq.send(null);

     });

     return promise;
  },

  loadAllPerformanceRecords: function() {
    var self = this;
    var promise = new Promise(function(resolve, reject) {
      var recordsReq = new XMLHttpRequest();
      recordsReq.onload = function() {
        self.updateAllPerformanceRecords(JSON.parse(recordsReq.responseText));
        resolve();
      };
      recordsReq.onerror = function() {
        reject(Error("It broke"));
      };
      recordsReq.open('GET', 'records?user='+self.Username);
      recordsReq.send(null);

     });

     return promise;
  },


  //After submitting the response
  //Update user with new history etc.
  submitResponse: function(promptId, phaseId, jsonResponse, renderCallback) {
    var self = this;
    // var text = value;
    var question = this.CurrentUIPrompt.Texts;
    var jsonQuestion = JSON.stringify(question); // Turns the texts array into json
    var phaseId = this.CurrentPhaseId;
    var promptId = this.CurrentUIPrompt.PromptId;

    var formData = new FormData();

    formData.append("user", this.Username);
    formData.append("questionText", jsonQuestion);
    formData.append("promptId", promptId);
    formData.append("phaseId", phaseId);
    formData.append("jsonResponse", jsonResponse);

    var responsePromise = new Promise(function(resolve, reject) {
      var xhr = new XMLHttpRequest();
      xhr.onload = function() {
        self.updateUser(JSON.parse(xhr.responseText));
        resolve();
      };
      xhr.error = function() {
        reject();
      }; 
      xhr.open('POST', 'sendresponse');
      xhr.send(formData);
    });

    responsePromise.then(renderCallback, function(error) {
                                               console.error("Failed to submit a response!", error);
                                           });
  },
};



