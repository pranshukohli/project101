import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./Menu.scss";
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';


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
	  };
	}

	handleChange = (e) =>{
          this.setState({[e.target.name]: e.target.value});
        }


	fetchMenu = () => {
	  axios.get('/menu')
	    .then(
	    (repos) => {
		    console.log(repos.data)
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

	addToOrder = (name, dish_id, price, quantity) => {
	    var item =[ 
		     name,
		     dish_id,
		     price,
		     quantity
	    ];
	    this.props.parentCallback(item);
	}

	addToMenu = () => {
	  axios.post('/menu', {
	    "name": this.state.itemName,
	    "price": parseInt(this.state.itemPrice),
	    "description": this.state.itemDescription
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
	  let subset = items.slice(1);
	  return(
            <div className="menu">
              <p>Menu</p>
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
						onClick={() => this.addToOrder(
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
		  <p>
              Name:<input name="itemName" type="text" value={this.state.itemName} onChange={this.handleChange}/>
              </p>
		  <p>Price:<input name="itemPrice" type="number" value={this.state.itemPrice} onChange={this.handleChange}/>
              </p>
		  <p>Description:<input name="itemDescription" type="text" value={this.state.itemDescription} onChange={this.handleChange}/>
              </p><button onClick={this.addToMenu}>Add</button>
            </div>
          );
        }
	}
};


export default Menu;
