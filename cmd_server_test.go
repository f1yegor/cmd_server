package main

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestParseCommand(t *testing.T) {
	jsonStr := `["date"]`
	_, err := parseCommand(jsonStr)
	if err != nil {
		t.Fail()
	}

	jsonStr = `["sqlcmd", "-Slocalhost", "-U", "guest", "-P", "guest", "-d", "test_db", "-h", "-1", "-w", "65535", "-s^", "-Q", "SELECT [name], [server_timestamp], [id] FROM note_category", "-o", "test_sqlcmd_bcp.csv"]`
	_, err = parseCommand(jsonStr)
	if err != nil {
		t.Fail()
	}
}

func TestExecute(t *testing.T) {
	relPath := ""
	array, _ := parseCommand(`["sqlcmd", "-Slocalhost", "-U", "guest", "-P", "guest", "-d", "test_db", "-h", "-1", "-w", "65535", "-s^", "-Q", "SELECT [name], [server_timestamp], [id] FROM note_category", "-o", "test_sqlcmd_bcp.csv"]`)
	w := httptest.NewRecorder()

	execute(w, relPath, array...)
}

func TestEnsureDirectory(t *testing.T) {

	ensureDirectory("")

	ensureDirectory("./tenant1")

	ensureDirectory("./tenant1/project1")
}

func TestReadConfig(t *testing.T) {
	config := readConfig()
	fmt.Printf("config %q \n", config)
}
