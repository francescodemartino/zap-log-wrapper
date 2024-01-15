package writer_custom_log

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WriterCustomLog struct {
	muxLogFiles      *sync.RWMutex
	lifespanLogFile  time.Duration
	lifetimeLogFiles time.Duration
	printInConsole   bool
	basePathLog      string
	baseNameLog      string
	file             *os.File
}

func New(basePathLog string, baseNameLog string, printInConsole bool, lifespanLogFile time.Duration, lifetimeLogFiles time.Duration) *WriterCustomLog {
	if !strings.HasSuffix(basePathLog, "/") {
		basePathLog = basePathLog + "/"
	}
	wcl := &WriterCustomLog{
		muxLogFiles:      &sync.RWMutex{},
		lifespanLogFile:  lifespanLogFile,
		lifetimeLogFiles: lifetimeLogFiles,
		printInConsole:   printInConsole,
		basePathLog:      basePathLog,
		baseNameLog:      baseNameLog,
		file:             createLogFile(basePathLog, baseNameLog),
	}
	go wcl.substituteLogFile()
	go wcl.deleteOldLogFiles()
	return wcl
}

func (w *WriterCustomLog) Write(p []byte) (n int, err error) {
	if w.printInConsole {
		os.Stdout.Write(p)
	}
	w.muxLogFiles.RLock()
	nWrite, errWrite := w.file.Write(p)
	w.muxLogFiles.RUnlock()
	return nWrite, errWrite
}

func (w *WriterCustomLog) Sync() error {
	return w.file.Sync()
}

func (w *WriterCustomLog) substituteLogFile() {
	for {
		time.Sleep(w.lifespanLogFile)
		w.muxLogFiles.Lock()
		w.file.Sync()
		w.file.Close()
		w.file = createLogFile(w.basePathLog, w.baseNameLog)
		w.muxLogFiles.Unlock()
	}
}

func (w *WriterCustomLog) deleteOldLogFiles() {
	for {
		now := time.Now()
		filesLog, err := os.ReadDir(w.basePathLog)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, fileEntry := range filesLog {
			if !fileEntry.IsDir() && strings.HasSuffix(fileEntry.Name(), ".log") {
				nameFile := fileEntry.Name()
				unixTimeFile, err := strconv.ParseInt(nameFile[strings.LastIndex(nameFile, "_")+1:len(nameFile)-4], 10, 64)
				if err != nil {
					continue
				}
				dateFile := time.Unix(unixTimeFile, 0)
				dateFile = dateFile.Add(w.lifetimeLogFiles)
				if dateFile.Before(now) {
					os.Remove(w.basePathLog + nameFile)
				}
			}
		}
		time.Sleep(w.lifespanLogFile)
	}
}

func createLogFile(basePathLog string, baseNameLog string) *os.File {
	file, _ := os.Create(basePathLog + baseNameLog + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".log")
	return file
}
