/*
@Time : 2019-10-10 14:33
@Author : zr
*/
package goutils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

func MacStringFromBytes(macBytes []byte, splitStr string) (macStr string, err error) {

	if macLen := len(macBytes); macLen != 6 {
		err = errors.New(fmt.Sprintf("mac data len err: (%d)", macLen))
		return
	}

	var mac []string

	for _, v := range macBytes {
		mac = append(mac, hex.EncodeToString([]byte{v}))
	}
	macStr = strings.Join(mac, splitStr)
	return
}

func MacStringToBytes(macStr string, splitStr string) (macBytes []byte, err error) {
	macList := strings.Split(macStr, splitStr)
	if macLen := len(macList); macLen != 6 {
		err = errors.New(fmt.Sprintf("mac len not 6: (%d)", macLen))
		return
	}
	var b []byte
	for _, v := range macList {
		b, err = hex.DecodeString(v)
		macBytes = append(macBytes, b...)
	}

	return
}
