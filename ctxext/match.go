package ctxext

type ListGetter interface {
	List() []string
}

// ValueInList 判断参数是否在列表中
func ValueInList[Ctx any](val string, list ListGetter) func(Ctx) bool {
	return func(ctx Ctx) bool {
		for _, v := range list.List() {
			if val == v {
				return true
			}
		}
		return false
	}
}
