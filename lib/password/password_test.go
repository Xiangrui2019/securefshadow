package password

import (
	"testing"
)

func TestRandPassword(t *testing.T) {
	pw := RandPassword()
	t.Log(pw)

	if pw == "" {
		t.Error("生成密码失败.")
	}
}

func TestParsePassword(t *testing.T) {
	pw := RandPassword()
	t.Log(pw)

	if pw == "" {
		t.Error("生成密码失败.")
	}

	_, err := ParsePassword(pw)

	if err != nil {
		t.Error("解析失败.")
	}
}
