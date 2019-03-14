import React, { Component } from 'react';
import { Button, TabContent, TabPane, Nav, NavItem, NavLink, Container, Row, Col } from 'reactstrap'
import classnames from 'classnames';

const jobIcon = require('../images/job.png')
const locationIcon = require('../images/location.png')

const info = {
    textAlign: 'center',
    width: '100%',
  }
const style = { display: 'block', justifyContent: 'center', alignItems: 'center', width: '..', height: '..'}

class ProfileView extends Component {
    constructor(props) {
        super(props);
    
        this.toggle = this.toggle.bind(this);
        this.state = {
          activeTab: '2'
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
                <div style={{backgroundColor:'#D8D8D8'}}>
                <Container style={{paddingTop: '12vh'}}>
                    <Row>
                            <h1 style={{textAlign: 'center', width: '100%'}}>Sam Johnson</h1>
                    </Row>
                </Container>
                <Container style={{paddingTop: '20px', paddingBottom:'10vh'}}>
                    <Row style={{paddingLeft:'20%', paddingRight: '20%'}}>
                        <Col>
                            <img src={jobIcon} alt="test" style={{height:'50px', horizontalAlign: 'middle'}}/>
                            <span style={{paddingLeft:'10px', fontSize:'2vh'}}>Software Engineer</span>
                        </Col>
                        <Col >
                            <img src={locationIcon} alt="test" style={{height:'50px', horizontalAlign: 'middle'}}/>
                            <span style={{paddingLeft:'10px', fontSize:'2vh'}}>United State, WA</span>
                        </Col>
                    </Row>
                </Container>
            </div>

            
            <div style={{paddingTop:'15px'}}>
                <Nav tabs style={{paddingLeft:'60px'}}>
                    <NavItem>
                        <NavLink
                        className={classnames({ active: this.state.activeTab === '1' })}
                        onClick={() => { this.toggle('1');}}
                        >
                        Personal Profile
                        </NavLink>
                    </NavItem>
                    <NavItem style={{paddingRight:'59%'}}>
                        <NavLink
                        className={classnames({ active: this.state.activeTab === '2' })}
                        onClick={() => { this.toggle('2'); }}
                        >
                        Publication
                        </NavLink>
                    </NavItem>

                    <Button variant="link " onClick={() => { this.toggle('3'); }}>Edit Profile</Button>

                </Nav>
                <TabContent activeTab={this.state.activeTab} style={{paddingLeft:'60px', paddingRight:'60px',backgroundColor:'#F6F6F6'}}>
                    <TabPane tabId="1" >
                        <Row>
                        <Col sm="12">
                            <div style={{paddingLeft:'27px'}}>
                                <Row style={{paddingTop:'30px'}}>
                                    <h4 style={{paddingBottom:'15px'}}>Profile</h4>
                                    <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id 
                                        facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et 
                                        malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante, 
                                        tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                </Row>
                            </div>
                        </Col>
                        </Row>
                    </TabPane>
                    <TabPane tabId="2">
                        <div style={{paddingLeft:'27px'}}>
                            <Row style={{paddingTop:'30px'}}>
                                <h4 style={{paddingBottom:'15px'}}>Publication 1 Title</h4>
                                <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id 
                                    facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et 
                                    malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante, 
                                    tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                <text style={{textAlign: 'right', width: '100%', paddingTop:'15px'}} >Sept. 2018</text>
                            </Row>
                            <Row >
                                <h4 style={{paddingBottom:'15px'}}>Publication 2 Title</h4>
                                <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id 
                                    facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et 
                                    malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante, 
                                    tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                <text style={{textAlign: 'right', width: '100%', paddingTop:'15px'}} >Sept. 2018</text>
                            </Row>
                            <Row >
                                <h4 style={{paddingBottom:'15px'}}>Publication 3 Title</h4>
                                <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id 
                                    facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et 
                                    malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante, 
                                    tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                <text style={{textAlign: 'right', width: '100%', paddingTop:'15px'}} >Sept. 2018</text>
                            </Row>
                            <Row >
                                <h4 style={{paddingBottom:'15px'}}>Publication 4 Title</h4>
                                <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id 
                                    facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et 
                                    malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante, 
                                    tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                <text style={{textAlign: 'right', width: '100%', paddingTop:'15px'}} >Sept. 2018</text>
                            </Row>
                        </div>
                    </TabPane>
                    <TabPane tabId="3" >
                        <Row>
                        <Col sm="12">
                            <div style={{paddingLeft:'27px'}}>
                                <Row style={{paddingTop:'30px'}}>
                                    <h4 style={{paddingBottom:'15px'}}>Edit Profile</h4>
                                    <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id 
                                        facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et 
                                        malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante, 
                                        tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                </Row>
                            </div>
                        </Col>
                        </Row>
                    </TabPane>
                </TabContent>
            </div>
        </div>
        );
    }
}

export default ProfileView