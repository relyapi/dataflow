package utils

import (
	"fmt"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	encrypted, err := EncryptAES("123456", "FaC16jUwZo&DNZbj_6S$#lqhoU6l0%XW")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(encrypted)
}

func TestDecryptAES(t *testing.T) {
	result, err := DecryptAES("QneqJ3jqt4YbzLkheK87BDdTro5Nkg==", "FaC16jUwZo&DNZbj_6S$#lqhoU6l0%XW")
	if err != nil {
		return
	}
	fmt.Println(result)
}
