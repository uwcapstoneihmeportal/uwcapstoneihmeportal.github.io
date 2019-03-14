import React, { Component } from 'react';
import { FormGroup, Label, Input, Button, FormFeedback } from 'reactstrap';

const ButtonStyle = {
    backgroundColor: '#26a146',
    borderRadius: '25px',
    margin: '0 auto',
    marginTop: '100px',
    width: '85%'
}

const InputStyle = {
    borderRadius: '25px',
    margin: '0 auto',
    height: '50px',
    width: '85%'
}

const FormGroupStyle = {
    marginTop: '30px'
}

export class AuthForm extends Component {
    render() {
        return (
            <FormGroup style={FormGroupStyle}>
                <Input type={ this.props.type } placeholder={this.props.labelText} style={InputStyle} />
            </FormGroup>
        );
    }
}

export class AuthButton extends Component {
    render() {
        return (
            <FormGroup>
                <Button onClick={this.props.onClick} style={ButtonStyle} size="lg" block>
                    {this.props.labelText}
                </Button>
            </FormGroup>
        )
    }
}
