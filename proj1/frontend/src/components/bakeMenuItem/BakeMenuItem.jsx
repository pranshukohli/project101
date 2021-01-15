import React, { Component } from "react";
import axios from 'axios';
import { connect, sendMsg } from "../../api";
import "./BakeMenuItem.scss";


class BakeMenuItem extends Component {
	constructor(props) { 
	  super(props); 
          this.state = {
	    error: null,
            isLoaded: false,
	    bakeOrders: [],
            menu: '',
	    order_number: props.match.params.ordernumber,
	  }
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }
        componentDidMount() {
                connect((msg) => {
                        this.fetchBakeMenu();
                });
                this.fetchBakeMenu();
        }


	fetchBakeMenu = () => {
	console.log(this.state.order_number)
	  axios.get('/bakemenu/'+this.state.order_number)
	    .then(
	    (repos) => {
	      this.setState({
	        isLoaded: true,
	        bakeOrders: repos.data
	      });
	    },
	    (error) => {
	      this.setState({
	        isLoaded: true,
	        error
	      });
	    }
	   )
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
              <p>BakeMenu</p>
              <ul>
                {bakeOrders.map(itemlist => (
			<li key={itemlist.OrderNumber}>
			<ul>	
				<div>
					<span>
						#{itemlist.OrderNumber}
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
			<button>Set Completed</button>
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


export default BakeMenuItem;
