import React, { Component } from 'react';
import { FormGroup, InputGroup, InputGroupAddon, Input } from 'reactstrap';

const InputStyle = {
    borderRadius: '25px',
    margin: '0 auto',
    height: '50px',
    width: '85%',
    paddingLeft: '50px',
    backgroundRepeat: 'no-repeat',
    backgroundSize: '30px',
    backgroundPosition: '2%'
}

const FormGroupStyle = {
    marginTop: '30px'
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
