package encryption

import (
	"github.com/xiangrui2019/securefshadow/lib/password"
)

type Cipher struct {
	// 编码用的密码
	encodePassword *password.Password
	// 解码用的密码
	decodePassword *password.Password
}

// 加密原数据
func (cipher *Cipher) Encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encodePassword[v]
	}
}

// 解码加密后的数据到原数据
func (cipher *Cipher) Decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decodePassword[v]
	}
}

// 新建一个编码解码器
func NewCipher(encodePassword *password.Password) *Cipher {
	decodePassword := &password.Password{}
	for i, v := range encodePassword {
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}
	return &Cipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}
