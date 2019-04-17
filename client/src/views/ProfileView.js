import React, { Component } from 'react';
import { Button, TabContent, TabPane, Nav, NavItem, NavLink, Container, Row, Col } from 'reactstrap'
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

                {/* <Nav tabs>
                    <NavItem>

                    </NavItem>
                    <NavItem>

                    </NavItem>
                </Nav> */}


                <Nav tabs style={{ paddingLeft: '60px', alignItems: 'center' }}>
                    <NavItem>
                        <NavLink
                            className={classnames({ active: this.state.activeTab === '1' })}
                            onClick={() => { this.toggle('1'); }}
                        >
                            Personal Profile
                    </NavLink>
                    </NavItem>
                    <NavItem style={{marginRight: '600px'}}>
                        <NavLink
                            className={classnames({ active: this.state.activeTab === '2' })}
                            onClick={() => { this.toggle('2'); }}
                        >
                            Publications
                    </NavLink>
                    </NavItem>

                    <Button style={EditProfileButtonStyle} variant="link">
                        Edit Profile
                    </Button>
                </Nav>
                <TabContent activeTab={this.state.activeTab} style={{ paddingLeft: '60px', paddingRight: '60px', backgroundColor: '#F6F6F6' }}>
                    <TabPane tabId="1" >
                        <Row>
                            <Col sm="12">
                                
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