package ingestion

import (
	"encoding/csv"
	"os"
)

func ReadCSV(path string) ([]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, err
	}

	// First row is headers
	headers := records[0]
	var out []map[string]string

	for _, row := range records[1:] {
		rowMap := map[string]string{}
		for i, val := range row {
			if i < len(headers) {
				rowMap[headers[i]] = val
			}
		}
		out = append(out, rowMap)
	}

	return out, nil
}
