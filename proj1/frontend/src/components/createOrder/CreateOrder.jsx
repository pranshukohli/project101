import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./CreateOrder.scss";


class CreateOrder extends Component {
	constructor(props) {
		super(props);
		var d = new Date();
		this.state = {
			error: null,
			isLoaded: false,
			orderList: [],
			orderNumber: d.getTime(),
		}
	}

	handleChange = (e) => {
		this.setState({
			[e.target.name]: e.target.value
		});
	}

	createNewOrder = () => {
		console.log(this.props.newOrderItem[0][1]);
		var orderItems = this.props.newOrderItem;
		for (var i=0;i<orderItems.length;i++){
			axios.post('/order', {
				"dish_id": parseInt(orderItems[i][1]),
				"order_number": this.state.orderNumber,
				"quantity": parseInt(orderItems[i][3]), 
			}).then(function (response) {
				console.log(response);
			})
			sendMsg("update Menu");
		}
	}

	render() {
		return(
			<div className="createorder">
				<h3>New Order #{this.state.orderNumber}</h3>
				<ul>
				   {this.props.newOrderItem.map(item=>(  
					<li key={item[1]}>
					   {item[0]}
					   &nbsp;&nbsp;&nbsp;
					   &#x20B9;{item[2]}
					   &nbsp;&nbsp;&nbsp;
					   {item[3]}
					</li>
				   ))}
				</ul>
				<h3>Total: &#x20B9;{this.props.totalSum}</h3>
				<button name="createOrder" onClick={this.createNewOrder}>
					Create Order
				</button>
			</div>
		);
	}
};


export default CreateOrder;
