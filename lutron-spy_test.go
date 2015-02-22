package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

type MockTransport struct {
	needFailNoResponder bool
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestMain(m *testing.M) {
	cf, _ := os.Open("example-config.json")
	parseConfig(cf)

	client = &http.Client{Transport: &MockTransport{}}
	os.Exit(m.Run())
}

func Example_spy_r1_on() {
	tf, _ := os.Open("traces/AEB551_on.trace")
	spy(tf)
	// Output:
	// serial:  AEB551
	// button: 2
	// nickname:  bedside
	// PUT {"on":true} http://192.168.1.1/api/blah/lights/2/state
	//
	// EOF reached.
}

func Example_spy_r1_up() {
	tf, _ := os.Open("traces/AEB551_up.trace")
	spy(tf)
	// Output:
	// serial:  AEB551
	// button: 5
	// nickname:  bedside
	//
	// EOF reached.
}

func Example_spy_r1_select() {
	tf, _ := os.Open("traces/AEB551_select.trace")
	spy(tf)
	// Output:
	// serial:  AEB551
	// button: 3
	// nickname:  bedside
	//
	// EOF reached.
}

func Test_parseLine_r1_on(t *testing.T) {
	Helper_parseLine(t, "AEB551", "on")
}

func Test_parseLine_r2_on(t *testing.T) {
	Helper_parseLine(t, "AFA592", "on")
}

func Test_parseLine_r1_off(t *testing.T) {
	Helper_parseLine(t, "AEB551", "off")
}

func Test_parseLine_r1_up(t *testing.T) {
	Helper_parseLine(t, "AEB551", "up")
}

func Test_parseLine_r1_down(t *testing.T) {
	Helper_parseLine(t, "AEB551", "down")
}

func Test_parseLine_r1_select(t *testing.T) {
	Helper_parseLine(t, "AEB551", "select")
}

func Helper_parseLine(t *testing.T, remote string, button string) {
	fh, err := os.Open(fmt.Sprintf("traces/%s_%s.trace", remote, button))
	if err != nil {
		t.Error("Missing testing trace")
	}
	reader := bufio.NewReader(fh)
	// First string is a button press
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Error("Corrupted testing trace")
	}
	serial, pressed, err := parseLine(line)
	assert.Equal(t, remote, serial)
	assert.Equal(t, button, BUTTONS[pressed])
	assert.Nil(t, err)
}
