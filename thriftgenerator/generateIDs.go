package main

import (
	"fmt"
	"strings"

	"github.com/codegp/game-object-types/types"
)

const (
	idsTemplateFile = "ids.thrift.tmpl"
	idsFileOutput   = "../thrift/ids.thrift"
	idTemplate      = "const i64 %s = %d;"
)

type idsTemplateData struct {
	BotTypeIDs     string
	AttackTypeIDs  string
	ItemTypeIDs    string
	TerrainTypeIDs string
	MoveTypeIDs    string
}

func generateIDs() error {
	botTypes, err := getBotTypes()
	if err != nil {
		return err
	}

	botIDs := getBotIDDefinitions(botTypes)
	itemIDs, err := getItemIDDefinitions()
	if err != nil {
		return err
	}

	terrainIDs, err := getTerrainIDDefinitions()
	if err != nil {
		return err
	}

	attackIDs, err := getAttackIDDefinitions(botTypes)
	if err != nil {
		return err
	}

	moveIDs, err := getMoveIDDefinitions(botTypes)
	if err != nil {
		return err
	}

	tmpl := idsTemplateData{
		BotTypeIDs:     botIDs,
		AttackTypeIDs:  attackIDs,
		ItemTypeIDs:    itemIDs,
		TerrainTypeIDs: terrainIDs,
		MoveTypeIDs:    moveIDs,
	}
	return generate(idsTemplateFile, idsFileOutput, tmpl)
}

func getIDDef(name string, ID int64) string {
	noSpaces := strings.Replace(name, " ", "_", -1)
	appendID := fmt.Sprintf("%s_ID", noSpaces)
	return fmt.Sprintf(idTemplate, appendID, ID)
}

func getBotTypes() ([]*types.BotType, error) {
	botTypes := []*types.BotType{}
	for _, ID := range gameType.BotTypes {
		botType, err := cp.GetBotType(ID)
		if err != nil {
			return botTypes, err
		}
		botTypes = append(botTypes, botType)
	}
	return botTypes, nil
}

func getBotIDDefinitions(botTypes []*types.BotType) string {
	names := []string{}
	for _, botType := range botTypes {
		names = append(names, getIDDef(botType.Name, botType.ID))
	}
	return strings.Join(names, "\n")
}

func getAttackIDDefinitions(botTypes []*types.BotType) (string, error) {
	names := []string{}
	found := map[int64]bool{}

	for _, botType := range botTypes {
		for _, attackTypeID := range botType.AttackTypeIDs {
			if _, exists := found[attackTypeID]; exists {
				continue
			}

			attackType, err := cp.GetAttackType(attackTypeID)
			if err != nil {
				return "", err
			}

			names = append(names, getIDDef(attackType.Name, attackTypeID))
			found[attackTypeID] = true
		}
	}
	return strings.Join(names, "\n"), nil
}

func getMoveIDDefinitions(botTypes []*types.BotType) (string, error) {
	names := []string{}
	found := map[int64]bool{}

	for _, botType := range botTypes {
		for _, moveTypeID := range botType.MoveTypeIDs {
			if _, exists := found[moveTypeID]; exists {
				continue
			}

			moveType, err := cp.GetMoveType(moveTypeID)
			if err != nil {
				return "", err
			}

			names = append(names, getIDDef(moveType.Name, moveTypeID))
			found[moveTypeID] = true
		}
	}
	return strings.Join(names, "\n"), nil
}

func getTerrainIDDefinitions() (string, error) {
	names := []string{}

	for _, terrainTypeID := range gameType.TerrainTypes {
		terrainType, err := cp.GetTerrainType(terrainTypeID)
		if err != nil {
			return "", err
		}

		names = append(names, getIDDef(terrainType.Name, terrainTypeID))
	}

	return strings.Join(names, "\n"), nil
}

func getItemIDDefinitions() (string, error) {
	names := []string{}

	for _, itemTypeID := range gameType.ItemTypes {
		itemType, err := cp.GetItemType(itemTypeID)
		if err != nil {
			return "", err
		}

		names = append(names, getIDDef(itemType.Name, itemTypeID))
	}

	return strings.Join(names, "\n"), nil
}
