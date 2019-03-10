import React, { Component } from 'react';
import { FormGroup, Label, Input, Button, FormFeedback } from 'reactstrap';

const ButtonStyle = {
    backgroundColor: '#26a146'
}

export class AuthForm extends Component {
    render() {
        return (
            <FormGroup>
                <Label >{this.props.labelText}</Label>
                <Input />
            </FormGroup>
        );
    }
}

export class AuthButton extends Component {
    render() {
        return (
            <FormGroup>
                <Button style={ButtonStyle} size="lg" block>
                    {this.props.labelText}
                </Button>
            </FormGroup>
        )
    }
}
