package iteration

func Repeat(char string, num int) string {
	var repeated string
	for i := 0; i < num; i++ {
		repeated += char
	}
	return repeated
}
