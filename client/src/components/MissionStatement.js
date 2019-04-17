import React, { Component } from 'react';

const StatementStyle = {
    backgroundColor: "#cbe2a0",
    padding: "30px 50px 30px 50px",
    width: '60%'
}

class MissionStatement extends Component {
    render() {
        return (
            <div className="d-none d-md-block" style={StatementStyle}>
                <h2>
                    Measuring what matters
                </h2>
                <p>
                    Our mission is to improve the health of the world's populations
                    by providing the best information on population health
                </p>
            </div>
        )
    }
}

export default MissionStatement
