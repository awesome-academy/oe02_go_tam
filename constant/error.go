package constant

import (
	"errors"
)

var (
	EmailAlreadyExists      = errors.New("auth.email.exists")
	ErrUserNotFound         = errors.New("auth.user.not_found")
	LoginFailed             = errors.New("auth.login.failed")
	UnexpectedSigningMethod = errors.New("auth.token.unexpected_method")
	TokenExpired            = errors.New("auth.token.expired")
	InvalidToken            = errors.New("auth.token.invalid")
	InvalidTokenClaim       = errors.New("auth.token.invalid_claim")
	UserIdNotFoundInToken   = errors.New("auth.token.user_id_missing")
	RoleNotFoundInToken     = errors.New("auth.token.role_missing")
	PasswordMismatch        = errors.New("auth.password.mismatch")
	ErrLoadEnvFailed        = errors.New("config.env.load_failed")
	ErrJWTSecretNotSet      = errors.New("config.env.jwt_secret_missing")
)
