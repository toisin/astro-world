// @flow

import React from 'react';
import './App.css';
import Graph from './graph.js';
import type {Data} from './data.js';
import {Link, Route, BrowserRouter as Router} from 'react-router-dom';

function App({data}: {data: Data}) {
  return (
    <Router>
      <div>
        <Route exact path="/" render={() => <Graphs data={data} />} />
        <Route
          path="/users/:username"
          render={({match}) => <UserGraph data={data} match={match} />}
        />
      </div>
    </Router>
  );
}

function UserGraph({data, match}) {
  return <Graphs data={data} username={match.params.username} />;
}

function Graphs({data, username: user = ''}: {data: Data, username?: string}) {
  const sections = [];

  for (const [codingFieldName, {map, maxX, maxY}] of data) {
    const cs = [];
    for (const [username, data] of map) {
      if (user === '' || user === username) {
        cs.push(
          <div key={`${username}-${codingFieldName}`} className={`${codingFieldName} coding-field`}>
            {user === '' ? (
              <h4>
                <Link to={`/users/${username}`}>{username}</Link>
              </h4>
            ) : (
              ''
            )}
            <Graph
              data={data}
              width={400}
              height={120}
              prop={codingFieldName}
              maxX={maxX}
              maxY={maxY}
            />
          </div>
        );
      }
    }

    sections.push(
      <div className={`${codingFieldName}s coding-fields`} key={codingFieldName}>
        <h3>{codingFieldName}</h3>
        <div className="coding-fields-container">{cs}</div>
      </div>
    );
  }

  return sections;
}

export default App;
