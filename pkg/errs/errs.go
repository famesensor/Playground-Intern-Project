package errs

import (
	"errors"
	"net/http"

	"github.com/HangoKub/Hango-service/internal/core/domain"
)

var (
	DocumentNotFound      = errors.New("NO_DATA_FOUND")
	CannotParseData       = errors.New("CANNOT_PARSE_DATA")
	InternalServerError   = errors.New("INTERNAL_SERVER_ERROR")
	UserAlreadyExsiting   = errors.New("User_Already_Exsiting")
	EmailAlreadyExsiting  = errors.New("Email_Already_Exsiting")
	InvalidPhoneNumber    = errors.New("INVALID_PHONE_NUMBER")
	InvalidToken          = errors.New("Invalid_Token")
	OtpExpires            = errors.New("OTP_is_Expires")
	InvalidOTP            = errors.New("Invalid_OTP")
	FailedReqOTP          = errors.New("Failed_Request_OTP")
	OperationIsnotAllowed = errors.New("Operation_is_not_Allowed")
	UserOutofReduis       = errors.New("User_Out_of_Reduis")
	PictureExceed         = errors.New("The_number_of_images_was_exceeded_the_limit")
	UserNotAllowed        = errors.New("You_are_Not_Allowed")
	YouCannotCheckInApart = errors.New("You_cannot_check_in_apart_of_office_hour")
)

var (
	Msgs_CannotParseData       string = "data cannot be parsed"
	Msgs_DocumentNotFound      string = "data not found"
	Msgs_UserAlreadyExsiting   string = "User already exsiting"
	Msgs_EmailAlreadyExsiting  string = "Email already exsiting"
	Msgs_InvalidPhoneNumber    string = "Invalid phone number"
	Msgs_InvalidToken          string = domain.TOKEN_INVALID
	MSgs_OtpExpires            string = "OTP is Expires"
	Msgs_InvalidOTP            string = "Invalid OTP or Wrong Type OTP"
	Msgs_FailedReqOTP          string = "Failed to request otp(cool down)"
	Msgs_OperationIsnotAllowed string = "Operation isn't Allowed"
	Msgs_UserOutofReduis       string = "User out of restaurant reduis"
	Msgs_PictureExceed         string = "The number of images was exceeded the limit"
	Msgs_UserNotAollwed        string = "You are not allowed"
	Msgs_YouCannotCheckInApart string = "You cannot check in apart of office hour"
)

func ErrorDetails(err error) (int, string, []string) {
	switch err {
	case DocumentNotFound:
		return http.StatusNotFound, DocumentNotFound.Error(), []string{Msgs_DocumentNotFound}
	case CannotParseData:
		return http.StatusBadRequest, CannotParseData.Error(), []string{Msgs_CannotParseData}
	case UserAlreadyExsiting:
		return http.StatusBadRequest, UserAlreadyExsiting.Error(), []string{Msgs_UserAlreadyExsiting}
	case InvalidToken:
		return http.StatusBadRequest, Msgs_InvalidToken, []string{Msgs_InvalidToken}
	case InvalidPhoneNumber:
		return http.StatusBadRequest, InvalidPhoneNumber.Error(), []string{Msgs_InvalidPhoneNumber}
	case OtpExpires:
		return http.StatusBadRequest, OtpExpires.Error(), []string{MSgs_OtpExpires}
	case InvalidOTP:
		return http.StatusBadRequest, InvalidOTP.Error(), []string{Msgs_InvalidOTP}
	case FailedReqOTP:
		return http.StatusBadRequest, FailedReqOTP.Error(), []string{Msgs_FailedReqOTP}
	case OperationIsnotAllowed:
		return http.StatusBadRequest, OperationIsnotAllowed.Error(), []string{Msgs_OperationIsnotAllowed}
	case UserOutofReduis:
		return http.StatusBadRequest, UserOutofReduis.Error(), []string{Msgs_UserOutofReduis}
	case EmailAlreadyExsiting:
		return http.StatusBadRequest, EmailAlreadyExsiting.Error(), []string{Msgs_EmailAlreadyExsiting}
	case PictureExceed:
		return http.StatusBadRequest, PictureExceed.Error(), []string{Msgs_PictureExceed}
	case UserNotAllowed:
		return http.StatusBadRequest, UserNotAllowed.Error(), []string{Msgs_UserNotAollwed}
	case YouCannotCheckInApart:
		return http.StatusBadRequest, YouCannotCheckInApart.Error(), []string{Msgs_YouCannotCheckInApart}
	default:
		return http.StatusInternalServerError, InternalServerError.Error(), []string{err.Error()}
	}
}
