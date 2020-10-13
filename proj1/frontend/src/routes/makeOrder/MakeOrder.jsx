import React, { Component } from 'react';
//import Category from '../../components/category/Category';
import Menu from '../../components/menu/Menu';
import CreateOrder from '../../components/createOrder/CreateOrder';
import "./MakeOrder.scss";
import { connect } from "../../api";

class MakeOrder extends Component {
	constructor(props){
		super(props);
		this.state = {
			newOrderItem: [],
			itemFromMenu: '',
		};
	}
 	callbackFunction = (childData) => {
		var items = this.state.newOrderItem;
		var isPresent = false;
		items.map((item, index) => {
			console.log("dd"+item);
			if(item[1] == childData[1]) {
				console.log("ddima"+item);
				item[3]=item[3]+1;
				isPresent = true;
			}
		});
		console.log("ff");
		if(!isPresent)
			items.push(childData);
      		this.setState({newOrderItem: items})
	}	
	componentDidMount() {
		connect((msg) => {
			this.refs.child.fetchMenu();
		});
		this.refs.child.fetchMenu();
	}
	render() {
		return (
			<div className="makeorder">
				<h1>Make Order</h1>
				<Menu ref="child" parentCallback = {this.callbackFunction}/>
				<CreateOrder newOrderItem={this.state.newOrderItem} />
			</div>
		);
	}
}
 
export default MakeOrder;
