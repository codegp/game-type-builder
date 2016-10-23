package main

import (
	"fmt"
	"strings"
)

/*
  This file generates a thrift file so thrift can generate gametype constants, such meta
*/

const (
	apiTemplateFile = "api.thrift.tmpl"
	apiFileOutput   = "../thrift/api.thrift"
)

type apiTemplateData struct {
	Signatures string
}

func generateAPI() error {
	sigMap := signatureMap()
	signatures := []string{}
	for _, apiFunc := range gameType.ApiFuncs {
		sig, exists := sigMap[apiFunc]
		if !exists {
			return fmt.Errorf("No signature exists for function %s", apiFunc)
		}
		signatures = append(signatures, sig)
	}

	tmpl := apiTemplateData{
		Signatures: strings.Join(signatures, ",\n"),
	}
	return generate(apiTemplateFile, apiFileOutput, tmpl)
}

func signatureMap() map[string]string {
	return map[string]string{
		"me":            "gameObjects.Bot me()",
		"spawn":         "gameObjects.Bot spawn(1:gameObjects.Direction d, 2:i64 botTypeID)",
		"canMove":       "bool canMove(1:gameObjects.Direction d, 2:i64 moveTypeID)",
		"move":          "void move(1:gameObjects.Direction d, 2:i64 moveTypeID) throws (1: gameObjects.InvalidMove err)",
		"canAttack":     "bool canAttack(1:gameObjects.Location l, 2:i64 attackTypeID)",
		"attack":        "void attack(1:gameObjects.Location l, 2:i64 attackTypeID) throws (1: gameObjects.InvalidMove err)",
		"botAtLocation": "gameObjects.Bot botAtLocation(1: gameObjects.Location l) throws (1: gameObjects.InvalidMove err)",
	}
}
