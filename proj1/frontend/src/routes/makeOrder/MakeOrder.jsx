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
			totalSum: 0,
			itemFromMenu: '',
		};
	}
 	callbackFunction = (childData) => {
		var items = this.state.newOrderItem;
		var isPresent = false;
		var total = 0;
		items.map((item, index) => {
			console.log("dd"+item);
			if(item[1] == childData[1]) {
				console.log("ddima"+item);
				item[3]=item[3]+1;
				isPresent = true;
			}
			total = total + item[2]*item[3];
		});
		console.log("ff");
		if(!isPresent){
			items.push(childData);
			total = total + childData[2]*childData[3];
		}
		console.log("ss" + total);
      		this.setState({
			newOrderItem: items,
			totalSum: parseInt(total),
		})
		console.log("dd" + this.state.totalSum);
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
				<CreateOrder newOrderItem={this.state.newOrderItem} totalSum={this.state.totalSum} />
			</div>
		);
	}
}
 
export default MakeOrder;
