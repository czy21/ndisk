package util

type Object struct {
}

func (Object) Default(source, obj any) {
	if source == nil {
		source = obj
	}
}
