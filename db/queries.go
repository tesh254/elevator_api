package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

type Queries struct {
	database *sql.DB
}

type Elevator struct {
	ID          int    `json:"id"`
	BuildingID  int    `json:"building_id"`
	Direction   string `json:"direction"`
	Floor       int    `json:"floor"`
	Moving      bool   `json:"moving"`
	IsDoorsOpen bool   `json:"is_door_open"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (p *Queries) GetAllRecords(table string) ([]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := p.database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	structType := reflect.TypeOf(reflect.New(reflect.TypeOf(Elevator{})).Elem().Interface())
	results := []interface{}{}

	for rows.Next() {
		result := reflect.New(structType).Interface()
		columns := make([]interface{}, structType.NumField())
		for i := 0; i < structType.NumField(); i++ {
			columns[i] = reflect.ValueOf(result).Elem().Field(i).Addr().Interface()
		}
		if err := rows.Scan(columns...); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *Queries) GetById(table string, id int) (interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table)
	row := p.database.QueryRow(query, id)

	// Get the type of the struct associated with the table
	structType := reflect.TypeOf(reflect.New(reflect.TypeOf(Elevator{})).Elem().Interface())

	// Create a new instance of the struct
	result := reflect.New(structType).Interface()

	// Extract columns from the result interface
	columns := make([]interface{}, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		columns[i] = reflect.ValueOf(result).Elem().Field(i).Addr().Interface()
	}

	err := row.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Queries) UpdateById(table string, id int, jsonData []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	var updateColumns []string
	var updateValues []interface{}

	for key, value := range data {
		updateColumns = append(updateColumns, fmt.Sprintf("%s = $%d", key, len(updateColumns)+1))
		updateValues = append(updateValues, value)
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(updateColumns, ", "), len(updateColumns)+1)
	updateValues = append(updateValues, id)

	_, err := p.database.Exec(updateQuery, updateValues...)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated record
	return p.GetById(table, id)
}

func (p *Queries) Insert(table string, data map[string]interface{}, resultStruct interface{}) (interface{}, error) {
	keys := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for key, value := range data {
		keys = append(keys, key)
		values = append(values, value)
	}

	columns := strings.Join(keys, ", ")
	placeholders := make([]string, len(keys))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	valuesQuery := strings.Join(placeholders, ", ")

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *", table, columns, valuesQuery)

	structType := reflect.TypeOf(resultStruct)

	resultInstance := reflect.New(structType).Interface()

	columnsToScan := make([]interface{}, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		columnsToScan[i] = reflect.ValueOf(resultInstance).Elem().Field(i).Addr().Interface()
	}

	err := p.database.QueryRow(insertQuery, values...).Scan(columnsToScan...)
	if err != nil {
		return nil, err
	}

	return resultInstance, nil
}
