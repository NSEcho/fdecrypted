package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/frida/frida-go/frida"
)

//go:embed script.js
var sc string

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s DATA\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s Gadget\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s Gadget:[BDL]:some_filename.json\n", os.Args[0])
		os.Exit(1)
	}

	var target, filename, targetDir string
	isFile := false

	if strings.Contains(os.Args[1], ":") {
		splitted := strings.Split(os.Args[1], ":")
		if len(splitted) != 3 {
			fmt.Fprintln(os.Stderr, "filename should be in format APP_NAME:DIR:FILENAME")
			os.Exit(1)
		}
		isFile = true
		target = splitted[0]
		targetDir = splitted[1]
		filename = splitted[2]
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if isFile {
		err := script.ExportsCallWithContext(ctx, "download_file", targetDir, filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching file")
			os.Exit(1)
		}
	} else {
		err := script.ExportsCallWithContext(ctx, "download_bin")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching file")
			os.Exit(1)
		}
	}
	<-done
	fmt.Printf("[*] Saved \"%s\" (%d bytes)\n", name, length)
}
