import React from 'react'
import PropTypes from 'prop-types'
import logo from '../images/IHMElogo.png';

const NavBarStyle = {
    backgroundColor: "#26a146",
    padding: '1rem',
    position: 'fixed',
    top: '0',
    left: '0',
    width: '100%',
}

const NavItemStyle = (shouldShowItems) => {
    return ({
        color: 'white',
        display: shouldShowItems ? 'table' : 'none',
        float: 'right',
        fontSize: '1.2em',
        marginRight: '2rem',
        marginTop: '1rem',
        verticalAlign: 'middle',
        textAlign: 'center',
        textDecoration: 'none'
    });
}

const LogoStyle = {
    width: '140px',
    // height: 'auto'
}

const NavigationBar = ({ shouldShowNavItems }) => (
    <nav style={NavBarStyle}>
        <a href=""><img src={logo} alt="" style={LogoStyle}/></a>
        <a href="" style={NavItemStyle(shouldShowNavItems)}>My Profile</a>
        <a href="" style={NavItemStyle(shouldShowNavItems)}>Home</a>
        
    </nav>
)

NavigationBar.propTypes = {
    shouldShowItems: PropTypes.bool.isRequired
}

export default NavigationBar