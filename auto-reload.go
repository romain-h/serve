package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type autoReload struct {
	fs    http.FileSystem
	dir   string
	port  int
	index string
}

var re = regexp.MustCompile("</body>")

func (al autoReload) OpenIndex(name string) (http.File, error) {
	// If the file doesn't exist, create it, or append to the file
	orig, _ := ioutil.ReadFile(al.dir + name)
	// append socketScript to body
	withSocket := re.ReplaceAll(orig, []byte(fmt.Sprintf(SocketScript, al.port)+"</body>"))

	if err := ioutil.WriteFile(al.dir+"/.hrc/index.html", withSocket, 0644); err != nil {
		log.Fatal(err)
	}

	return al.fs.Open("/.hrc/index.html")
}

func (al autoReload) Open(name string) (http.File, error) {
	if name == fmt.Sprintf("/%s", al.index) {
		return al.OpenIndex(name)
	}
	return al.fs.Open(name)
}
