package lib

import (
	"errors"
	"strconv"
)

//
//func GenerateDeleteSqlQuery(field string) (string, error) {
//	var sqlQuery string
//	switch field {
//	case "rate":
//		sqlQuery = `
//			DELETE
//		`
//	}
//	return sqlQuery, nil
//}

func GenerateChangeSqlQuery(field, rawValue string) (string, any, error) {
	var value any
	sqlQuery := ``
	switch field {
	case "rate":
		intvalue, err := strconv.Atoi(rawValue)
		if err != nil {
			return "", "", err
		}
		sqlQuery = `
			UPDATE movies
			SET rate=$1
			WHERE id=$2
		`
		value = intvalue
	case "date":
		sqlQuery = `
			UPDATE movies
			SET date=$1
			WHERE id=$2
		`
		value = rawValue
	case "description":
		sqlQuery = `
			UPDATE movies
			SET description=$1
			WHERE id=$2
		`
		value = rawValue
	case "name":
		sqlQuery = `
			UPDATE movies
			SET name=$1
			WHERE id=$2
		`
		value = rawValue
	default:
		return "", "", errors.New("unknown field")
	}
	return sqlQuery, value, nil
}
