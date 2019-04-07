// actions.js

// The three possible states for the login process
export const LOGIN_REQUEST = "LOGIN_REQUEST"
export const LOGIN_SUCCESS = "LOGIN_SUCCESS"
export const LOGIN_FAILURE = "LOGIN_FAILURE"

function requestLogin(credentials) {
    return {
        type: LOGIN_REQUEST,
        isFetching: true,
        isAuthenticated: false,
        credentials
    }
}

function successfulLogin(user) {
    return {
        type: LOGIN_SUCCESS,
        isFetching: false,
        isAuthenticated: true,
        userID: user.userId
    }
}

function loginError(message) {
    return {
        type: LOGIN_SUCCESS,
        isFetching: false,
        isAuthenticated: true,
        message
    }
}
