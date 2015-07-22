/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

function User(name) {
  this.username = username;
  this.oldCart = null;
  this.newCart = null;
}

User.prototype = {

  loadAllUserData: function(username, callback) {
    debugger;
    this.loadUserCartData(username).then(callback, function(error) {
                                                console.error("Failed!", error);
                                              });
  },

  loadUserCartData: function(username) {
    var self = this;
    var promise = new Promise(function(resolve, reject) {
      var cartsReq = new XMLHttpRequest();
      cartsReq.onload = function() {
        self.results = JSON.parse(xhr.responseText);
        resolve("Stuff worked!");
      };
      cartsReq.onerror = function() {
        reject(Error("It broke"));
      };
      cartsReq.open('GET', '/usercart/' + this.username + '/findallcarts');
      cartsReq.send(null);

    });

    return promise;
  },

  getUserData: function(username, callback) {
    var self = this;
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
      self.results = JSON.parse(xhr.responseText);
      callback();
    };
    xhr.open('GET', '/usercart/' + this.username + '/findallcarts');
    xhr.send(null);
  },

  updateCart: function(result) {
    if (this.oldCart == null) {
      this.oldCart = result;
      return;
    }
    var latestCart = this.oldCart;
    if (this.newCart != null) {
      latestCart = this.newCart;
    }
    var ivnames = variableModels.iVariables.map(function(iv) {
      return iv.name;
    });
    for (var i = 0; i < ivnames.length; i++) {
      if (result[ivnames[i]] != latestCart[ivnames[i]]) {
        this.oldCart = latestCart;
        this.newCart = result;
        return;
      }
    }
  },

  addResult: function(result, callback) {
    this.updateCart(result);
    var xhr = new XMLHttpRequest();
    var self = this;
    xhr.onload = function() {
      self.loadAllUserData(username, callback);
    };
    xhr.open('POST', '/usercart/' + this.username + '/addcartdata');
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(result));
  },

  enterChallenge: function(callback) {
    if (!this.results) {
      // if the user data are empty, receive it
      this.loadAllUserData(this.username, callback);
    } else {
      callback();
    }
  }




};

var UserResultData = React.createClass({displayName: 'UserResultData',
  render: function() {
    var user = this.props.user;
    var variableModels = this.props.variableModels;
    var ivnames = variableModels.iVariables.map(function(iv) {
      return iv.name;
    });
    var resultsDisplay = [];
    var allDisplay = [user.results.length];

    switch (this.props.mode) {
/*      case 'all':
        var results = user.results;
        resultsDisplay = results.map(function(result) {
          var index = result[variableModels.dvResultCount];
          return <UserResult variableModels={variableModels} data={result} index={'#' + index}/>;
        });
        break;*/
      case 'notebook':
        var newDisplay = null;
        var oldDisplay = null;
        for (var j=0; j < user.results.length; j++) {
          var result = user.results[j];
          var isNew = true;
          var isOld = true;
          for (var i = 0; i < ivnames.length; i++) {
            if ((!user.newCart) || (result[ivnames[i]] != user.newCart[ivnames[i]])) {
              isNew = false;
            }
            if ((!user.oldCart) || (result[ivnames[i]] != user.oldCart[ivnames[i]])) {
              isOld = false;
            }
          }
          var index = result[variableModels.dvResultCount];
          allDisplay[index-1] = UserResult( {variableModels:variableModels, data:result, index:'#' + index});

          if (isNew) {
            newDisplay = UserResult( {variableModels:variableModels, data:result, index:'#' + index + ' (Newly Saved)'});
          } else if (isOld) {
            oldDisplay = UserResult( {variableModels:variableModels, data:result, index:'#' + index + ' (Last Saved)'});
          }
        }
        if (newDisplay) {
          resultsDisplay.push(newDisplay);
        }
        if (oldDisplay) {
          resultsDisplay.push(oldDisplay);
        }
        break;
    }

    var headers = variableModels.iVariables.map(function(iv) {
      //return <th><VariableImage name={iv.name}/>{iv.label}</th>;
      return React.DOM.th(null, iv.label);
    });

    return React.DOM.table( {className:"result"}, 
      React.DOM.thead(null, 
        React.DOM.tr(null, 
          React.DOM.th(null
          ),
          headers,
          React.DOM.th(null, 
            variableModels.dvLabel
          )
        )
      ),
      React.DOM.tbody(null, 
        resultsDisplay,
        React.DOM.tr(null, React.DOM.td( {co:headers.length}, "All Results")),
        allDisplay
      )
    );
  },
});

var UserResult = React.createClass({displayName: 'UserResult',
  render: function() {
    var variableModels = this.props.variableModels;
    var data = this.props.data;
    var dvValues = data[variableModels.dvName].join(', ');
    var index = this.props.index;

    var variables = variableModels.iVariables.map(function(variable) {
      return UserResultSelection( {iv:variable, value:data[variable.name]});
    });

    return React.DOM.tr(null, 
      React.DOM.td(null, 
        " Cart ", index, " : "
      ),
      variables,
      React.DOM.td(null, 
        dvValues
      )
      );
  }
});

var UserResultSelection = React.createClass({displayName: 'UserResultSelection',
  getDisplayValue: function(value) {
    var options = this.props.iv.options;
    for (var i = 0; i < options.length; i++) {
      if (options[i].value == value) {
        return options[i].label;
      }
    }
    return null;
  },

  render: function() {
    var iv = this.props.iv;
    var ivValue = this.getDisplayValue(this.props.value);
    return React.DOM.td(null, ivValue);
  }
});


