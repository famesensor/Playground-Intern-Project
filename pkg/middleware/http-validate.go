package middleware

import (
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/HangoKub/Hango-service/pkg/errs"
	"github.com/HangoKub/Hango-service/pkg/reponseHandler"
	"github.com/HangoKub/Hango-service/pkg/validators"
	"github.com/gofiber/fiber/v2"
)

var validImageTypes = map[string]bool{
	"jpeg": true,
	"png":  true,
	"jpg":  true,
}

func ValidateBody(newInstance func() interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		doc := newInstance()
		if err := c.BodyParser(doc); err != nil {
			return reponseHandler.ReponseMsg(c, 400, "failed", err.Error(), errs.CannotParseData.Error())
		}

		isValid, err := validators.ValidateStructErrors(doc)
		if !isValid {
			return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", err)
		}

		c.Locals("bodyData", doc)
		return c.Next()
	}
}

func ValidateQuery(newInstance func() interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		doc := newInstance()
		if err := c.QueryParser(doc); err != nil {
			return reponseHandler.ReponseMsg(c, 400, "failed", err.Error(), errs.CannotParseData.Error())
		}

		isValid, err := validators.ValidateStructErrors(doc)
		if !isValid {
			return reponseHandler.ReponseMsg(c, fiber.StatusBadRequest, "failed", "Validation Errors", err)
		}

		c.Locals("queryData", doc)
		return c.Next()
	}
}

func ValidatePicture() fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return reponseHandler.ReponseMsg(c, 400, "failed", err.Error(), nil)
		}

		files := form.File["files"]

		if len(files) == 0 {
			return reponseHandler.ReponseMsg(c, 400, "failed", "Validate errors", "Images is required!")
		}

		var pictures []model.PictureDoc
		for _, file := range files {
			var filetype []string
			if filetype = strings.Split(file.Header["Content-Type"][0], "/"); filetype[0] != "image" {
				return reponseHandler.ReponseMsg(c, 400, "failed", file.Filename+" type is not support", nil)
			}

			_, exists := validImageTypes[filetype[1]]
			if !exists {
				return reponseHandler.ReponseMsg(c, 400, "failed", file.Filename+" type is not support", nil)
			}

			img, err := file.Open()
			if err != nil {
				return err
			}
			defer img.Close()

			var imgTmp image.Image
			switch filetype[1] {
			case "jpeg":
				imgTmp, err = jpeg.Decode(img)
				if err != nil {
					return err
				}
				break
			case "png":
				imgTmp, err = png.Decode(img)
				if err != nil {
					return err
				}
				break
			}

			pictureDoc := model.PictureDoc{
				File: imgTmp,
				Type: filetype[1],
			}
			pictures = append(pictures, pictureDoc)
		}
		if len(pictures) > 0 {
			c.Locals("pictureData", pictures)
		}

		return c.Next()
	}
}
