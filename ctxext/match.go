package ctxext

type ListGetter interface {
	List() []string
}

// ValueInList 判断参数是否在列表中
func ValueInList[Ctx any](getval func(Ctx) string, list ListGetter) func(Ctx) bool {
	return func(ctx Ctx) bool {
		val := getval(ctx)
		for _, v := range list.List() {
			if val == v {
				return true
			}
		}
		return false
	}
}
