/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

const rectSize = 8;
const toolBoxSizeHeight = 200;
const toolBoxSizeWidth = 250;
const spacingFactor = 2.2;
const rowHeight = 100;
const elementsPerRow = 5;
const columnWidth = 225;
const columnLabelHeight = 20;
const paddingLeft = 50;
const paddingBottom = 50;
const paddingTop = 20;
const paddingRight = 20;
const rowColors = ['transparent', 'rgba(0,0,0,0.05)'];
// const data = [{
//   label: 'Average',
//   rcount: [5, 7, 3, 28, 1]
// }, {
//   label: 'Good',
//   rcount: [15, 2, 13, 8, 11]
// }];


var Chart = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  render: function() {
    var state = this.state;
    var user = this.props.user;
    var app = this.props.app;
    var yLabels = [];
    var performanceLabels = [{
                              grade: 'A',
                              label: 'A (very well)'
                            },
                            {
                              grade: 'B',
                              label: 'B (well)'
                            },
                            {
                              grade: 'C',
                              label: 'C (so so)'
                            },
                            {
                              grade: 'D',
                              label: 'D (poorly)'
                            },
                            {
                              grade: 'E',
                              label: 'E (very poorly)'
                            }];
    var data = [];
    var pRecords = user.getAllPerformanceRecords();

    var colFilters = [];
    if (!this.props.showTargetFactorRecords) {
      colFilters = ["all"];
    } else {
      var factors = user.getContentFactors();
      var targetFactorId;
      if (user.getState().TargetFactor) {
        targetFactorId = user.getState().TargetFactor.FactorId;
        var fkey = Object.keys(factors)
        for (var i = 0; i < fkey.length; i++) {
          if (factors[fkey[i]].FactorId == targetFactorId) {
            for (var colIndex = 0; colIndex < factors[fkey[i]].Levels.length; colIndex++) {
              colFilters[colIndex] = targetFactorId + ":" + factors[fkey[i]].Levels[colIndex].FactorLevelId;
            }
            break;
          }
        }
      }
    }
    for (var colIndex = 0; colIndex < colFilters.length; colIndex++) {
      data[colIndex] = { label: '', rcount: [] };
      data[colIndex].label = colFilters[colIndex];
      for (var gindex = 0; gindex < pRecords.length; gindex++) {
        if (pRecords[gindex].Records[colFilters[colIndex]]) {
          data[colIndex].rcount[gindex] = pRecords[gindex].Records[colFilters[colIndex]].length;
        } else {
          data[colIndex].rcount[gindex] = 0;
        }
        // pRecords[gindex].Grade should be the same as performanceLabels[gindex].grade
      }
    }
    return  <div>
              <Graph user={user} colFilters={colFilters} data={data} allowToolbox={this.props.allowToolbox} yTitle="Performance" xTitle="All" yLabels={performanceLabels}/>
            </div>;
  }
});

function Diamond(props) {
  var h = rectSize / 2;

  var toggleToolbox = function() {
    props.onDiamondClick(props.col, props.grade, props.rIndex, true);
  };

  if (props.allowToolbox) {
    return <rect onClick={toggleToolbox} width={rectSize} height={rectSize}
      transform={`translate(${props.x},${props.y}) rotate(45) translate(-${h},-${h})`}
      style={{stroke: 'green', fill: 'green'}}/>
  }
  return <rect width={rectSize} height={rectSize}
    transform={`translate(${props.x},${props.y}) rotate(45) translate(-${h},-${h})`}
    style={{stroke: 'green', fill: 'green'}}/>
}

function Diamonds(props) {
  var size = rectSize * spacingFactor;
  var diamonds = [];
  var ePerRow = elementsPerRow;
  if (props.singleColumn) {
    ePerRow = elementsPerRow * 2;
  }
  for (let i = 0; i < props.count; i++) {
    let y = Math.floor(i / ePerRow);
    let x = i % ePerRow;
    diamonds.push(<Diamond x={x * size} y={y * size} allowToolbox={props.allowToolbox} onDiamondClick={props.onDiamondClick} col={props.col} grade={props.grade} rIndex={i} key={i}/>);
  }
  return <g transform={`translate(${size / 2}, -${size / 2}) scale(1, -1)`}>{diamonds}</g>;
}

