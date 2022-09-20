package psqldb

import (
	"pongpongi.com/osdep"
)

func DbQuery(sql string, args ...interface{}) []dbQueryType {
	results := []dbQueryType{}
	rows, err := db.Queryx(sql, args...)
	osdep.Check(err)
	for rows.Next() {
		result := make(map[string]interface{})
		osdep.Check(rows.MapScan(result))
		results = append(results, result)
	}
	return results
}
