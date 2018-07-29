package actioncontext

func assert(test bool) {
	if !test {
		panic("")
	}
}

func assertEx(test bool, msg string) {
	if !test {
		panic(msg)
	}
}
