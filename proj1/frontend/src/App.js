import React from 'react';
import './App.css';
import axios from 'axios';
import { connect, sendMsg } from "./api";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      items: [],
      soe: '',
      menu: ''
    };
  }

  handleChange = (e) =>{
    this.setState({[e.target.name]: e.target.value});
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
  addToMenu = () => {
    axios.post('/menu', {
      "name": this.state.menu,
      "description": this.state.menu
    })
    .then(function (response) {
       console.log("dd");
       console.log(response);
    })
    sendMsg("update Menu"); 
  }
  componentDidMount() {
    connect((msg) => {
    	this.fetchMenu();
    });
    this.fetchMenu();
  }

  add = (dish_id) => {
    console.log("Add it" + dish_id);
  }

  render() {
    const { error, isLoaded, items} = this.state;
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else {
      return (
        <div>
          <h1>StarManager</h1>
          <ul>
            {items.map(item => (
              <li key={item.dish_id}>
                <p>
		    {item.name}&nbsp;&nbsp;&nbsp; 
	            <button onClick={() => this.add(item.dish_id)}>+</button>
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
}

export default App;
