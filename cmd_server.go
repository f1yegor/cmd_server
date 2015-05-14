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
	http.HandleFunc("/generate", CommandHandler)
	config := readConfig()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Listen commands on %s \n", config.Port)

	http.ListenAndServe(":"+config.Port, nil)
}

func CommandHandler(w http.ResponseWriter, req *http.Request) {
	relPath := req.FormValue("path")
	log.Printf("Path %q \n", relPath)
	jsonStr := req.FormValue("cmd")

	// bcp variant
	// jsonStr = `["BCP", "select top 555 * from issue_history",
	// "queryout", "temp_bcp.csv", "-c", "-t,", "-Slocalhost", "-U", "sa", "-P", "adm1n", "-d",
	// "lat_fs", "-a", "65535"]`
	// sqlcmd variant
	// jsonStr = `["sqlcmd", "-Slocalhost", "-U", "sa", "-P", "adm1n", "-d", "lat_fs_test", "-h", "-1", "-w", "65535", "-s#", "-Q", "SELECT * FROM note_category", "-o", "test_sqlcmd_bcp.csv"]`

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
	err := cmd.Run()

	if err != nil {
		log.Printf("Fail execute %q, out %s \n", err, out)
		failure(w, err)

	} else {
		fmt.Fprintf(w, "ok\n")
		log.Printf("Result %s \n", out)
	}
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
