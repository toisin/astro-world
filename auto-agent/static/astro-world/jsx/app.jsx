/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var App = React.createClass({
  getInitialState: function() {
    return {mode: 0, actionReady: false};
  },

  showAction: function(show) {
    this.setState({mode: 0, actionReady: true});
  },

  changeState: function() {
    this.setState({mode: 0});
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    if (!this.state.actionReady) {
      return  <div className="content">
                  <Dialog user={user} app={this}/>
              </div>;
    } else {
      this.state.actionReady = false;
      return  <div className="content">
                  <Dialog user={user} app={this}/>
                  <Action user={user} app={this}/>
              </div>;
    }
  }

});

