package web

type ValidationFn func(Request) error

//Request represents a POST or PUT HttpRequest body
type Request interface {
	Validate(ValidationFn) error
}

//func CheckVal[T Request](t T) bool {
//
//	rootDir, _ := os.Getwd()
//	file, _ := os.OpenFile(fmt.Sprintf("%s/test.json", rootDir), os.O_RDONLY, 0777)
//
//	_ = utils.FromJson(&t, file)
//
//	return t.Validate(func(t Request) bool {
//		return utils.IsValid(t) == nil
//	})
//}
