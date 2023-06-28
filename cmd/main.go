package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

func GetFolderSize(folderPath string) (uint64, error) {
	var size uint64 = 0

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += uint64(info.Size())
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return size, nil
}

func DeleteNodeModules(folderPath string, daysThreshold int, wg *sync.WaitGroup, bytesSavedChan chan uint64) {
	defer wg.Done()

	parentInfo, err := os.Stat(filepath.Dir(folderPath))
	if err != nil {
		fmt.Println(err)
		return
	}

	lastModifiedHrs := time.Since(parentInfo.ModTime()).Abs().Hours()
	if lastModifiedHrs < float64(daysThreshold)*24 {
		return // too recently modified
	}

	size, err := GetFolderSize(folderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.RemoveAll(folderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Deleted %s, saving %s\n", folderPath, humanize.Bytes(size))

	savedTillNow := <-bytesSavedChan
	bytesSavedChan <- size + savedTillNow
}

func main() {
	basePath := flag.String("path", "", "Specify the base path")
	days := flag.Int("days", 0, "Specify the number of days")

	flag.Parse()

	if *basePath == "" || *days == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s -path <path> -days <days>\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	var wg sync.WaitGroup
	bytesSavedChan := make(chan uint64, 1)
	bytesSavedChan <- 0

	err := filepath.WalkDir(*basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == "node_modules" {
			wg.Add(1)
			go DeleteNodeModules(path, *days, &wg, bytesSavedChan)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	close(bytesSavedChan)

	totalSaved := <-bytesSavedChan

	fmt.Println("Done! Freed up", humanize.Bytes(totalSaved))
}
