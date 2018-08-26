package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/gamelife1314/go-note/config"
	"io"
	"strconv"
)

func EncryptModelId(id uint) string {
	idString := string(id)
	return base64.StdEncoding.EncodeToString([]byte(idString))
}

func DecodeModelIdString(id string) uint {

	if data, err := base64.StdEncoding.DecodeString(id); err != nil {
		return 0
	} else {
		if modelId, err := strconv.Atoi(string(data)); err != nil {
			return 0
		} else {
			return uint(modelId)
		}
	}
}

func CryptPassword(password string) string {
	h := hmac.New(sha256.New, []byte(config.Configuration.Other["AppKey"].(string)))
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
