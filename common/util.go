package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gamelife1314/go-note/config"
	"io"
	"math/rand"
	"time"
)

func Struct2Map(in interface{}) map[string]interface{} {
	var output map[string]interface{}
	inJson, _ := json.Marshal(in)
	json.Unmarshal(inJson, &output)
	return output
}

func InStringArray(search string, list []string) bool {
	for _, b := range list {
		if b == search {
			return true
		}
	}
	return false
}

func Md5(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func GenerateRandomString(length int64) string {
	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+-,.;':|"
	var sequences []byte
	sequences = make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range sequences {
		sequences[i] = letters[r.Intn(len(letters))]
	}
	return string(sequences)
}

func Sha256(s string) string {
	h := hmac.New(sha256.New, []byte(config.Configuration.Other["AppKey"].(string)))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
