/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

var TwoRecordSelection = React.createClass({
  getInitialState: function() {
    return {enabled: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  // return an array of selected levels for each factor
  // f.FactorId : the id of a factor
  // f.SelectedLevelId: the id of the level selected for the factor
  getSelectedFactors: function(record) {
    var user = this.props.user;
    var prompt = user.CurrentUIPrompt;
    var form = document.getElementById("covactionForm");
    var selectedFactors = prompt.Factors.map(
      function(factor, i) {
        var fid = form.elements[factor.FactorId+record];
        var f = {};
        f.FactorId = factor.FactorId;
        f.SelectedLevelId = fid ? fid.value : "";
        return f;
      });
    return selectedFactors;
  },

  handleChange: function(event) {
    var selectedFactors = this.getSelectedFactors("1");
    for (i=0; i < selectedFactors.length; i++) {
      if (selectedFactors[i].SelectedLevelId == "") {
        return;
      }
    }    
    selectedFactors = this.getSelectedFactors("2");
    for (i=0; i < selectedFactors.length; i++) {
      if (selectedFactors[i].SelectedLevelId == "") {
        return;
      }
    }    
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault();

    var user = this.props.user;
    var prompt = user.CurrentUIPrompt;
    var onComplete = this.props.onComplete;
    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var r1selectedFactors = this.getSelectedFactors("1");
    var r2selectedFactors = this.getSelectedFactors("2");

    var response = {};
    response.RecordNoOne = r1selectedFactors;
    response.RecordNoTwo = r2selectedFactors;    

    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
    this.setState({mode: 0, enabled:false});
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();
    var promptId = prompt.PromptId;
    var phaseId = user.CurrentPhaseId;

    var recordOneFactors = prompt.Factors.map(
      function(factor, i) {
        return <FactorSelection factor={factor} key={i} record="1"/>;
      });
    var recordTwoFactors = prompt.Factors.map(
      function(factor, i) {
        return <FactorSelection factor={factor} key={i} record="2"/>;
      });

    return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
      <div className ="hbox">
        <div className="frame">
            <table>
              <tbody>
              <tr>
                <td>&nbsp;</td>
                <td colSpan="3" className="question">First Record</td>
              </tr>
              {recordOneFactors}
              </tbody>
            </table>
        </div>
        <div className="frame">
          <table>
            <tbody>
            <tr>
              <td>&nbsp;</td>
              <td colSpan="3" className="question">Second Record</td>
            </tr>
            {recordTwoFactors}
            </tbody>
          </table>
        </div>
      </div>
      <p>
        <input type="hidden" id="promptId" value={promptId}/>
        <input type="hidden" id="phaseId" value={phaseId}/>
        <button type="submit" disabled={!this.isEnabled()}>Enter</button>
      </p>
      </form>;
  }
});

var FactorSelection = React.createClass({
  render: function() {
    var state = this.state;
    var factor = this.props.factor;
    var record = this.props.record;

    var levels = factor.Levels.map(
      function(level, i) {
        return <FactorLevelSelection factor={factor} level={level} key={i} record={record}/>;
      });


    return <tr>
            <td>{factor.Text}</td>
            {levels}
          </tr>;
  }
});

var FactorLevelSelection = React.createClass({
  render: function() {
    var state = this.state;
    var factor = this.props.factor;
    var record = this.props.record;
    var level = this.props.level;
    var imgPath = "/img/"+level.ImgPath;
    var factorId = factor.FactorId+record;

    return <td><label>
            <input type="radio" name={factorId} value={level.FactorLevelId}><img src={imgPath}/><br/>{level.Text}</input></label></td>;
  }
});


var OneRecordPerformance = React.createClass({

  // getInitialState: function() {
  //   return {mode: 0};
  // },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;
      var prompt = user.getPrompt();

  		return  <div className="frame">
        <table className="record">
          <tbody>
            <tr>
              <td colSpan="3" className="robot">Record #18 <b>Daisy Smith</b></td>
              <td className="robot">Gender: F</td>
            </tr>
          </tbody>
        </table>
        <table className="xxx">
          <tbody>
          <tr>
            <td>Fitness</td>
            <td><label className="dimmed"><img src="/img/excellent fitness.jpg"/><br/>
            Excellent</label></td>
            <td>&nbsp;</td>
            <td><label><img src="/img/average fitness.jpg"/><br/>
            Average</label></td>
          </tr>
          <tr>
            <td>Parents health</td>
            <td><label><img src="/img/excellent parents.jpg"/><br/>
            Excellent</label></td>
            <td>&nbsp;</td>
            <td><label className="dimmed"><img src="/img/fair parents.jpg"/><br/>
            Fair</label></td>
          </tr>
          <tr>
            <td>Family size</td>
            <td><label><img src="/img/large family.jpg"/><br/>
            Large</label></td>
            <td>&nbsp;</td>
            <td><label className="dimmed"><img src="/img/small family.jpg"/><br/>
            Small</label></td>
          </tr>
          <tr>
            <td>Education</td>
            <td><label className="dimmed"><img src="/img/college.jpg"/><br/>
            College</label></td>
            <td><label className="dimmed"><img src="/img/some college.jpg"/><br/>
            Some College</label></td>
            <td><label><img src="/img/no college.jpg"/><br/>
            No College</label></td>
          </tr>
          </tbody>
        </table>
        <p className="performance-level">Performance Level:
          <span className="grade">D</span>
        </p>
        </div>;
  }
});
