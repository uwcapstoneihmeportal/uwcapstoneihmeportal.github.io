import React, { Component } from 'react';
import { Container, Row, Col } from 'reactstrap'
import { AuthForm, AuthButton } from '../components/AuthForm'
import { withRouter, Redirect } from 'react-router-dom'

import { DotLoader } from 'react-spinners'


const imagePath = require('../images/placeholderTwo.jpg')

const imageStyle = {
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'contain',
    width: '100%',
    height: '100vh'
}

const H1Style = {
    marginTop: '80px',
    textAlign: 'center',
    fontSize: '32px',
    fontWeight: 'bold'
}

const loadingStyle = {
    display: 'block',
    margin: '0 auto',
    marginTop: '10px',
    marginBottom: '0'
}

class SignInView extends Component {
    constructor(props) {
        super(props);
        this.state = {
            navigate: false,
            referrer: null,
            loading: false
        };
    }

    handleClick = () => {
        this.setState({ loading: true })

        setTimeout(() => {
            this.setState({ referrer: '/home' })
            this.setState({ loading: false })
        }, 2 * 1000)
    }

    render() {
        const { referrer } = this.state;
        if (referrer) return <Redirect to={referrer} />;

        return (
            <Container style={{ maxWidth: '100%' }}>
                <Row >
                    <Col sm="6" className='d-none d-sm-block' style={{ paddingLeft: '0' }}>
                        <img src={imagePath} alt="test"
                            style={imageStyle} />
                    </Col>
                    <Col sm="6" >
                        <img src={require("../images/ihme_logo.png")} alt="IHME logo" style={{ paddingTop: '10px', height: '80px' }} />

                        <h1 style={H1Style}>Sign in</h1>
                        <form style={{}}>
                            <AuthForm labelText="Email" />
                            <AuthForm type="password" labelText="Password" />
                            <DotLoader loading={this.state.loading} color="#26a146" css={loadingStyle} />
                            <AuthButton labelText="Sign in" onClick={this.handleClick} />
                        </form>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default withRouter(SignInView)