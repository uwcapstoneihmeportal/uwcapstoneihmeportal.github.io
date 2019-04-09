import React, { Component } from 'react';
import { FormGroup, Input } from 'reactstrap';

const InputStyle = {
    borderRadius: '25px',
    margin: '0 auto',
    height: '50px',
    width: '85%'
}

const FormGroupStyle = {
    marginTop: '30px'
}

export default class AuthForm extends Component {
    render() {
        return (
            <FormGroup style={FormGroupStyle}>
                <Input type={ this.props.type } placeholder={this.props.labelText} style={InputStyle} />
            </FormGroup>
        );
    }
}
