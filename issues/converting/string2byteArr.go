package converting

/**
In general, a []byte is not allowed to write to the underlying buffer of a string.
Converting a byte or rune into a string. This operation is actually a conversion from a Unicode code point into a single-character string
*/

func Write(buf []byte) {
	for _, c := range buf {
		WriteByte(c)
	}
}

// WriteString not cause a heap allocation
func WriteString(s string) {
	Write([]byte(s))
}

func WriteByte(c byte) {

}

/* ---- */
// this causes a heap allocation.
var global *int

func Foo1() {
	i := 3
	global = &i
}

// this does not cause a heap allocation.
func Foo2() {
	i := 3
	bar(&i)
}

func bar(i *int) {
	println(*i)
}
