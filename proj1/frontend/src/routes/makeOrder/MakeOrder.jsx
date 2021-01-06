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
		var deleteIndex = -1;
		var total = 0;
		items.map((item, index) => {
			if(item[1] == childData[1]) {
				item[3]=item[3]+childData[3];
				isPresent = true;

			}
			if (item[3] >=0) { 
				total = total + item[2]*item[3];
			}
			if(item[3] <= 0){
				deleteIndex = index;	
			}
		});
		if(!isPresent && childData[3]>0){
			items.push(childData);
			total = total + childData[2]*childData[3];
		}

		if(deleteIndex >=0){
			items.splice(deleteIndex,1);
		}
      		this.setState({
			newOrderItem: items,
			totalSum: parseInt(total),
		})
		console.log("dd" + this.state.totalSum);
	}	
	componentDidMount() {
		connect((msg) => {
			if(msg.data != null)
				if(JSON.parse(msg.data).body == "update_menu")
					this.refs.child.fetchMenu();
		});
		this.refs.child.fetchMenu();
	}
	render() {
		return (
			<div className="makeorder">
			<h1>&nbsp;</h1>
				<Menu ref="child" parentCallback = {this.callbackFunction}/>
				<CreateOrder newOrderItem={this.state.newOrderItem} totalSum={this.state.totalSum} />
			</div>
		);
	}
}
 
export default MakeOrder;
