package documentFields

type Field struct {
	Id   string
	Name string
	Lvl  int
}

func FilterByLvl(fields []Field, lvl int) (result []Field) {

	for _, field := range fields {

		if lvl >= field.Lvl {
			result = append(result, field)
		}
	}

	return
}

func FilterByIds(fields []Field, ids []string) (result []Field) {

	exist := make(map[string]bool)

	for _, v := range ids {
		exist[v] = true
	}

	for _, field := range fields {

		if _, ok := exist[field.Id]; ok == true {
			result = append(result, field)
		}
	}

	return
}

func FieldWithIdExist(fields []Field, id string) bool {

	for _, field := range fields {

		if id == field.Id {
			return true
		}
	}

	return false
}
