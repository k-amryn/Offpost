package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
)

type instance struct {
	Name              string            `json:"Name"`
	ImgFolders        []string          `json:"ImgFolders"`
	QueueDelay        string            `json:"QueueDelay"`
	PostDelay         string            `json:"PostDelay"`
	StartupPostDelay  string            `json:"StartupPostDelay"`
	NextPostTime      int64             `json:"NextPostTime"`
	Platforms         map[string]string `json:"Platforms"`
	Caption           string            `json:"Caption"`
	ItemsInQueue      int               `json:"ItemsInQueue"`
	restartMonitoring chan int
}

type allInstances struct {
	c         []*instance
	readySend chan int
	authComm  chan string
	mu        sync.Mutex
}

var wsSend = make(chan string, 1)
var guiOpen = false

func loadInstances() []*instance {
	jsonBlob, err := ioutil.ReadFile("./userdata/offpost.json")
	if err != nil {
		log.Panic("offpost.json not found.")
	}
	instancesRaw := make([]*instance, 0)
	_ = json.Unmarshal(jsonBlob, &instancesRaw)
	return instancesRaw
}

func (instance *instance) initQueue() {
	for _, folder := range instance.ImgFolders {
		// read files from folder
		fileStatus := make(map[string]string) // filename:p for posted, filename:q for queued
		files, _ := ioutil.ReadDir(folder)

		for _, file := range files {
			dotIndex := strings.LastIndex(file.Name(), ".")
			// if the directory entry is a folder
			if dotIndex == -1 {
				continue
			}

			filepath := folder + "/" + file.Name()

			filetype := file.Name()[dotIndex:]
			if filetype == ".jpg" || filetype == ".png" || filetype == ".webp" || filetype == ".txt" || filetype == ".mp4" {
				fileStatus[filepath] = "n"
			}
		}

		readFromFile := instance.readTxtFile("queue", false)
		for _, val := range readFromFile {
			fileStatus[val[0]] = "q"
		}

		readFromFile = instance.readTxtFile("posted", false)
		for _, val := range readFromFile {
			fileStatus[val[0]] = "p"
		}

		var newQueue [][]string
		for key, value := range fileStatus {
			if value == "n" {
				newQueue = append(newQueue, []string{key})
			}
		}

		newQueue = groupOrganize(newQueue)
		instance.appendTxtFile(newQueue, "queue")
	}
	instance.ItemsInQueue = instance.countQueueItems()
}

// instance.QueueDelay and instance.PostDelay are stored as strings,
// func queueDelay() and func postDelay() convert string to time duration
func (instance *instance) queueDelay() time.Duration {
	return processTime(instance.QueueDelay)
}

func (instance *instance) postDelay() time.Duration {
	return processTime(instance.PostDelay)
}

func processTime(stringTime string) time.Duration {
	dur, _ := time.ParseDuration(stringTime)
	return dur
}

func (instance *instance) countQueueItems() int {
	r, err := os.Open("./userdata/" + instance.Name + "_queue.txt")
	if err != nil {
		return 0
	}

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			return count
		}
	}
}

func numberAtEnd(filename string) int {
	filename = filename[:strings.LastIndex(filename, ".")]
	rgx, err := regexp.Compile(`-\d{1,3}$`)
	if err != nil {
		return -1
	}
	num := rgx.FindString(filename)
	if num == "" {
		return -1
	}
	num = num[1:]
	numInt, err := strconv.Atoi(string(num))
	if err != nil {
		return -1
	}
	return numInt
}

// get the filename minus end numbers and file extension
func getBaseName(filename string) string {
	slashIndex := strings.LastIndex(filename, "/")
	if slashIndex == -1 {
		filename = filename[:strings.LastIndex(filename, ".")]
	} else {
		filename = filename[slashIndex+1 : strings.LastIndex(filename, ".")]
	}
	rgx, err := regexp.Compile(`-\d{1,3}$`)
	if err != nil {
		return ""
	}
	index := rgx.FindStringIndex(filename)
	if index == nil {
		return filename
	}
	return filename[:index[0]]
}

