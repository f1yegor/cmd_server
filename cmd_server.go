package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"net/http"
	"encoding/json"
)

func main() {
	http.HandleFunc("/generate", CommandHandler)
    http.ListenAndServe(":4000", nil)
}

func CommandHandler(w http.ResponseWriter, req *http.Request) {
	// command := []string{req.FormValue("cmd")}
	// fmt.Printf("Request %q \n", command)
 	command := []string{"BCP", "select top 555 * from issue_history", 
	"queryout", "D:\\sync_dumps\\temp_bcp.csv", "-c", "-t,", "-Slocalhost", "-U", "sa", "-P", "adm1n", "-d", "lat_fs", "-a", "65535"}

    go Execute(w, command...)
}

func Execute(w http.ResponseWriter, command ...string) {
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

