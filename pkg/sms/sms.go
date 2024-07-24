package sms

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

type Config struct {
	Debug     bool   // flag we are in tests/dev
	DebugText string // simulate SMS
	MaxAge    int    // seconds
}

func (c *Config) injectStdin() {
	// create a pipe
	r, w, _ := os.Pipe()

	// write test input to the pipe
	fmt.Fprintln(w, c.DebugText)

	// close the write end of the pipe
	w.Close()

	// set os.Stdin to the read end of the pipe
	os.Stdin = r
}

type Sms struct {
	text  string // full SMS text
	token string // extracted token
	ts    int    // timestamp in MacOS messages
	age   int    // how old is the message in seconds
}

func (s *Sms) readStdin() {
	// scan stdin from sqlite3 cli
	scanner := bufio.NewScanner(os.Stdin)

	// save all lines to arr
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// array -> multiline string
	s.text = strings.Join(lines, "\n")
}

func (s *Sms) readSmsFromStdin() error {
	// run func read_stdin() for 1 second or it times out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // executes after end of func

	result := make(chan int)
	go func() {
		s.readStdin()
		result <- 0
	}()

	select {
	case <-ctx.Done():
		return errors.New("No stdin (SMS) recieved for 1 sec. Arent you running it directly without '-install' first?")
	case <-result:
		// good, we have sms text
		return nil
	}
}

func (s *Sms) copyTokenToClipboard() error {
	err := clipboard.WriteAll(s.token + "\n")
	if err != nil {
		return err
	}
	return nil
}

func (s *Sms) extractTokenAndTimestamp() error {
	var token string
	var ts int
	if strings.Contains(s.text, Microsoft) {
		token, ts = extractMicrosoftSms(s.text)
	} else if strings.Contains(s.text, Other) {
		token, ts = extractOtherSms(s.text)
	} else {
		return errors.New("SMS format not supported")
	}
	s.token = token
	s.ts = ts
	return nil
}

func (s *Sms) getAgeFromTimestamp() error {
	const Fmt string = "2006-01-02"
	const Origin string = "2001-01-01"

	// time obj which will be a baseline for calculating message age
	originDatetime, err := time.Parse(Fmt, Origin)
	if err != nil {
		return err
	}

	// this is how timestamps are stored in MacOS Messages...
	n := s.ts / 1_000_000_000
	delta := time.Duration(n) * time.Second
	msgTimestamp := originDatetime.Add(delta)

	// messages dates are in UTC
	now := time.Now()
	age := now.Sub(msgTimestamp).Seconds()

	s.age = int(age)
	return nil
}

func (s *Sms) checkAge(secs int) error {
	if s.age > secs {
		msg := fmt.Sprintf(" :: Token is too old (%d minutes), try again.", s.age/60)
		return errors.New(msg)
	}
	return nil
}

func (s *Sms) getSuccessMsg() string {
	var ageString string

	switch {
	case s.age == 1:
		ageString = "1 second"
	case s.age < 120:
		ageString = fmt.Sprintf("%d seconds", s.age)
	default:
		ageString = fmt.Sprintf("%d minutes", s.age/60)
	}

	return fmt.Sprintf(" :: token '%s' copied to clipboard (%s old)", s.token, ageString)
}

func ReadToken(config Config) (Sms, error) {
	if config.Debug {
		config.injectStdin()
	}

	sms := Sms{}

	err := sms.readSmsFromStdin()
	if err != nil {
		return sms, err
	}

	err = sms.extractTokenAndTimestamp()
	if err != nil {
		return sms, err
	}

	err = sms.getAgeFromTimestamp()
	if err != nil {
		return sms, err
	}

	err = sms.checkAge(config.MaxAge)
	if err != nil {
		return sms, err
	}

	err = sms.copyTokenToClipboard()
	if err != nil {
		return sms, err
	}

	// build informative text for user
	println(sms.getSuccessMsg())

	return sms, nil
}
