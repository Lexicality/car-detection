package main

import (
	"fmt"
	// "github.com/sorcix/irc"
	"io/ioutil"
	"os"
	"strings"
)

const (
	nsFileName = "nspass.txt"
)

func getNSPass() (string, error) {
	_, err := os.Stat(nsFileName)
	// No NS password
	if err != nil {
		log.Warning("No NickServ password specified. Please make a file called %s with your bot's password or blank to disable.", nsFileName)
		return "", nil
	}
	contents, err := ioutil.ReadFile(nsFileName)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}

func setupNickserv(conn *Session) (err error) {
	pass, err := getNSPass()
	if err != nil {
		return fmt.Errorf("Could not load NS password: %s", err)
	} else if pass == "" {
		log.Info("Nickserv disabled")
		return nil
	}
	return conn.Privmsg("nickserv", fmt.Sprintf("IDENTIFY %s", pass))
}
