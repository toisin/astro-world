// @flow

import React from 'react';
import './App.css';
import {scaleLinear} from 'd3-scale';
import type {RowData} from './data.js';

type Props = {
  data: RowData[],
  width: number,
  height: number,
  prop: string,
  maxX: number,
  maxY: number,
};

export default function Graph({data, width, height, prop, maxX, maxY}: Props) {
  const strokeWidth = 1;
  const margin = {top: 10, right: 10, bottom: 20, left: 20};
  const w = width - margin.left - margin.right;
  const h = height - margin.top - margin.bottom;

  const x = scaleLinear()
    .range([0, w])
    .domain([0, maxX]);
  const y = scaleLinear()
    .range([0, h])
    .domain([maxY, 0]);

  return (
    <svg width={width} height={height}>
      <g transform={`translate(${margin.left},${margin.top})`}>
        <XAxis x={x} y={y} max={maxX} />
        <YAxis x={x} y={y} max={maxY} />

        <path d={computePath(data, x, y, prop)} strokeWidth={strokeWidth} />
        {renderTaskIdLines(data, x, y, maxY)}
      </g>
    </svg>
  );
}

function computePath(data, x, y, prop) {
  let p = '';
  let index = 0;
  for (const d of data) {
    // const {index} = d;

    switch (index) {
      case 0:
        p += 'M ';
        break;
      case 1:
        p += 'L ';
        break;
      default:
        break;
    }
    p += `${x(index)},${y(d[prop])} `;
    index++;
  }
  return p;
}

function renderTaskIdLines(
  data: RowData[],
  x: number => number,
  y: number => number,
  maxY: number
) {
  let lastTaskId = '';
  const lines = [];
  let index = 0;
  for (const {taskId} of data) {
    if (taskId !== lastTaskId) {
      lines.push(
        <g className="task-id" key={index}>
          <line className="task-id-line" x1={x(index)} y1={y(0)} x2={x(index)} y2={y(maxY)} />
          <text
            className="task-id-label"
            transform={`translate(${x(index) + 3},${y(maxY)}) rotate(90)`}
          >
            {taskId}
          </text>
        </g>
      );
      lastTaskId = taskId;
    }
    index++;
  }
  return lines;
}

type AxisProps = {x: number => number, y: number => number, max: number};

function XAxis({x, y, max}: AxisProps) {
  const ticks = [];
  const labels = [];
  for (let i = 0; i <= max; i++) {
    ticks.push(
      <line
        className="x-axis-tick axis-label"
        key={`t${i}`}
        x1={x(i)}
        y1={y(0)}
        x2={x(i)}
        y2={y(0) + 3}
      />
    );
    labels.push(
      <text
        className="x-axis-label axis-label"
        key={`l${i}`}
        transform={`translate(${x(i)},${y(0)})`}
        textAnchor="middle"
        alignmentBaseline="hanging"
        y="5"
      >
        {i}
      </text>
    );
  }
  return (
    <g className="x-axis axis">
      <line x1={x(0)} y1={y(0)} x2={x(max)} y2={y(0)} />
      {ticks}
      {labels}
    </g>
  );
}

function YAxis({x, y, max}: AxisProps) {
  const ticks = [];
  const labels = [];
  for (let i = 0; i <= max; i++) {
    ticks.push(
      <line
        className="y-axis-tick axis-label"
        key={`t${i}`}
        x1={x(0)}
        y1={y(i)}
        x2={x(0) - 3}
        y2={y(i)}
      />
    );
    labels.push(
      <text
        className="y-axis-label axis-label"
        key={`l${i}`}
        transform={`translate(${x(0)},${y(i)})`}
        textAnchor="end"
        alignmentBaseline="middle"
        x="-5"
      >
        {i}
      </text>
    );
  }
  return (
    <g className="y-axis axis">
      <line x1={x(0)} y1={y(0)} x2={x(0)} y2={y(max)} />
      {ticks}
      {labels}
    </g>
  );
}
