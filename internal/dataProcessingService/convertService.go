package dataProcessingService

import (
	"strings"
)

func convertDataToSend(dataType string, attr []string, params []string) string {
	if len(params) == 0 {
		params = []string{"NA"}
	}
	return "#" + dataType + "#" + strings.Join(attr, ";") + ";" + listToSrt(params, ",")
}

func listToSrt(params []string, delim string) string {
	if len(params) == 0 {
		return ""
	}
	msg := ""
	for i := 0; i < len(params)-1; i++ {
		msg = params[i] + delim
	}
	return msg + params[len(params)-1]
}

func makeSlices(i int, list []string) [][]string {
	if len(list) <= i {
		return [][]string{list}
	}
	return append(makeSlices(i, list[i:]), list[:i])
}
