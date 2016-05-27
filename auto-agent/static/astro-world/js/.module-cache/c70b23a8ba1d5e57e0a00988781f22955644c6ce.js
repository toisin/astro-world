/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var variableModels = {

  dvLabel: 'Number of trips per hour',
  dvName: 'trips',
  dvResultCount: 'cartNumber',


  iVariables: [
    {
      name: 'handleLength',
      label: 'Handle length',
      options: [
        {value: 'Long', label: 'Long'},
        {value: 'Short', label: 'Short'},
      ]
    },
    {
      name: 'wheelSize',
      label: 'Wheel Size',
      options: [
        {value: 'Large(4)', label: 'Large(4)'},
        {value: 'Small(3)', label: 'Small(3)'}
      ],
    },
    {
      name: 'bucketSize',
      label: 'Bucket Size',
      options: [
        {value: 'Big(13)', label: 'Big(13)'},
        {value: 'Small(10)', label: 'Small(10)'},
      ]
    },
    {
      name: 'bucketPlacement',
      label: 'Bucket Placement',
      options: [
        {value: 'Far', label: 'Far'},
        {value: 'Near', label: 'Near'},
      ]
    }
  ]
};

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

  user.loadAllUserData(callback = function() {
                         React.render(
                           React.createElement(App, {user: user}),
                         document.body);});
}

// TODO
// window.onbeforeunload = function() {
//   return "";
// };

