import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./BakeMenu.scss";


class BakeMenu extends Component {
	constructor(props) { 
	  super(props); 
          this.state = {
	    error: null,
            isLoaded: false,
	    bakeOrders: [],
            menu: '',
	  }
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }


	fetchBakeMenu = () => {
	  axios.get('/bakemenubyorder')
	    .then(
	    (repos) => {
	      console.log(repos.data)
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


export default BakeMenu;
