/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var DELAY_PROMPT_TIME_SHORT = 1000;
var DELAY_PROMPT_TIME_LONG = 3000;
var LONG_PROMPT_SIZE = 150;

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
  this.ContentFactors = [];
  this.State = {};
  this.ArchiveHistoryLength = 0;
}

User.prototype = {

  loadAllUserData: function(renderCallback) {
    var self = this;
    var historyPromise = self.loadHistory();

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

  updateUser: function(j) {
    var self = this;
    self.Screenname = j.Screenname;
    self.History = j.History;
    self.CurrentUIPrompt = j.CurrentUIPrompt;
    self.CurrentUIAction = j.CurrentUIAction;
    self.CurrentPhaseId = j.CurrentPhaseId;
    self.ContentFactors = j.ContentFactors;
    self.State = j.State;
    self.ArchiveHistoryLength = j.ArchiveHistoryLength;
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


  //After submitting the response
  //Update user with new history etc.
  submitResponse: function(promptId, phaseId, jsonResponse, renderCallback) {
    var self = this;
    // var text = value;
    var question = self.CurrentUIPrompt.Texts;
    var jsonQuestion = JSON.stringify(question); // Turns the texts array into json
    var phaseId = self.CurrentPhaseId;
    var promptId = self.CurrentUIPrompt.PromptId;

    var formData = new FormData();

    formData.append("user", self.Username);
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

  // passing in self because otherwise, the scope can be screwed up if
  //     this is called from Promise
  // loadUserChallengeData: function() {
  //   var self = this;
  //   var promise = new Promise(function(resolve, reject) {
  //     var challengeReq = new XMLHttpRequest();
  //     challengeReq.onload = function() {
  //       //self.results = JSON.parse(challengeReq.responseText);
  //       resolve(self);
  //     };
  //     challengeReq.onerror = function() {
  //       reject(Error("It broke"));
  //     };
  //     challengeReq.open('GET', '/userchallenge/' + self.username + '/findallchallenges');
  //     challengeReq.send(null);

  //   });

  //   return promise;
  // },

  // passing in self because otherwise, the scope can be screwed up if
  //     this is called from Promise
  // loadUserResultData: function() {
  //   var self = this;
  //   var promise = new Promise(function(resolve, reject) {
  //     var resultsReq = new XMLHttpRequest();
  //     resultsReq.onload = function() {
  //       debugger;
  //       self.results = JSON.parse(resultsReq.responseText);
  //       resolve();
  //     };
  //     resultsReq.onerror = function() {
  //       reject(Error("It broke"));
  //     };
  //     resultsReq.open('GET', '/usercart/' + self.username + '/findallcarts');
  //     resultsReq.send(null);

  //   });

  //   return promise;
  // },

  // DELETE:Replaced by loadUserResultData using Promise
  // getUserData: function(username, callback) {
  //   var self = this;
  //   var xhr = new XMLHttpRequest();
  //   xhr.onload = function() {
  //     self.results = JSON.parse(xhr.responseText);
  //     callback();
  //   };
  //   xhr.open('GET', '/usercart/' + this.username + '/findallcarts');
  //   xhr.send(null);
  // },

  // updateCart: function(result) {
  //   if (this.oldCart == null) {
  //     this.oldCart = result;
  //     return;
  //   }
  //   var latestCart = this.oldCart;
  //   if (this.newCart != null) {
  //     latestCart = this.newCart;
  //   }
  //   var ivnames = variableModels.iVariables.map(function(iv) {
  //     return iv.name;
  //   });
  //   for (var i = 0; i < ivnames.length; i++) {
  //     if (result[ivnames[i]] != latestCart[ivnames[i]]) {
  //       this.oldCart = latestCart;
  //       this.newCart = result;
  //       return;
  //     }
  //   }
  // },

  // addResult: function(result, renderCallback) {
  //   var self = this;
    
  //   self.updateCart(result);

  //   var addCartPromise = new Promise(function(resolve, reject) {
  //     var xhr = new XMLHttpRequest();
  //     xhr.onload = function() {
  //       resolve(self);
  //     };
  //     xhr.error = function() {
  //       reject();
  //     }; 
  //     xhr.open('POST', '/usercart/' + self.username + '/addcartdata');
  //     xhr.setRequestHeader('Content-Type', 'application/json');
  //     xhr.send(JSON.stringify(result));
  //   });

  //   var loadUserCartPromise = addCartPromise.then(function() {
  //     self.loadUserResultData();
  //   });

  //   loadUserCartPromise.then(renderCallback);
  // },

  // enterChallenge: function(renderCallback) {
  //   // if (!this.currentChallenge) {
  //   //   // if the user data are empty, receive it
  //   //   this.loadAllUserData(renderCallback);
  //   // } else {
  //     renderCallback();
  //   // }
  // }




};



