package domain

// Domain event type names published by the settings context onto the event bus.
const (
	EventSettingsUpdated      = "SettingsUpdated"
	EventPrivacyChanged       = "PrivacySettingsChanged"
	EventNotificationsChanged = "NotificationSettingsChanged"
)
