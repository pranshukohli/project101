import React from "react";
import "./Header.scss";
import { NavLink } from 'react-router-dom';

const Header = () => (
       <div className="header">
          <NavLink className="navlink" to="/">Home</NavLink>
          <NavLink className="navlink" to="/makeorder">MakeOrder</NavLink>
          <NavLink className="navlink" to="/bakeorder">BakeOrder</NavLink>
       </div>
)

export default Header;
