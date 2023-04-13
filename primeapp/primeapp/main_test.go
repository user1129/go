package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func TestIntro(t *testing.T) {
	tt := []struct {
		text string
	}{
		{"Is it Prime?"},
		{"------------"},
		{"Enter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit."},
		{"-> "},
	}

	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()

	r, w, err := os.Pipe()
	defer w.Close()
	if err != nil {
		t.Fatalf("%s Test, pipeline creates IO files: %v", failed, err)
	}
	t.Logf("%s Test, pipeline creates IO files.", succeed)

	os.Stdout = w
	os.Stderr = w
	log.SetOutput(w)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buff bytes.Buffer
		wg.Done()
		io.Copy(&buff, r)
		out <- buff.String()
	}()
	wg.Wait()
	intro()
	w.Close()

	str := strings.Split(<-out, "\n")
	for i, test := range tt {
		if str[i] != test.text {
			t.Errorf("%s Test, output should be \"%s\" but it's \"%s\"", failed, test.text, str[i])
		}
		t.Logf("%s Test, output should be \"%s\" but it's \"%s\"", succeed, test.text, str[i])
	}
}

func TestCheckNumbers(t *testing.T) {
	tt := []struct {
		input string
		text  string
		done  bool
	}{
		{"q", "", true},
		{"a", "Please enter a whole number!", false},
		{"2", "2 is a prime number!", false},
	}

	r, w, _ := os.Pipe()
	str := ""
	for _, test := range tt {
		str += test.input + "\n"
	}
	w.Write([]byte(str))
	scanner := bufio.NewScanner(r)
	w.Close()
	for _, test := range tt {
		text, done := checkNumbers(scanner)
		if text != test.text {
			t.Errorf("%s Test, response should be \"%s\" but it's \"%s\"", failed, test.text, text)
		} else {
			t.Logf("%s Test, response should be \"%s\" but it's \"%s\"", succeed, test.text, text)
		}
		if done != test.done {
			t.Errorf("%s Test, quit command should be %t but it's %t", failed, test.done, done)
		} else {
			t.Logf("%s Test, quit command should be %t but it's %t", succeed, test.done, done)
		}
	}

}

func TestReadUserInput(t *testing.T) {
	expectedRes := [2]string{"Please enter a whole number!", "-> "}

	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	r, w, _ := os.Pipe()
	r1, w1, _ := os.Pipe()
	defer w.Close()
	defer w1.Close()

	os.Stdout = w1
	os.Stderr = w1
	log.SetOutput(w1)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buff bytes.Buffer
		wg.Done()
		io.Copy(&buff, r1)
		out <- buff.String()
	}()
	wg.Wait()
	w.Write([]byte("a\nq\n"))
	doneChan := make(chan bool)
	go func() {
		<-doneChan
	}()
	readUserInput(r, doneChan)
	w1.Close()
	w.Close()
	str := strings.Split(<-out, "\n")

	for i, test := range str {
		if test != expectedRes[i] {
			t.Errorf("%s Test, output should be \"%s\" but it's \"%s\"", failed, expectedRes[i], test)
		}
		t.Logf("%s Test, output should be \"%s\" but it's \"%s\"", succeed, expectedRes[i], test)
	}
}
