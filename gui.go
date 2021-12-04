package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/gorilla/websocket"
)

// --------------------------------tray functions-------------------------------
// open() is a cross-platform way to open URLs in the default browser
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

// onReady() defines settings for the tray icon
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

// --------------------------- end tray functions ------------------------------

// -------------------------- svelte gui functions -----------------------------
func (instances *allInstances) createWebSocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _ := upgrader.Upgrade(w, r, nil)
	clientClosed := make(chan bool)

	instances.mu.Lock()
	guiOpen = true
	instances.mu.Unlock()
	fmt.Print("GUI opening\n\n")

	wsSend <- ""

	// sends to the GUI, when wsSend is fed a string
	go func() {
		for {
			select {
			case <-wsSend:
				err := conn.WriteJSON(instances.c)
				if err != nil {
					fmt.Println(err)
				}
			case <-clientClosed:
				// exit writer function when ReadMessage says client is closed
				return
			}

		}
	}()

	// reads from the GUI
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			instances.mu.Lock()
			// ReadMessage returns err when client closes
			clientClosed <- true
			fmt.Print("GUI closing\n\n")
			guiOpen = false
			instances.mu.Unlock()
			return
		}

		// this switch is used to identify the incoming message by the first 3
		// characters in the string
		switch string(p)[:3] {
		case "s, ":
			fmt.Println("we received a message to save the settings")
			fmt.Println(string(p)[3:])
			var data []*instance
			json.Unmarshal(p[3:], &data)
			instances.saveGuiSettings(data)
		default:
			fmt.Println("Invalid socket message received.")
		}

		// for i := range data {
		// 	// fmt.Println(instances.c[i].Caption)
		// 	fmt.Println(data[i].Caption)
		// 	instances.c[i].Caption = data[i].Caption
		// }

		fmt.Println("---")
		// fmt.Println("Received:", string(p))
	}
}

func (instances *allInstances) saveGuiSettings(newInstances []*instance) {
	for i := range instances.c {
		instances.c[i].Caption = newInstances[i].Caption
	}
}

// ------------------------ end svelte gui functions ---------------------------
