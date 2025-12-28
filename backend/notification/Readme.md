# Notification Service

## Requirements
- Support for multiple notification channels (e.g., email, SMS, push notifications)
- Support Templates with localizations


## Schema

templates
- id
- event
- channel
- locale
- subject
- body
- version


## Event
{
  "event": "USER_SIGNUP_OTP",
  "channel": ["email", "sms"],
  "recipient": {
    "email": "user@example.com",
    "phone": "+919876543210"
  },
  "data": {
    "otp": "123456",
    "expires_in": "5 minutes"
  },
  "locale": "en-IN"
}