function Column(props) {
  var totalHeight = props.rcount.length * rowHeight
  var colWidth = columnWidth;
  var ePerRow = elementsPerRow;
  if (props.singleColumn) {
    colWidth = columnWidth * 2;
    ePerRow = elementsPerRow * 2;
  }
  const x = (colWidth - ePerRow * rectSize * spacingFactor) / 2;
  return <g transform={`translate(${x},${0})`}>{
    props.rcount.map((count, i) =>
      <g transform={`translate(0,${(i+1) * rowHeight})`} key={i}>
        <Diamonds singleColumn={props.singleColumn} allowToolbox={props.allowToolbox} onDiamondClick={props.onDiamondClick} count={props.rcount[i]} col={props.col} grade={i}/>
      </g>)
    }</g>;
}

function Toolbox(props) {
  var rcount = props.data[props.toolboxCol].rcount
  var totalHeight = rcount.length * rowHeight
  var colWidth = columnWidth;
  var ePerRow = elementsPerRow;
  if (props.singleColumn) {
    colWidth = columnWidth * 2;
    ePerRow = elementsPerRow * 2;
  }
  var colX = (colWidth - ePerRow * rectSize * spacingFactor) / 2;

  var size = rectSize * spacingFactor;
  var x = props.toolboxIndex % ePerRow * size;
  var y = Math.floor(props.toolboxIndex / ePerRow) * size;
  var h = rectSize / 2;

  var allRecords = props.user.getAllPerformanceRecords()
  var record = allRecords[props.toolboxGrade].Records[props.colFilters[props.toolboxCol]][props.toolboxIndex];

  var toggleToolbox = function(){
    props.onDiamondClick(props.toolboxCol, props.toolboxGrade, props.toolboxIndex, false)
  }

  return  <g transform={`translate(${props.toolboxCol * colWidth}, 0)`}>
            <g transform={`translate(${colX},0)`}>
              <g transform={`translate(0,${(props.toolboxGrade+1) * rowHeight})`}>
                <g transform={`translate(${size / 2}, -${size / 2}) scale(1, -1)`}>
                 <rect width={rectSize+2} height={rectSize+2}
                    transform={`translate(${x},${y}) rotate(45) translate(-${h},-${h})`}
                    style={{stroke: 'black', fill: 'darkgreen'}}/>
                </g>
                <g transform={`translate(${size / 2}, -${size / 2}) scale(1, -1)`}>
                  <rect onClick={toggleToolbox} width={toolBoxSizeWidth} height={toolBoxSizeHeight}
                    transform={`translate(${x},${y}) translate(${h*3},0)`}
                    style={{stroke: 'white', fill: 'lightgrey'}}/>
                  <g transform={`translate(${x + (toolBoxSizeWidth+h*3)/2}, ${y+columnLabelHeight}) scale(1, -1)`}>
                    {
                      Object.keys(record.FactorLevels).map((l, i) =>
                        <text x={0} y={-columnLabelHeight*(i+1)} textAnchor='middle' className='axis-label' key={i}>{record.FactorLevels[l].FactorName}: {record.FactorLevels[l].SelectedLevel}</text>)
                    }
                    <text x={0} y={-columnLabelHeight*(Object.keys(record.FactorLevels).length+2)} textAnchor='middle' className='axis-title' key={Object.keys(record.FactorLevels).length}>Record #{record.RecordNo} {record.RecordName}</text>
                    <text x={0} y={0} textAnchor='middle' className='axis-title' key={Object.keys(record.FactorLevels).length+1}>Performance: {record.Performance}</text>
                  </g>
                </g>
              </g>
            </g>
          </g>;
}

function XAxis(props) {
  var colWidth = columnWidth;
  if (props.singleColumn) {
    colWidth = columnWidth * 2;
  }
  return <g className='axis'>
    <line x1='0' y1='0' x2={colWidth * props.labels.length} y2='0' stroke='black' strokeWidth='1'/>
    {
      props.labels.map((l, i) =>
        <text key={i} x={(i + .5) * colWidth} y={columnLabelHeight}
          textAnchor='middle' className='axis-label'>{l}</text>)
    }
    <text x={colWidth * props.labels.length / 2} y={2 * columnLabelHeight} textAnchor='middle' className='axis-title'>{props.title}</text>
  </g>;
}

