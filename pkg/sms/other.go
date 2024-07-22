package sms

import "strings"

// string which must be present in SMS text so we know to use this particular extractor
const Other string = "text = Token code:"

func extractOtherSms(token_msg string) (string, int) {
	// sample sms:
	// Token code: 123456

	// get token
	line_with_token := strings.Split(token_msg, "\n")[0]
	token_text := strings.Split(line_with_token, ": ")[1]
	token := strings.TrimSpace(token_text)

	// get message timestamp
	words := strings.Fields(token_msg)
	var ts_str string
	for i, word := range words {
		if word == "date" {
			ts_str = words[i+2]
		}
	}
	return token, strToInt(ts_str)
}
