import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./BakeMenu.scss";
const baseBackendURL = "http://ec2-65-0-12-62.ap-south-1.compute.amazonaws.com:8080"

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
	  }
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }


	fetchBakeMenu = (isNew, isUpdate) => {
	  axios.get(baseBackendURL + '/v1/bakemenubyorder')
	    .then(
	    (repos) => {
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
	   if(isNew) {
	   	console.log("new order")
		this.setState({
			showAlert: true,
			war: "Got New Order!!",
		})
	   }else if(isUpdate) {
	   	console.log("update order")
		this.setState({
			showAlert: true,
			war: "Completed Order!!",
		})
	   }
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

	updateOrder = async(order_number) => {
                axios
                        .put(baseBackendURL + "/v1/bakemenuupdate/" + order_number)
                        .then(
				(response) => {
					console.log(response);
                                	sendMsg("update_bakemenu_com");
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
          const { error, isLoaded, bakeOrders} = this.state;
          if (error) {
            return <div>Error: {error.message}</div>;
          } else if (!isLoaded) {
            return <div className="bakemenu">Loading BakeMenu</div>;
          } else {

	  return(
            <div className="bakemenu">
		  {this.alertNewOrder()}
	      <p>BakeMenu</p>
	      <p>Pending orders:  {this.state.bakeOrderInProgress}</p>
              <ul>
                {bakeOrders.map(itemlist => (
			<li key={itemlist.OrderNumber}>
			<ul>	
				<div>
					<span>
						#{itemlist.OrderNumber}
						@{this.getDate(itemlist.OrderNumber)}
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
