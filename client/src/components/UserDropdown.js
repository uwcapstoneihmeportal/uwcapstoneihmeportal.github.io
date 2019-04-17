import React, { Component } from 'react';

import { UncontrolledDropdown, DropdownToggle, DropdownMenu, DropdownItem } from 'reactstrap'

const DropdownToggleStyle = {
    backgroundColor: '#cbe2a0',
    border: 'none',
    color: 'black',
    marginRight: '10px',
    textAlign: 'right'
}

const BaseImageStyle = {
    height: '30px', 
    marginBottom: '5px',
    marginRight: '5px'
}

const DropdownImageStyle = {
    marginLeft: '10px'
}

const userImage = require('../images/user.png')
const dropdownImage = require('../images/dropdown.png')

class UserDropdown extends Component {
    render() {
        return (
            <UncontrolledDropdown style={{textAlign: 'right'}}>
                <DropdownToggle style={DropdownToggleStyle}>
                    <img src={userImage} alt="image of user" style={{...BaseImageStyle}}/>
                    Welcome, Sam
                    <img src={dropdownImage} alt="image of user" style={{...BaseImageStyle, ...DropdownImageStyle}}/>
                </DropdownToggle>
                <DropdownMenu>
                    <DropdownItem header style={{fontSize: '20px', color: '#2F4F4F', textAlign: 'right'}}>
                        Sam Johnson
                    </DropdownItem>
                    <DropdownItem>
                        View Profile
                    </DropdownItem>
                    <DropdownItem divider />
                    <DropdownItem style={{ color: 'red' }}>
                        Logout
                    </DropdownItem>
                </DropdownMenu>
            </UncontrolledDropdown>
        )
    }
}

export default UserDropdown
