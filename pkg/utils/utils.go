package utils

func GetSubInterfaceId(customerId int, outer int, inner int) int {
	return (customerId * 10000000) + (outer * 10000) + inner
}

func ConcatVlanPrio(vlanId int, prio int) int {
	return ((prio << 13) | vlanId)
}
