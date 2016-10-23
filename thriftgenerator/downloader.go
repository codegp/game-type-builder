package main

import "io/ioutil"

const filename = "../../game-runner/gamerunner/gametypedef.go"

func downloadSource() error {
	bytes, err := cp.ReadGameTypeCode(gameType.ID)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}
