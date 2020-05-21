package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"

	notepad "github.com/jchv/go-notepad"
	"golang.org/x/sys/windows"
)

func startNotepad() error {
	cmd := exec.Command("notepad")
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	writer, err := notepad.NewWriter()
	if err != nil {
		if err := startNotepad(); err != nil {
			windows.MessageBox(0, windows.StringToUTF16Ptr(fmt.Sprintf("Could not find an existing Notepad window, nor start a new Notepad instance: %s", err)), windows.StringToUTF16Ptr("Error starting Notepad logger"), windows.MB_ICONERROR)
			os.Exit(1)
		}
		time.Sleep(1 * time.Second)
		writer, err = notepad.NewWriter()
		if err != nil {
			windows.MessageBox(0, windows.StringToUTF16Ptr(fmt.Sprintf("Unable to find a notepad instance after starting one: %s.", err)), windows.StringToUTF16Ptr("Error starting Notepad logger"), windows.MB_ICONERROR)
			os.Exit(1)
		}
	}

	log.SetHandler(json.New(writer))
	for {
		log.Info("Logging to notepad is fun.")
		log.Warn("This is an awful idea.")
		log.Error("Nobody should ever use this.")
		time.Sleep(1 * time.Second)
	}
}
