package sms

import (
	"fmt"
	"strconv"
)

func strToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		panic(" !! Yikes")
	}
	return num
}
