package qrcode

import (
	"strings"

	"github.com/skip2/go-qrcode"
)

// GenerateQrCode genetare a qrcode in a png image
func GenerateQrCode(data string, quality string, size uint) ([]byte, error) {
	var recoveryLevel qrcode.RecoveryLevel
	switch strings.ToLower(quality) {
	case "low":
		recoveryLevel = qrcode.Low
	case "medium":
		recoveryLevel = qrcode.Medium
	case "high":
		recoveryLevel = qrcode.High
	case "highest":
		recoveryLevel = qrcode.Highest
	default:
		recoveryLevel = qrcode.Medium
	}

	png, err := qrcode.Encode(data, recoveryLevel, int(size))
	if err != nil {
		return nil, err
	}

	return png, nil
}
