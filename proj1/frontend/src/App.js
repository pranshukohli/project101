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

  componentDidMount() {
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


  send = () => {
    sendMsg(this.state.soe);
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
            {this.state.items.map(item => (
              <li key={item.dish_id}>
                <p>{item.name}</p> {item.description}
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
