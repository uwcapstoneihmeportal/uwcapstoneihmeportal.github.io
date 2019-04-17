import React, { Component } from 'react';
import { Collapse, Navbar, NavbarToggler, NavbarBrand, Nav } from 'reactstrap';

import NavigationItem from '../components/NavigationItem'
import TopNavigationBar from '../components/TopNavigationBar'

const backgroundColor = "#cbe2a0"

const NavBarStyle = {
    backgroundColor: backgroundColor,
    paddingTop: '10px'
}

const brandImage = require('../images/ihme_logo.png')

class NavigationBar extends Component {
    render() {
        const NavBrandHeight = '80px'

        return (
            <div>
                <TopNavigationBar />

                <Navbar style={NavBarStyle} expand="md">
                    <NavbarBrand href="/home">
                        <img src={brandImage} alt="IHME logo" style={{ height: NavBrandHeight }} />
                    </NavbarBrand>

                    <Collapse isOpen={false} navbar>
                        <Nav className="nav-fill w-100" navbar>
                            <NavigationItem label="Home" />
                            <NavigationItem label="Research" />
                            <NavigationItem label="News &amp; Events" />
                            <NavigationItem label="Projects" />
                            <NavigationItem label="Get Involved" />
                            <NavigationItem label="About" />
                        </Nav>
                    </Collapse>
                </Navbar>
            </div>
        )
    }
}

export default NavigationBar