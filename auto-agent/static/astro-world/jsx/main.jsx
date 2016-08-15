/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

function getQueryStringValue (key) {
  return unescape(window.location.search.replace(new RegExp("^(?:.*[&\\?]" +
      encodeURIComponent(key) + "(?:\\=([^&]*))?)?.*$", "i"), "$1"));  
}  

// Begin side-effect
var username = getQueryStringValue("user");

if (!username) {
  window.location = "index.html";
} else {

  var user = new User(username);

  // React.renderComponent(<App variableModels={variableModels} user={user}/>,
  //                   document.body);

  user.loadAllUserData(function() {
                         React.render(
                           <App user={user}/>,
                         document.body);});
}

// TODO
// window.onbeforeunload = function() {
//   return "";
// };

