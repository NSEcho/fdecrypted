package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frida/frida-go/frida"
)

//go:embed script.js
var sc string

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s DATA\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s Gadget\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s Gadget:some_filename.json\n", os.Args[0])
		os.Exit(1)
	}

	var target, filename string
	isFile := false

	if strings.Contains(os.Args[1], ":") {
		splitted := strings.Split(os.Args[1], ":")
		if len(splitted) != 2 {
			fmt.Fprintln(os.Stderr, "filename should be in format APP_NAME:FILENAME")
			os.Exit(1)
		}
		isFile = true
		target = splitted[0]
		filename = splitted[1]
	} else {
		target = os.Args[1]
	}

	d := frida.USBDevice()
	session, err := d.Attach(target, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error attaching to target: %v\n", err)
		os.Exit(1)
	}

	script, err := session.CreateScript(sc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating script: %v\n", err)
		os.Exit(1)
	}

	var name string
	var length int
	done := make(chan struct{})

	script.On("message", func(message string, data []byte) {
		if len(data) > 0 {
			unmarshalled := make(map[string]string)
			json.Unmarshal([]byte(message), &unmarshalled)

			name = filepath.Base(unmarshalled["payload"])
			length = len(data)
			if err := os.WriteFile(name, data, os.ModePerm); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving binary: %v\n", err)
				os.Exit(1)
			}
			go func() {
				done <- struct{}{}
			}()
		}
	})
	script.Load()

	if isFile {
		script.ExportsCall("download_file", filename)
	} else {
		script.ExportsCall("download_bin")
	}
	<-done
	fmt.Printf("[*] Saved \"%s\" (%d bytes)\n", name, length)
}
