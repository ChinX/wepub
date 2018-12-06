package module

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/chinx/morph"
)

func SaveScale(r io.Reader, outDir string, width float64) (string, error) {
	byteArr, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	buff := bytes.NewReader(byteArr)
	filename := outDir + SHA256File(buff)
	buff.Seek(0, 0)
	return morph.Scale(buff, filename, width)
}

func SHA256File(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return fmt.Sprintf("%x", h.Sum(nil))
}