func groupOrganize(shortQueue [][]string) [][]string {
	lookingFor := 1
	var foundBaseName string
	var foundIndex int
	// checkedOne stores whether items ending in 1 have been checked, to prevent re-checking
	checkedOne := make(map[int]bool)
	for i := 0; i < len(shortQueue); i++ {
		if lookingFor == 1 && !checkedOne[i] {
			if numberAtEnd(shortQueue[i][0]) == lookingFor &&
				len(shortQueue[i]) == 1 {
				foundBaseName = getBaseName(shortQueue[i][0])
				foundIndex = i

				checkedOne[i] = true

				i = -1
				lookingFor++
			}
		} else {
			if getBaseName(shortQueue[i][0]) == foundBaseName &&
				numberAtEnd(shortQueue[i][0]) == lookingFor &&
				len(shortQueue[i]) == 1 &&
				!checkedOne[i] {
				shortQueue[foundIndex] = append(shortQueue[foundIndex], shortQueue[i][0])
				shortQueue[i][0] = "x."

				i = -1
				lookingFor++
			} else if i == len(shortQueue)-1 && !checkedOne[i] {
				i = -1
				lookingFor = 1
			}
		}
	}
	for i := 0; i < len(shortQueue); i++ {
		if shortQueue[i][0] == "x." {
			copy(shortQueue[i:], shortQueue[i+1:])
			shortQueue = shortQueue[:len(shortQueue)-1]
			i--
		}
	}
	return shortQueue
}

func (instance *instance) isQueueEmpty() bool {
	queueInfo, err := os.Stat("./userdata/" + instance.Name + "_queue.txt")
	if err != nil {
		log.Panic(instance.Name, "_queue.txt not found")
	}
	return queueInfo.Size() == 0
}

func (instance *instance) readTxtFile(queueOrPost string, grouped bool) [][]string {
	f, err := os.OpenFile("./userdata/"+instance.Name+"_"+queueOrPost+".txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// read the whole file into a variable
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("error:", err)
	}

	// put each line in the file into an array item
	lines := strings.Split(string(data), "\n")
	var readFromFile [][]string
	for i := range lines {
		// skip empty lines
		if lines[i] == "" {
			continue
		}

		// Windows txt lines end with \r\n, lines retaining \r cause issues
		re := regexp.MustCompile(`\r`)
		lines[i] = re.ReplaceAllString(lines[i], "")

		if grouped {
			readFromFile = append(readFromFile, strings.Split(lines[i], "***"))
		} else {
			for _, val := range strings.Split(lines[i], "***") {
				readFromFile = append(readFromFile, []string{val})
			}
		}

	}
	return readFromFile
}

