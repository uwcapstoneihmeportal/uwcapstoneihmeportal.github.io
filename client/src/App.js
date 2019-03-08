import React from 'react';
import NavigationBar from './components/NavigationBar'
import LoginForm from './components/LoginForm'
import './App.css';

const App = () => (
  <div>
    <NavigationBar shouldShowNavItems={ true } />
    <LoginForm />
  </div>
)

export default App;
