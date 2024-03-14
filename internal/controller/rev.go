package controller

func ReverseName(name string) string {
	rev := []byte(name)
	for i := 0; i < len(rev)/2; i++ {
		rev[i], rev[len(rev)-i-1] = rev[len(rev)-i-1], rev[i]
	}
	return string(rev)
}
