/** @jsx React.DOM */
"use strict"

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
    for (var i = 0; i < options.length; i++) {
      if (options[i].ResponseId == value) {
        text = options[i].Text;
        id = value;
        break;
      }
    }

    var response = {};
    response.text = text;
    response.id = id;
    var jsonResponse = JSON.stringify(response);
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
                      <tr><td className="question">Select the factor to investigate</td></tr>
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
              <input type="radio" name="covactioninput" value={option.ResponseId}/><br/>{option.Text}</label></td></tr>;
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
        f.IsBeliefCausal = fid.value == "true" ? true : false;
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

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var response = {};
    response.BeliefFactors = this.getSelectedFactors();

    var jsonResponse = JSON.stringify(response);
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
                <td className="factorNameFront">{factor.Text}</td>
                <td><label>
                  <input type="radio" name={factorId} value={true}/><br/>Yes
                </label></td>
                <td><label>
                  <input type="radio" name={factorId} value={false}/><br/>No
                </label></td>
              </tr>;
      });


    return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
      <div className ="hbox">
        <div className="frame">
            <table>
              <tbody>
              <tr>
                <td colSpan="4" className="question">Select "Yes" for factor that you think makes a difference.</td>
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
            <input type="radio" name={factorId} value={level.FactorLevelId}/><img src={imgPath}/><br/>{level.Text}
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
        f.IsBeliefCausal = factor.IsBeliefCausal;
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

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var response = {};
    response.BeliefFactors = this.getSelectedFactors();

    var jsonResponse = JSON.stringify(response);
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
                  <td className="factorNameFront">{factor.Text}</td>
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
        if (fid) {
          var f = {};
          f.FactorId = factor.FactorId;
          f.SelectedLevelId = fid ? fid.value : "";
          return f;
        }
      });
    return selectedFactors;
  },

  handleChange: function(event) {
    var doubleRecord = this.props.doubleRecord;
    var comparePrevious = this.props.comparePrevious;

    var selectedFactors;
    if (!comparePrevious) {
      selectedFactors = this.getSelectedFactors("1");
      for (var i = 0; i < selectedFactors.length; i++) {
        if (selectedFactors[i].SelectedLevelId == "") {
          return;
        }
      }
    }
    if (doubleRecord) {
      selectedFactors = this.getSelectedFactors("2");
      for (var i = 0; i < selectedFactors.length; i++) {
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
    var doubleRecord = this.props.doubleRecord;
    var comparePrevious = this.props.comparePrevious;

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById("covactionForm");

    var response = {};

    var r1selectedFactors
    if (!comparePrevious) {
      r1selectedFactors = this.getSelectedFactors("1");
    } else {
      response.UseDBRecordNoOne = true;
    }
    var r2selectedFactors
    if (doubleRecord) {
      r2selectedFactors = this.getSelectedFactors("2");
    }
    if (r1selectedFactors && (r1selectedFactors.length > 0)) {
      response.RecordNoOne = r1selectedFactors;
    }
    if (r2selectedFactors && (r2selectedFactors.length > 0)) {
      response.RecordNoTwo = r2selectedFactors;    
    }
    var jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var doubleRecord = this.props.doubleRecord;
    var comparePrevious = this.props.comparePrevious;
    var prompt = user.getPrompt();
    var factors = user.getContentFactors();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();
    var recordOneFactors = {}
    if (!comparePrevious) {
      recordOneFactors = factors.map(
        function(factor, i) {
          return <FactorSelection factor={factor} key={i} record="1"/>;
        });
    }
    var recordTwoFactors = factors.map(
      function(factor, i) {
        return <FactorSelection factor={factor} key={i} record="2"/>;
      });

    if (!doubleRecord) {
      return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
        <div className ="hbox">
          <div className="frame">
              <table className="record">
                <tbody>
                <tr>
                  <td colSpan="4" className="recordTitle">First Record</td>
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
    } else if (!comparePrevious) {
      return <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
        <div className ="hbox">
          <div className="frame">
              <table className="record">
                <tbody>
                <tr>
                  <td colSpan="4" className="recordTitle">First Record</td>
                </tr>
                {recordOneFactors}
                </tbody>
              </table>
          </div>
          <div className="frame">
            <table className="record">
              <tbody>
              <tr>
                <td colSpan="4" className="recordTitle">Second Record</td>
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
    } else {
      return <div className ="hbox">
              <div className="frame">
                <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
                <table className="record">
                  <tbody>
                  <tr>
                    <td colSpan="4" className="recordTitle">Second Record</td>
                  </tr>
                    {recordTwoFactors}
                  </tbody>
                </table>
                <p>
                  <input type="hidden" id="promptId" value={promptId}/>
                  <input type="hidden" id="phaseId" value={phaseId}/>
                  <button type="submit" disabled={!this.isEnabled()} key={"RecordSelection"}>Enter</button>
                </p>
                </form>
              </div>
              <RecordPerformance user={user} app={app} recordOneOnly/>
            </div>;
    }
  }
});

var FactorSelection = React.createClass({
  render: function() {
    var state = this.state;
    var factor = this.props.factor;
    var record = this.props.record;

    var size = factor.Levels.length;
    if (size == 2) {
      factor.Levels[2]=factor.Levels[1];
      factor.Levels[1]="_";
    }

    var levels = factor.Levels.map(
      function(level, i) {
        if (level == "_") {
          return <td>&nbsp;</td>;
        }
        return <FactorLevelSelection factor={factor} level={level} key={i} record={record} size={size}/>;        
      });


    return <tbody>
            <tr>
              <td colSpan="3" className="factorNameRow">{factor.Text}</td>
            </tr>
            <tr>
              {levels}
            </tr>
          </tbody>;
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
            <input type="radio" name={factorId} value={level.FactorLevelId}/><img src={imgPath}/><br/>{level.Text}</label></td>;
  }
});


