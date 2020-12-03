package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

// B64ToImage accepts a b64 image string (having content-type specified) and returns a
// base64 encoded string.
func B64ToImage(b64 string) (string, image.Image, error) {
	i := strings.Index(b64, ",")
	if i < 0 {
		return "", nil, fmt.Errorf("no comma in image")
	}

	start := strings.Index(b64, ":")
	end := strings.Index(b64, ";")

	mime := b64[start+1 : end]
	// pass reader to NewDecoder
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64[i+1:]))

	var err error
	var img image.Image

	switch mime {
	case "image/png":
		img, err = png.Decode(dec)
		if err != nil {
			break
		}
		return "image/png", img, err
	case "image/jpg":
		fallthrough
	case "image/jpeg":
		img, err = jpeg.Decode(dec)
		if err != nil {
			break
		}
		return "image/jpg", img, err
	case "image/gif":
		img, err = gif.Decode(dec)
		if err != nil {
			break
		}
		return "image/gif", img, err
	default:
		return "", nil, fmt.Errorf("image format not allowed: %s", mime)
	}

	log.Error(err)
	return "", nil, err
}

//ImageToTypedImage converts an image into a type Image of type png, jpg, gif
func ImageToTypedImage(mime string, img image.Image) (io.Reader, error) {
	var f = bytes.NewBuffer([]byte{})
	var err error

	switch mime {
	case "image/png":
		log.Printf("encoding as %s", mime)
		err = png.Encode(f, img)
		if err != nil {
			log.Error(err)
			break
		}
		return f, nil
	case "image/jpg":
		fallthrough
	case "image/jpeg":
		log.Printf("encoding as %s", mime)
		err = jpeg.Encode(f, img, nil)
		if err != nil {
			log.Error(err)
			break
		}
		return f, nil
	case "image/gif":
		log.Printf("encoding as %s", mime)
		err = gif.Encode(f, img, nil)
		if err != nil {
			log.Error(err)
			break
		}
		return f, nil
	default:
		return f, fmt.Errorf("image format not allowed: %s", mime)
	}

	log.Printf("image encoded as %s, %v", mime, f)
	return f, nil
}
