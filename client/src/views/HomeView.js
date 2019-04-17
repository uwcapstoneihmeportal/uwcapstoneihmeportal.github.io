import React, { Component } from 'react';
import { Container, Row, Col } from 'reactstrap'

import NavigationBar from '../components/NavigationBar'
import MissionStatement from '../components/MissionStatement'

class HomeView extends Component {
    render() {
        return (
            <div>
                <NavigationBar />

                <div style={{
                    backgroundColor: '#ADD8E6',
                    height: '300px'
                }}>
                </div>

                <MissionStatement />

                <Container>
                    <Row>
                        <Col sm="12">
                            <div style={{}}>
                                <Row style={{ paddingTop: '30px' }}>
                                    <h4 style={{ paddingBottom: '15px' }}>News Article</h4>
                                    <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id
                    facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et
                    malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante,
                    tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                </Row>
                            </div>
                        </Col>
                    </Row>
                    <Row>
                        <Col sm="12">
                            <div style={{}}>
                                <Row style={{ paddingTop: '30px' }}>
                                    <h4 style={{ paddingBottom: '15px' }}>Another News Article</h4>
                                    <text>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quis scelerisque lectus, id
                    facilisis massa. Integer nec euismod ex. Donec vel ante et lectus mollis varius. Interdum et
                    malesuada fames ac ante ipsum primis in faucibus. Aliquam vulputate auctor elementum. Etiam nibh ante,
                    tincidunt id arcu sit amet, pharetra tincidunt nisi. Praesent in pellentesque massa, eget pellentesque leo.</text>
                                </Row>
                            </div>
                        </Col>
                    </Row>
                </Container>
            </div>
        )
    }
}

export default HomeView