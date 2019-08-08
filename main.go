package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var port = flag.Int("port", 8000, "port")
var localOnly = flag.Bool("local-only", false, "limit access to localhost")

func GetIntParam(u *url.URL, name string) (int, error) {
	params, ok := u.Query()[name]
	if !ok || len(params) == 0 {
		return 0, errors.New("not found")
	}
	return strconv.Atoi(params[len(params)-1])
}

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(args) > 0 {
		if filepath.IsAbs(args[0]) {
			dir = args[0]
		} else {
			dir = filepath.Join(dir, args[0])
		}
	}

	fs := http.FileServer(http.Dir(dir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		delay := 0
		if d, err := GetIntParam(r.URL, "delay"); err == nil {
			delay = d
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.Header().Set("Cache-Control", "no-cache")
		fs.ServeHTTP(w, r)
	})

	log.Println("Listening...", *port, dir)
	addr := fmt.Sprintf(":%d", *port)
	if *localOnly {
		addr = "localhost" + addr
	}
	log.Fatal(http.ListenAndServe(addr, nil))
}
