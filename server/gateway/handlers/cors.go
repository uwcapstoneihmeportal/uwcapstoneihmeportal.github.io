package handlers

import (
	"net/http"
)

/* TODO: implement a CORS middleware handler, as described
  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/

type CorsHandler struct {
	http.Handler
}

func NewCorsHandler(handlerToWrap http.Handler) http.Handler {
	return &CorsHandler{handlerToWrap}
}

func (ch *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(headerAccessControlAllowOrigin, originAny)
	w.Header().Add(headerAccessControlAllowMethods, allowMethods)
	w.Header().Add(headerAccessControlAllowHeaders, allowHeaders)
	w.Header().Add(headerAccessControlExposeHeaders, authHeader)
	w.Header().Add(headerAccessControlMaxAge, maxAge)
	if r.Method != http.MethodOptions {
		ch.Handler.ServeHTTP(w, r)
	}
}