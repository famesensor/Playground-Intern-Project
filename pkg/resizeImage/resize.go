package resizeImage

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	model "github.com/HangoKub/Hango-service/internal/core/domain"
	"github.com/nfnt/resize"
)

func ResizeImage(PictureDoc []model.PictureDoc) (model.PictureResize, error) {
	var resizeThumbnail, resizeLarge image.Image
	var pictureResizes model.PictureResize
	for _, element := range PictureDoc {
		if element.File.Bounds().Dy() > element.File.Bounds().Dx() {
			resizeThumbnail = resize.Resize(0, 1000, element.File, resize.Lanczos3)
			resizeLarge = resize.Resize(0, 2000, element.File, resize.Lanczos3)
		} else {
			resizeThumbnail = resize.Resize(1000, 0, element.File, resize.Lanczos3)
			resizeLarge = resize.Resize(2000, 0, element.File, resize.Lanczos3)
		}

		var encodedThumbnail, encodedLarge bytes.Buffer
		switch element.Type {
		case "jpeg":
			if err := jpeg.Encode(&encodedThumbnail, resizeThumbnail, nil); err != nil {
				return model.PictureResize{}, err
			}
			if err := jpeg.Encode(&encodedLarge, resizeLarge, nil); err != nil {
				return model.PictureResize{}, err
			}
			break
		case "png":
			if err := png.Encode(&encodedThumbnail, resizeThumbnail); err != nil {
				return model.PictureResize{}, err
			}
			if err := png.Encode(&encodedLarge, resizeLarge); err != nil {
				return model.PictureResize{}, err
			}
			break
		}

		pictureResizes.ImageThumbnail = append(pictureResizes.ImageThumbnail, encodedThumbnail)
		pictureResizes.ImageLarge = append(pictureResizes.ImageLarge, encodedLarge)
	}

	return pictureResizes, nil
}
