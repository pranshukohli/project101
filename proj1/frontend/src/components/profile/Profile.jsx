import React, { Component } from "react";
import axios from 'axios';
import { sendMsg } from "../../api";
import "./Profile.scss";

const baseBackendURL = "http://ec2-65-0-12-62.ap-south-1.compute.amazonaws.com:8080" 


class Profile extends Component {
	constructor(props) { 
		super(props); 
		this.state = {
			error: null,
			isLoaded: false,
			profileData: "",
		};
	}

	handleChange = (e) =>{
		this.setState({[e.target.name]: e.target.value});
	}

	fetchProfile = () => {
		axios.get(baseBackendURL + '/v1/profile',{withCredentials: true})
			.then(
				(repos) => {
					console.log("fetched profile data"+repos.data)
					this.setState({
						isLoaded: true,
						profileData: repos.data
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
	componentDidMount() {
		this.fetchProfile()
	}

	gotoAuth = (isAuth) => {
		if (isAuth) {
			window.location.replace(baseBackendURL+"/v1/auth/google")
		}else {
			window.location.replace(baseBackendURL+"/v1/logout/google")
		}
	}
	render() {
		const { error, isLoaded, profileData} = this.state;
		if (error) {
			return <div>Error: {error.message}</div>;
		} else if (!isLoaded) {
			return <div className="profile">Loading Profile</div>;
		} else {
			if(this.state.profileData == ""){
				return(
					<div className="profile">
						<button onClick={ () => this.gotoAuth(true)}>Login</button>
					</div>
				);
			}else{
				return(
					<div className="profile">
					{profileData}
							<img src=""/>	
							<button onClick={ () => this.gotoAuth(false)}>Logout</button>
					</div>
				);
			}
		}
	}
};


export default Profile;
