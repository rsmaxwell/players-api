package model

import (
	"fmt"
	"strings"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionBuild = debug.NewFunction(pkg, "Build")
)

// Condition type
type Condition struct {
	Operation string      `json:"operation"`
	Value     interface{} `json:"value"`
}

// Query type
type Query = map[string]Condition

var (
	supportedOperations = make(map[string]string)
)

func init() {
	supportedOperations["="] = "="
	supportedOperations[">"] = ">"
	supportedOperations["<"] = "<"
	supportedOperations[">="] = ">="
	supportedOperations["<="] = "<="
	supportedOperations["<>"] = "<>"
	supportedOperations["BETWEEN"] = "BETWEEN"
	supportedOperations["LIKE"] = "LIKE"
	supportedOperations["IN"] = "IN"
}

// CheckOperation function
func checkOperation(operation string) (string, error) {
	if val, ok := supportedOperations[operation]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Operation not supported: '%s'", operation)
}

// BuildQuery function
func BuildQuery(table string, allFields []string, returnedFields []string, q *Query) (string, []interface{}, error) {
	f := functionBuild

	var values []interface{}
	sqlStatement := "SELECT " + strings.Join(returnedFields, ", ") + " FROM " + table

	if q != nil {
		count := 0
		var where []string
		for _, k := range allFields {
			if c, ok := (*q)[k]; ok {
				count++
				values = append(values, c.Value)
				op, err := checkOperation(c.Operation)
				if err != nil {
					f.DumpError(err, "query string is not valid")
					return "", nil, err
				}
				where = append(where, fmt.Sprintf("%s%s$%d", k, op, count))
			}
		}

		if count > 0 {
			sqlStatement = sqlStatement + " WHERE " + strings.Join(where, " AND ")
		}
	}

	return sqlStatement, values, nil
}
