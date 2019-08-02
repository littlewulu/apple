package manager

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// 生成签名
func CalcuSignature(accessToken, secret string)string{
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(accessToken))
	sig := mac.Sum(nil)
	return hex.EncodeToString(sig)
}

