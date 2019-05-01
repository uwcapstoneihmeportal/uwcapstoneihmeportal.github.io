import React, { Component } from 'react';
import { Button, Container, Row, Col } from 'reactstrap'
import CardContainer from '../components/CardContainer'

class ProfileInformationView extends Component {
    render() {
        return (
            <Container>
                <Row>
                    <Col>
                        <CardContainer title="Contact Details"/>
                    </Col>
                    <Col>
                        <CardContainer title="Contact Information"/>
                    </Col>
                </Row>
            </Container>
        )
    }
}

export default ProfileInformationView
