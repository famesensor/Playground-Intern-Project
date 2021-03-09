package handlers

import (
	"fmt"
	"log"
	"time"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	interfaces "github.com/HangoKub/Hango-service/internal/core/ports"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/reponseHandler"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUc interfaces.AuthenUsecase
	userUc interfaces.UserUsecase
}

const (
	acExp time.Duration = time.Minute * 10
	rfExp time.Duration = time.Hour * 168
)

func reponseError(c *fiber.Ctx, err error) error {
	status, code, messages := errs.ErrorDetails(err)
	return reponseHandler.ReponseMsg(c, status, "failed", code, messages)
}

func NewAuthHandler(authUc interfaces.AuthenUsecase, userUc interfaces.UserUsecase) interfaces.AuthHanler {
	return &AuthHandler{
		authUc,
		userUc,
	}
}

func (h *AuthHandler) ReqToken(c *fiber.Ctx) error {
	user := c.Locals("bodyData").(*model.LoginUser)

	userDoc, errUser := h.userUc.GetUserById(c.Context(), user.ID)
	// case login user platform
	switch user.Platform {
	case "social":
		if errUser != nil {
			log.Printf("Error Get User By ID : %v\n", errUser)
			return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "", &fiber.Map{"isRegistered": false})
		}
		// generate token and create token document into database
		tkDoc, err := h.authUc.ReqTokenDocument(c.Context(), acExp, rfExp, "Bearer", &model.AuthClaim{UserId: userDoc.UserId, Email: userDoc.Email, HgId: userDoc.HgId})
		if err != nil {
			return reponseError(c, err)
		}
		if err := h.authUc.CreateAuthDoc(c.Context(), tkDoc.AuthDocument); err != nil {
			return reponseError(c, err)
		}
		return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", tkDoc.TokenCard)
	case "otp":
		if len(user.ID) != 10 {
			return reponseError(c, errs.InvalidPhoneNumber)
		}

		// generate otp, create otp document and send otp to user
		otpT := time.Minute * 5
		if err := h.authUc.ReqOtpDocument(c.Context(), user.ID, otpT, model.VERIFY_AUTH_OTP); err != nil {
			return reponseError(c, err)
		}

		if errUser != nil {
			return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "OTP is send", &fiber.Map{"isRegistered": false})
		}
		return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", &fiber.Map{"isRegistered": true})
	default:
		return reponseError(c, errs.InternalServerError)
	}
}

// TODO: Upload Picture Profile
func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	userDoc := c.Locals("bodyData").(*model.RegisterUser)

	if userDoc.Platform == "otp" {
		log.Println("Run")
		if err := h.authUc.VerifyOtp(c.Context(), userDoc.Otp, userDoc.ID, "own-verification", model.VERIFY_AUTH_OTP); err != nil {
			return reponseError(c, err)
		}
	}

	// create user into database
	hgId, err := h.userUc.CreateUser(c.Context(), *userDoc)
	if err != nil {
		fmt.Printf("Error Create : %v\n", err)
		return reponseError(c, err)
	}

	// generate token and create token document into database
	tkDoc, err := h.authUc.ReqTokenDocument(c.Context(), acExp, rfExp, "Bearer", &model.AuthClaim{UserId: userDoc.ID, Email: userDoc.Email, HgId: hgId})
	if err != nil {
		return reponseError(c, err)
	}
	if err := h.authUc.CreateAuthDoc(c.Context(), tkDoc.AuthDocument); err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusCreated, "success", "", tkDoc.TokenCard)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	rfToken := c.Locals("bodyData").(*model.RefreshToken)

	// check refresh token, generate new token and update authen document after refresh token
	tkDoc, err := h.authUc.RefreshToken(c.Context(), rfToken.RefreshToken)
	if err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", tkDoc)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	rfToken := c.Locals("bodyData").(*model.RefreshToken)

	// get token form header
	tk := c.Locals("user").(*jwt.Token).Raw
	// verify user token
	_, err := h.authUc.VerifyToken(tk)
	if err != nil {
		return reponseHandler.ReponseMsg(c, fiber.StatusUnauthorized, "failed", err.Error(), nil)
	}

	// logout
	if err := h.authUc.Logout(c.Context(), rfToken.RefreshToken); err != nil {
		return reponseError(c, err)
	}
	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", nil)
}

func (h *AuthHandler) Otp(c *fiber.Ctx) error {
	otpInf := c.Locals("bodyData").(*model.Otp)

	// check user already existing
	userDoc, err := h.userUc.GetUserById(c.Context(), otpInf.Phone)
	if err != nil {
		if err.Error() == "NO_DATA_FOUND" {
			// verify otp
			if err := h.authUc.VerifyOtp(c.Context(), otpInf.Otp, otpInf.Phone, model.REGISTER_OTP, model.VERIFY_AUTH_OTP); err != nil {
				return reponseError(c, err)
			}
			return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "OTP verify success", nil)
		}
		return reponseError(c, err)
	}

	// verify otp
	if err := h.authUc.VerifyOtp(c.Context(), otpInf.Otp, otpInf.Phone, model.LOGIN_OTP, model.VERIFY_AUTH_OTP); err != nil {
		return reponseError(c, err)
	}

	// generate token and create token document into database
	tkDoc, err := h.authUc.ReqTokenDocument(c.Context(), acExp, rfExp, "Bearer", &model.AuthClaim{UserId: userDoc.UserId, Email: userDoc.Email, HgId: userDoc.HgId})
	if err != nil {
		return reponseError(c, err)
	}
	if err := h.authUc.CreateAuthDoc(c.Context(), tkDoc.AuthDocument); err != nil {
		return reponseError(c, err)
	}
	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", tkDoc.TokenCard)
}

func (h *AuthHandler) RefreshOtp(c *fiber.Ctx) error {
	phone := c.Locals("bodyData").(*model.Phone)

	// check condition otp and generate, send otp to user
	if err := h.authUc.RefreshOtp(c.Context(), phone.Phone, model.VERIFY_AUTH_OTP); err != nil {
		return reponseError(c, err)
	}

	return reponseHandler.ReponseMsg(c, fiber.StatusOK, "success", "", nil)
}
