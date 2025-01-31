package parser

import (
	"log"
	"reflect"
)

func BuildWhereFilters[U any](
	filters Filter,
	customFilters []U,
	tm TypeMapper,
	mapper string,
) []U {
	prismaFilters := customFilters
	// Map generic filters to Prisma filters
	if whereFilter, ok := filters["where"]; ok {
		filterMappers := WhereMapper(tm.GetDomainType(mapper), whereFilter.(Filter), tm.GetWhereClause(mapper), tm.GetDBInstance(mapper), tm)
		if len(filterMappers) > 0 {
			for _, filterMapper := range filterMappers {
				prismaFilters = append(prismaFilters, filterMapper.(U))
			}
		}
	}

	return prismaFilters
}

func WhereMapper[T any, K any, C any](model K, filter Filter, returnType T, where C, tm TypeMapper) []T {
	var prismaFilters []T

	// Handle the "where" key explicitly at the top level
	if whereValue, exists := filter["where"]; exists {
		if whereMap, ok := whereValue.(Filter); ok {
			// Recursively process the "where" filter
			return WhereMapper(model, whereMap, returnType, where, tm)
		}
	}

	// Process filters after unwrapping "where"
	for key, value := range filter {
		if ok := tm.GetDomainType(key); ok != nil {
			fieldType, _ := getFieldTypeByName(model, capitalizeFirstLetter(key))
			if fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Struct {
				// Handle Slice field with Slice.Where() method
				sliceFilters := (WhereMapper(tm.GetDomainType(key), value.(Filter), tm.GetWhereClause(key), tm.GetDBInstance(key), tm))
				if len(sliceFilters) > 0 {
					// Apply Category.Where() to the LibraryItem query
					method := reflect.ValueOf(where).FieldByName(capitalizeFirstLetter(key)).MethodByName("Where")

					if method.IsValid() {
						for _, categorFilter := range sliceFilters {
							prismaFilters = append(prismaFilters, method.Call([]reflect.Value{reflect.ValueOf(categorFilter)})[0].Interface().(T))
						}
					}
				}
			} else if fieldType.Kind() == reflect.Slice {
				if relationFilters, ok := value.(Filter); ok {
					for relationOperator, relationValue := range relationFilters {
						switch relationOperator {
						case "some":
							subFilters := WhereMapper(tm.GetDomainType(key), relationValue.(Filter), tm.GetWhereClause(key), tm.GetDBInstance(key), tm)
							method := reflect.ValueOf(where).FieldByName(capitalizeFirstLetter(key)).MethodByName("Some")
							if method.IsValid() {
								for _, subFilter := range subFilters {
									prismaFilters = append(prismaFilters, method.Call([]reflect.Value{reflect.ValueOf(subFilter)})[0].Interface().(T))
								}
							}
						case "every":
							subFilters := WhereMapper(tm.GetDomainType(key), relationValue.(Filter), tm.GetWhereClause(key), tm.GetDBInstance(key), tm)
							method := reflect.ValueOf(where).FieldByName(capitalizeFirstLetter(key)).MethodByName("Every")
							if method.IsValid() {
								for _, subFilter := range subFilters {
									prismaFilters = append(prismaFilters, method.Call([]reflect.Value{reflect.ValueOf(subFilter)})[0].Interface().(T))
								}
							}
						case "none":
							subFilters := WhereMapper(tm.GetDomainType(key), relationValue.(Filter), tm.GetWhereClause(key), tm.GetDBInstance(key), tm)
							method := reflect.ValueOf(where).FieldByName(capitalizeFirstLetter(key)).MethodByName("None")
							if method.IsValid() {
								for _, subFilter := range subFilters {
									prismaFilters = append(prismaFilters, method.Call([]reflect.Value{reflect.ValueOf(subFilter)})[0].Interface().(T))
								}
							}
						default:
							log.Printf("Unsupported relation operator %s for field %s\n", relationOperator, key)
						}
					}
				}
			}

			continue
		}
		// Use reflection to check if the field exists in the model via JSON tags
		_, err := getFieldByJSONTag(model, key)
		if err != nil {
			log.Printf("Field %s not found in model %T, skipping...\n", key, model)
			continue
		}
		// Handle logical operators like "OR", "NOT"
		switch key {
		case "OR":
			if conditions, ok := value.(Filter); ok {
				prismaFilters = append(prismaFilters, processFilterOperation(model, conditions, returnType, where, tm, "Or")...)
			} else {
				log.Println("value", value)
			}
		case "NOT":
			if conditions, ok := value.(Filter); ok {
				prismaFilters = append(prismaFilters, processFilterOperation(model, conditions, returnType, where, tm, "Not")...)
			} else {
				log.Println("value", value)
			}
		case "AND":
			if conditions, ok := value.(Filter); ok {
				log.Println("value", value)
				prismaFilters = append(prismaFilters, processFilterOperation(model, conditions, returnType, where, tm, "And")...)
			} else {
				log.Println("value", value)
			}
		// Handle fields dynamically
		default:
			if subFilters, ok := value.(Filter); ok {
				for operator, v := range subFilters {
					strValue := convertToString(v) // Convert any type to string
					// Dynamically map the operator to Prisma method
					prismaMethod, err := mapOperatorToPrismaMethod(model, key, operator, strValue, returnType, where, tm)
					if err == nil {
						prismaFilters = append(prismaFilters, prismaMethod)
					} else {
						log.Printf("Unsupported operator %s for field %s\n", operator, key)
					}
				}
			} else {
				strValue := convertToString(value) // Convert any type to string
				// Dynamically map the operator to Prisma method
				prismaMethod, err := mapOperatorToPrismaMethod(model, key, "Equals", strValue, returnType, where, tm)
				if err == nil {
					prismaFilters = append(prismaFilters, prismaMethod)
				} else {
					log.Printf("Unsupported operator %s for field %s\n", "Equals", key)
				}
			}
		}
	}

	return prismaFilters
}

