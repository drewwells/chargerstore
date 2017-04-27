import React, { Component } from 'react';
import './App.css';
import Status from './status.jsx';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <h2>Welcome to Volt Charge</h2>
        </div>
				<div className="app-body">
					<Status />
				</div>
      </div>
    );
  }
}

export default App;
