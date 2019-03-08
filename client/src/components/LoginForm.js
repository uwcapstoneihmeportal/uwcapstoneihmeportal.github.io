import React from 'react'
import PropTypes from 'prop-types'

const LoginFormStyle = {
    width: '300px',
    margin: '200px auto'
}

const LoginForm = () => (
    <form style={LoginFormStyle}>
        <label>
            Username:
            <input type="text" name="name" />
        </label>
        <input type="submit" value="Submit" />
    </form>
)

LoginForm.propTypes = {
}

export default LoginForm