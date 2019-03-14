import { SIGN_IN_REQUEST, SIGN_IN_SUCCESS, SIGN_IN_FAILURE, SIGN_OUT } from '../actions/authActions'

const initialState = { signedIn: false }

export function authReducer(state = initialState, action) {
    switch (action.type) {
      case SIGN_IN_REQUEST:
        return {
          signedIn: false,
        };
      case SIGN_IN_SUCCESS:
        return {
          signedIn: true,
        };
      case SIGN_IN_FAILURE:
        return {};
      case SIGN_OUT:
        return {};
      default:
        return state
    }
  }
