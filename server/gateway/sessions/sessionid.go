package sessions

import (
	"crypto/sha256"
	"errors"
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"encoding/base64"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//Represents a valid and digitally signed ID for session.
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) < 1 {
		return InvalidSessionID, errors.New("error in creating new session ID: " +
			"signing key must have at least one character")
	}
	//a byte slice where the first idLength of bytes are cryptographically random bytes
	buffer := make([]byte, idLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return InvalidSessionID, errors.New(fmt.Sprintf("Error in generating random 32 bytes: %v", err))
	}
	//use HMAC hash of the first 32 bytes and the provided signingKey for the remaining bytes.
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(buffer)
	hashedBuffer := h.Sum(nil)
	preBaseEncoded := make([]byte, signedLength)
	preBaseEncoded = append(buffer, hashedBuffer...)
	//encode using base64 URL encoding.
	baseEncoded := SessionID(base64.URLEncoding.EncodeToString(preBaseEncoded))
	return baseEncoded, nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
	//Decode base64 encoded session ID
	decodedId, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, errors.New(fmt.Sprintf("Error in decoding given ID: %v", err))
	}

	//create a new hmac hash of the id portion of the byte slice.
	h := hmac.New(sha256.New, []byte(signingKey))
	buffer := []byte(decodedId)[0:idLength]
	h.Write(buffer)
	hashedBuffer := h.Sum(nil)
	remaining := []byte(decodedId)[idLength:]
	//compare to see if match. Invalid if it doesn't match.
	if hmac.Equal(hashedBuffer, remaining) {
		return SessionID(id), nil
	}
	return InvalidSessionID, ErrInvalidID
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
