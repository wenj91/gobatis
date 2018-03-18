package tagutil

func IsContains(name string) bool {
	if(name == "insert" ||
		name == "delete" ||
		name == "update" ||
		name == "select" ||
		name == "if" ||
		name == "foreach"){
		return true
	}

	return false
}
