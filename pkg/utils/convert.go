package utils

func ToStringArray(source interface{}) []string {
	stringArray := make([]string, 0)
	if source == nil {
		return stringArray
	}
	ifArray, flag := source.([]interface{})
	if !flag {
		return stringArray
	}
	if len(ifArray) == 0 {
		return stringArray
	}
	for _, item := range ifArray {
		stringArray = append(stringArray, item.(string))
	}
	return stringArray
}
