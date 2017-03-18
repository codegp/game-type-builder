package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codegp/cloud-persister"
)

func main() {
	cp, err := cloudpersister.NewCloudPersister()
	if err != nil {
		log.Fatalf("Failed to start cloud persister: %v", err)
	}

	projIDStr := os.Getenv("PROJECT_ID")
	if projIDStr == "" {
		log.Fatalf("No project id provided in environment")
	}

	projID, err := strconv.ParseInt(projIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Failed to cast projID, %v", err)
	}

	project, err := cp.GetProject(projID)
	if err != nil {
		log.Fatalf("Failed to get project %d, %v", projID, err)
	}

	if project.Directory == "" {
		log.Println("Directory is set, not downloading source.")
		return
	}

	for _, filename := range project.FileNames {
		bytes, err := cp.ReadProjectFile(projID, filename)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		tmpfn := filepath.Join("/source", filename)
		if err = ioutil.WriteFile(tmpfn, bytes, 0666); err != nil {
			log.Fatalf("Error writing file to disk: %v", err)
		}
	}
}
