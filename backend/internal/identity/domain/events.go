package domain

// Domain event type names published by the identity context onto the event bus.
const (
	EventUserRegistered  = "UserRegistered"
	EventEmailVerified   = "EmailVerified"
	EventUserLoggedIn    = "UserLoggedIn"
	EventPasswordReset   = "PasswordReset"
	EventPasswordChanged = "PasswordChanged"
	EventMFAEnabled      = "MFAEnabled"
	EventMFADisabled     = "MFADisabled"
	EventOAuthLinked     = "OAuthLinked"
)

const EventAccountDeactivated = "AccountDeactivated"

// EventUserRolesUpdated is published when a user changes their own roles via the
// self-serve PUT /users/me/roles endpoint.
const EventUserRolesUpdated = "UserRolesUpdated"
