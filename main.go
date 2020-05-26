package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	port := flag.Int("p", 3009, "port")
	dir := flag.String("d", ".", "directory to serve")
	index := flag.String("i", "index.html", "HTML index entry")

	flag.Parse()

	c := make(chan os.Signal, 1)
	tmpDir := *dir + "/.hrc"
	os.Mkdir(tmpDir, os.ModePerm)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}()

	fs := autoReload{
		fs:    http.Dir(*dir),
		dir:   *dir,
		port:  *port,
		index: *index,
	}
	buildDone := make(chan bool, 1)

	http.Handle("/", http.FileServer(fs))
	http.Handle("/websocket", getSocketHandler(buildDone))

	watcher := watchDir(*dir, buildDone)
	defer watcher.Close()

	log.Println(fmt.Sprintf("Listening on :%d, serving %s with %s as HTML index", *port, *dir, *index))
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
