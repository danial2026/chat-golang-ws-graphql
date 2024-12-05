package decoders

import (
	"errors"
	"time"
)

const (
	ErrKidNotExists            = "expecting JWT header to have string kid"
	ErrKidNotExistsInJWK       = "key with specified kid is not present in jwks"
	ErrPublicKeyParseException = "could not parse pubkey"
	ErrInvalidSigningMethod    = "unexpected signing method: "
	ErrInvalidAudience         = "token has invalid audience"
	ErrAudienceNotExists       = "failed to get claims[client_id]"
	ErrInvalidIssuer           = "token has invalid issuer"
	ErrIssuerNotExists         = "failed to get claims[iss]"
	ErrTokenExpired            = "token is expired"
	ErrInvalidTokenType        = "invalid token type"
	ErrTokenTypeNotExists      = "failed to get claims[token_use]"
	ErrTokenExpiryNotExists    = "failed to get claims[exp]"
	ErrInvalidToken            = "token is invalid"
)

var ErrTokenExpire = errors.New(ErrTokenExpired)

const DefaultTimeout = 120 * time.Second
