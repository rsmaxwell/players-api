package basic

import (
	"database/sql"
	"strconv"
)

// My Interface
type My interface {
	Check(c *Condition) bool
}

// Condition type
type Condition struct {
	Field     string `json:"field"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

// WhereClause type
type WhereClause []Condition

// MyInt type
type MyInt struct {
	Value int
}

// MyNullString type
type MyNullString struct {
	Value sql.NullString
}

// Check function
func (v MyInt) Check(c *Condition) bool {

	value, err := strconv.Atoi(c.Value)
	if err != nil {
		return false
	}

	if c.Operation == "=" {
		if v.Value != value {
			return false
		}
	} else {
		return false
	}

	return true
}

// Check function
func (ns MyNullString) Check(c *Condition) bool {

	if c.Operation == "nil" {
		if ns.Value.Valid {
			return false
		}
	} else if c.Operation == "=" {
		if ns.Value.Valid {
			if ns.Value.String != c.Value {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}

	return true
}
