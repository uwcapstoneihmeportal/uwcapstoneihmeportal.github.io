import React, { Component } from 'react'; //import React Component
import { FormGroup, Label, Input, Button, FormFeedback, Row, Col } from 'reactstrap';
import { Link } from 'react-router-dom';

class AuthForm extends Component {
    constructor(props) {
        super(props);
        this.state = { value: "" }
    }
    handleChange(event) {
        this.props.handleChange(event);
        this.setState({ value: event.target.value });
    }
    render() {
        return (
            <FormGroup>
                <Label htmlFor={this.props.id}>{this.props.name}</Label>
                <Input id={this.props.type}
                    type={this.props.type}
                    name={this.props.id}
                    valid={this.props.valid}
                    onChange={(event) => this.handleChange(event)}
                />
                {this.props.valid !== undefined &&
                    this.props.errors.map((error) => {
                        return <FormFeedback key={error}>{error}</FormFeedback>;
                    })}
            </FormGroup>
        );
    }
}

export default class SignInForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            email: undefined,
            password: undefined,
        };
    }

    handleChange(event) {
        let newState = {};
        newState[event.target.name] = event.target.value;
        this.setState(newState);
    }

    /**
     * A helper function to validate a value based on an object of validations
     * Second parameter has format e.g., 
     *    {required: true, minLength: 5, email: true}
     * (for required field, with min length of 5, and valid email)
     */
    validate(value, validations) {
        let errors = [];

        if (value !== undefined) { //check validations
            //handle required
            if (validations.required && value === '') {
                errors.push('Required field.');
            }

            //handle minLength
            if (validations.minLength && value.length < validations.minLength) {
                errors.push(`Must be at least ${validations.minLength} characters.`);
            }

            //handle email type
            if (validations.email) {
                //pattern comparison from w3c
                //https://www.w3.org/TR/html-markup/input.email.html#input.email.attrs.value.single
                let valid = /^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(value)
                if (!valid) {
                    errors.push('Not an email address.')
                }
            }
            return errors; //report the errors
        }
        return undefined; //no errors defined (because no value defined!)
    }

    handleSignIn(event) {
        event.preventDefault(); //don't submit
        this.props.signInCallback(this.state.email, this.state.password);
    }

    render() {
        let emailError = this.validate(this.state.email, { required: true, email: true });
        let passwordError = this.validate(this.state.password, { required: true, minLength: 6 });
        let validations = [];
        validations.push((emailError === undefined) ? undefined : (emailError.length === 0));
        validations.push((passwordError === undefined) ? undefined : (passwordError.length === 0));
        return (
            <div className="container mt-5">
                <Row>
                    <Col sm={{ size: 10, offset: 1 }} md={{size: 4, offset:4}}>
                        <form>
                            <AuthForm type="email" id="email" name="Email" valid={validations[0]} errors={emailError} handleChange={(event) => this.handleChange(event)} />
                            <AuthForm type="password" id="password" name="Password" valid={validations[1]} errors={passwordError} handleChange={(event) => this.handleChange(event)} />
                            <AuthButton type="signin" isValid={validations} click={(event) => this.handleSignIn(event)} />
                        </form>
                        {/* <p>Don't have an account already? Then <Link to="/join">Sign Up</Link></p> */}
                    </Col>
                </Row>
            </div>
        );
    }
}

class AuthButton extends Component {
    render() {
        if (this.props.type === "signup") {
            return (
                <FormGroup>
                    <Button className="mr-2" color="primary" onClick={(e) => this.props.click(e)} disabled={!this.props.isValid[0] && !this.props.isValid[1] && !this.props.isValid[2]}>
                        Sign-up
                    </Button>
                </FormGroup>
            );
        } else {
            return (
                <FormGroup>
                    <Button color="primary" onClick={(e) => this.props.click(e)} disabled={!this.props.isValid[0] && !this.props.isValid[1]}>
                        Sign-in
                    </Button>
                </FormGroup>
            );
        }
    }
}

export class SignUpForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            email: undefined,
            handle: undefined,
            password: undefined,
            passwordVerify: undefined
        };
    }

    handleChange(event) {
        let newState = {};
        newState[event.target.name] = event.target.value;
        this.setState(newState);
    }

    /**
     * A helper function to validate a value based on an object of validations
     * Second parameter has format e.g., 
     *    {required: true, minLength: 5, email: true}
     * (for required field, with min length of 5, and valid email)
     */
    validate(value, validations) {
        let errors = [];

        if (value !== undefined) { //check validations
            //handle required
            if (validations.required && value === '') {
                errors.push('Required field.');
            }

            //handle minLength
            if (validations.minLength && value.length < validations.minLength) {
                errors.push(`Must be at least ${validations.minLength} characters.`);
            }

            if (validations.passwordVerify) {
                if (this.state.password !== value) {
                    errors.push('Password does not match.');
                }
            }

            //handle email type
            if (validations.email) {
                //pattern comparison from w3c
                //https://www.w3.org/TR/html-markup/input.email.html#input.email.attrs.value.single
                let valid = /^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(value)
                if (!valid) {
                    errors.push('Not an email address.')
                }
            }
            return errors; //report the errors
        }
        return undefined; //no errors defined (because no value defined!)
    }

    handleSignUp(event) {
        event.preventDefault(); //don't submit
        this.props.signUpCallback(this.state.email, this.state.password, this.state.handle);
    }

    render() {
        let emailError = this.validate(this.state.email, { required: true, email: true });
        let passwordError = this.validate(this.state.password, { required: true, minLength: 6 });
        let passwordConfirmError = this.validate(this.state.passwordVerify, { required: true, passwordVerify: true });
        let handleError = this.validate(this.state.handle, { required: true });
        let validations = [];
        validations.push((emailError === undefined) ? undefined : (emailError.length === 0));
        validations.push((passwordError === undefined) ? undefined : (passwordError.length === 0));
        validations.push((passwordConfirmError === undefined) ? undefined : (passwordConfirmError.length === 0));
        validations.push((handleError === undefined) ? undefined : (handleError.length === 0));
        return (
            <div className="container mt-5">
                <Row>
                    <Col sm={{ size: 10, offset: 1 }} md={{size: 4, offset:4}}>
                        <form>
                            <AuthForm type="email" id="email" name="Email" valid={validations[0]} errors={emailError} handleChange={(event) => this.handleChange(event)} />
                            <AuthForm type="password" id="password" name="Password" valid={validations[1]} errors={passwordError} handleChange={(event) => this.handleChange(event)} />
                            <AuthForm type="password" id="passwordVerify" name="Confirm Password" valid={validations[2]} errors={passwordConfirmError} handleChange={(event) => this.handleChange(event)} />
                            <AuthForm type="handle" id="handle" name="Nick Name" valid={validations[3]} errors={handleError} handleChange={(event) => this.handleChange(event)} />
                            <p>If you want to add a profile picture, create an account at <a href="https://en.gravatar.com/">https://en.gravatar.com/</a> with the Email used above!</p>
                            <AuthButton type="signup" isValid={validations} click={(event) => this.handleSignUp(event)} />
                        </form>
                        {/* <p>Already have an account? Then <Link to="/login">Sign In</Link></p> */}
                    </Col>
                </Row>
            </div>
        );
    }
}

