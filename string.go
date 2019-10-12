/*
@Time : 2019-10-12 15:10
@Author : zr
*/
package goutils

func StringChineseMaxLen(str string, maxLen int) string {
	outStr := ""
	testStr := ""
	nameRune := []rune(str)
	for _, v := range nameRune {
		testStr += string(v)
		testLen := len(testStr)
		if testLen > maxLen {
			return outStr
		}
		outStr = testStr
	}

	return outStr
}
