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
                            <h4 style={{ marginTop: '30px', marginBottom: '15px' }}>Welcome to the GBD Collaborator Portal!</h4>
                            <p>
                                GBD collaborators are crucial for the production, analysis, and improvement of the Global Burden of Disease. 
                                We’ve put together this page of information and resources for members of the Global Burden of Disease collaborative network.
                            </p>     
                        </Col>
                    </Row>
                    <Row>
                        <Col sm="12">
                            <h4 style={{ marginTop: '30px', marginBottom: '15px' }}>Resources, Tools, &amp; Instructional Videos</h4>
                            <p>


                            </p>
                        </Col>
                    </Row>
                </Container>
            </div>
        )
    }
}

export default HomeView
