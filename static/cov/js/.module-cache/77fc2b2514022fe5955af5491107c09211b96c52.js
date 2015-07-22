var Challenge = React.createClass({
  getInitialState: function() {
    return {};
  },

  handleChange: function(event) {
    var state = {};
    state[event.target.id] = event.target.value;
    this.setState(state);
  },

  handleSubmit: function(event) {
    event.preventDefault();

    this.post(this.state);
  },

  post: function(data) {
/*    if (!this.isEnabled())
      return;

    var xhr = new XMLHttpRequest();
    var self = this;
    xhr.onload = function() {
      if (self.props.onComplete) {
        self.props.onComplete(JSON.parse(xhr.responseText));
      }
    };
    xhr.open('POST', '/carts/gettrips');
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(data));
*/
  },

  isEnabled: function() {
    // if (!('findout' in this.state)) {
    //   return false;
    // }
    return true;
  },

  render: function() {
    var user = this.props.user;
    var variableModels = this.props.variableModels;
    var ivnames = variableModels.iVariables.map(function(iv) {
      return iv.name;
    });

    var variables = this.props.variableModels.iVariables.map(function(variable) {
      return <IndependentVariable iv={variable}/>;
    });

    switch (this.state.mode) {
      default:
        return <form onSubmit={this.handleSubmit} onChange={this.handleChange}
                className="request">
          <table><tbody>
            <tr>
              <td>
                What did you find out about whether the Handle Length makes a difference?
              </td>
              <td>
                <textarea id='handlelength'></textarea>
              </td>
            </tr>
            <tr>
              <td>
                What results show you are right?
              </td>
              <td>
                <textarea id='results'></textarea>
              </td>
            </tr>
          </tbody></table>
          <button type="submit" disabled={!this.isEnabled()}>Enter</button>
        </form>;
    }      
  
  }
});
