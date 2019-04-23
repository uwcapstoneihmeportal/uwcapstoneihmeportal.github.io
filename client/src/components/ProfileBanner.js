import React, { Component } from 'react';

const jobIcon = require('../images/job.png')
const locationIcon = require('../images/location.png')

const BannerStyle = {
    background: 'linear-gradient(#cbe2a0, #26a146)',
    marginBottom: '15px',
    paddingTop: '12vh',
    paddingBottom: '12vh',
    textAlign: 'center'
}

class ProfileBanner extends Component {
    render() {
        return (
            <div style={BannerStyle}>
                <h1>Sam Johnson</h1>
                <div style={{ display: 'inline-block', marginRight: '10px' }}>
                    <img src={jobIcon} alt="job image icon" style={{ height: '30px', marginBottom: '5px'}} />
                    <span style={{ marginLeft: '5px', fontSize: '20px' }}>
                        Health Specialist
                    </span>
                </div>
                <div style={{ display: 'inline-block', marginLeft: '10px' }}>
                    <img src={locationIcon} alt="location image icon" style={{ height: '40px'}} />
                    <span style={{ marginLeft: '-5px', fontSize: '20px' }}>
                        United States, WA
                    </span>
                </div>
            </div>
        )
    }
}

export default ProfileBanner
