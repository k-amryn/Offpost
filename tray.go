package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"
)

// open is a cross-platform way to open URLs in the default browser
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func onReady() {
	ico, _ := ioutil.ReadFile("favicon.ico")

	systray.SetTemplateIcon(ico, ico)
	systray.SetTitle("Offpost")
	systray.SetTooltip("Offpost")
	mURL := systray.AddMenuItem("Open UI", "my home")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mURL.ClickedCh:
				open("http://localhost:8081/")
			case <-mQuit.ClickedCh:
				log.Panic("User clicked quit button.")
			}
		}
	}()
}

func onExit() {
	fmt.Println(time.Now())
	// ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}
