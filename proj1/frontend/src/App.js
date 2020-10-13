import React from 'react';
import './App.css';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

import MakeOrder from './routes/makeOrder/MakeOrder.jsx';
import BakeOrder from './routes/bakeOrder/BakeOrder.jsx';
import Header from './components/header/Header';

class App extends React.Component {
  render() {
    return (
      <BrowserRouter>
        <div className="app">
         <Header /> 
            <Switch>
             <Route path="/" component={MakeOrder} exact/>
             <Route path="/makeorder" component={MakeOrder}/>
             <Route path="/bakeorder" component={BakeOrder}/>
             <Route component={Error}/>
            </Switch>
        </div> 
      </BrowserRouter>
    );
  }
}

export default App;
