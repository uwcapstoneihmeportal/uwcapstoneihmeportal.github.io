export const SIGN_IN_REQUEST = "SIGN_IN_REQUEST";
export const SIGN_IN_SUCCESS = "SIGN_IN_SUCCESS";
export const SIGN_IN_FAILURE = "SIGN_IN_FAILURE";
export const SIGN_OUT = "SIGN_OUT";

// TODO: ADD user as parameter to appropriate functions

export function signInRequest() {
    return {
        type: SIGN_IN_REQUEST
    }
}

export function signInSuccess(user) {
    return {
        type: SIGN_IN_SUCCESS
    }
}

export function signInFailure(error) {
    return {
        type: SIGN_IN_FAILURE
    }
}

export function signOut(user) {
    return {
        type: SIGN_OUT
    }
}

