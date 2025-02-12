package parser

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// mapOperatorToPrismaMethod dynamically maps the operator and field to the corresponding Prisma method.
func mapOperatorToPrismaMethod[T any, C any, K any](model K, field, operator, value string, returnType T, where C, tm TypeMapper) (T, error) {
	// Use reflection to find the method for the given field in the model
	_, err := getFieldByJSONTag(model, field)
	if err != nil {
		log.Printf("Field %s not found in model %T, skipping...\n", field, model)
		return returnType, err
	}

	// Get the method name dynamically
	field, operator = getPrismaMethodName(field, operator)
	log.Println("field", field, operator, value)
	// Use reflection to find the method by name
	if value == "<nil>" {
		method := reflect.ValueOf(where).FieldByName(field).MethodByName("IsNull")
		if method.IsValid() {
			// Call the method dynamically and pass the value
			return method.Call([]reflect.Value{})[0].Interface().(T), nil
		}
	}
	method := reflect.ValueOf(where).FieldByName(field).MethodByName(operator)
	if method.IsValid() {
		// Call the method dynamically and pass the value
		return method.Call([]reflect.Value{reflect.ValueOf(convertInput(value, operator, tm))})[0].Interface().(T), nil
	}

	log.Printf("Method %s not found for field %s\n", operator, field)
	return returnType, nil
}

func getFieldTypeByName(structType interface{}, fieldName string) (reflect.Type, error) {
	// Get the reflection type of the struct
	structReflectType := reflect.TypeOf(structType)
	// Ensure the input is a pointer to a struct (to allow FieldByName to work)
	if structReflectType.Kind() == reflect.Ptr {
		structReflectType = structReflectType.Elem()
	}

	// Get the field by name
	field, found := structReflectType.FieldByName(fieldName)
	if !found {
		return nil, fmt.Errorf("field '%s' not found", fieldName)
	}

	// Return the type of the field
	return field.Type, nil
}

// convertInput converts a string value into its appropriate type (bool, nil, etc.)
func convertInput(value, operator string, tm TypeMapper) interface{} {
	lower := strings.ToLower(value)
	switch lower {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return nil
	case "asc":
		if operator == "Order" {
			return tm.GetASC()
		}
		return value
	case "desc":
		if operator == "Order" {
			return tm.GetDESC()
		}
		return value
	default:
		if operator == "Mode" {
			return tm.GetMode(value)
		} else if operator != "" {
			return tm.GetValue(value)
		}
		return value // Default return as string
	}
}

// setNestedValue sets a nested value in a map based on a dot-separated path
func setNestedValue(obj Filter, path string, value interface{}) {
	keys := strings.Split(path, ".")
	currentObj := obj

	for i, key := range keys {
		isLastKey := i == len(keys)-1
		if _, exists := currentObj[key]; !exists {
			if isLastKey {
				currentObj[key] = value
			} else {
				currentObj[key] = make(Filter)
			}
		}

		if isLastKey {
			if existing, ok := currentObj[key].([]interface{}); ok {
				currentObj[key] = append(existing, value)
			} else {
				currentObj[key] = value
			}
		} else {
			if nested, ok := currentObj[key].(Filter); ok {
				currentObj = nested
			} else {
				newObj := make(Filter)
				currentObj[key] = newObj
				currentObj = newObj
			}
		}
	}
}

// getFieldByJSONTag retrieves the field by JSON tag name from the struct
func getFieldByJSONTag(model interface{}, jsonTag string) (reflect.Value, error) {
	// Skip checking for "OR" and "NOT" operators
	if jsonTag == "OR" || jsonTag == "NOT" || jsonTag == "AND" {
		return reflect.Value{}, nil // Return nil to signify no need to check further
	}

	// Use reflection to access the model's fields
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer if necessary
	}

	// Iterate through all fields in the struct to match the JSON tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		jsonTagValue := field.Tag.Get("json")
		if jsonTagValue == jsonTag {
			return val.Field(i), nil
		}
	}

	return reflect.Value{}, fmt.Errorf("field with json tag %s not found", jsonTag)
}

// getPrismaMethodName generates the Prisma method name dynamically based on the field and operator
func getPrismaMethodName(field, operator string) (string, string) {
	// The Prisma method name is typically in the form of "<Field><Operator>"
	return replaceId(capitalizeFirstLetter(field)), capitalizeFirstLetter(operator)
}

// capitalizeFirstLetter capitalizes the first letter of a string
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	// Capitalize the first letter
	capitalized := strings.ToUpper(s[:1]) + s[1:]
	return capitalized
}

// Replace "Id" with "ID" if it appears at the end
func replaceId(capitalized string) string {
	if strings.HasSuffix(capitalized, "Id") {
		capitalized = strings.TrimSuffix(capitalized, "Id") + "ID"
	}
	return capitalized
}

// convertToString safely converts any type to a string
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}
