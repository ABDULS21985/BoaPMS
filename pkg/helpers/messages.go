package helpers

// Notification message constants mirroring .NET NotificationMessages.
const (
	MsgNoRecordFound            = "No Record Found"
	MsgOperationCompleted       = "Operation Completed successfully"
	MsgInvalidDateFormat        = "Invalid Date Format Expected Format (yyyyMMdd)"
	MsgTryAgainLater            = "Oops! Something went wrong! Please try again later..."
	MsgInvalidCredentials       = "Invalid email or password"
	MsgPasswordChangeSuccessful = "Password changed successfully"
	MsgInvalidPassword          = "Invalid password"
	MsgLoginLockedOut           = "Your account is locked because of too many invalid login attempts. Please try again later."
	MsgTwoFACodeInvalid         = "Invalid code"
	MsgTwoFAEnabled             = "Two-factor authentication enabled successfully"
	MsgTwoFADisabled            = "Two-factor authentication disabled successfully"
	MsgTwoFADisableError        = "An error occured while disabling two-factor authentication"
	MsgSessionExpired           = "Session has expired. Please log in again."
	MsgEmailVerifLinkResent     = "Account activation link sent successfully"
	MsgEmailVerifSuccess        = "Account activated. You can now log in."
	MsgEmailVerifFailed         = "An error occured while verifying your email. Please try again."
	MsgEmailAlreadyVerified     = "Your email is already verified. You can log in."
	MsgEmailVerifRequired       = "Your account is not activated. Please check your email for the activation link."
	MsgGenericException         = "Process could not be completed."
)
