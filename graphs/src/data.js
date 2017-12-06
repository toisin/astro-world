// @flow

import {csvParseRows as parse} from 'd3-dsv';
import invariant from './invariant.js';

// These are the fields we should make graphs for.
const codingFieldIndexes = {
  Planning: 12,
  Regulating: 13,
  'Monitoring & Evaluating': 14,
  'Epistemological Thinking': 15,
  'Cognitive COV': 16,
  'Cognitive Chart': 17,
  'Cognitive Prediction': 18,
  // 'Cognitive Select Team': 19,
  'Argumentation skill': 20,
};

const usernameIdx = 0;
const dateIdx = 9;
const taskIdIdx = 11;

export type RowData = {
  username: string,
  taskId: string,
  date: Date,
};

// These docIds can be found by opening the spreadsheets and looking at the URL.
const docIds = [
  // rm2g1
  '1WDc6jxw-Hj-NXWBM48KysEQfOi9lKyXdAbo6RFMfxjI',
  // rm2g2
  '1EcOGMY3l3eEtAkIsOzmC68ASVFv3Ih26j2AtW7qyZUA',
  // rm10g1
  '1M9hnIUyG21Tda40tL-8jwMOnrNWHxUYrl0p0_E7oc7c',
  // rm10g2
  '1TyRukwPf-OeRMJuGCo9MmX8Ogv4LIOo0uz_h888fA6s',
];

function getCsvUrl(docId) {
  return `https://docs.google.com/spreadsheets/d/${docId}/gviz/tq?tqx=out:csv`;
}

async function loadCsv(docId: string): Promise<RowData[]> {
  const url = getCsvUrl(docId);
  const s = await (await fetch(url)).text();
  return parse(s)
    .slice(1)
    .map(coalesce);
}

export type SectionData = {map: Map<string, RowData[]>, maxX: number, maxY: number};
export type Data = Map<string, SectionData>;

export default async function getData(): Promise<Data> {
  const rowsOfRows = await Promise.all(docIds.map(loadCsv));
  const rows = removeDuplicates([].concat(...rowsOfRows));
  return splitRowsByCodingFields(rows, Object.keys(codingFieldIndexes));
}

function isEmpty(row, index) {
  const s = row[index];
  return s === '' || s === '?';
}

function coalesce(row: string[]): RowData {
  const data = {
    username: row[usernameIdx],
    date: new Date(row[dateIdx]),
    taskId: row[taskIdIdx],
    ...codingFieldIndexes, // value updated below
  };

  for (const key in codingFieldIndexes) {
    if (isEmpty(row, codingFieldIndexes[key])) {
      data[key] = undefined;
    } else {
      data[key] = +row[codingFieldIndexes[key]];
    }
  }

  return data;
}

function removeDuplicates(rows: RowData[]): RowData[] {
  const seen = new Set();
  return rows.filter(row => {
    const key = `${row.username}${row.date.getTime()}${row.taskId}`;
    if (seen.has(key)) {
      return false;
    }
    seen.add(key);
    return true;
  });
}

function getSectionData(data: RowData[], codingFieldName: string): SectionData {
  const usernameToRows = new Map();
  let maxX = 0;
  let maxY = 0;
  for (const row of data) {
    if (row[codingFieldName] === undefined) {
      continue;
    }

    const {username} = row;
    let rows = usernameToRows.get(username);
    if (rows === undefined) {
      rows = [];
      usernameToRows.set(username, rows);
    }

    maxX = Math.max(maxX, rows.length);
    rows.push(row);

    invariant(typeof row[codingFieldName] === 'number');
    maxY = Math.max(maxY, row[codingFieldName]);
  }
  return {map: usernameToRows, maxX, maxY};
}

function splitRowsByCodingFields(rows: RowData[], codingFieldNames: string[]): Data {
  const codingMap = new Map();
  for (const fieldName of codingFieldNames) {
    codingMap.set(fieldName, getSectionData(rows, fieldName));
  }
  return codingMap;
}
