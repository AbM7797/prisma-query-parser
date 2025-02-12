package parser

import (
	"net/url"
	"strconv"
	"strings"
)

// processArgs processes query parameters into a structured filter map
func ProcessArgs(args url.Values) Filter {
	where := make(Filter)
	orderBy := make(Filter)
	skip := -1
	take := -1

	for key, values := range args {
		if len(values) == 0 {
			continue
		}

		value := values[0]

		// Handle WHERE clauses
		if strings.HasPrefix(key, "where[") {
			decodedKey := strings.ReplaceAll(key, `\]`, "]")
			decodedKey = strings.ReplaceAll(decodedKey, `\[`, "[")
			field := decodedKey[6 : len(decodedKey)-1]
			field = strings.ReplaceAll(field, "][", ".")
			setNestedValue(where, field, convertInput(value, "", nil))
			// Handle ORDER BY
		} else if strings.HasPrefix(key, "orderBy[") {
			decodedKey := strings.ReplaceAll(key, `\\]`, "]")
			decodedKey = strings.ReplaceAll(decodedKey, `\\[`, "[")
			field := decodedKey[8 : len(decodedKey)-1]
			field = strings.ReplaceAll(field, "][", ".")

			// Parse the value as db.ASC or db.DESC
			var direction interface{}
			if strings.EqualFold(value, "asc") {
				direction = "asc"
			} else if strings.EqualFold(value, "desc") {
				direction = "desc"
			} else {
				direction = value // Keep as-is for unexpected values
			}

			setNestedValue(orderBy, field, direction)

			// Handle SKIP
		} else if key == "skip" {
			if i, err := strconv.Atoi(value); err == nil {
				skip = i
			} else {
				skip = -1
			}

			// Handle TAKE
		} else if key == "take" {
			if i, err := strconv.Atoi(value); err == nil {
				take = i
			} else {
				take = -1
			}
		}
	}

	// Handle OR condition
	if orValue, exists := where["OR"]; exists {
		if orMap, ok := orValue.(map[string]interface{}); ok {
			orArray := []map[string]interface{}{}
			for _, v := range orMap {
				if condition, ok := v.(map[string]interface{}); ok {
					orArray = append(orArray, condition)
				}
			}
			where["OR"] = orArray
		}
	}

	// Handle NOT condition
	if notValue, exists := where["NOT"]; exists {
		if notMap, ok := notValue.(map[string]interface{}); ok {
			notArray := []map[string]interface{}{}
			for _, v := range notMap {
				if condition, ok := v.(map[string]interface{}); ok {
					notArray = append(notArray, condition)
				}
			}
			where["NOT"] = notArray
		}
	}

	// Final structured output
	result := make(Filter)
	if len(where) > 0 {
		result["where"] = where
	}
	if len(orderBy) > 0 {
		result["orderBy"] = orderBy
	}
	if skip >= 0 {
		result["skip"] = skip
	}
	if take >= 0 {
		result["take"] = take
	}

	return result
}
