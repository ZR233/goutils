/*
@Time : 2019-10-12 15:13
@Author : zr
*/
package goutils

import "testing"

func TestStringChineseMaxLen(t *testing.T) {
	type args struct {
		str    string
		maxLen int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"我是好人", 6}, "我是"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringChineseMaxLen(tt.args.str, tt.args.maxLen); got != tt.want {
				t.Errorf("StringChineseMaxLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
