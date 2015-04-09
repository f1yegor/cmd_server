package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"net/http"
	"encoding/json"
)

func main() {
	http.HandleFunc("/generate", CommandHandler)
    http.ListenAndServe(":4000", nil)
}

func CommandHandler(w http.ResponseWriter, req *http.Request) {
	relPath := req.FormValue("path")
	fmt.Printf("Path %q \n", relPath)
	jsonStr := req.FormValue("cmd")
	fmt.Printf("Json %q \n", jsonStr)

	// command = `["BCP", "select top 555 * from issue_history", 
		// "queryout", "D:\\sync_dumps\\temp_bcp.csv", "-c", "-t,", "-Slocalhost", "-U", "sa", "-P", "adm1n", "-d", 
		// "lat_fs", "-a", "65535"]`

	array := make([]string, 0)
	err := json.Unmarshal([]byte(jsonStr), &array)
	if err != nil {
		fmt.Printf("Fail %q \n", err)	
		return
	}
	fmt.Printf("Request %q \n", array)
 	
    go Execute(w, relPath, array...)
}

func Execute(w http.ResponseWriter, relPath string, command ...string) {
	
	EnsureDirectory(relPath)
	//os.Chdir(relPath); defer os.Chdir(".")

	prog, args := command[0], command[1:]
	cmd := exec.Command(prog, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	
	if err != nil {
		fmt.Fprintf(w, "fail")
		fmt.Printf("Fail %q, out %s \n", err, out)
	} else {
		fmt.Fprintf(w, "ok")
		fmt.Printf("Result %s \n", out)
	}
}

func EnsureDirectory(relPath string) {
	err := os.MkdirAll(relPath, 0666)
	if err != nil {
		fmt.Printf("Fail dir %q \n", err)	
	}
}


