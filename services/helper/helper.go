package helper

import (
	"database/sql"
	"fmt"
	"log"

	"morris-backend.com/main/services/models"
)

var DB *sql.DB

func PostPart(part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code string) (uint, error) {
	// Connect to the database
	var id uint

	err := DB.QueryRow("INSERT INTO parts (part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code).Scan(&id)
	if err != nil {
		return 0, err
	}

	fmt.Println("Post Successful")

	return id, nil
}

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

func GetPartByPartNumber(partNumber string) ([]models.Part, error) {
	var parts []models.Part

	var query string
	var args []interface{}

	if len(partNumber) >= 3 {
		// Check if the full part number is provided
		fullPartQuery := "SELECT id, part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code FROM parts WHERE part_number = $1"
		fullPartRows, err := DB.Query(fullPartQuery, partNumber)
		if err != nil {
			log.Println("Error executing full part query:", err)
			return nil, err
		}
		defer fullPartRows.Close()

		// Check for exact match
		for fullPartRows.Next() {
			var part models.Part
			err := fullPartRows.Scan(&part.ID, &part.PartNumber, &part.RemainPartNumber, &part.PartDescription, &part.FgWisonPartNumber, &part.SuperSSNumber, &part.Weight, &part.Coo, &part.HsCode)
			if err != nil {
				log.Println("Error scanning row:", err)
				return nil, err
			}
			parts = append(parts, part)
		}

		if len(parts) > 0 {
			// If exact match is found, return these results
			return parts, nil
		}

		// If no exact match, perform prefix search
		query = "SELECT id, part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code FROM parts WHERE part_number LIKE $1"
		args = append(args, partNumber[:3]+"%")
	} else {
		// Handle as an exact match if partNumber is shorter than 3 characters
		query = "SELECT id, part_number, remain_part_number, part_description, fg_wison_part_number, super_ss_number, weight, coo, hs_code FROM parts WHERE part_number = $1"
		args = append(args, partNumber)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var part models.Part
		err := rows.Scan(&part.ID, &part.PartNumber, &part.RemainPartNumber, &part.PartDescription, &part.FgWisonPartNumber, &part.SuperSSNumber, &part.Weight, &part.Coo, &part.HsCode)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		parts = append(parts, part)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error processing rows:", err)
		return nil, err
	}

	return parts, nil
}
