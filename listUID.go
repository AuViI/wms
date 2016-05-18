package main

import "fmt"
import "math/rand"
import "time"

func getRUID(length int) (uid string) {
	n := int(time.Now().Unix() % 777)
	for x := 0; x < length; x++ {
		uid += order(rand.Intn(61) + n)
	}
	return
}

func order(num int) string {
	num = num % 61
	switch {
	case num < 0:
		return "_"
	case num < 10:
		return fmt.Sprintf("%d", num)
	case num < 36:
		return string(rune(65 + num - 10))
	case num < 61:
		return string(rune(98 + num - 36))
	}
	return "#"
}
