package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/sorcix/irc"
	// "io"
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
	var m *irc.Message
	for {
		m, err = conn.Decode()
		if err != nil {
			log.Error("Could not read message: %s", err)
			break
		}
		log.Info("Got message: %+v", m)
		if m.Command == "PING" {
			m.Command = "PONG"
			err = conn.Encode(m)
			if err != nil {
				log.Error("Couldn't pong: %s", err)
			}
		}
	}
	conn.Close()
	log.Info("It worked!")
}
