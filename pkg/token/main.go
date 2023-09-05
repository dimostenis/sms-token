package token

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

// Simulate stdin, its for debugging and testing only
func inject_stdin(s string) {

	// create a pipe
	r, w, _ := os.Pipe()

	// write test input to the pipe
	fmt.Fprintln(w, s)

	// close the write end of the pipe
	w.Close()

	// set os.Stdin to the read end of the pipe
	os.Stdin = r
}

func read_stdin() string {
	// check if there is any data available to read from standard input
	// stat, _ := os.Stdin.Stat()
	// size := stat.Size()
	// if size == 0 {
	// 	fmt.Println("No stdin (SMS) recieved. Arent you running it directly without '-install' first?")
	// 	os.Exit(1)
	// }

	// scan stdin from sqlite3 cli
	scanner := bufio.NewScanner(os.Stdin)

	// save all lines to arr
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, " !! Err while reading standard input:", err)
	}

	// array -> multiline string
	return strings.Join(lines, "\n")
}

func str_to_int(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		panic(" !! Yikes")
	}
	return num
}

func copy_to_clipboard(s string) {
	err := clipboard.WriteAll(s)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func get_age(ts int) int {
	const FMT string = "2006-01-02"
	const ORIGIN string = "2001-01-01"

	// time obj which will be a baseline for calculating message age
	origin_dt, err := time.Parse(FMT, ORIGIN)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// this is how timestamps are stored in MacOS Messages...
	n := ts / 1_000_000_000
	delta := time.Duration(n) * time.Second
	msg_ts := origin_dt.Add(delta)

	// messages dates are in UTC
	now := time.Now()
	age := now.Sub(msg_ts).Seconds()

	return int(age)
}

func get_sms() string {
	if debug := os.Getenv("DEBUG"); debug != "" {
		println("DEBUG MODE ON, reading fake SMS")
		if microsoft := os.Getenv("MICROSOFT"); microsoft != "" {
			inject_stdin(microsoft)
		} else if other := os.Getenv("OTHER"); other != "" {
			inject_stdin(other)
		} else {
			fmt.Println("Set some debug env var (MICROSOFT | OTHER) with SMS text")
			os.Exit(1)
		}
	}

	// run func read_stdin() for 1 second or it times out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // executes after end of func

	var sms string
	result := make(chan int)
	go func() {
		sms = read_stdin()
		result <- 0
	}()

	select {
	case <-ctx.Done():
		fmt.Println("No stdin (SMS) recieved for 1 sec. Arent you running it directly without '-install' first?")
		os.Exit(1)
	case <-result:
		// good, we have sms text
	}

	return sms
}

func GetToken() {
	// this can be influeced by DEBUG env var
	sms_with_token := get_sms()

	// extract token and timestamp from sms
	var token string
	var ts int
	if strings.Contains(sms_with_token, MICROSOFT) {
		token, ts = extract_microsoft_sms(sms_with_token)
	} else if strings.Contains(sms_with_token, OTHER) {
		token, ts = extract_other_sms(sms_with_token)
	} else {
		// when there is no message at all
		fmt.Println(" :: no token")
		return
	}

	// token age in seconds
	age := get_age(ts)

	// tokens older than 10 mins are useless
	if age > 10*60 {
		fmt.Printf(" :: token is too old, try again")
		os.Exit(0)
	}

	var unit_name string
	var unit_num int
	if age == 1 {
		unit_num = age
		unit_name = "second"
	} else if age < 120 {
		unit_num = age
		unit_name = "seconds"
	} else {
		unit_num = age * 60
		unit_name = "minutes"
	}
	copy_to_clipboard(token)
	fmt.Printf(" :: token '%s' copied to clipboard (%d %s old)", token, unit_num, unit_name)
}
