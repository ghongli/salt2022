package utils

import (
	"strings"
)

const (
	// CloudEventContentType the content type for cloud event.
	CloudEventContentType = "application/cloudevents+json; charset=utf-8"

	// PrettyCloudEventContentType the content type for pretty cloud event.
	PrettyCloudEventContentType = "application/cloudevents+json+pretty; charset=utf-8"

	// JSONContentType the content type for JSON.
	JSONContentType = "application/json; charset=utf-8"

	// PrettyJSONContentType the content type for pretty JSON.
	PrettyJSONContentType = "application/json+pretty; charset=utf-8"

	// XMLContentType the content type for XML.
	XMLContentType = "application/xml; charset=utf-8"

	// GRPCContentType the content type for gRPC.
	GRPCContentType = "application/grpc; charset=utf-8"

	// TextPlainContentType the content type for plan.
	TextPlainContentType = "text/plain; charset=utf-8"

	// TextHtmlContentType the content type for plan.
	TextHtmlContentType = "text/html; charset=utf-8"

	// FormUrlencodedContentType the content type for form-urlencoded.
	FormUrlencodedContentType = "application/x-www-form-urlencoded; charset=utf-8"

	// BinaryContentType the content type for binary.
	BinaryContentType = "application/octet-stream; charset=utf-8"

	// HeaderContentTypeKey the header content-type key.
	HeaderContentTypeKey = HeaderContentType
)

const (
	DefaultContentTypeSeparated    = ";"
	DefaultContentTypeStringPrefix = "text/"
)

func PrettyContentType(contentType string) string {
	lowerContentType := strings.ToLower(contentType)
	semiColonPos := strings.Index(lowerContentType, DefaultContentTypeSeparated)
	if semiColonPos >= 0 {
		return lowerContentType[:semiColonPos]
	}
	return lowerContentType
}

func IsContentType(contentType string, expected string) bool {
	lowerContentType := strings.ToLower(contentType)
	if lowerContentType == expected {
		return true
	}

	semiColonPos := strings.Index(lowerContentType, DefaultContentTypeSeparated)
	if semiColonPos >= 0 {
		return lowerContentType[:semiColonPos] == expected[:semiColonPos]
	}

	semiColonPos = strings.Index(expected, DefaultContentTypeSeparated)
	if semiColonPos >= 0 {
		return lowerContentType == expected[:semiColonPos]
	}

	return false
}

func IsJSONContentType(contentType string) bool {
	return IsContentType(contentType, JSONContentType)
}

func IsGRPCContentType(contentType string) bool {
	return IsContentType(contentType, GRPCContentType)
}

// IsStringContentType determines whether content type is string
func IsStringContentType(contentType string) bool {
	if strings.HasPrefix(strings.ToLower(contentType), DefaultContentTypeStringPrefix) {
		return true
	}

	return IsContentType(contentType, XMLContentType)
}

// IsBinaryContentType determines whether content type is byte[].
func IsBinaryContentType(contentType string) bool {
	return IsContentType(contentType, BinaryContentType)
}
