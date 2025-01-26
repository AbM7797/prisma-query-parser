package parser

import (
	"log"
)

func BuildOrderFilters[U any](
	filters Filter,
	customFilters []U,
	tm TypeMapper,
	mapper string,
) []U {
	prismaFilters := customFilters
	// Map generic filters to Prisma filters
	if whereFilter, ok := filters["orderBy"]; ok {
		filterMappers := OrderByMapper(whereFilter.(Filter), tm.GetDomainType(mapper), tm.GetOrderByClause(mapper), tm.GetDBInstance(mapper), tm)
		if len(filterMappers) > 0 {
			for _, filterMapper := range filterMappers {
				prismaFilters = append(prismaFilters, filterMapper.(U))
			}
		}
	}

	return prismaFilters
}

func OrderByMapper[T any, C any, U any](filter Filter, model U, returnType T, order C, tm TypeMapper) []T {
	var orderBy []T

	for key, value := range filter {
		strValue := convertToString(value) // Convert any type to string
		// Dynamically map the operator to Prisma method
		prismaMethod, err := mapOperatorToPrismaMethod(model, key, "Order", strValue, returnType, order, tm)
		if err == nil {
			orderBy = append(orderBy, prismaMethod)
		} else {
			log.Printf("Unsupported operator %s for field %s\n", "Equals", key)
		}
	}

	return orderBy
}
