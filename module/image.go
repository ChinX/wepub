package module

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/gif"
	"io"
	"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
)

func ScaleImage(r io.Reader, width, height int) (image.Image, error) {
	img, err := imaging.Decode(r)
	if err != nil {
		return img, err
	}
	if width == 0 && height == 0 {
		return img, err
	}
	return imaging.Resize(img, width, height, imaging.Lanczos), nil
}

func saveScaleGif(byteArr []byte, outDir string, width, height int) (string, error) {
	buff := bytes.NewBuffer(byteArr)
	img, err := gif.DecodeAll(buff)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.gif", SHA256File(buff))
	f, err := os.Create(outDir + filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return filename, gif.EncodeAll(f, img)
}

func SaveScale(r io.Reader, outDir string, width, height int) (string, error) {
	byteArr, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBuffer(byteArr)
	config, fm, err := image.DecodeConfig(buff)
	if err != nil {
		return "", err
	}

	if fm == "gif" {
		return saveScaleGif(byteArr, outDir, width, height)
	}
	buff.Reset()
	buff.Write(byteArr)
	var img image.Image
	if config.Width > width {
		img, err = ScaleImage(buff, width, height)
	} else {
		img, err = imaging.Decode(buff)
	}
	if err != nil {
		return "", err
	}

	buff.Reset()
	buff.Write(byteArr)

	filename := fmt.Sprintf("%s.%s", SHA256File(buff), fm)
	return filename, imaging.Save(img, outDir+filename)
}

func SHA256File(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return fmt.Sprintf("%x", h.Sum(nil))
}
