package dataProcessingService

import (
	"strings"
)

func convertDataToSend(dataType string, attr []string, params []string) string {
	if len(params) == 0 {
		params = []string{"NA"}
	}
	// log.Print("converter")
	return "#" + dataType + "#" + strings.Join(attr, ";") + ";" + strings.Join(params, ",")
}

// func makeSlices(i int, list []string) [][]string {
// 	if len(list) <= i {
// 		return [][]string{list}
// 	}
// 	return append(makeSlices(i, list[i:]), list[:i])
// }
