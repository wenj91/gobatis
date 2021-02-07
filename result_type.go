package gobatis

type ResultType string

const (
	ResultTypeMap     ResultType = "map"     // Result set is a map: map[string]interface{}
	ResultTypeMaps    ResultType = "maps"    // Result set is a slice, item is map: []map[string]interface{}
	ResultTypeStruct  ResultType = "struct"  // Result set is a struct
	ResultTypeStructs ResultType = "structs" // Result set is a slice, item is struct
	ResultTypeSlice   ResultType = "slice"   // Result set is a value slice, []interface{}
	ResultTypeSlices  ResultType = "slices"  // Result set is a value slice, item is value slice, []interface{}
	ResultTypeArray   ResultType = "array"   //
	ResultTypeArrays  ResultType = "arrays"  // Result set is a value slice, item is value slice, []interface{}
	ResultTypeValue   ResultType = "value"   // Result set is single value
)
