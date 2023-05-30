package constant

import (
	"net/textproto"
)

var (
	XRequestIdHeader        = textproto.CanonicalMIMEHeaderKey("x-request-id")
	XServiceNameHeader      = textproto.CanonicalMIMEHeaderKey("x-service-name")
	XUserType               = textproto.CanonicalMIMEHeaderKey("x-user-type")
	XUserDetail             = textproto.CanonicalMIMEHeaderKey("x-user-detail")
	XApiKeyHeader           = textproto.CanonicalMIMEHeaderKey("x-api-key")
	AuthorizationHeader     = textproto.CanonicalMIMEHeaderKey("authorization")
	ContextBackground       = textproto.CanonicalMIMEHeaderKey("ContextBackground")
	XRequestSignatureHeader = textproto.CanonicalMIMEHeaderKey("x-request-signature")
	XRequestAtHeader        = textproto.CanonicalMIMEHeaderKey("x-request-at")
	ContentTypeHeader       = textproto.CanonicalMIMEHeaderKey("content-type")
)
