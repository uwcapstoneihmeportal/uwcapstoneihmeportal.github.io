import React, { Component } from 'react';
import { Container, Row, Col } from 'reactstrap'
import { AuthForm, AuthButton } from '../components/AuthForm'

import { withRouter, Redirect } from 'react-router-dom'

const imagePath = require('../images/placeholderTwo.jpg')

const imageStyle = {
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'contain',
    width: '100%',
    height: '100vh'
}

const H1Style = {
    marginTop: '100px',
    textAlign: 'center',
    fontSize: '32px',
    fontWeight: 'bold'
}

class SignInView extends Component {
    constructor(props) {
        super(props);
        this.state = {
            navigate: false,
            referrer: null,
        };
    }

    handleClick = () => {
        this.setState({referrer: '/profile'});
    }

    render() {
        const {referrer} = this.state;
        if (referrer) return <Redirect to={referrer} />;

        return (
            <Container style={{ maxWidth: '100%' }}>
                <Row >
                    <Col sm="6" className='d-none d-sm-block' style={{ paddingLeft: '0' }}>
                        <img src={imagePath} alt="test"
                            style={imageStyle} />
                    </Col>
                    <Col sm="6" >
                        <h1 style={H1Style}>Sign in</h1>
                        <form style={{ }}>
                            <AuthForm labelText="Email" />
                            <AuthForm type="password" labelText="Password" />
                            <AuthButton labelText="Sign in" onClick={this.handleClick}/>
                        </form>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default withRouter(SignInView)