package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type instance struct {
	name               string
	imgFolders         []string
	timeToQueue        time.Duration
	postInterval       time.Duration
	postDelayAtStartup string
	platforms          map[string]string
}

type postQ [][]string

func processTime(stringTime string) time.Duration {
	// turn "10 minutes" into ["10"], ["minutes"]
	splitted := strings.Split(stringTime, " ")

	quantity, _ := strconv.Atoi(splitted[0])

	if splitted[1] == "seconds" || splitted[1] == "second" {
		return time.Duration(quantity) * time.Second
	} else if splitted[1] == "minutes" || splitted[1] == "minute" {
		return time.Duration(quantity) * time.Minute
	} else if splitted[1] == "hours" || splitted[1] == "hour" {
		return time.Duration(quantity) * time.Hour
	}

	return time.Minute
}

func loadInstances() map[string]*instance {
	jsonBlob, err := ioutil.ReadFile("offpost.json")
	if err != nil {
		log.Panic("offpost.json not found.")
	}
	instancesRaw := make(map[string]struct {
		ImgFolders         []string
		TimeToQueue        string
		PostInterval       string
		PostDelayAtStartup string
		Platforms          map[string]string
	})
	_ = json.Unmarshal(jsonBlob, &instancesRaw)

	realInstances := make(map[string]*instance)
	for key := range instancesRaw {
		realInstances[key] = &instance{
			name:               key,
			imgFolders:         instancesRaw[key].ImgFolders,
			timeToQueue:        processTime(instancesRaw[key].TimeToQueue),
			postInterval:       processTime(instancesRaw[key].PostInterval),
			postDelayAtStartup: instancesRaw[key].PostDelayAtStartup,
			platforms:          instancesRaw[key].Platforms}
	}

	return realInstances
}

func numberAtEnd(filename string) int {
	filename = filename[:strings.LastIndex(filename, ".")]
	rgx, err := regexp.Compile("-\\d{1,3}$")
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
	filename = filename[:strings.LastIndex(filename, ".")]
	rgx, err := regexp.Compile("-\\d{1,3}$")
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
	queueInfo, _ := os.Stat(instance.name + "_queue.txt")
	if queueInfo.Size() == 0 {
		return true
	}
	return false
}

func (instance instance) readTxtFile(queueOrPost string, grouped bool) [][]string {
	f, err := os.OpenFile(instance.name+"_"+queueOrPost+".txt", os.O_RDONLY|os.O_CREATE, 0666)
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

func (instance *instance) initQueue() {

	for _, folder := range instance.imgFolders {
		// read files from folder
		var readFromFolder [][]string
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
				readFromFolder = append(readFromFolder,
					[]string{filepath})
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
}

func (instance instance) appendTxtFile(shortQueue [][]string, queueOrPost string) {
	f, err := os.OpenFile(instance.name+"_"+queueOrPost+".txt", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for i := range shortQueue {
		line := strings.Join(shortQueue[i], "***")
		f.WriteString(line + "\n")
	}
}

func (instance *instance) monitorFolder() {
	// monitor the folder and add new images to the queue
	// fmt.Printf("\n%v LongQueue: %v\n", instance.name, instance.initQueue())
	instance.initQueue()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	rename := -1

	done := make(chan int)
	go func() {
		var shortQueue [][]string

		queueTimer := time.NewTimer(instance.timeToQueue)
		queueTimer.Stop()

		var postTimer *time.Timer
		var postTimerCheck *time.Timer
		switch instance.postDelayAtStartup {
		case "random":
			rand.Seed(time.Now().UnixNano())
			randSecondsStr := fmt.Sprint(rand.Intn(int(instance.postInterval.Seconds())))
			randSecondsDur, _ := time.ParseDuration(randSecondsStr + "s")
			postTimer = time.NewTimer(randSecondsDur)
			// postTimerCheck allows timer.Stop check without stopping main timer
			postTimerCheck = time.NewTimer(instance.postInterval)
		case "full":
			postTimer = time.NewTimer(instance.postInterval)
			postTimerCheck = time.NewTimer(instance.postInterval)
		default:
			postTimer = time.NewTimer(0 * time.Second)
			postTimerCheck = time.NewTimer(0 * time.Second)
		}

		for {

			// this select waits for either a new file, or a shortQueue timer to expire
			select {
			case event, _ := <-watcher.Events:
				tweetlink := "Tweet error/n"

				// event.Name returns full filepath, this isolates the filename
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
					// fmt.Printf("%v timer stopped\n", instance.name)
					queueTimer.Reset(instance.timeToQueue)
					// fmt.Printf("%v timer reset\n", instance.name)

					fmt.Printf("%v New image found\n", instance.name)

					filetype := filename[strings.LastIndex(filename, "."):]

					if filetype == ".jpg" || filetype == ".png" || filetype == ".webp" || filetype == ".txt" || filetype == ".mp4" {
						shortQueue = append(shortQueue, []string{filename})
					}

					fmt.Printf("%v Current Images: %v\n\n", instance.name, shortQueue)
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
							fmt.Printf("%v Current Images: %v\n\n", instance.name, shortQueue)
							break
						}
					}
					if len(shortQueue) != 0 {
						queueTimer.Stop()
						fmt.Printf("%v timer stopped\n", instance.name)
						queueTimer.Reset(instance.timeToQueue)
						fmt.Printf("%v timer reset\n", instance.name)
					}
				} // end of event switch

				if tweetlink == "Tweet error\n" {
					fmt.Println(tweetlink)
				}
			case <-queueTimer.C:
				if len(shortQueue) != 0 {
					shortQueue = groupOrganize(shortQueue)

					isEmpty := instance.isQueueEmpty()

					instance.appendTxtFile(shortQueue, "queue")
					fmt.Printf("%v queued the Current Images\n", instance.name)

					if isEmpty && !postTimerCheck.Stop() {
						instance.makePost()
						postTimer.Stop()
						postTimer.Reset(instance.postInterval)
						postTimerCheck.Stop()
						postTimerCheck.Reset(instance.postInterval)
					}

					queueTimer = time.NewTimer(instance.timeToQueue)
					queueTimer.Stop()
					shortQueue = [][]string{}

				}
			case <-postTimer.C:
				if !instance.isQueueEmpty() {
					instance.makePost()
					postTimer.Stop()
					postTimer.Reset(instance.postInterval)
					postTimerCheck.Stop()
					postTimerCheck.Reset(instance.postInterval)
					break
				}
				fmt.Printf("%v Post Timer done, but the queue is empty.\nNext thing added to queue will be posted immediately.\n", instance.name)

			} // end of timer switch

		}
	}()

	for _, folder := range instance.imgFolders {
		err = watcher.Add(folder)
	}

	<-done
}

func main() {
	instances := loadInstances()
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

	for key := range instances {
		fmt.Println(instances[key].name, "-", instances[key].imgFolders[0])
		if len(instances[key].imgFolders) > 1 {
			for i := 1; i < len(instances[key].imgFolders); i++ {
				for _ = range instances[key].name {
					fmt.Print(" ")
				}
				fmt.Println(" -", instances[key].imgFolders[i])
			}
		}
		fmt.Print("\n")
	}

	fmt.Println("\nType anything to start working.")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	fmt.Println("Initializing post queue and monitoring your folders.")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Print("\n")
	for i := range instances {
		go instances[i].monitorFolder()
	}

	time.Sleep(10000 * time.Hour)
}
