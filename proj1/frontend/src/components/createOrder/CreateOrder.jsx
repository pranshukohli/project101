import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./CreateOrder.scss";
import { Redirect } from 'react-router-dom';
import empty_cart from '../../empty_cart.png';

const baseBackendURL="http://"+process.env.REACT_APP_BASE_BACKEND_URI+":"+process.env.REACT_APP_BASE_BACKEND_PORT;

axios.defaults.headers.post['Content-Type'] ='application/json;charset=utf-8';
axios.defaults.headers.post['Access-Control-Allow-Origin'] = '*';
axios.defaults.headers.post['Access-Control-Allow-Headers'] = '*';

class CreateOrder extends Component {
	constructor(props) {
		super(props);
		this.state = {
			error: null,
			isLoaded: false,
			orderList: [],
			orderNumber: '',
		}
	}

	ColoredLine = (color) => (
		<hr
			style={{
				color: color,
				backgroundColor: color,
				height: 1,
				margin: 0,
			}}
		/>
	)

	handleChange = (e) => {
		this.setState({
			[e.target.name]: e.target.value
		});
	}
	popUpMenu = () => {
		return(
			<div>{this.state.orderNumber}</div>
		)
	
	}

	alterOrder = (name, dish_id, price, quantity) => {
            var item =[
                     name,
                     dish_id,
                     price,
                     quantity
            ];
            this.props.parentCallback(item);
        }

	createNewOrder = async () => {
		var orderItems = this.props.newOrderItem;
		var d = new Date();
		var orderNumber = d.getTime();
		this.setState({orderNumber: orderNumber});
		var responses = [];
		for (var i=0;i<orderItems.length;i++){
			responses.push(
			axios.post(baseBackendURL + '/v1/order', {
				"dish_id": parseInt(orderItems[i][1]),
				"order_number": orderNumber,
				"quantity": parseInt(orderItems[i][3]),
				"status": "new_order",
			}))
		}
		await axios
		  .all(responses)
		  .then(
			axios.spread((...response) => {
			console.log(response)
		  	sendMsg(JSON.stringify({msg:"update_bakemenu_new",orderNumber:orderNumber}));
		  }))
		console.log("Done");
		//ADD MODEL for "ORDER PLACED proceed to :NEW ORDER or VIEW ORDER"
	}

	render() {
		var noi = this.props.newOrderItem;
		if (noi.length == 0){
			return(
			<div className="createOrderParent">
				<h3 className="top-sticky">
					Cart	
					{this.ColoredLine("red")}	
				</h3>
				<div className="starting">
					<img src = {empty_cart} />
				</div>
			</div>
			)
		}
		else {
		return(
			<div className="createOrderParent">
			<div className="createOrder">
				<h3 className="top-sticky">
					Cart	
					{this.ColoredLine("red")}
				</h3>
				{this.popUpMenu()}
				<ul className="menu-ul">
				   {this.props.newOrderItem.map(item=>(  
					<li className="menu-li" key={item[1]}>
						<ul>
							<li>
					  			<span className="">
								   	{item[0]}
							   	</span>
								<span className="float-right">
									Full&nbsp;
					   				<span className="vnv">
					   					&#9679;
					   				</span>
								</span>
						   	</li>
							<li>
					   <button className="" onClick={() => this.alterOrder(
						   item[0],
						   item[1],
						   item[2], -1)
					   }>
					   -
					   </button>
					   			&nbsp;&nbsp;&nbsp;		
					   			{item[3]}
					   			&nbsp;&nbsp;&nbsp;		
					   <button className="" onClick={() => this.alterOrder(
						   item[0],
						   item[1],
						   item[2], 1)
					   }>
					   +
					   </button>
					   			&nbsp;&nbsp;&nbsp;
					   			X
					   			&nbsp;&nbsp;&nbsp;
								&#x20B9;{item[2]}
					   			&nbsp;&nbsp;&nbsp;		
					   			=
					   			&nbsp;&nbsp;&nbsp;		
					   			&#x20B9;{item[2]*item[3]}
						   	</li>
						</ul> 
					</li>
				   ))}
				</ul>
			</div>
			<div className="fixed-bottom">	
			<h3>
					<button name="createOrder" onClick={this.createNewOrder}>
						Order
					</button>
					&nbsp;&nbsp;&nbsp;		
					Total: &#x20B9;{this.props.totalSum}
				</h3>
			</div>
		</div>
		)}
	}
};


export default CreateOrder;