func (instance *instance) appendTxtFile(shortQueue [][]string, queueOrPost string) {
	f, err := os.OpenFile("./userdata/"+instance.Name+"_"+queueOrPost+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for i := range shortQueue {
		line := strings.Join(shortQueue[i], "***")
		f.WriteString(line + "\n")
	}
}

// monitors folders, manages queueing and posting, *allInstances is being sent
// here to access its Mutex to prevent data races
func (instance *instance) monitorFolder(postDelayReset bool, all *allInstances) {
	instance.restartMonitoring = make(chan int)
	instance.initQueue()
	timeToExit := false
	exitChan := make(chan int)

	rename := -1

	queueTimer := time.NewTimer(instance.queueDelay())

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, folder := range instance.ImgFolders {
		err = watcher.Add(folder)
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		var shortQueue [][]string

		queueTimer.Stop()

		var postTimer *time.Timer
		var postTimerCheck *time.Timer
		guiSend := make(chan string, 1)

		// if NextPostTime has passed, then schedule a new NextPostTime
		if time.Now().UnixMilli() > instance.NextPostTime || postDelayReset {
			if time.Now().UnixMilli() > instance.NextPostTime {
				fmt.Println(instance.Name + ": Scheduled post time has passed. Scheduling new post time.")
			} else {
				fmt.Println(instance.Name + ": Post Delay reset. Scheduling new post time.")
			}
			switch instance.StartupPostDelay {
			case "random":
				// my funny seed algorithm, mixes part of instance name with current time
				fiveLetters := fmt.Sprint([]byte(instance.Name[:5]))
				fiveLetters = strings.ReplaceAll(fiveLetters, " ", "")
				fiveLetters = strings.ReplaceAll(fiveLetters, "[", "")
				fiveLetters = strings.ReplaceAll(fiveLetters, "]", "")
				fiveLetters64, _ := strconv.Atoi(fiveLetters)
				rand.Seed(time.Now().UnixNano() + int64(fiveLetters64))
				//-----seed algorithm finished-----

				randSecondsStr := fmt.Sprint(rand.Intn(int(instance.postDelay().Seconds())))
				randSecondsDur, _ := time.ParseDuration(randSecondsStr + "s")
				postTimer = time.NewTimer(randSecondsDur)
				// postTimerCheck allows timer.Stop check without stopping main timer
				postTimerCheck = time.NewTimer(randSecondsDur)
				instance.NextPostTime = time.Now().Add(randSecondsDur).UnixMilli()
			case "full":
				postTimer = time.NewTimer(instance.postDelay())
				postTimerCheck = time.NewTimer(instance.postDelay())
				instance.NextPostTime = time.Now().Add(instance.postDelay()).UnixMilli()
			case "none": // StartupPostDelay = "none"
				postTimer = time.NewTimer(0 * time.Second)
				postTimerCheck = time.NewTimer(0 * time.Second)
				instance.NextPostTime = time.Now().UnixMilli()
			default:
				log.Panic("StartupPostDelay value must be \"random\", \"full\", or \"none\". Check your offpost.json.")
			}
		} else { // use the saved NextPostTime
			fmt.Println(instance.Name + ": Using scheduled post time.")
			timeUntilNextPost := time.Until(time.UnixMilli(instance.NextPostTime))
			postTimer = time.NewTimer(timeUntilNextPost)
			postTimerCheck = time.NewTimer(timeUntilNextPost)
		}

		guiSend <- ""
		all.readySend <- 0

		for !timeToExit {

			// this select waits for either a new file, or a shortQueue timer to expire
			select {
			case event := <-watcher.Events:

				// event.Name returns full filepath
				filename := strings.ReplaceAll(event.Name, "\\", "/") //[strings.LastIndex(event.Name, "\\")+1:]
				switch event.Op.String() {
				case "CREATE":
					// rename is set to the index of the renamed file in shortQueue, is -1 when there is no rename
					if rename != -1 {
						fmt.Printf("%v renamed to ", shortQueue[rename][0])
						shortQueue[rename][0] = filename
						fmt.Printf("%v\n", shortQueue[rename][0])
						rename = -1
						fmt.Printf("Current Images: %v\n\n", shortQueue)
						break
					}

					queueTimer.Stop()
					// fmt.Printf("%v timer stopped\n", instance.Name)
					queueTimer.Reset(instance.queueDelay())
					// fmt.Printf("%v timer reset\n", instance.Name)

					fmt.Printf("%v New image found\n", instance.Name)

					filetype := filename[strings.LastIndex(filename, "."):]

					if filetype == ".jpg" || filetype == ".png" || filetype == ".webp" || filetype == ".txt" || filetype == ".mp4" {
						shortQueue = append(shortQueue, []string{filename})
					}

					fmt.Printf("%v Current Images: %v\n\n", instance.Name, shortQueue)
				case "RENAME":
					for i := range shortQueue {
						if shortQueue[i][0] == filename {
							rename = i
							// cannot perform the rename here because a CREATE event
							// with the new filename occurs after the RENAME event
							break
						}
					}
				case "REMOVE":
					for i := range shortQueue {
						if shortQueue[i][0] == filename {
							fmt.Printf("\nremoving %v\n", filename)
							copy(shortQueue[i:], shortQueue[i+1:])
							shortQueue = shortQueue[:len(shortQueue)-1]
							fmt.Printf("%v Current Images: %v\n\n", instance.Name, shortQueue)
							break
						}
					}
					if len(shortQueue) != 0 {
						queueTimer.Stop()
						fmt.Printf("%v timer stopped\n", instance.Name)
						queueTimer.Reset(instance.queueDelay())
						fmt.Printf("%v timer reset\n", instance.Name)
					}
				}
				// end of event switch

			case <-queueTimer.C:
				all.mu.Lock()
				if len(shortQueue) != 0 {
					shortQueue = groupOrganize(shortQueue)

					isEmpty := instance.isQueueEmpty()

					instance.appendTxtFile(shortQueue, "queue")
					instance.ItemsInQueue = instance.countQueueItems()
					fmt.Printf("%v queued the Current Images, %v items in queue\n\n", instance.Name, instance.ItemsInQueue)

					if isEmpty && !postTimerCheck.Stop() {
						instance.makePost()
						postTimer.Stop()
						postTimer.Reset(instance.postDelay())
						postTimerCheck.Stop()
						postTimerCheck.Reset(instance.postDelay())
						instance.NextPostTime = time.Now().Add(instance.postDelay()).UnixMilli()
						all.saveSettings(false, all.c)
					}

					guiSend <- ""

					queueTimer = time.NewTimer(instance.queueDelay())
					queueTimer.Stop()
					shortQueue = [][]string{}

				}
				all.mu.Unlock()
			case <-postTimer.C:
				all.mu.Lock()
				if !instance.isQueueEmpty() {
					instance.makePost()
					postTimer.Stop()
					postTimer.Reset(instance.postDelay())
					postTimerCheck.Stop()
					postTimerCheck.Reset(instance.postDelay())
					instance.NextPostTime = time.Now().Add(instance.postDelay()).UnixMilli()
					all.saveSettings(false, all.c)

					guiSend <- ""
					all.mu.Unlock()

					break
				}
				fmt.Printf("%v Post Timer done, but the queue is empty.\nNext thing added to queue will be posted immediately.\n\n", instance.Name)
				all.mu.Unlock()
			case <-guiSend:
				all.mu.Lock()
				if guiOpen {
					wsSend <- ""
				}
				all.mu.Unlock()
			case <-exitChan:
				timeToExit = true

			} // end of timer switch
		}
	}()

	// send to exit chan to close watcher goroutine
	why := <-instance.restartMonitoring

	queueTimer.Stop()
	queueTimer.Reset(0 * time.Second)

	exitChan <- why
	if why == 0 {
		fmt.Println(instance.Name + ": refreshing folder monitoring.")
	} else {
		fmt.Println(instance.Name + ": deleting.\n")
	}
}

func main() {
	instances := allInstances{c: loadInstances()}
	fmt.Print(` _______  _______  _______  _______  _______  _______ _________
(  ___  )(  ____ \(  ____ \(  ____ )(  ___  )(  ____ \\__   __/
| (   ) || (    \/| (    \/| (    )|| (   ) || (    \/   ) (
| |   | || (__    | (__    | (____)|| |   | || (_____    | |
| |   | ||  __)   |  __)   |  _____)| |   | |(_____  )   | |
| |   | || (      | (      | (      | |   | |      ) |   | |
| (___) || )      | )      | )      | (___) |/\____) |   | |
(_______)|/       |/       |/       (_______)\_______)   )_(
by coekuss

offpost.json settings loaded.

`)

	fmt.Println("Your instances:")
	// for i := range instances {
	// 	fmt.Printf("  - %v\n", instances[i].name)
	// }

	for i := range instances.c {
		fmt.Println(instances.c[i].Name, "-", instances.c[i].ImgFolders[0])
		if len(instances.c[i].ImgFolders) > 1 {
			for i2 := 1; i2 < len(instances.c[i].ImgFolders); i2++ {
				for range instances.c[i].Name {
					fmt.Print(" ")
				}
				fmt.Println(" -", instances.c[i].ImgFolders[i2])
			}
		}
		fmt.Print("\n")
	}

	fmt.Println("\nType anything to start working.")

	// reader := bufio.NewReader(os.Stdin)
	// _, _ = reader.ReadString('\n')
	go systray.Run(onReady, onExit)

	fmt.Println("Initializing post queue and monitoring your folders.")
	fmt.Println("GUI on http://localhost:14859/")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Print("\n")

	// readySend ensures the instance data is gathered before sending to GUI
	instances.readySend = make(chan int)
	for i := range instances.c {
		go instances.c[i].monitorFolder(false, &instances)
		<-instances.readySend
	}
	instances.refreshUsernames()
	fmt.Print("\n")
	// saving at this point to save rescheduled post times
	instances.saveSettings(false, instances.c)

	go instances.handleWebServer()

	// <-stayOpen makes the program stay open since no value is sent to the channel
	stayOpen := make(chan int)
	<-stayOpen
}
