package web

func TryCatch(try func(), catch func(interface{})) (r bool) {
	defer func() {
		if x := recover(); x != nil {
			r = true
			catch(x)
		}
	}()
	try()
	return
}
