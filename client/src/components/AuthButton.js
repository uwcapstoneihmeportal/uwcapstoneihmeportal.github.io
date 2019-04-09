import React, { Component } from 'react';
import { FormGroup, Button } from 'reactstrap';

const ButtonStyle = {
    backgroundColor: '#26a146',
    borderRadius: '25px',
    margin: '0 auto',
    marginTop: '20px',
    width: '85%'
}

export default class AuthButton extends Component {
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
