package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/codegp/cloud-persister"
)

var cp *cloudpersister.CloudPersister
var gameTypeID int64

func init() {
	var err error
	cp, err = cloudpersister.NewCloudPersister()
	if err != nil {
		panic(err)
	}

	gameTypeIDStr := os.Getenv("GAME_TYPE_ID")
	if gameTypeIDStr == "" {
		panic(fmt.Errorf("No gametype id found in environment"))
	}
	gameTypeID, err = strconv.ParseInt(gameTypeIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
}

func main() {
	api, err := ioutil.ReadFile("gen-html/api.html")
	fatalOnErr(err)
	ids, err := ioutil.ReadFile("gen-html/ids.html")
	fatalOnErr(err)
	fatalOnErr(cp.WriteDocs(fmt.Sprintf("docs-api-%d", gameTypeID), api))
	fatalOnErr(cp.WriteDocs(fmt.Sprintf("docs-ids-%d", gameTypeID), ids))
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}
