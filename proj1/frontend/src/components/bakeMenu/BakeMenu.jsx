import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./BakeMenu.scss";
const baseBackendURL="http://"+process.env.REACT_APP_BASE_BACKEND_URI+":"+process.env.REACT_APP_BASE_BACKEND_PORT; 

console.log("Backend Server IP:" + baseBackendURL);

class BakeMenu extends Component {
	constructor(props) { 
	  super(props); 
          this.state = {
	    error: null,
            isLoaded: false,
	    bakeOrders: [],
            menu: '',
	    bakeOrderInProgress:  0,
	    showAlert: false,
	    war: "",
	    unak: [],
	  }
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }


	fetchBakeMenu = (msg,orderNumber) => {
	   if(msg == "update_bakemenu_new") {
	   	console.log("new order")
		this.setState({
			showAlert: true,
			war: "!!Got New Order!! "+orderNumber,
		})
	   }else if(msg == "update_bakemenu_com") {
	   	console.log("update order")
		this.setState({
			showAlert: true,
			war: "Last Completed Order!! "+orderNumber,
		})
	   }
	  axios.get(baseBackendURL + '/v1/bakemenubyorder')
	    .then(
	    (repos) => {
		    if(msg == "update_bakemenu_new"){
		    	repos.data[0].type="<span class='badge badge-secondary'>New</span>"
			var _unak = this.state.unak
			_unak.push(orderNumber)
			    this.setState({
			    	unak:_unak,
			    })
			    console.log("S"+_unak)
			    console.log("a"+this.state.unak)
		    }
		    console.log(repos.data[0])
	      this.setState({
	        isLoaded: true,
	        bakeOrders: repos.data,
		bakeOrderInProgress: repos.data.length,
	      });
	    },
	    (error) => {
	      this.setState({
	        isLoaded: true,
	        error
	      });
	    }
	   );
	}

	closeAlert = () =>{
	   	console.log("release new order")
		this.setState({
			showAlert: false,
			war: "",
		})
	
	}

	alertNewOrder = () => {
		const show = this.state.showAlert;
		var war = this.state.war;
		if (show) {
			return (
				<div className="alert alert-warning alert-dismissible show" role="alert">
				<strong>{war}</strong>
				<button type="button" className="close" data-dismiss="alert" 
					aria-label="Close" onClick = {() => this.closeAlert()}>
					<span aria-hidden="true">&times;</span>
				</button>
				</div>
			);
		}
	}	

	viewOrder = async(order_number) =>{
		alert(order_number)
	}

	updateOrder = async(order_number) => {
                axios
                        .put(baseBackendURL + "/v1/bakemenuupdate/" + order_number)
                        .then(
				(response) => {
					console.log(response);
                                	sendMsg(JSON.stringify(
						{msg:"update_bakemenu_com",orderNumber:order_number}
					));
				},
				(error) => {
					console.log(error);
				}
			);
                console.log("Done");
        }
	
	getDate = (utcSec) => {
		var d = new Date(0); // The 0 there is the key, which sets the date to the epoch
		d.setUTCSeconds(utcSec/1000);
		return d.toLocaleString();
	}
	render() {
          const { error, isLoaded, bakeOrders, unak} = this.state;
          if (error) {
            return <div>Error: {error.message}</div>;
          } else if (!isLoaded) {
            return <div className="bakemenu">Loading BakeMenu</div>;
          } else {

	  return(

		  <div className="bakemenu">
		
		  {this.alertNewOrder()}
	      <p>BakeMenu</p>
		  <p>{process.env.REACT_APP_BASE_BACKEND_URL}</p>
	      <p>Pending orders:  {this.state.bakeOrderInProgress}</p>
	      <p>Unaknowledged orders: {unak.map(ulist=>(
	      					<span>{ulist}, </span>
	      				))}</p>
              <ul>
                {bakeOrders.map(itemlist => (
			<li key={itemlist.OrderNumber} 
			    onClick={() => this.viewOrder(itemlist.OrderNumber)}>
			<ul>	
				<div>
					<span>
						#{itemlist.OrderNumber}
						@{this.getDate(itemlist.OrderNumber)}
						&nbsp;&nbsp;
						<span dangerouslySetInnerHTML={
							{__html:itemlist.type}
						}></span>
					</span>
					<span className="float-right">
						{itemlist.OrderList[0].OrderStatus}
					</span>
                		</div>
				{itemlist.OrderList.map(item => (
                			<li key={item.OrderId}>
			                  <div>
						{item.DishName}
						&nbsp;&nbsp;&nbsp;
						{item.OrderQuantity}
						&nbsp;&nbsp;&nbsp;
						{item.OrderType}
						&nbsp;&nbsp;&nbsp;
						<span className="status">
							{item.OrderStatus}
						</span>
						&nbsp;&nbsp;&nbsp;
			                  </div> 
					</li>
		                ))}
			<div>
			<button className="Button"
				onClick={() => this.updateOrder(itemlist.OrderNumber)}	
			>Set Completed</button>
			</div>
			</ul>
			</li>
		))}
              </ul>
            </div>
          );
        }
	}
};


export default BakeMenu;
