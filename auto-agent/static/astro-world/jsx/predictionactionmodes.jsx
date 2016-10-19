/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

function FactorsRequestForm(props) {
  var question = "Check the box for up to four factors that you would like to know about an applicant.";
  var formName = "predictionactionForm";
  return <div>
          <MultipleFactorsSelect formName={formName} question={question} user={props.user} onComplete={props.onComplete} app={props.app}/>
          <ChartButtons user={props.user} app={props.app}/>
        </div>;
}

function ContributingFactorsForm(props) {
  var question = "Which of the four factors you have data on mattered to your prediction?";
  var formName = "predictionactionForm";
  var factors = props.user.getState().DisplayFactors;
  return <div>
          <MultipleFactorsSelect formName={formName} factors={factors} question={question} user={props.user} onComplete={props.onComplete} app={props.app}/>
          <ChartButtons user={props.user} app={props.app}/>
          <PredictionRecord user={props.user} onComplete={props.onComplete} app={props.app}/>
        </div>;
}

var MultipleFactorsSelect = React.createClass({

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
    var formName = this.props.formName;

    var prompt = user.getPrompt();
    var form = document.getElementById(formName);

    var factorOrder = [];
    var tempfactorIdsMap;

    if (this.props.factors) {
      tempfactorIdsMap = this.props.factors.map(
        function(v, i) {return v.FactorId});
    } else {
      tempfactorIdsMap = Object.keys(user.getContentFactors());
    }

    var tempfactors = tempfactorIdsMap.map(
      function(fkey, i) {
        var factor = user.getContentFactors()[fkey];
        factorOrder[i] = factor.Order;
        var fid = form.elements[factor.FactorId];
        if (fid) {
          var f = {};
          f.FactorId = factor.FactorId;
          f.IsBeliefCausal = fid.checked;
          return f;
        }
      });

    var selectedFactors = [];
    for (var i = 0; i < tempfactors.length; i++) {
      selectedFactors[factorOrder[i]] = tempfactors[i];
    }
    return selectedFactors;
  },

  handleChange: function(event) {
    this.setState({enabled:true});
  },

  handleSubmit: function(event) {
    event.preventDefault();

    var user = this.props.user;
    var formName = this.props.formName;

    var prompt = user.getPrompt();
    var onComplete = this.props.onComplete;

    var e = document.getElementById("promptId");
    var promptId = e ? e.value : "";
    var e = document.getElementById("phaseId");
    var phaseId = e ? e.value : "";
    var f = document.getElementById(formName);

    var response = {};
    response.BeliefFactors = this.getSelectedFactors();

    var count = 0;
    for (var i = 0; i < response.BeliefFactors.length; i++) {
      if (response.BeliefFactors[i].IsBeliefCausal) {
        count++;
      }
    }

    if (count > 4) {
      alert("You have selected more than 4 factors. Please remove at least 1 and try again.");
    } else {
      var jsonResponse = JSON.stringify(response);
      user.submitResponse(promptId, phaseId, jsonResponse, onComplete);
    }
  },

  render: function() {
    var user = this.props.user;
    var formName = this.props.formName;
    var question = this.props.question;

    var prompt = user.getPrompt();

    var promptId = prompt.PromptId;
    var phaseId = user.getCurrentPhaseId();

    var factorOrder = [];
    var tempfactorIdsMap;

    if (this.props.factors) {
      tempfactorIdsMap = this.props.factors.map(
        function(v, i) {return v.FactorId});
    } else {
      tempfactorIdsMap = Object.keys(user.getContentFactors());
    }

    var tempfactors = tempfactorIdsMap.map(
      function(fkey, i) {
        var factor = user.getContentFactors()[fkey];
        var factorId = factor.FactorId;
        factorOrder[i] = factor.Order;

        return <tr  key={i}>
                <td><label>
                  <input type="checkbox" name={factorId}/>
                </label></td>
                <td className="factorNameFront">{factor.Text}</td>
              </tr>;
      });

    var factors = [];
    for (var i = 0; i < tempfactors.length; i++) {
      factors[factorOrder[i]] = tempfactors[i];
    }

    return <form id={formName} onSubmit={this.handleSubmit} onChange={this.handleChange}>
      <div className ="hbox">
        <div className="frame">
            <table>
              <tbody>
              <tr>
                <td colSpan="4" className="question">{question}</td>
              </tr>
              {factors}
              </tbody>
            </table>
        </div>
      </div>
      <p>
        <input type="hidden" id="promptId" value={promptId}/>
        <input type="hidden" id="phaseId" value={phaseId}/>
        <button type="submit" disabled={!this.isEnabled()} key={"MultipleFactorsSelect"}>Enter</button>
      </p>
      </form>;
  },
});

var PredictionRecord = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;

      var prompt = user.getPrompt();
      var promptId = prompt.PromptId;
      var phaseId = user.getCurrentPhaseId();
      var record = user.getState().TargetPrediction;

      var recordDetails = function(r) {
        var factorOrder = [];
        var tempfactors = user.getState().DisplayFactors.map(
          function(v, i) {
            var factor = v;
            factorOrder[i] = factor.Order;
            var fid = factor.FactorId;
            var selectedf = r.FactorLevels[fid];
            var SelectedLevelName = selectedf.SelectedLevel;

            var size = factor.Levels.length;
            if (size == 2) {
              factor.Levels[2]=factor.Levels[1];
              factor.Levels[1]="_";
            }
            var levels = factor.Levels.map(
              function(level, j) {
                if (level.ImgPath) {
                  var imgPath = "/img/"+level.ImgPath;
                  if (level.Text == SelectedLevelName) {
                    return <td key={j}><label>
                        <img src={imgPath}/><br/>{level.Text}</label></td>;
                  }
                }
                return <td key={j}><label className="dimmed">
                      <img src={imgPath}/><br/>{level.Text}</label></td>;
              });

            return <tr key={i}>
                    <td className="factorNameFront">{factor.Text}</td>
                    {levels}
                  </tr>;
          });
        var factors = [];
        for (var i = 0; i < tempfactors.length; i++) {
          factors[factorOrder[i]] = tempfactors[i];
        }

        return r ? <div className="frame" key={r.RecordNo}>
                <table className="record">
                  <tbody>
                    <tr>
                      <td colSpan="4" className="robot">Applicant #{r.RecordNo} <b>{r.RecordName}</b></td>
                    </tr>
                    {factors}
                  </tbody>
                </table>
              </div> : null;};
              
      var recordDetails
      recordDetails = recordDetails(record);
      return <div className = "hbox">
                {recordDetails}
              </div>;
  }
});

