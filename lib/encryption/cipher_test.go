package encryption

import (
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/xiangrui2019/securefshadow/lib/password"
)

const (
	MB = 1024 * 1024
)

// 测试 Cipher 加密解密
func TestCipher(t *testing.T) {
	password_r := password.RandPassword()
	t.Log(password_r)
	p, _ := password.ParsePassword(password_r)
	cipher := NewCipher(p)
	// 原数据
	org := make([]byte, password.PasswordLength)
	for i := 0; i < password.PasswordLength; i++ {
		org[i] = byte(i)
	}
	// 复制一份原数据到 tmp
	tmp := make([]byte, password.PasswordLength)
	copy(tmp, org)
	t.Log(tmp)
	// 加密 tmp
	cipher.Encode(tmp)
	t.Log(tmp)
	// 解密 tmp
	cipher.Decode(tmp)
	t.Log(tmp)
	if !reflect.DeepEqual(org, tmp) {
		t.Error("解码编码数据后无法还原数据，数据不对应")
	}
}

func BenchmarkEncode(b *testing.B) {
	password_r := password.RandPassword()
	p, _ := password.ParsePassword(password_r)
	cipher := NewCipher(p)
	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs)
	cipher.Encode(bs)
}

func BenchmarkDecode(b *testing.B) {
	password_r := password.RandPassword()
	p, _ := password.ParsePassword(password_r)
	cipher := NewCipher(p)
	bs := make([]byte, MB)
	b.ResetTimer()
	rand.Read(bs)
	cipher.Decode(bs)
}
