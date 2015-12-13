package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/sorcix/irc"
	"io"
	"io/ioutil"
	"strings"
)

const (
	server   = "irc.mindfang.org"
	port     = 6667
	username = "pcc31"
	nickname = "randomTesting"
)

var log = logging.MustGetLogger("Random Encounters")

func getNSPass() (string, error) {
	contents, err := ioutil.ReadFile("nspass")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}

func handshake(conn *irc.Conn) (err error) {
	err = conn.Encode(&irc.Message{
		Command: "NICK",
		Params:  []string{nickname},
	})
	if err != nil {
		return err
	}
	return conn.Encode(&irc.Message{
		Command:  "USER",
		Params:   []string{username, "0", "*"},
		Trailing: "Hello world",
	})
}

func handlePing(conn *irc.Conn, message *irc.Message) (err error) {
	// Ha ha ha this is so dodgy
	message.Command = "PONG"
	return conn.Encode(message)
}

func readPump(conn *irc.Conn) (err error) {
	var toIgnore = [...]string{"001", "002", "003", "005", "251", "252", "254", "255", "265", "266"}
	var shouldIgnore = make(map[string]bool, len(toIgnore))
	for _, num := range toIgnore {
		shouldIgnore[num] = true
	}
	var m *irc.Message

	for {
		m, err = conn.Decode()
		if err != nil {
			return err
		}
		// Ignore informative spam
		if shouldIgnore[m.Command] {
			continue
		}
		// Start of glorious message type switches
		if m.Command == "PING" {
			err = handlePing(conn, m)
			if err != nil {
				log.Error("Couldn't pong: %s", err)
			}
		} else {
			log.Info("Got message: %+v", m)
		}
	}
}

func main() {
	fmt.Println("Hello world")
	// nspass, err := getNSPass()
	// if err != nil {
	// 	log.Fatalf("Unable to get Nickserv password: %s", err)
	// }
	// db stuff probably
	conn, err := irc.Dial(fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		log.Fatalf("Could not connect to the server: %s", err)
	}
	// say hello
	err = handshake(conn)
	if err != nil {
		conn.Close()
		log.Fatalf("Could not handshake: %s", err)
	}

	// Blocking
	err = readPump(conn)

	if err == io.EOF {
		log.Info("Server hung up")
	} else if err != nil {
		log.Critical("Could not read message: %s", err)
	}

	conn.Close()
}
