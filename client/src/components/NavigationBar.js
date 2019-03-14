import React, { Component } from 'react';
import {
    Collapse,
    Navbar,
    NavbarToggler,
    NavbarBrand,
    Nav,
    NavItem,
    NavLink,
    UncontrolledDropdown,
    DropdownToggle,
    DropdownMenu,
    DropdownItem
} from 'reactstrap';

const NavBarStyle = {
    backgroundColor: 'gray',
}

const TopNavigationStyle = {
    backgroundColor: 'gray',
    marginRight: '0'
}

const NavStyle = {
    fontSize: '11pt',
    textAlign: 'right',
    justifyContent: 'right',
    float: 'right',
    alignItems: 'right' 
}

const brandImage = require('../images/ihme_logo.png')

class NavigationBar extends Component {
    render() {
        return (
            <div>
                <TopNavigationItems />
                <Navbar style={NavBarStyle} light expand="md">
                    <NavbarBrand href="/">
                        <img src={brandImage} alt="IHME logo" style={{ height: '70px' }} />
                    </NavbarBrand>
                    <Collapse isOpen={false} navbar style={{ float: 'right' }}>
                        <Nav style={NavStyle} className="" navbar>
                            <NavItem style={{marginLeft: '15px'}}>
                                <NavLink style={{color: 'white'}} href="/home">Home</NavLink>
                            </NavItem>
                            <NavItem style={{marginLeft: '15px'}}>
                                <NavLink style={{color: 'white'}} href="">Research</NavLink>
                            </NavItem>
                            <NavItem style={{marginLeft: '15px'}}>
                                <NavLink style={{color: 'white'}} href="">News & Events</NavLink>
                            </NavItem>
                            <NavItem style={{marginLeft: '15px'}}>
                                <NavLink style={{color: 'white'}} href="">Projects</NavLink>
                            </NavItem>
                            <NavItem style={{marginLeft: '15px'}}>
                                <NavLink style={{color: 'white'}} href="">Get Involved</NavLink>
                            </NavItem>
                            <NavItem style={{marginLeft: '15px'}}>
                                <NavLink style={{color: 'white'}} href="">About</NavLink>
                            </NavItem>
                        </Nav>
                    </Collapse>
                </Navbar>
            </div>
        )
    }
}

class TopNavigationItems extends Component {
    render() {
        return (
            <div style={TopNavigationStyle}>
                <UserDropdown />
            </div>
        )
    }
}

class UserDropdown extends Component {

    render() {
        return (
            <UncontrolledDropdown style={{textAlign: 'right', marginRight: '10px', color: 'white'}}>
                <DropdownToggle nav caret style={{color: 'white', paddingTop: '20px'}}>
                    Welcome, Sam
                </DropdownToggle>
                <DropdownMenu right style={{backgroundColor: 'green'}}>
                    <DropdownItem>
                    <NavLink href='/profile' style={{color: 'black'}}>
                            View Profile
                    </NavLink>
                    </DropdownItem>
                    <DropdownItem divider />
                    <DropdownItem>
                        <NavLink href='/signin' style={{color: 'black'}}>
                            Logout
                        </NavLink>    
                </DropdownItem>
                </DropdownMenu>
            </UncontrolledDropdown>
        )
    }
}




export default NavigationBar