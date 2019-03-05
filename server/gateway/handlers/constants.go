package handlers

//header names
const (
	headerContentType = "Content-Type"
	headerAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	headerAccessControlAllowMethods = "Access-Control-Allow-Methods"
	headerAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	headerAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	headerAccessControlMaxAge = "Access-Control-Max-Age"
)

//content types
const (
	contentTypeJSON = "application/json"
	contentTypeText = "text/plain; charset=utf-8"
)

//header values

const (
	originAny    = "*"
	allowMethods = "GET, PUT, POST, PATCH, DELETE"
	allowHeaders = "Content-Type, Authorization"
	maxAge       = "600"
	authHeader   = "Authorization"
	schemeBearer = "Bearer "
)