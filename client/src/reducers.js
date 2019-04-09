import { combineReducers } from 'redux'
import { LOGIN_REQUEST, LOGIN_SUCCESS, LOGIN_FAILURE } from "./actions"

const authDefaultState = {
    isFetching: false,
    isAuthenticated: false
}

function auth(state = authDefaultState, action) {
    switch(action.type) {
        case LOGIN_REQUEST:
            return Object.assign({}, state, {
                isFetching: true,
                isAuthenticated: false,
                user: action.credentials
            })
        case LOGIN_SUCCESS:
            return Object.assign({}, state, {
                isFetching: false,
                isAuthenticated: true,
                errorMessage: ''
            })
        case LOGIN_FAILURE:
            return Object.assign({}, state, {
                isFetching: false,
                isAuthenticated: true,
                errorMessage: action.message
            })
        default:
            return state
    }
}

const rootReducer = combineReducers({auth})

export default rootReducer