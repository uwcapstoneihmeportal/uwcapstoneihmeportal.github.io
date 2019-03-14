import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch, Redirect } from 'react-router-dom';

// import Views
import SignInView from './views/SignInView'
// import HomeView from './views/HomeView'
import ProfileView from './views/ProfileView'

// Style related imports
import './styles/App.css'

class App extends Component {
  render() {
    return (
      <div className="App">
        <Router>
          <Switch>
            <Route exact path='/' component={SignInView} />
            <Route path='/profile' component={ProfileView} />
            <Redirect to='/' component={SignInView} />
          </Switch>
        </Router>
      </div>
    );
  }
}

export default App