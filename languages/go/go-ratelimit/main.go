package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func actualDelete(path string, olderThan time.Duration) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	// age := time.Since(info.ModTime())

	fmt.Printf("Removing %s\n", path)
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}

func deleteOldDirs(root string, olderThan time.Duration) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	var directories []string
	for _, file := range files {
		if file.IsDir() { // TODO: we can also remove this check
			directories = append(directories, filepath.Join(root, file.Name()))
		}
	}

	var wg sync.WaitGroup
	const concurrentLimit = 30
	sem := make(chan struct{}, concurrentLimit)

	for _, dir := range directories {
		sem <- struct{}{}
		wg.Add(1)

		// recreate loop variable before passing
		// to a goroutine
		dir := dir
		go func() {
			defer wg.Done()
			defer func() {
				// release the limiter
				<-sem
			}()
			err := actualDelete(dir, olderThan)
			if err != nil {
				fmt.Printf("error occured in deleting %s: %v \n", dir, err)
				return
			}
		}()
	}

	fmt.Println("waiting for the goroutines to delete the directories")
	wg.Wait()
	fmt.Println("deleted all directories")

	return nil
}

func removeWorkspaceDir(workspacePath string) error {
	olderThan := 14 * 24 * time.Hour
	return deleteOldDirs(workspacePath, olderThan)
}
func main() {
	agent := os.Args[1]
	if agent == "" {
		logrus.Fatalln("agent can't be empty")
		os.Exit(1)
		return
	}
	if agent == "agent23" { // TODO: rename the agent in Jenkins
		agent = "agent023"
	}
	rootPattern := fmt.Sprintf("/tmp/workspace/")
	if err := removeWorkspaceDir(rootPattern); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
		return
	}
	fmt.Println("Exiting fn")
}
