package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "localhost:8080"
const defaultChunkSize = 8
const defaultWait = 10

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "PANIC: %v\n", err)
			os.Exit(1)
		}
	}()

	t := newTask()
	flag.StringVar(&t.flags.d, "d", "response.dat", "response body")
	flag.StringVar(&t.flags.t, "t", "text/html; charset=UTF-8", "content type")
	flag.IntVar(&t.flags.c, "c", defaultChunkSize, "chunk size")
	flag.IntVar(&t.flags.w, "w", defaultWait, "chunk delay (ms)")
	flag.StringVar(&t.flags.s, "s", defaultPort, "listening server")
	flag.BoolVar(&t.flags.disable, "disable", false, "disable chunk mode")
	flag.Parse()

	if err := t.main(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

type flags struct {
	d       string
	t       string
	c       int
	w       int
	s       string
	disable bool
}

type task struct {
	Logger *log.Logger
	flags  flags
}

func fatalError(msg string, arg ...interface{}) error {
	return fmt.Errorf("%s: fatal error: %s", os.Args[0], fmt.Sprintf(msg, arg...))
}

const logFlag = log.Ldate | log.Ltime | log.Lshortfile

func newTask() *task {
	return &task{
		Logger: log.New(os.Stderr, "chunked-server ", logFlag),
	}
}

func (t *task) main() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Logger.Print("request /")

		w.Header().Set("Content-Type", t.flags.t)

		flusher, ok := w.(http.Flusher)
		if !ok {
			fatalError("expected http.ResponseWriter to be http.Flusher")
		}

		file, err := os.Open(t.flags.d)
		if err != nil {
			fatalError("%v\n", err)
		}
		defer file.Close()

		bufSize := t.flags.c
		if t.flags.disable {
			bufSize = 4096
		}

		buf := make([]byte, bufSize)
		for {
			n, err := file.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				fatalError("%v\n", err)
			}

			fmt.Fprintf(w, string(buf[:n]))

			if !t.flags.disable {
				flusher.Flush()
				time.Sleep(time.Duration(t.flags.w) * time.Millisecond)
			}
		}
	})

	t.Logger.Printf("Listening on %s", t.flags.s)
	t.Logger.Fatal(http.ListenAndServe(t.flags.s, nil))

	return nil
}
