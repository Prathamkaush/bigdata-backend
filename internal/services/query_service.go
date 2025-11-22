package services

import (
	"fmt"
	"strings"

	"bigdata-api/internal/models"
)

var allowedFields = map[string]bool{
	"dedupe_key": true, "source": true, "source_id": true, "customer_id": true,
	"first_name": true, "last_name": true, "email": true, "phone": true,
	"city": true, "state": true, "country": true,
	"ingest_ts": true, "record_ts": true,
}

func BuildSelectQuery(
	selectCols []string,
	filters map[string]interface{},
	ranges map[string]models.RangeFilter,
	fuzzy map[string]models.FuzzyFilter,
	sort string,
	limit, offset int,
) (string, []interface{}) {

	where := []string{"1=1"}
	args := []interface{}{}

	// --- exact filters ---
	for field, value := range filters {
		if allowedFields[field] {
			where = append(where, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
	}

	// --- range filters ---
	for field, r := range ranges {
		if allowedFields[field] {
			if r.From != "" {
				where = append(where, fmt.Sprintf("%s >= ?", field))
				args = append(args, r.From)
			}
			if r.To != "" {
				where = append(where, fmt.Sprintf("%s <= ?", field))
				args = append(args, r.To)
			}
		}
	}

	// --- fuzzy filters ---
	for field, f := range fuzzy {
		if allowedFields[field] && f.Query != "" {
			where = append(where, fmt.Sprintf("lower(%s) LIKE lower(?)", field))
			args = append(args, "%"+f.Query+"%")
		}
	}

	cols := strings.Join(selectCols, ", ")
	query := fmt.Sprintf(
		"SELECT %s FROM default.master_records WHERE %s",
		cols,
		strings.Join(where, " AND "),
	)

	// safe sorting
	if sort != "" {
		parts := strings.Fields(sort)
		if len(parts) >= 1 && allowedFields[parts[0]] {
			dir := "ASC"
			if len(parts) == 2 && strings.ToUpper(parts[1]) == "DESC" {
				dir = "DESC"
			}
			query = fmt.Sprintf("%s ORDER BY %s %s", query, parts[0], dir)
		}
	}

	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
	return query, args
}

func BuildCountQuery(
	filters map[string]interface{},
	ranges map[string]models.RangeFilter,
	fuzzy map[string]models.FuzzyFilter,
) (string, []interface{}) {

	where := []string{"1=1"}
	args := []interface{}{}

	for field, value := range filters {
		if allowedFields[field] {
			where = append(where, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
	}

	for field, r := range ranges {
		if allowedFields[field] {
			if r.From != "" {
				where = append(where, fmt.Sprintf("%s >= ?", field))
				args = append(args, r.From)
			}
			if r.To != "" {
				where = append(where, fmt.Sprintf("%s <= ?", field))
				args = append(args, r.To)
			}
		}
	}

	for field, f := range fuzzy {
		if allowedFields[field] && f.Query != "" {
			where = append(where, fmt.Sprintf("lower(%s) LIKE lower(?)", field))
			args = append(args, "%"+f.Query+"%")
		}
	}

	query := fmt.Sprintf(
		"SELECT count() FROM default.master_records WHERE %s",
		strings.Join(where, " AND "),
	)

	return query, args
}
