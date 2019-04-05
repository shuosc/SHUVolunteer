package service

import (
	"github.com/otiai10/gosseract"
	"math/rand"
	"time"
)

var OcrClient *gosseract.Client

func init() {
	rand.Seed(time.Now().UnixNano())
	OcrClient = gosseract.NewClient()
	_ = OcrClient.SetWhitelist("0123456789")
}

func Ocr(bytes []byte) (string, error) {
	_ = OcrClient.SetImageFromBytes(bytes)
	return OcrClient.Text()
}
