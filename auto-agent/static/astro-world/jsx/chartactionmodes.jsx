/** @jsx React.DOM */
"use strict"

// npm install -g react-tools
// jsx -w -x jsx public/js public/js

const rectSize = 4;
const spacingFactor = 2.2;
const rowHeight = 100;
const elementsPerRow = 5;
const columnWidth = 300;
const columnLabelHeight = 20;
const paddingLeft = 50;
const paddingBottom = 50;
const paddingTop = 20;
const paddingRight = 20;
const rowColors = ['transparent', 'rgba(0,0,0,0.05)'];
const data = [{
  label: 'Average',
  data: [5, 7, 3, 28, 1]
}, {
  label: 'Good',
  data: [15, 2, 13, 8, 11]
}];


var Chart = React.createClass({

  getInitialState: function() {
    return {mode: 0};
  },

  render: function() {
      var state = this.state;
      var user = this.props.user;
      var app = this.props.app;
      var yLabels= ['A', 'B', 'C', 'D', 'E'] ;

      return  <div>
                <Graph data={data} yTitle="Performance" xTitle="Fitness" yLabels={yLabels}/>
              </div>;
  }
});

function Diamond(props) {
  const h = rectSize / 2;
  return <rect width={rectSize} height={rectSize}
    transform={`translate(${props.x},${props.y}) rotate(45) translate(-${h},-${h})`}
    style={{stroke: 'green', fill: 'green'}}/>
}


function Diamonds(props) {
  const size = rectSize * spacingFactor;
  const diamonds = [];
  const rowCount = Math.ceil(props.count / elementsPerRow);
  for (let i = 0; i < props.count; i++) {
    let y = Math.floor(i / elementsPerRow);
    let x = i % elementsPerRow;
    diamonds.push(<Diamond x={x * size} y={y * size} key={i}/>);
  }
  return <g transform={`translate(${size / 2},${size / 2})`}>{diamonds}</g>;
}

function Column(props) {
  const x = (columnWidth - elementsPerRow * rectSize * spacingFactor) / 2;
  return <g transform={`translate(${x},${props.data.length * rowHeight}) scale(1,-1)`}>{
    props.data.map((count, i) =>
      <g transform={`translate(0,${i * rowHeight})`} key={i}>
        <Diamonds count={props.data[i]}/>
      </g>)
    }</g>;
}

function XAxis(props) {
  return <g className='axis'>
    <line x1='0' y1='0' x2={columnWidth * props.labels.length} y2='0' stroke='black' strokeWidth='1'/>
    {
      props.labels.map((l, i) =>
        <text key={i} x={(i + .5) * columnWidth} y={columnLabelHeight}
          textAnchor='middle' className='axis-label'>{l}</text>)
    }
    <text x={columnWidth * props.labels.length / 2} y={2 * columnLabelHeight} textAnchor='middle' className='axis-title'>{props.title}</text>
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
          <text textAnchor='middle' className='axis-label'>{l}</text>
        </g>)
    }
    <g transform={`translate(${3 * x},${rowHeight * props.labels.length / 2}) rotate(-90)`} >
      <text textAnchor='middle' className='axis-title'>{props.title}</text>
    </g>
  </g>;
}

function Graph(props) {
  const labels = props.data.map(v => v.label);
  const columns = props.data.map((v, i) => <g transform={`translate(${i * columnWidth}, 0)`} key={i}>
    <Column data={v.data}/>
  </g>);

  const rowBackground = props.data[0].data.map((_, i) => {
    return <rect width={props.data.length * columnWidth} height={rowHeight} y={i * rowHeight}
      fill={rowColors[i % rowColors.length]} key={i}/>
  });


  return <svg className='graph'style={{
      width: paddingLeft + props.data.length * columnWidth + paddingRight,
      height: paddingBottom + props.yLabels.length * rowHeight + paddingTop,
    }}>
    <g transform={`translate(${paddingLeft}, ${paddingTop})`}>
      {rowBackground}
      {columns}
      <YAxis labels={props.yLabels} title={props.yTitle}/>
      <g transform={`translate(0,${props.data[0].data.length * rowHeight})`}>
        <XAxis labels={labels} title={props.xTitle}/>
      </g>
    </g>
  </svg>;
}