function YAxis(props) {
  // TODO: Why?
  const x = -columnLabelHeight / 2;
  return <g className='axis'>
    <line x1='0' y1='0' x2='0' y2={rowHeight * props.labels.length} stroke='black' strokeWidth='1'/>
    {
      props.labels.map((l, i) =>
        <g transform={`translate(${x},${(i + .5) * rowHeight}) rotate(-90)`} key={i}>
          <text textAnchor='middle' className='axis-label'>{l.label}</text>
        </g>)
    }
    <g transform={`translate(${3 * x},${rowHeight * props.labels.length / 2}) rotate(-90)`} >
      <text textAnchor='middle' className='axis-title'>{props.title}</text>
    </g>
  </g>;
}

var Graph = React.createClass({

  getInitialState: function() {
    return {mode: 0, showToolbox: false, toolboxCol: -1, toolboxGrade: -1, toolboxIndex: -1};
  },

  toggleToolbox: function(col, grade, index, show) {
    if (col == this.state.toolboxCol &&
      grade == this.state.toolboxGrade &&
      index == this.state.toolboxIndex) {
      this.state.showToolbox = false;
    } else {
      this.state.showToolbox = show;
    }
    if (!this.state.showToolbox) {
      this.state.toolboxCol = -1;
      this.state.toolboxGrade = -1;
      this.state.toolboxIndex = -1;
    } else {
      this.state.toolboxCol = col;
      this.state.toolboxGrade = grade;
      this.state.toolboxIndex = index;
    }
    this.setState(this.state);
  },

  render: function() {
    var props = this.props;
    var singleColumn = props.data.length > 1 ? false : true;
    var drawingAreaH, drawingAreaW; 
    var colWidth = columnWidth;
    if (singleColumn) {
      colWidth = columnWidth * 2;
    }
    const labels = props.data.map(v => v.label);
    const columns = props.data.map((v, i) => <g transform={`translate(${i * colWidth}, 0)`} key={i}>
      <Column singleColumn={singleColumn} onDiamondClick={this.toggleToolbox} allowToolbox={props.allowToolbox} rcount={v.rcount} col={i}/>
    </g>);

    const rowBackground = props.data[0].rcount.map((_, i) => {
      return <rect width={props.data.length * colWidth} height={rowHeight} y={i * rowHeight}
        fill={rowColors[i % rowColors.length]} key={i}/>
    });

    var toolbox
    if (this.state.showToolbox) {
      var key = "k"+this.state.toolboxCol+":"+this.state.toolboxGrade+":"+this.state.toolboxIndex;
      toolbox = <Toolbox user={props.user} colFilters={props.colFilters} singleColumn={singleColumn} onDiamondClick={this.toggleToolbox} toolboxCol={this.state.toolboxCol} toolboxGrade={this.state.toolboxGrade} toolboxIndex={this.state.toolboxIndex} data={props.data} key={key}/>
    }

    if (props.allowToolbox) {
      drawingAreaW = "100%";
      drawingAreaH = paddingBottom + props.yLabels.length * rowHeight + paddingTop;
    } else {
      drawingAreaW = paddingLeft + props.data.length * colWidth + paddingRight;
      drawingAreaH = paddingBottom + props.yLabels.length * rowHeight + paddingTop;
    }

    return <svg className='graph'style={{
        width: drawingAreaW,
        height: drawingAreaH,
      }}>
      <g transform={`translate(${paddingLeft}, ${paddingTop})`}>
        {rowBackground}
        {columns}
        <YAxis labels={props.yLabels} title={props.yTitle}/>
        <g transform={`translate(0,${props.data[0].rcount.length * rowHeight})`}>
          <XAxis labels={labels} title={props.xTitle} singleColumn={singleColumn}/>
        </g>
        {toolbox}
      </g>
    </svg>;
  },
});


var ChartSelectTargetFactor = React.createClass({

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
    var f = document.getElementById("chartactionForm");
    e = f.elements['chartactioninput'];
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
        return <ChartFactorPromptOption option={option} key={i}/>;
      });

    return   <form id="chartactionForm" onSubmit={this.handleSubmit} onChange={this.handleChange}>
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
                <button type="submit" disabled={!this.isEnabled()} key={"ChartSelectTargetFactor"}>Enter</button>
              </p>
              </form>;
  },
});

var ChartFactorPromptOption = React.createClass({

  render: function() {
    var option = this.props.option;
      return <tr><td><label>
              <input type="radio" name="chartactioninput" value={option.ResponseId}/><br/>{option.Text}</label></td></tr>;
  },
});

