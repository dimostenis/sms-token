package sms

import (
	"testing"
)

func TestInjectAndReadStdin(t *testing.T) {
	// my custom stdin
	const myText = "test input"

	// Inject!
	config := Config{Debug: true, DebugText: myText}
	config.injectStdin()

	// call the method that reads from stdin
	sms := Sms{}
	err := sms.readSmsFromStdin()
	if err != nil {
		t.Error("noooo")
	}

	// check if the result is correct
	if sms.text != myText {
		t.Errorf("read_stdin() = %q, want %s", sms.text, myText)
	}
}

func TestReadToken(t *testing.T) {
	// Setup debug config
	const years int = 10000 * 10000
	config := Config{
		Debug:     true,
		DebugText: "text = Token code: 123456\n  date = 715469616040455040",
		MaxAge:    years,
	}

	s, err := ReadToken(config)
	if err != nil {
		t.Error(err)
	}
	if s.text != "text = Token code: 123456\n  date = 715469616040455040" {
		t.Error("Error parsing SMS text")
	}
	if s.token != "123456" {
		t.Error("Error parsing SMS token")
	}
	if s.ts != 715469616040455040 {
		t.Error("Error parsing SMS timestamp")
	}
	if s.age < 1 {
		t.Error("Error getting SMS age")
	}
}
