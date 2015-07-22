/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js


var PROMPT_TEXT = "TEXT";
var PROMPT_YES_NO = "YES_NO";
var PROMPT_NO_RESPONSE = "NO_RESPONSE";
var PROMPT_MC = "MC";


function User(name) {
  this.Username = name;
  this.Screenname = "";
  this.History = [];
  this.CurrentPrompt = {};
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

  getPrompt: function() {
    return this.CurrentPrompt;
  },

  getScreenname: function() {
    return this.Screenname;
  },

  updateUser: function(j) {
    var self = this;
    self.Screenname = j["Screenname"];
    self.History = j["History"];
    self.CurrentPrompt = j["CurrentPrompt"];
    //self.CurrentPrompt = {"prompt": {"text": "First Question", "workflowStateID": "2"};    
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
  submitResponse: function(workflowStateID, text, renderCallback) {
    var self = this;

    var formData = new FormData();

    formData.append("user", self.Username);
    formData.append("response", text);
    formData.append("workflowStateID", workflowStateID);

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



