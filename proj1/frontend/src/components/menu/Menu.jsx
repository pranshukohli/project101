import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./Menu.scss";
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';
import image1 from '../../view_image.png';

const baseBackendURL = "http://192.168.3.120:8080" 


class Menu extends Component {
	constructor(props) { 
	  super(props); 
          this.state = {
	    error: null,
            isLoaded: false,
	    items: [],
            itemName: '',
            itemPrice: 0,
            itemDescription: '',
	    view: "list",
	  };
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }


	fetchMenu = () => {
	  axios.get(baseBackendURL + '/v1/menu')
	    .then(
	    (repos) => {
		    console.log("fetched data"+repos.data)
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

	alterOrder = (name, dish_id, price, quantity) => {
	    var item =[ 
		     name,
		     dish_id,
		     price,
		     quantity
	    ];
	    this.props.parentCallback(item);
	}

	addToMenu = () => {
	  axios.post(baseBackendURL + '/v1/menu', {
	    "name": this.state.itemName,
	    "price": parseInt(this.state.itemPrice),
	    "description": this.state.itemDescription
	  })
	  .then(function (response) {
	     console.log(response);
	  })
	  sendMsg("update Menu");
	}

	MenuWithCards = (items) => {
	  let subset = items.slice(1);
	  return(
          <Container className="menu-container">
              <Row>
                {subset.map(item => (
                <Col md={4} sm={12} key={item.dish_id}>
			<div className="card-menu">
				<div className="card-menu-header">
					<span className="card-top-left">
						&#x20B9;{item.price}
					</span>
					<span className="float-right">
						{item.name}
					</span>
				</div>
				<div className="card-menu-body">
					{item.description}
				</div>
				<div className="card-menu-footer">
					<button className="button"
						onClick={() => this.alterOrder(
								item.name,
								item.dish_id,
								item.price, 1)
						}>
						+
					</button>
				</div>
			</div>
                </Col>
                ))}
              </Row>
	      </Container>
	  )
	}

	MenuWithList = (items) => {
	  let subset = items.slice(1);
	  return(
          <Container className="list-menu-container">
              <ul>
                {subset.map(item => (
                <li key={item.dish_id}>
			<ul>
				<li className="">
					<span className="float-right">
						&#x20B9;{item.price}
					</span>
					<span className="">
						{item.name}
					</span>
				</li>
				<li className="">
					{item.description}
				</li>
				<li className="">
					<button className="" onClick={() => this.alterOrder(
								item.name,
								item.dish_id,
								item.price, -1)
					}>
						-	
					</button>
					&nbsp;&nbsp;&nbsp;1&nbsp;&nbsp;&nbsp;
					<button className="" onClick={() => this.alterOrder(
								item.name,
								item.dish_id,
								item.price, 1)
					}>
						+
					</button>
				</li>
			</ul>
                </li>
                ))}
              </ul>
	      </Container>
	  )
	}

	AddMenu = () => {
	  return(
	  <div>
		<p> Name: <input name="itemName" type="text" 
			   value={this.state.itemName} onChange={this.handleChange}/>
		</p>
		<p> Price: <input name="itemPrice" type="number"
			    value={this.state.itemPrice} onChange={this.handleChange}/>
		</p>
		<p> Description: <input name="itemDescription" type="text"
				  value={this.state.itemDescription} onChange={this.handleChange}/>
		</p>
		<button onClick={this.addToMenu}>Add</button>
	  </div>)
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


	render() {
          const { error, isLoaded, items} = this.state;
          if (error) {
            return <div>Error: {error.message}</div>;
          } else if (!isLoaded) {
            return <div className="menu">Loading Menu</div>;
          } else {
		  if(this.state.view == "cards"){
		 
	  return(
	    <div className="menu">
              <h3 className="top-sticky">
		  Menu
		  <span className="view-button" 
		  	onClick={() => this.setState({view: "list"})}>
		  	<img src={image1}/>	
		  </span>
		  {this.ColoredLine("red")}
	      </h3>
	      {this.MenuWithCards(items)}
            </div>
          );
		  }
		  else if(this.state.view == "list"){
	  return(
	    <div className="menu">
              <h3 className="top-sticky">
		  Menu
		  <span className="view-button" 
		  	onClick={() => this.setState({view: "cards"})}>
		  	<img src={image1}/>	
		  </span>
		  {this.ColoredLine("red")}
	      </h3>
	      {this.MenuWithList(items)}
            </div>
          );
		  }
        }
	}
};


export default Menu;
