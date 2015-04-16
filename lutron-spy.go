// TODO Reimplment slsnif?

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type RemoteConfig struct {
	Remotes map[string]struct {
		Nickname string `json:"nickname"`
		Buttons  map[string]struct {
			Data   interface{} `json:"data"`
			URL    string      `json:"url"`
			Method string      `json:"method"`
		  Type   string      `json:"type"`
		} `json:"buttons"`
	} `json:"remotes"`
}

var BUTTONS = map[byte]string{
	0x02: "on",
	0x04: "off",
	0x05: "up",
	0x06: "down",
	0x03: "select",
}

const DEBUG = false

var config RemoteConfig
var client *http.Client

func main() {
	f, _ := os.Open("remote-config.json")
	parseConfig(f)

	client = &http.Client{}

	os.Exit(spy(os.Stdin))
}

func parseConfig(input *os.File) {
	dec := json.NewDecoder(input)
	dec.Decode(&config)
	if DEBUG {
		fmt.Println(config)
	}
}

func spy(input *os.File) int {
	reader := bufio.NewReader(input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("EOF reached.")
			return 0
		}
		if len(line) > 0 {
			if DEBUG {
				fmt.Println(line)
			}
			serial, button, err := parseLine(line)
			if err == nil {
				handleButtonPress(serial, button)
			}
		}
	}
}

func parseLine(line string) (string, byte, error) {
	re := regexp.MustCompile("[(]([0-9a-f]+)[)]")
	matches := re.FindAllStringSubmatch(line, -1)
	if len(matches) > 0 {
		var buf []byte
		if DEBUG {
			fmt.Println(matches)
		}
		for _, s := range matches {
			temp, err := strconv.ParseUint(s[1], 16, 8)
			if err != nil {
				fmt.Printf("Fatal Error, unable to convert %s to byte\n", s[1])
				os.Exit(1)
			}
			buf = append(buf, byte(temp))
		}
		if DEBUG {
			fmt.Println(buf)
		}
		if buf[0] == 0x7e && buf[len(buf)-1] == 0x7e {
			return parseMessage(buf)
		}
		// TODO Support message split across lines
	}
	return "", 0x00, errors.New("Not a button press")
}

func parseMessage(message []byte) (string, byte, error) {
	// Message type is button press
	if len(message) > 5 && message[4] == 0x00 && message[5] == 0x81 {
		return fmt.Sprintf("%X%X%X", message[7], message[8], message[9]), message[11], nil
	}
	return "", 0x00, errors.New("Not a button press")
}

func handleButtonPress(serial string, button byte) {
	fmt.Println("serial: ", serial)
	fmt.Printf("button: %x\n", button)
	if val, ok := config.Remotes[serial]; ok {
		fmt.Println("nickname: ", val.Nickname)
		if val2, ok := val.Buttons[BUTTONS[button]]; ok {
			b, err := json.Marshal(val2.Data)
			if err != nil {
				fmt.Println("Fatal Error, data could not be remarshalled")
				os.Exit(1)
			}
			fmt.Printf("%s %s %s %s\n", val2.Method, string(b), val2.Type, val2.URL)

			req, err := http.NewRequest(val2.Method, val2.URL, bytes.NewBuffer(b))
			if err != nil {
				fmt.Println("Configured request is invalid")
			}

			req.Header.Set("Content-Type", val2.Type)

			_, err = client.Do(req)
			if err != nil {
				fmt.Println("The configured request failed")
			}
		}
	}
	fmt.Printf("\n")
}
