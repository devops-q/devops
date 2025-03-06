package utils

import (
	"fmt"
	"os"
	"strconv"
)

func UpdateLatest(parsedCommandId string) {

	if parsedCommandIdInt, err := strconv.Atoi(parsedCommandId); err != nil {
		fmt.Println("Couldn't convert value to Integer")
		return
	} else {
		if parsedCommandIdInt != -1 {
			if err := os.WriteFile("internal/api/handlers/latest_processed_sim_action_id.txt", []byte(parsedCommandId), 0644); err != nil {
				fmt.Println("Error with file ${}")
			}

		}
	}

}
