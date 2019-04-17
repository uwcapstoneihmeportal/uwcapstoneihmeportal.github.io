import React, { Component } from 'react';
import { FormGroup, FormFeedback, Input } from 'reactstrap';

const InputStyle = {
    backgroundRepeat: 'no-repeat',
    backgroundSize: '30px',
    backgroundPosition: '2%',
    border: '1px solid #cbe2a0',
    boxShadow: '0 0 5px #26a146',
    borderRadius: '25px',
    height: '50px',
    margin: '0 auto',
    paddingLeft: '50px',
    width: '85%'
}

const FormGroupStyle = {
    marginTop: '30px'
}

const FormFeedbackStyle = {
    marginLeft: '50px'
}

export default class AuthForm extends Component {
    render() {
        const image = 'url(' + this.props.imagePath + ')'

        return (
            <FormGroup style={FormGroupStyle}>
                <Input type={ this.props.type } placeholder={this.props.labelText} style={{...InputStyle, ...{backgroundImage: image}}} />
            </FormGroup>
        );
    }
}
