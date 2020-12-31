import React from 'react';
import './App.css';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

import MakeOrder from './routes/makeOrder/MakeOrder.jsx';
import BakeOrder from './routes/bakeOrder/BakeOrder.jsx';
import Header from './components/header/Header';
import BakeMenuItem from './components/bakeMenuItem/BakeMenuItem';

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
						<Route path="/bakeorderitem/:ordernumber" component={BakeMenuItem}/>
						<Route component={Error}/>
					</Switch>
				</div>
			</BrowserRouter>
		);
	}
}

export default App;
