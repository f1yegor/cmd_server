package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/generate", CommandHandler)
	config := readConfig()
	fmt.Printf("Listen commands on %s \n", config.Port)

	http.ListenAndServe(":"+config.Port, nil)
}

func CommandHandler(w http.ResponseWriter, req *http.Request) {
	relPath := req.FormValue("path")
	fmt.Printf("Path %q \n", relPath)
	jsonStr := req.FormValue("cmd")
	fmt.Printf("Json %q \n", jsonStr)

	// command = `["BCP", "select top 555 * from issue_history",
	// "queryout", "temp_bcp.csv", "-c", "-t,", "-Slocalhost", "-U", "sa", "-P", "adm1n", "-d",
	// "lat_fs", "-a", "65535"]`

	array := make([]string, 0)
	err := json.Unmarshal([]byte(jsonStr), &array)
	if err != nil {
		fmt.Printf("Fail %q \n", err)
		return
	}
	//fmt.Printf("Request %q \n", array)

	Execute(w, relPath, array...)
}

func Execute(w http.ResponseWriter, relPath string, command ...string) {
	EnsureDirectory(relPath)

	prog, args := command[0], command[1:]
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Fprintf(w, "fail\n")
		fmt.Printf("Fail %q, out %s \n", err, out)
	} else {
		fmt.Fprintf(w, "ok\n")
		fmt.Printf("Result %s \n", out)
	}
}

func EnsureDirectory(relPath string) {
	if relPath == "" || relPath == "./" {
		return
	}

	err := os.MkdirAll(relPath, 0666)
	if err != nil {
		fmt.Printf("Fail dir %q \n", err)
	}
}

type Config struct {
	Port string // port listen to
}

func readConfig() Config {
	file, err := os.Open("settings.txt")
	if err != nil {
		fmt.Printf("Use default port 4000\n")
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
