package constants

const(
	// result set is a map: map[string]interface{}
	RESULT_TYPE_MAP = "Map"
	// result set is a slice, item is map: []map[string]interface{}
	RESULT_TYPE_MAPS = "Maps"
	// result set is a struct
	RESULT_TYPE_STRUCT = "Struct"
	// result set is a slice, item is struct
	RESULT_TYPE_STRUCTS = "Structs"
	// result set is a value slice, []interface{}
	RESULT_TYPE_SLICE = "Slice"
	// result set is a value slice, item is value slice, []interface{}
	RESULT_TYPE_SLICES = "Slices"
)