func processFilterOperation[T any, K any, C any](
	model K,
	conditions Filter,
	returnType T,
	where C,
	tm TypeMapper,
	operation string, // The name of the operation: "And", "Or", or "Not"
) []T {
	// Step 1: Create filters dynamically
	var filters []T
	if operation == "Or" {
		for _, val := range conditions {
			// Map the condition using FilterMapper
			filters = append(filters, WhereMapper(model, val.(Filter), returnType, where, tm)...)
		}
	} else {
		filters = WhereMapper(model, conditions, returnType, where, tm)
	}

	// Step 2: Find the specified method dynamically
	method := reflect.ValueOf(where).MethodByName(operation)
	if !method.IsValid() {
		log.Fatalf("The '%s' method does not exist for type: %T", operation, where)
	}

	// Step 3: Dynamically create a slice of the required type
	filtersValue := reflect.ValueOf(filters)
	if filtersValue.Kind() != reflect.Slice {
		log.Fatalf("filters is not a slice, got: %T", filters)
	}

	// Dynamically construct a new slice of the type expected by the method
	filtersType := method.Type().In(0) // Get the type of the first argument of the method
	dynamicSlice := reflect.MakeSlice(filtersType, filtersValue.Len(), filtersValue.Len())

	// Populate the new slice with the elements of `filters`, casting to the correct type
	for i := 0; i < filtersValue.Len(); i++ {
		element := filtersValue.Index(i).Interface()
		typedElement, ok := element.(T) // Explicit type assertion to T
		if !ok {
			log.Fatalf("Element %v is not of type %T", element, returnType)
		}
		dynamicSlice.Index(i).Set(reflect.ValueOf(typedElement))
	}

	// Step 4: Call the specified method with `CallSlice`
	results := method.CallSlice([]reflect.Value{dynamicSlice})

	// Step 5: Convert the result back to `[]T`
	var prismaFilters []T
	for _, result := range results {
		prismaFilters = append(prismaFilters, result.Interface().(T))
	}

	return prismaFilters
}
