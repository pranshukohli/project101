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
	  axios.get('/bakemenufull')
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
                {bakeOrders.map(item => (
                <li key={item.OrderId}>
                  <p>
			{item.OrderNumber}
			&nbsp;&nbsp;&nbsp;
			{item.OrderStatus}
			&nbsp;&nbsp;&nbsp;
			{item.DishName}
			&nbsp;&nbsp;&nbsp;
			{item.OrderQuantity}
			&nbsp;&nbsp;&nbsp;
			{item.DishType}
			&nbsp;&nbsp;&nbsp;
			{item.OrderPaymentMode}
			&nbsp;&nbsp;&nbsp;
			{item.OrderType}
			&nbsp;&nbsp;&nbsp;
                  </p> 
		</li>
                ))}
              </ul>
            </div>
          );
        }
	}
};


export default BakeMenu;
