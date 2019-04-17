import React, { Component } from 'react';
import { Container, Row, Col } from 'reactstrap'

import UserDropdown from '../components/UserDropdown'

const backgroundColor = "#cbe2a0"
const TopNavigationStyle = {
    backgroundColor: backgroundColor,
    paddingTop: '20px'
}

class TopNavigationBar extends Component {
    render() {
        return (
            <div style={TopNavigationStyle}>
                <Container style={{ maxWidth: '100%', margin: '0', padding: '0' }}>
                    <Row>
                        <Col sm="12">
                            <UserDropdown />
                        </Col>
                    </Row>
                </Container>
            </div>

            


            



        )
    }
}

export default TopNavigationBar
