package jsonUtils

import "encoding/json"

//
// MapToJson
//  @Description: 将结构体转化成json格式
//  @param Mapstruct
//  @return string
//
func MapToJson(Mapstruct interface{}) string {
	// map转 json str
	jsonBytes, _ := json.Marshal(Mapstruct)
	jsonStr := string(jsonBytes)
	return jsonStr
}
