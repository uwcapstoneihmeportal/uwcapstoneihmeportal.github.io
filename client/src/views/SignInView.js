import React, { Component } from 'react';
import { Container, Row, Col } from 'reactstrap'
import SignInBanner from '../components/SignInBanner'
import CustomForm from '../components/CustomForm'
import AuthButton from '../components/AuthButton'
import LoadingIcon from '../components/LoadingIcon'
import { withRouter, Redirect } from 'react-router-dom'

const ihme_logo = require("../images/ihme_logo.png")

const H1Style = {
    textAlign: 'center',
    fontSize: '32px',
    fontWeight: 'bold'
}

const FormContainerStyle = {
    margin: 'auto', 
    position: 'absloute', 
    transform: 'translate(0%, 50%)'
}

class SignInView extends Component {
    constructor(props) {
        super(props);
        this.state = {
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
                <Row>
                    <Col sm="6" className='d-none d-sm-block' style={{ paddingLeft: '0' }}>
                        <SignInBanner />
                    </Col>
                    <Col xs="12" sm="6">
                        {<img src={ihme_logo} alt="IHME logo" className="d-sm-none d-xs-block" style={{ paddingTop: '10px', height: '80px' }} />}

                    
                        <div style={FormContainerStyle}>
                            <h1 style={H1Style}>Sign in</h1>
                            <form>
                                <CustomForm labelText="Email" imagePath={require("../images/green_user.png")}/>
                                <CustomForm labelText="Password" type="password" imagePath={require("../images/password.png")}/>
                            </form>
                            <div style={{ marginTop: '60px'}}>
                                <LoadingIcon loading={this.state.loading} />
                                <AuthButton labelText="Sign in" onClick={this.handleClick} />
                            </div>
                        </div>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default withRouter(SignInView)
