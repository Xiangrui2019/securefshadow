package password

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const PasswordLength = 256

type Password [PasswordLength]byte

func init() {
	rand.Seed(time.Now().Unix())
}

// 采用base64编码把密码转换为字符串
func (password *Password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}

// 解析采用base64编码的字符串获取密码
func ParsePassword(passwordString string) (*Password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != PasswordLength {
		return nil, errors.New("不合法的密码")
	}
	password := Password{}
	copy(password[:], bs)
	bs = nil
	return &password, nil
}

// 产生 256个byte随机组合的 密码，最后会使用base64编码为字符串存储在配置文件中
// 不能出现任何一个重复的byte位，必须又 0-255 组成，并且都需要包含
func RandPassword() string {
	// 随机生成一个由  0~255 组成的 byte 数组
	intArr := rand.Perm(PasswordLength)
	password := &Password{}
	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return password.String()
}
