package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/generate", commandHandler)
	http.Handle("/", gzipHandler(http.FileServer(http.Dir("."))))

	config := readConfig()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Listen commands on %s \n", config.Port)

	http.ListenAndServe(":"+config.Port, nil)
}

func commandHandler(w http.ResponseWriter, req *http.Request) {
	relPath := req.FormValue("path")
	log.Printf("Path %q \n", relPath)
	jsonStr := req.FormValue("cmd")

	array, err := parseCommand(jsonStr)
	if err != nil {
		failure(w, err)
		return
	}

	execute(w, relPath, array...)
}

func execute(w http.ResponseWriter, relPath string, command ...string) {
	ensureDirectory(relPath)

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

func ensureDirectory(relPath string) {
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

type config struct {
	Port string
}

func readConfig() config {
	file, err := os.Open("settings.txt")
	if err != nil {
		log.Printf("Use default port 4000\n")
		return config{Port: "4000"}
	}
	defer file.Close()
	config := config{}
	configScanner := bufio.NewScanner(file)
	if configScanner.Scan() {
		config.Port = configScanner.Text()
	}
	return config
}
