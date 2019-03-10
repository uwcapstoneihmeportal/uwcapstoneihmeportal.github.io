import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch, Redirect } from 'react-router-dom';

// import Views
import LogInView from './views/LogInView'
import HomeView from './views/HomeView'
import { Provider } from 'react-redux'

// Style related imports
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Provider>
          <Router>
            <Switch>
              <Route path='/' component={LogInView} />
              <Route path='/home' component={HomeView} />
              <Redirect to="/" />
            </Switch>
          </Router>
        </Provider>
      </div>
    );
  }
}

export default App