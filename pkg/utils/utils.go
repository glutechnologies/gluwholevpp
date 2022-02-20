package utils

func GetSubInterfaceId(customerId int, outer int, inner int) int {
	return (customerId * 10000) + (outer * 1000) + inner
}
