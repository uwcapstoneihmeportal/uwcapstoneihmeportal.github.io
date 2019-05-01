import React, { Component } from 'react';
import { Card, Button, CardHeader, CardFooter, CardBody, CardTitle, CardText } from 'reactstrap';
import classnames from 'classnames';

const HeaderStyle = {
    backgroundColor: 'grey',
    color: 'white',
    fontSize: '15pt'
}

class CardContainer extends Component {
    render() {
        return (
            <div style={{ marginTop: '20px', marginBottom: '20px' }}>
                <Card>
                    <CardHeader style={HeaderStyle}>
                        {this.props.title}
                    </CardHeader>
                    <CardBody>
                        <CardTitle>Name</CardTitle>
                        <CardText>
                            With supporting text below as a natural lead-in to additional content.
                        </CardText>
                    </CardBody>
                </Card>
            </div>
        )
    }
}

export default CardContainer
