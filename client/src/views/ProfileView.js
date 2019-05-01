import React, { Component } from 'react';
import { Button, TabContent, TabPane, Nav, NavItem, NavLink, Container, Row, Col } from 'reactstrap'
import ProfileInformationView from './ProfileInformationView'


import classnames from 'classnames';

import ProfileBanner from '../components/ProfileBanner'

const EditProfileButtonStyle = {
    borderRadius: '25px', 
    border: '1px solid #26a146', 
    backgroundColor: '#26a146', 
    color: 'white',
    marginBottom: '10px', 
}

class ProfileView extends Component {
    constructor(props) {
        super(props);

        this.toggle = this.toggle.bind(this);
        this.state = {
            activeTab: '1'
        };
    }

    toggle(tab) {
        if (this.state.activeTab !== tab) {
            this.setState({
                activeTab: tab
            });
        }
    }

    render() {
        return (
            <div>
                <ProfileBanner />
                <Nav tabs style={{ paddingLeft: '50px', alignItems: 'center' }}>
                    <NavItem>
                        <NavLink className={classnames({ active: this.state.activeTab === '1' })} onClick={() => { this.toggle('1'); }}>
                            Personal Profile
                        </NavLink>
                    </NavItem>
                    <NavItem>
                        <NavLink className={classnames({ active: this.state.activeTab === '2' })} onClick={() => { this.toggle('2'); }}>
                            Publications
                        </NavLink>
                    </NavItem>
                </Nav>
                <TabContent activeTab={this.state.activeTab} style={{ backgroundColor: '#F6F6F6' }}>
                    <TabPane tabId="1" >
                        <Row>
                            <Col sm="12">
                                <ProfileInformationView />
                            </Col>
                        </Row>
                    </TabPane>
                    <TabPane tabId="2">
                        

                    </TabPane>
                </TabContent>
            </div>
        );
    }
}

export default ProfileView