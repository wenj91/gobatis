package gobatis

type ResultType string

const (
	// result set is a map: map[string]interface{}
	resultTypeMap ResultType = "Map"
	// result set is a slice, item is map: []map[string]interface{}
	resultTypeMaps ResultType = "Maps"
	// result set is a struct
	resultTypeStruct ResultType = "Struct"
	// result set is a slice, item is struct
	resultTypeStructs ResultType = "Structs"
	// result set is a value slice, []interface{}
	resultTypeSlice ResultType = "Slice"
	// result set is a value slice, item is value slice, []interface{}
	resultTypeSlices ResultType = "Slices"
	resultTypeValue  ResultType = "Value"
)
