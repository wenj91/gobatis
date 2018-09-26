package gobatis

func GetFieldMapperByTagName(t interface{}, tagName string) (map[string]string, error){
	res := make(map[string]string)

	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Ptr {
		return nil, errors.New("t must not be ptr")
	}

	tt := reflect.TypeOf(t)
	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		fieldMapperKey := field.Tag.Get(tagName)
		if fieldMapperKey != "" {
			res[fieldMapperKey] = field.Name
		}else{
			res[field.Name] = field.Name
		}
	}

	return res, nil
}
