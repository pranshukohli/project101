import React from 'react';
import './App.css';
import { connect, sendMsg } from "./api";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      items: [],
      soe: ''
    };
    connect();
  }

  handleChange = (e) =>{
    this.setState({soe: e.target.value});
  }

  fetchMenu = () => {
    fetch("/menu")
//      .then(res => res.text())          // convert to plain text
//      .then(text => console.log(text)) 
      .then(res => res.json())
      .then(
        (result) => {
         this.setState({
            isLoaded: true,
            items: result
          });
        },
        // Note: it's important to handle errors here
        // instead of a catch() block so that we don't swallow
        // exceptions from actual bugs in components.
        (error) => {
          this.setState({
            isLoaded: true,
            error
          });
        }
      )
  }
  componentDidMount() {
    this.fetchMenu();
  }

  send = () => {
    sendMsg(this.state.soe);
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
	<p>SOE: {this.state.soe}</p>
	<input type="text" value={this.state.soe} onChange={this.handleChange}/>
	<button onClick={this.send}>Hit</button>
        </div>
      );
    }
  }
}

export default App;
