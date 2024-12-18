package app

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Labellable interface {
	GetId() int
	GetLabel() string
	GetValue() string
	GetDescription() string
}

func StringifyHardSkillTypesIntoQueryValues[T Labellable](items []T) string {
	valuesArray := make([]string, len(items))

	for i, item := range items {
		valuesArray[i] = fmt.Sprintf("('%s', '%s', '%s')",
			item.GetLabel(),
			item.GetValue(),
			item.GetDescription(),
		)

	}
	return strings.ReplaceAll(strings.Trim(strings.Join(valuesArray, ","), ""), " ", "")
}

func StringifyHardSkillTypeValueIntoQueryValues[T Labellable](items []T) string {
	valuesArray := make([]string, len(items))

	for i, item := range items {
		valuesArray[i] = fmt.Sprintf("'%s'", item.GetValue())
	}
	return strings.Trim(strings.Join(valuesArray, ","), "")
}

func createReplaceQuery[T Labellable](tableName string, items []T) string {
	valuePlaceholder := StringifyHardSkillTypesIntoQueryValues(items)
	query := fmt.Sprintf("REPLACE INTO %s (label, value, description) VALUES %s", tableName, valuePlaceholder)
	fmt.Println("query", query)
	return query
}

func executeReplaceQuery(query string) (sql.Result, error) {
	db := DatabaseConnection()
	defer db.Close()
	results, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	log.Print("RowsAffected:")
	log.Print(results.RowsAffected())
	return results, nil
}

func createReadQuery[T Labellable](tableName string, items []T) string {
	valuePlaceholder := StringifyHardSkillTypeValueIntoQueryValues(items)
	query := fmt.Sprintf("SELECT id, label, value, description FROM %s WHERE value IN (%s)", tableName, valuePlaceholder)
	fmt.Println("query", query)
	return query
}

func CreateAndExecuteReadQuery[T Labellable](tableName string, items []T) (*sql.Rows, error) {
	query := createReadQuery(tableName, items)
	rows, err := DatabaseConnection().Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	return rows, nil
}

func CreateAndExecuteReplaceQuery[T Labellable](tableName string, items []T) (sql.Result, error) {
	query := createReplaceQuery(tableName, items)
	return executeReplaceQuery(query)
}
