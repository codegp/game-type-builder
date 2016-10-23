package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/codegp/cloud-persister"
	"github.com/codegp/cloud-persister/models"
)

/*
  This file generates a thrift file so thrift can generate gametype constants, such meta
*/
var cp *cloudpersister.CloudPersister
var gameType *models.GameType

func init() {
	var err error
	cp, err = cloudpersister.NewCloudPersister()
	if err != nil {
		panic(err)
	}

	gameType, err = getGameType()
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := generateAPI(); err != nil {
		panic(err)
	}
	if err := generateIDs(); err != nil {
		panic(err)
	}
	if err := downloadSource(); err != nil {
		panic(err)
	}
}

func generate(templatePath string, outputPath string, templateData interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	defer func() {
		writer.Flush()
		f.Close()
	}()

	return tmpl.Execute(writer, templateData)
}

func getGameType() (*models.GameType, error) {
	gameTypeIDStr := os.Getenv("GAME_TYPE_ID")
	if gameTypeIDStr == "" {
		return nil, fmt.Errorf("No gametype id found in environment")
	}

	gameTypeID, err := strconv.ParseInt(gameTypeIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return cp.GetGameType(gameTypeID)
}
