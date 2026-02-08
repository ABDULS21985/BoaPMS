package notifications

// Notification message constants converted from .NET NotificationMessages.
const (
	NoRecordFound                      = "No Record Found"
	OperationCompleted                 = "Operation Completed successfully"
	InvalidDateFormat                  = "Invalid Date Format Expected Format (yyyyMMdd)"
	TryAgainLaterError                 = "Oops! Something went wrong! Please try again later. If this error continues to occur, please contact our support center"
	InvalidCredentials                 = "Invalid email or password"
	PasswordChangeSuccessful           = "Password changed successfully"
	InvalidPassword                    = "Invalid password"
	LoginLockedOut                     = "Your account is locked because of too many invalid login attempts. Please try again in a few minutes."
	TwoFactorAuthCodeInvalid           = "Invalid code"
	TwoFactorAuthEnabled               = "Two-factor authentication enabled successfully"
	TwoFactorAuthDisabled              = "Two-factor authentication disabled successfully"
	TwoFactorAuthDisableError          = "An error occured while disabling two-factor authentication"
	SessionExpired                     = "Session has expired. Please log in again."
	EmailVerificationLinkResentSuccess = "Account activation link sent successfully"
	SuccessfulEmailVerification        = "Account activated. You can now log in."
	EmailVerificationFailed            = "An error occured while verifying your email. Please try again later and if this error continues to occur, contact our support center"
	EmailAlreadyVerified               = "Your email is already verified. You can log in."
	EmailVerificationRequired          = "Your account is not activated. If you have not received the activation email, you can request a new one on this page."
	GenericException                   = "Process could not be completed."
)
