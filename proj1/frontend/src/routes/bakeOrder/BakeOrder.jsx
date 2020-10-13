import React, { Component } from 'react';
import BakeMenu from '../../components/bakeMenu/BakeMenu';
import "./BakeOrder.scss";
import { connect } from "../../api";

class BakeOrder extends Component {
	constructor(props){
		super(props);
		this.state = {
			newOrderItem: [],
			itemFromBakeMenu: '',
		};
	}
 	callbackFunction = (childData) => {
		var items = this.state.newOrderItem;
		items.push(childData);
		console.log(childData);
		console.log(items);
      		this.setState({newOrderItem: items})
	}	
	componentDidMount() {
		connect((msg) => {
			this.refs.child.fetchBakeMenu();
		});
		this.refs.child.fetchBakeMenu();
	}
	render() {
		return (
			<div className="bakeorder">
				<h1>Bake Order</h1>
				<BakeMenu ref="child" parentCallback = {this.callbackFunction}/>
			</div>
		);
	}
}
 
export default BakeOrder;
