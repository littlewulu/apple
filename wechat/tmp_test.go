package wechat

import (
	"apple/utils"
	"math/rand"
	"testing"
	"time"
)

func TestGetWechatAccessToken(t *testing.T) {
	rand.Seed(time.Now().Unix())
	print(utils.RanString(32))

}



