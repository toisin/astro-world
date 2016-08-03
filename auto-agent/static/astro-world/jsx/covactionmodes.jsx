/** @jsx React.DOM */

// npm install -g react-tools
// jsx -w -x jsx public/js public/js
var SelectTargetFactor = React.createClass({

  getInitialState: function() {
    return {enabled: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  handleChange: function(event) {
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault(); // default might be to follow a link, instead, takes control over the event

    var user = this.props.user;
    var onComplete = this.props.onComplete;
    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");
    e = f.elements['covactioninput'];
    var value = e ? e.value : "";
    e.value = "";
    var text, id;

    var options = user.getPrompt().Options;
    for (i = 0; i < options.length; i++) {
      if (options[i].ResponseId == value) {
        text = options[i].Text;
        id = value;
        break;
      }
    }

    var response = {};
    response.text = text;
    response.id = id;
    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
    var user = this.props.user;
    var prompt = user.getPrompt();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();

    if (!prompt.Options) {
      console.error("Error: Select factor UI without options!");    
      return <div></div>;
    }
    var options = prompt.Options.map(
      function(option, i) {
        return <FactorPromptOption option={option} key={i}/>;
      });

    return   <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
              <div className ="hbox">
                <div className="frame">
                    <table>
                      <tbody>
                      <tr>
                        <td>{prompt.Text}</td>
                      </tr>
                      {options}
                      </tbody>
                    </table>
                </div>
              </div>
              <p>
                <input type="hidden" id="promptId" value={promptId}/>
                <input type="hidden" id="phaseId" value={phaseId}/>
                <button type="submit" disabled={!this.isEnabled()} key={"SelectTargetFactor"}>Enter</button>
              </p>
              </form>;
  },
});

var FactorPromptOption = React.createClass({

  render: function() {
    var option = this.props.option;
      return <tr><td><label>
              <input type="radio" name="covactioninput" value={option.ResponseId}><br/>{option.Text}</input></label></td></tr>;
  },
});

var PriorBeliefFactors = React.createClass({

  getInitialState: function() {
    return {enabled: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  // return an array of selected levels for each factor
  // f.FactorId : the id of a factor
  // f.SelectedLevelId: the id of the level selected for the factor
  getSelectedFactors: function() {
    var user = this.props.user;
    var prompt = user.getPrompt();
    var form = document.getElementById("covactionForm");
    var selectedFactors = user.getContentFactors().map(
      function(factor, i) {
        var fid = form.elements[factor.FactorId];
        var f = {};
        f.FactorId = factor.FactorId;
        f.IsCausal = fid.value == "true" ? true : false;
        return f;
      });
    return selectedFactors;
  },

  handleChange: function(event) {
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault();

    var user = this.props.user;
    var prompt = user.getPrompt();
    var onComplete = this.props.onComplete;
    var singleRecord = this.props.singleRecord;

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var response = {};
    response.CausalFactors = this.getSelectedFactors();;

    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
    var user = this.props.user;
    var prompt = user.getPrompt();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();

    var factors = user.getContentFactors().map(
      function(factor, i) {
        var factorId = factor.FactorId;

        return <tr  key={i}>
                <td>{factor.Text}</td>
                <td><label>
                  <input type="radio" name={factorId} value={true}><br/>Yes</input>
                </label></td>
                <td><label>
                  <input type="radio" name={factorId} value={false}><br/>No</input>
                </label></td>
              </tr>;
      });


    return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
      <div className ="hbox">
        <div className="frame">
            <table>
              <tbody>
              <tr>
                <td>&nbsp;</td>
                <td colSpan="3" className="question">Does it make a difference? (Select "Yes" or "No")</td>
              </tr>
              {factors}
              </tbody>
            </table>
        </div>
      </div>
      <p>
        <input type="hidden" id="promptId" value={promptId}/>
        <input type="hidden" id="phaseId" value={phaseId}/>
        <button type="submit" disabled={!this.isEnabled()} key={"PriorBeliefFactors"}>Enter</button>
      </p>
      </form>;
  },
});

var FactorLevelPriorBeliefSelection = React.createClass({
  render: function() {
    var factor = this.props.factor;
    var level = this.props.level;
    var imgPath = "/img/"+level.ImgPath;
    var factorId = factor.FactorId;

    return <td><label>
            <input type="radio" name={factorId} value={level.FactorLevelId}><img src={imgPath}/><br/>{level.Text}</input>
          </label></td>;
  }
});

var PriorBeliefLevels = React.createClass({

  getInitialState: function() {
    return {enabled: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  // return an array of selected levels for each factor
  // f.FactorId : the id of a factor
  // f.SelectedLevelId: the id of the level selected for the factor
  getSelectedFactors: function() {
    var user = this.props.user;
    var prompt = user.getPrompt();
    var form = document.getElementById("covactionForm");
    var selectedFactors = user.getContentFactors().map(
      function(factor, i) {
        var fid = form.elements[factor.FactorId];
        var f = {};
        f.FactorId = factor.FactorId;
        f.BestLevelId = fid ? fid.value : "";
        f.IsCausal = factor.IsBeliefCausal;
        return f;
      });
    return selectedFactors;
  },

  handleChange: function(event) {
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault();

    var user = this.props.user;
    var prompt = user.getPrompt();
    var onComplete = this.props.onComplete;
    var singleRecord = this.props.singleRecord;

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var response = {};
    response.CausalFactors = this.getSelectedFactors();;

    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
    var user = this.props.user;
    var prompt = user.getPrompt();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();

    var factors = user.getContentFactors().map(
      function(factor, i) {
        if (factor.IsBeliefCausal) {
          var factorId = factor.FactorId;

          var levels = factor.Levels.map(
            function(level, j) {
              return <FactorLevelPriorBeliefSelection factor={factor} level={level} key={j}/>;
            });

          return <tr key={i}>
                  <td>{factor.Text}</td>
                  {levels}
                </tr>;
        } else {
          return "";
        }
      });


    return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
      <div className ="hbox">
        <div className="frame">
            <table>
              <tbody>
              <tr>
                <td>&nbsp;</td>
                <td colSpan="3" className="question">Choose the level of the factor that you think would be best for performance.</td>
              </tr>
              {factors}
              </tbody>
            </table>
        </div>
      </div>
      <p>
        <input type="hidden" id="promptId" value={promptId}/>
        <input type="hidden" id="phaseId" value={phaseId}/>
        <button type="submit" disabled={!this.isEnabled()} key={"PriorBeliefFactors"}>Enter</button>
      </p>
      </form>;
  },
});

var RecordSelection = React.createClass({
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
    var prompt = user.getPrompt();
    var form = document.getElementById("covactionForm");
    var selectedFactors = user.getContentFactors().map(
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
    var singleRecord = this.props.singleRecord;

    var selectedFactors = this.getSelectedFactors("1");
    for (i=0; i < selectedFactors.length; i++) {
      if (selectedFactors[i].SelectedLevelId == "") {
        return;
      }
    }
    if (!singleRecord) {
      selectedFactors = this.getSelectedFactors("2");
      for (i=0; i < selectedFactors.length; i++) {
        if (selectedFactors[i].SelectedLevelId == "") {
          return;
        }
      }    
    }
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault();

    var user = this.props.user;
    var prompt = user.getPrompt();
    var onComplete = this.props.onComplete;
    var singleRecord = this.props.singleRecord;

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var r1selectedFactors = this.getSelectedFactors("1");
    var r2selectedFactors
    if (!singleRecord) {
      r2selectedFactors = this.getSelectedFactors("2");
    }
    var response = {};
    response.RecordNoOne = r1selectedFactors;
    response.RecordNoTwo = r2selectedFactors;    

    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var singleRecord = this.props.singleRecord;
    var prompt = user.getPrompt();
    var factors = user.getContentFactors();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();

    var recordOneFactors = factors.map(
      function(factor, i) {
        return <FactorSelection factor={factor} key={i} record="1"/>;
      });
    var recordTwoFactors = factors.map(
      function(factor, i) {
        return <FactorSelection factor={factor} key={i} record="2"/>;
      });

    if (singleRecord) {
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
        </div>
        <p>
          <input type="hidden" id="promptId" value={promptId}/>
          <input type="hidden" id="phaseId" value={phaseId}/>
          <button type="submit" disabled={!this.isEnabled()} key={"RecordSelection"}>Enter</button>
        </p>
        </form>;
    } else {
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
          <button type="submit" disabled={!this.isEnabled()} key={"RecordSelection"}>Enter</button>
        </p>
        </form>;
      }
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


var RecordPerformance = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  handleSubmit: function(event) {
    event.preventDefault();

    var user = this.props.user;
    var prompt = user.getPrompt();
    var onComplete = this.props.onComplete;
    var singleRecord = this.props.singleRecord;

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var response = {};

    jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;
      var prompt = user.getPrompt();
      var promptId = prompt.PromptId;
      var record1 = user.getState().RecordNoOne;
      var record2 = user.getState().RecordNoTwo;
      var factors = user.getContentFactors().map(
        function(factor, i) {
          var fid = factor.FactorId;
          var selectedf = record1.FactorLevels[fid];
          var SelectedLevelName = selectedf.SelectedLevel;

          var levels = factor.Levels.map(
            function(level, j) {
              var imgPath = "/img/"+level.ImgPath;
              if (level.Text == SelectedLevelName) {
                return <td key={j}><label>
                    <img src={imgPath}/><br/>{level.Text}</label></td>;
              }
              return <td key={j}><label className="dimmed">
                    <img src={imgPath}/><br/>{level.Text}</label></td>;
            });

          return <tr key={i}>
                  <td>{factor.Text}</td>
                  {levels}
                </tr>;
        });

      if (record1) {
      return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
          <div className="frame">
          <table className="record">
            <tbody>
              <tr>
                <td colSpan="3" className="robot">Record #{record1.RecordNo} <b>{record1.RecordName}</b></td>
                <td className="robot">Gender: ??</td>
              </tr>
            </tbody>
          </table>
          <table className="recorddetails">
            <tbody>
            {factors}
            </tbody>
          </table>
          <p className="performance-level">Performance Level:
            <span className="grade">{record1.Performance}</span>
          </p>
          </div>
          <p>
            <input type="hidden" id="promptId" value={promptId}/>
            <input type="hidden" id="phaseId" value={phaseId}/>
            <button type="submit" key={"RecordPerformance"}>Enter</button>
          </p>
          </form>;
      } else {
        return <div></div>;
      }
      // if (record1 && record2) {
      //   return  <div className="frame">
      //   <table className="record">
      //     <tbody>
      //       <tr>
      //         <td colSpan="3" className="robot">Record #{record1.RecordNo} <b>{record1.RecordName}</b></td>
      //         <td className="robot">Gender: ??</td>
      //       </tr>
      //     </tbody>
      //   </table>
      //   <table className="recorddetails">
      //     <tbody>
      //     <tr>
      //       <td>Fitness</td>
      //       <td><label className="dimmed"><img src="/img/excellent fitness.jpg"/><br/>
      //       Excellent</label></td>
      //       <td>&nbsp;</td>
      //       <td><label><img src="/img/average fitness.jpg"/><br/>
      //       Average</label></td>
      //     </tr>
      //     <tr>
      //       <td>Parents health</td>
      //       <td><label><img src="/img/excellent parents.jpg"/><br/>
      //       Excellent</label></td>
      //       <td>&nbsp;</td>
      //       <td><label className="dimmed"><img src="/img/fair parents.jpg"/><br/>
      //       Fair</label></td>
      //     </tr>
      //     <tr>
      //       <td>Family size</td>
      //       <td><label><img src="/img/large family.jpg"/><br/>
      //       Large</label></td>
      //       <td>&nbsp;</td>
      //       <td><label className="dimmed"><img src="/img/small family.jpg"/><br/>
      //       Small</label></td>
      //     </tr>
      //     <tr>
      //       <td>Education</td>
      //       <td><label className="dimmed"><img src="/img/college.jpg"/><br/>
      //       College</label></td>
      //       <td><label className="dimmed"><img src="/img/some college.jpg"/><br/>
      //       Some College</label></td>
      //       <td><label><img src="/img/no college.jpg"/><br/>
      //       No College</label></td>
      //     </tr>
      //     </tbody>
      //   </table>
      //   <p className="performance-level">Performance Level:
      //     <span className="grade">D</span>
      //   </p>
      //   </div>;
      // }
  }
});
