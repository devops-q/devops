package utils

import (
	"fmt"
	"os"
	"strconv"
)


func UpdateLatest(parsed_command_id string) {

	if parsed_command_id_int, err := strconv.Atoi(parsed_command_id); err != nil {
		fmt.Println("Couldn't convert value to Integer")
		return
	} else {
	if parsed_command_id_int != -1 {
		if err := os.WriteFile("./latest_processed_sim_action_id.txt", []byte(parsed_command_id),0644); err != nil {
			fmt.Println("Error with file ${}")
		}

	}
	}
	



}