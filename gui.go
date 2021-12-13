package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
			var data []*instance
			json.Unmarshal(p[3:], &data)
			instances.saveSettings(true, data)
		default:
			fmt.Println("Invalid socket message received.")
		}
	}
}

func (instances *allInstances) saveSettings(fromGui bool, new []*instance) {
	if fromGui {
		instances.mu.Lock()

		if len(new) > len(instances.c) {
			newOne := len(new) - 1
			instances.c = append(instances.c, new[newOne])
			go instances.c[newOne].monitorFolder(false, instances)
			<-instances.readySend
		} else {
			for i := range instances.c {
				instances.c[i].Name = new[i].Name
				instances.c[i].Caption = new[i].Caption
				instances.c[i].StartupPostDelay = new[i].StartupPostDelay
				new[i].NextPostTime = instances.c[i].NextPostTime

				postDelayReset := false
				needsRestart := false
				if !sliceEqual(instances.c[i].ImgFolders, new[i].ImgFolders) {
					instances.c[i].ImgFolders = new[i].ImgFolders
					needsRestart = true
				}
				if instances.c[i].PostDelay != new[i].PostDelay {
					instances.c[i].PostDelay = new[i].PostDelay
					needsRestart = true
					postDelayReset = true
				}
				if instances.c[i].QueueDelay != new[i].QueueDelay {
					instances.c[i].QueueDelay = new[i].QueueDelay
					needsRestart = true
				}
				if needsRestart {
					instances.c[i].restartMonitoring <- 0
					go instances.c[i].monitorFolder(postDelayReset, instances)
					<-instances.readySend
				}
			}
		}
		instances.mu.Unlock()
	}

	// defining this type separately from main "instance" type allows us to avoid
	// saving certain fields to the json file, e.g. "ItemsInQueue" doesn't need
	// to be saved because it's generated on startup
	typeToSave := make([]struct {
		Name             string            `json:"Name"`
		ImgFolders       []string          `json:"ImgFolders"`
		QueueDelay       string            `json:"QueueDelay"`
		PostDelay        string            `json:"PostDelay"`
		StartupPostDelay string            `json:"StartupPostDelay"`
		NextPostTime     int64             `json:"NextPostTime"`
		Platforms        map[string]string `json:"Platforms"`
		Caption          string            `json:"Caption"`
	}, len(new))
	for i := range new {
		typeToSave[i].Name = new[i].Name
		typeToSave[i].ImgFolders = new[i].ImgFolders
		typeToSave[i].QueueDelay = new[i].QueueDelay
		typeToSave[i].PostDelay = new[i].PostDelay
		typeToSave[i].StartupPostDelay = new[i].StartupPostDelay
		typeToSave[i].NextPostTime = new[i].NextPostTime
		typeToSave[i].Platforms = new[i].Platforms
		typeToSave[i].Caption = new[i].Caption
	}
	dataToSave, err := json.MarshalIndent(typeToSave, "", "\t")
	if err != nil {
		log.Panic(err)
	}
	os.WriteFile("./userdata/offpost.json", dataToSave, 0666)
}

func sliceEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// ------------------------ end svelte gui functions ---------------------------
