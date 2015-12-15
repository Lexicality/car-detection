package main

import (
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"strings"
)

var log = logging.MustGetLogger("Random Encounters")

func getNSPass() (string, error) {
	contents, err := ioutil.ReadFile("nspass")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}

func main() {
	var err error

	var session = &Session{
		Server:   "irc.mindfang.org",
		Port:     6667,
		UserName: "pcc31",
		NickName: "randomTesting",
		RealName: "HI MOPM",
	}

	err = session.Dial()
	if err != nil {
		log.Fatalf("Could not connect to the server: %s", err)
	}
	defer session.Close()

	// Blocking
	err = session.readPump()

	if err == io.EOF {
		log.Info("Server hung up")
	} else if err != nil {
		log.Critical("Could not read message: %s", err)
	}
}
