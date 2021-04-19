import React, { Component } from 'react';
import BakeMenu from '../../components/bakeMenu/BakeMenu';
import "./BakeOrder.scss";
import { connect } from "../../api";

class BakeOrder extends Component {
	constructor(props){
		super(props);
		this.state = {
			newOrderItem: [],
			itemFromBakeMenu: '',
			isConnected: "Not Connected",
		};
	}
 	callbackFunction = (childData) => {
		var items = this.state.newOrderItem;
		items.push(childData);
		console.log(childData);
		console.log(items);
      		this.setState({newOrderItem: items})
	}	
	componentDidMount() {
		connect((msg) => {
			console.log(msg)
			if(msg == "database_in_sync") {
				this.setDatabaseSync(true);
			} else if(msg == "database_out_of_sync") {
				this.setDatabaseSync(false);
			} else if (JSON.parse(msg.data).body != null){
				if(JSON.parse(JSON.parse(msg.data).body).msg != null) {
					this.refs.child.fetchBakeMenu(
						JSON.parse(JSON.parse(msg.data).body).msg,
						JSON.parse(JSON.parse(msg.data).body).orderNumber
					);
				} 
			}
		});
		this.refs.child.fetchBakeMenu(false, false);
	}

	setDatabaseSync(isConnected) {
                if(isConnected){
      			this.setState({isConnected: "In Sync"})
                }else{
      			this.setState({isConnected: "Not In Sync, Refresh Page!!"})
                }
        }

	render() {
		return (
			<div className="bakeorder">
				<h1>Bake Order</h1>
				<p id="database_conn">Database: {this.state.isConnected}</p>
				<BakeMenu ref="child" parentCallback = {this.callbackFunction}/>
			</div>
		);
	}
}
 
export default BakeOrder;
