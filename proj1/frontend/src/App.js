import React from 'react';
import './App.css';
import { HashRouter, Route, Switch } from 'react-router-dom';

import login from './routes/login/login.js';
import MakeOrder from './routes/makeOrder/MakeOrder.jsx';
import BakeOrder from './routes/bakeOrder/BakeOrder.jsx';
import Header from './components/header/Header';
import BakeMenuItem from './components/bakeMenuItem/BakeMenuItem';
import Profile from './components/profile/Profile.jsx';

class App extends React.Component {
	render() {
		return (
			<HashRouter>
				<div className="app">
					<Header /> 
					<Profile />
					<Switch>
						<Route path="/" component={MakeOrder} exact/>
						<Route path="/makeorder" component={MakeOrder}/>
						<Route path="/bakeorder" component={BakeOrder}/>
						<Route path="/bakeorderitem/:ordernumber" component={BakeMenuItem}/>
						<Route path="/login" component={login}/>
						<Route component={Error}/>
					</Switch>
				</div>
			</HashRouter>
		);
	}
}

export default App;
