package token

import "strings"

// string which must be present in SMS text so we know to use this particular extractor
const MICROSOFT string = "Microsoft authentication"

func extract_microsoft_sms(token_msg string) (string, int) {
	// sample sms:
	// Use verification code 123456 for Microsoft authentication.

	var token string
	var ts int
	words := strings.Fields(token_msg)
	for i, word := range words {
		// get token
		if word == "code" {
			token = words[i+1]
		}

		// get message timestamp
		if word == "date" {
			ts_str := words[i+2]
			ts = str_to_int(ts_str)
		}
	}
	return token, ts
}
