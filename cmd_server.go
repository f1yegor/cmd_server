package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/generate", CommandHandler)
	http.Handle("/", Gzip(http.FileServer(http.Dir("."))))

	config := readConfig()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Listen commands on %s \n", config.Port)

	http.ListenAndServe(":"+config.Port, nil)
}

func CommandHandler(w http.ResponseWriter, req *http.Request) {
	relPath := req.FormValue("path")
	log.Printf("Path %q \n", relPath)
	jsonStr := req.FormValue("cmd")

	array, err := parseCommand(jsonStr)
	if err != nil {
		failure(w, err)
		return
	}

	Execute(w, relPath, array...)
}

func Execute(w http.ResponseWriter, relPath string, command ...string) {
	EnsureDirectory(relPath)

	prog, args := command[0], command[1:]
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if err != nil {
		log.Printf("Fail execute [%q], out [%s] \n", err, printOut(&out))
		failure(w, err)

	} else {
		fmt.Fprintf(w, "ok\n")
		log.Printf("Result [%s] \n", printOut(&out))
	}
}

func printOut(buf *bytes.Buffer) string {
	return buf.String()
}

func EnsureDirectory(relPath string) {
	if relPath == "" || relPath == "./" {
		return
	}

	err := os.MkdirAll(relPath, 0666)
	if err != nil {
		log.Printf("Fail dir %q \n", err)
	}
}

func parseCommand(jsonStr string) ([]string, error) {
	log.Printf("Json %q \n", jsonStr)

	array := make([]string, 0)
	err := json.Unmarshal([]byte(jsonStr), &array)

	if err != nil {
		log.Printf("Fail %q with %s \n", err, jsonStr)
	}

	return array, err
}

func failure(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	fmt.Fprintf(w, "fail\n")
}

type Config struct {
	Port string // port listen to
}

func readConfig() Config {
	file, err := os.Open("settings.txt")
	if err != nil {
		log.Printf("Use default port 4000\n")
		return Config{Port: "4000"}
	}
	defer file.Close()
	config := Config{}
	configScanner := bufio.NewScanner(file)
	if configScanner.Scan() {
		config.Port = configScanner.Text()
	}
	return config
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		handler.ServeHTTP(gzw, r)
	})
}
