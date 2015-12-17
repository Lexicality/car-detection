package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func writeTempFile(tempFile *os.File, contents string) error {
	_, err := tempFile.WriteString(contents)
	if err != nil {
		return err
	}
	return tempFile.Close()
}

func TestGetNSPass(t *testing.T) {
	Convey("The Nickserv Stuff", t, func() {
		Convey("Disables itself if no file is present", func() {
			var temp = os.TempDir()
			var fname = "oaeuoaeuoaeu"
			ret, err := getNSPass(path.Join(temp, fname))
			So(err, ShouldBeNil)
			So(ret, ShouldBeBlank)
		})

		Convey("When passed a file", func() {
			tempFile, err := ioutil.TempFile("", "renspass")
			if err != nil {
				t.Fatal(err)
			}
			var fname = tempFile.Name()

			Convey("Returns the string in that file", func() {
				var str = "yo whatup"
				err := writeTempFile(tempFile, str)
				if err != nil {
					t.Fatal(err)
				}

				password, err := getNSPass(fname)
				So(err, ShouldBeNil)
				So(password, ShouldEqual, str)
			})

			Convey("Trims the contents of the file", func() {
				var str = "yo whatup"
				err := writeTempFile(tempFile, "\t\n  \t\r"+str+"\t\t\t\t\t \n")
				if err != nil {
					t.Fatal(err)
				}

				password, err := getNSPass(fname)
				So(err, ShouldBeNil)
				So(password, ShouldEqual, str)
			})

			Convey("Returns nothing if the file is empty", func() {
				password, err := getNSPass(fname)
				So(err, ShouldBeNil)
				So(password, ShouldBeBlank)
			})
		})
	})
}
