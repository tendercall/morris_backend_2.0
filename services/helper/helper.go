package helper

import (
	"database/sql"
	"fmt"

	"morris-backend.com/main/services/models"
)

var DB *sql.DB

func GetPart() ([]models.Part, error) {
	rows, err := DB.Query("SELECT id, part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code FROM parts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var part models.Part
		err := rows.Scan(&part.ID, &part.PartNumber, &part.RemainPartNumber, &part.PartDescription, &part.FgWisonPartNumber, &part.SuperSSNumber, &part.Weight, &part.Coo, &part.HsCode)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	fmt.Println("Get Successful")

	return parts, nil
}
