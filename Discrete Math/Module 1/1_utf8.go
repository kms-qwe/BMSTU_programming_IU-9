/*
Кодирование и декодирование UTF-8
Источники:
https://www.youtube.com/watch?v=uTJoJtNYcaQ
https://habr.com/ru/articles/138173/
*/
package main

func main() {
	// fmt.Println(string(decode(encode([]rune("Hello 世界")))))
}

func encode(utf32 []rune) []byte {
	var utf8Bytes []byte
	for _, codePoint := range utf32 {
		if codePoint <= 0x7F {
			utf8Bytes = append(utf8Bytes, byte(codePoint))
		} else if codePoint <= 0x7FF {
			utf8Bytes = append(utf8Bytes, byte(0xC0|codePoint>>6), byte(0x80|codePoint&0x3F))
		} else if codePoint <= 0xFFFF {
			utf8Bytes = append(utf8Bytes, byte(0xE0|codePoint>>12), byte(0x80|codePoint>>6&0x3F), byte(0x80|codePoint&0x3F))
		} else {
			utf8Bytes = append(utf8Bytes, byte(0xF0|codePoint>>18), byte(0x80|codePoint>>12&0x3F), byte(0x80|codePoint>>6&0x3F), byte(0x80|codePoint&0x3F))
		}
	}
	return utf8Bytes
}

func decode(utf8 []byte) []rune {
	var result []rune
	numLeadingOnes := func(b byte) int {
		for i := 7; i >= 0; i-- {
			if b&(1<<uint(i)) == 0 {
				return 7 - i
			}
		}
		return 8
	}
	for i := 0; i < len(utf8); {
		numBytes := numLeadingOnes(utf8[i])
		if numBytes == 0 {
			result = append(result, rune(utf8[i]))
			i++
		} else {
			codePoint := rune(utf8[i] & (1<<uint(8-numBytes) - 1))
			for j := 1; j < numBytes; j++ {
				codePoint = (codePoint << 6) | (rune(utf8[i+j]) & 0x3F)
			}
			result = append(result, codePoint)
			i += numBytes
		}
	}
	return result
}
