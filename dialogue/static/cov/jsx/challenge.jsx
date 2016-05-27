/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

// var Challenge = React.createClass({

//   render: function() {
//     var state = this.state;
//     var user = this.props.user;
//     var app = this.props.app;
//     return  <div className="action">
//               <p>
//                 <a href="chart.html" className="button">Show Charts</a>
//                 <br/><br/>
//               </p>
//               <table>
//                 <tr>
//                   <td colspan="2">What would you need to predict his/her performance?</td>
//                 </tr>
//                 <tr>
//                   <td colspan="2">&nbsp;</td>
//                 </tr>
//               </table>
//               <p>
//                 <a href="dialog26.html" className="button">OK</a>
//               </p>
//             </div>;
//   }
// });



// function ChallengeHandler(challenge) {
//   this.username = username;
//   this.oldCart = null;
//   this.newCart = null;
//   this.results = null;
//   this.currentChallengeID = null;
// }

// User.prototype = {

//   loadAllUserData: function(renderCallback) {
//     var self = this;
//     var cartPromise = self.loadUserResultData(self.username);

//     var challengePromise = cartPromise.then(function(username) {
//                                               return self.loadUserChallengeData(username);
//                                             });
//     challengePromise.then(renderCallback, function(error) {
//                                             console.error("Failed!", error);
//                                           });
//   },

// };














// var Challenge = React.createClass({
//   getInitialState: function() {
//     this.setState({enabled:false});
//     debugger;
//     return {};
//   },

//   handleChange: function(event) {
//     // var state = {};
//     // state[event.target.id] = event.target.value;
//     this.setState({enabled:true});
//   },

//   handleSubmit: function(event) {
//     event.preventDefault();

//     this.post(this.state);
//   },

//   post: function(data) {
// /*    if (!this.isEnabled())
//       return;

//     var xhr = new XMLHttpRequest();
//     var self = this;
//     xhr.onload = function() {
//       if (self.props.onComplete) {
//         self.props.onComplete(JSON.parse(xhr.responseText));
//       }
//     };
//     xhr.open('POST', '/carts/gettrips');
//     xhr.setRequestHeader('Content-Type', 'application/json');
//     xhr.send(JSON.stringify(data));
// */
//   },

//   isEnabled: function() {
//     return this.state.enabled;
//   },

//   render: function() {
//     var user = this.props.user;
//     var variableModels = this.props.variableModels;
//     var ivnames = variableModels.iVariables.map(function(iv) {
//       return iv.name;
//     });

//     var variables = this.props.variableModels.iVariables.map(function(variable) {
//       return <IndependentVariable iv={variable}/>;
//     });

//     switch (this.state.mode) {
//       default:
//         return <form onSubmit={this.handleSubmit} onChange={this.handleChange}
//                 className="request">
//           <table><tbody>
//             <tr>
//               <td>
//                 What did you find out about whether the Handle Length makes a difference?
//               </td>
//               <td>
//                 <textarea id='handlelength'></textarea>
//               </td>
//             </tr>
//             <tr>
//               <td>
//                 What results show you are right?
//               </td>
//               <td>
//                 <textarea id='results'></textarea>
//               </td>
//             </tr>
//           </tbody></table>
//           <button type="submit" disabled={!this.isEnabled()}>Enter</button>
//         </form>;
//     }      
  
//   }
// });
