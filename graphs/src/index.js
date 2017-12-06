// @flow

import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import getData from './data.js';
import invariant from './invariant.js';

async function main() {
  const dataP = getData();
  const root = document.querySelector('#root');
  invariant(root);
  ReactDOM.render(<Loading />, root);
  const data = await dataP;
  ReactDOM.render(<App data={data} />, root);
}

main();

function Loading() {
  return <h3 className="loading">Loading data...</h3>;
}
