import React, { Component } from 'react';

const logo = require('../images/logo.png')

const BannerStyle = {
    background: 'linear-gradient(#cbe2a0, #67b56b, #26a146)',
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'contain',
    height: '100vh',
    paddingLeft: '50px',
    paddingTop: '50px',
    width: '100%'
}

const ImageStyle = {
    height: '200px',
    width: '200px'
}

const PortalTextStyle = {
    
    color: 'white',
    fontSize: '30px',
    marginTop: '20px'
}

const IHMETextStyle = {
    color: '#505050',
    fontSize: '35px', 
}

const SloganTextStyle = {
    bottom: '25px',
    left: '50px',
    color: 'white',
    fontSize: '20px',
    position: 'absolute'
}

class SignInBanner extends Component {
    render() {
        return (
            <div style={BannerStyle}>
                <img src={logo} alt="IHME Logo" style={ImageStyle}/>
                <p style={PortalTextStyle}>
                    Welcome to the <br />
                    GBD Collaborator Portal
                </p>                
                <p style={SloganTextStyle}>
                    <span style={IHMETextStyle}>
                        IHME
                    </span>
                    <br/>
                    Measuring what matters
                </p>
            </div>
        )
    }
}

export default SignInBanner
