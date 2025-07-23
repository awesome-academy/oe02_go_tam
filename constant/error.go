package constant

import (
	"errors"
)

var (
	EmailAlreadyExists         = errors.New("auth.email.exists")
	ErrUserNotFound            = errors.New("auth.user.not_found")
	LoginFailed                = errors.New("auth.login.failed")
	UnexpectedSigningMethod    = errors.New("auth.token.unexpected_method")
	TokenExpired               = errors.New("auth.token.expired")
	InvalidToken               = errors.New("auth.token.invalid")
	InvalidTokenClaim          = errors.New("auth.token.invalid_claim")
	UserIdNotFoundInToken      = errors.New("auth.token.user_id_missing")
	RoleNotFoundInToken        = errors.New("auth.token.role_missing")
	PasswordMismatch           = errors.New("auth.password.mismatch")
	ErrLoadEnvFailed           = errors.New("config.env.load_failed")
	ErrJWTSecretNotSet         = errors.New("config.env.jwt_secret_missing")
	ErrValidation              = errors.New("validation failed")
	ErrGoogleIDMismatch        = errors.New(T("auth.google_id_mismatch"))
	ErrReviewNotFound          = errors.New("review not found")
	ErrParentCommentNotFound   = errors.New("parent comment not found")
	ErrAlreadyLiked            = errors.New("user already liked this review")
	ErrInvalidRating           = errors.New("Rating invalid")
	ErrEmptyContent            = errors.New("Empty content")
	ErrTourNotFound            = errors.New("tour not found")
	ErrNotEnoughSeats          = errors.New("not enough seats available")
	ErrAlreadyBooked           = errors.New("already booked")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrAlreadyCancelled        = errors.New("booking already cancelled")
	ErrBookingAlreadyProcessed = errors.New("booking already processed or cancelled")
	ErrForbidden               = errors.New("you are not allowed to access this booking")
)
