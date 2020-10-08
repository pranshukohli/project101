import React, { Component } from "react";
import axios from 'axios';
import { connect, sendMsg } from "../../api";
import "./Menu.scss";


class Menu extends Component {
	constructor(props) { 
	  super(props); 
          this.state = {
	    error: null,
            isLoaded: false,
	    items: [],
            menu: ''
	  }
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }

        componentDidMount() {
          connect((msg) => {
    	    this.fetchMenu();
          });
        this.fetchMenu();
        }

	fetchMenu = () => {
	  axios.get('/menu')
	    .then(
	    (repos) => {
	      this.setState({
	        isLoaded: true,
	        items: repos.data
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

	addToOrder = (dish_id) => {
	    console.log("Add it" + dish_id);
	}

	addToMenu = () => {
	  axios.post('/menu', {
	    "name": this.state.menu,
	    "description": this.state.menu
	  })
	  .then(function (response) {
	     console.log(response);
	  })
	  sendMsg("update Menu");
	}

	render() {
          const { error, isLoaded, items} = this.state;
          if (error) {
            return <div>Error: {error.message}</div>;
          } else if (!isLoaded) {
            return <div className="menu">Loading Menu</div>;
          } else {

	  return(
            <div className="menu">
              <p>Menu</p>
              <ul>
                {items.map(item => (
                <li key={item.dish_id}>
                  <p>
                    {item.name}&nbsp;&nbsp;&nbsp; 
                    <button onClick={() => this.addToOrder(item.dish_id)}>+</button>
                  </p> {item.description}
                </li>
                ))}
              </ul>
              <input name="menu" type="text" value={this.state.menu} onChange={this.handleChange}/>
              <button onClick={this.addToMenu}>Add</button>
            </div>
          );
        }
	}
};


export default Menu;
