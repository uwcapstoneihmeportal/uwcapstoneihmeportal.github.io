import React, { Component } from 'react';
import { Container, Row, Col } from 'reactstrap'
import { AuthForm, AuthButton } from './AuthForm'

const imagePath = require('../../images/placeholderLoginImage.jpg')
const imageStyle = {
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'contain',
    width: '100%',
    height: '100vh'
}

class SignInView extends Component {
    render() {
        return (
            <Container style={{ maxWidth: '100%' }}>
                <Row >
                    <Col sm="6" className='d-none d-sm-block' style={{ paddingLeft: '0' }}>
                        <img src={imagePath} alt="test"
                            style={imageStyle} />
                    </Col>
                    <Col sm="6" >
                        <form style={{marginTop: '150px'}}>
                            <AuthForm labelText="Email" />
                            <AuthForm labelText="Password" />
                            {<AuthButton labelText="Sign in" />}
                        </form>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default SignInView