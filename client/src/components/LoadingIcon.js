import React, { Component } from 'react';
import { DotLoader } from 'react-spinners'

const loadingStyle = {
    display: 'block',
    margin: '0 auto',
    marginTop: '10px',
    marginBottom: '0'
}

export default class LoadingIcon extends Component {
    render() {
        return (
            <DotLoader loading={this.props.loading} color="#26a146" css={loadingStyle} />
        )
    }
}