var RecordPerformance = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;
      var recordOneOnly = this.props.recordOneOnly;
      var recordTwoOnly = this.props.recordTwoOnly;
      var hidePerformance = this.props.hidePerformance;

      var prompt = user.getPrompt();
      var promptId = prompt.PromptId;
      var phaseId = user.getCurrentPhaseId();
      var record1 = user.getState().RecordNoOne;
      var record2 = user.getState().RecordNoTwo;

      var performance = function(r) {
        return !hidePerformance ? <p className="performance-level">Performance Level:
                    <span className="grade">{r.Performance}</span>
                  </p> : null;};

      var recordDetails = function(r) {
        var factors = user.getContentFactors().map(
          function(factor, i) {
            var fid = factor.FactorId;
            var selectedf = r.FactorLevels[fid];
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
                    <td className="factorNameFront">{factor.Text}</td>
                    {levels}
                  </tr>;
          });
        return r ? <div className="frame" key={r.RecordNo}>
                <table className="record">
                  <tbody>
                    <tr>
                      <td colSpan="4" className="robot">Record #{r.RecordNo} <b>{r.RecordName}</b></td>
                    </tr>
                    {factors}
                  </tbody>
                </table>
                {performance(r)}
              </div> : null;};
              
      var record1Details, record2Details
      if (record1 && !recordTwoOnly) {
        record1Details = recordDetails(record1);
      }
      if (record2 && !recordOneOnly) {
        record2Details = recordDetails(record2);
      }
      return <div className ="hbox">
                {record1Details}
                {record2Details}
              </div>
  }
});

var MemoForm = React.createClass({
  getInitialState: function() {
    return {enabled: false};
  },

  isEnabled: function() {
    return this.state.enabled;
  },

  handleChange: function(event) {
    var form = document.getElementById("covactionForm");
    var memo = form.elements["memo"];
    var evidence = form.elements["evidence"];
    if (memo.value && evidence.value) {
      this.setState({enabled:true});
    }
    return;
  },

  handleEnter: function(event) {
    if (!event.shiftKey) {
      if (event.which == 13) {  // "Enter" key was pressed.
        this.handleSubmit(event);
      }
    }
  },

  handleSubmit: function(event) {    
    if (event) {
      event.preventDefault();
    }

    var user = this.props.user;
    var targetFactorName, targetFactorId;
    if (user.getState().TargetFactor) {
      targetFactorName = user.getState().TargetFactor.FactorName;
      targetFactorId = user.getState().TargetFactor.FactorId;
    }
    var onComplete = this.props.onComplete;
    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var form = document.getElementById("covactionForm");
    var ask = form.elements["ask"];
    var memo = form.elements["memo"];
    var evidence = form.elements["evidence"];

    var response = {};
    response.Ask = ask ? ask.value : "";
    response.Memo = memo ? memo.value : "";
    response.Evidence = evidence ? evidence.value : "";
    response.Id = targetFactorId
    response.FactorName = targetFactorName
    
    var jsonResponse = JSON.stringify(response);
    user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();

    var targetFactorName;
    if (user.getState().TargetFactor) {
      targetFactorName = user.getState().TargetFactor.FactorName;
    }
    var investigatingFactorHeading;
    if (targetFactorName) {
      investigatingFactorHeading = <h3>Investigating Factor: <b>{targetFactorName}</b></h3>;
    }

    return <div>
            <div className="mbox">
              <h3>Memo to the foundation</h3>
              <form id="covactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
              <p>
                  We recommend that you &nbsp;
                  <input type="text" name="ask" size="20" autofocus className="con" placeholder="Enter ask/do not ask"/> &nbsp;
                  applicants about <u>{targetFactorName}</u> because &nbsp;
                  <input type="text" name="memo" autofocus className="con" placeholder="Enter if it does/does not make a difference."/><br/>
                  <br/>
                  Our evidence for claiming this is:<br/>
                  <textarea name="evidence" className="evid" onKeyDown={this.handleEnter} placeholder="Enter your answer here"></textarea>
                  <br/>
              </p>
              <p>
                <input type="hidden" id="promptId" value={promptId}/>
                <input type="hidden" id="phaseId" value={phaseId}/>
                <button type="submit" disabled={!this.isEnabled()} key={"MemoForm"}>Enter</button>
              </p>
              </form>
            </div>
            <div>
              {investigatingFactorHeading}
              <RecordPerformance user={user} app={app}/>
            </div>
           </div>;
  }
});

var Memo = React.createClass({
  getInitialState: function() {
    return {};
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var prompt = user.getPrompt();
    var targetFactorName, targetFactorId;
    if (user.getState().TargetFactor) {
      targetFactorName = user.getState().TargetFactor.FactorName;
      targetFactorId = user.getState().TargetFactor.FactorId;
    }


    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();
    var ask, memo, evidence;

    if (user.getState().LastMemo) {
      ask = user.getState().LastMemo.Ask;
      memo = user.getState().LastMemo.Memo;
      evidence = user.getState().LastMemo.Evidence;
    }

    return <div className="mbox">
              <h3>Memo to the foundation</h3>
              <p>
                  We recommend that you <u>{ask}</u> applicants about <u>{targetFactorName}</u> because <u>{memo}</u><br/>
                  <br/>
                  Our evidence for claiming this is:<br/>
                  <u>{evidence}</u>
                  <br/>
              </p>
            </div>;
  }
});
