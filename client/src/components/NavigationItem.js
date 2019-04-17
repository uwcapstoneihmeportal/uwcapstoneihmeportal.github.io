import React, { Component } from 'react';
import { NavItem, NavLink, } from 'reactstrap' 

const NavItemStyle = {


}

const NavLinkStyle = {
    color: 'black',
    fontSize: '16px'
}

class NavigationItem extends Component {
    render() {
        return (
            <NavItem style={{NavItemStyle}}>
                <NavLink style={NavLinkStyle} href="/home">{this.props.label}</NavLink>
            </NavItem>
        )
    }
}

export default NavigationItem
