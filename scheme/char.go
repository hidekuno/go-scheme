/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

var (
	whitespaceChar = map[string]byte{
		"#\\tab":     0x09,
		"#\\space":   0x20,
		"#\\newline": 0x0A,
		"#\\return":  0x0D,
	}
)

// Character Type
type Char struct {
	Expression
	Value rune
	exp   string
}

func NewChar(v string) *Char {
	b := new(Char)
	b.exp = v
	b.Value = []rune(v)[2]
	return b
}
func NewCharFromRune(c rune) *Char {
	b := new(Char)
	b.Value = c

	for k, v := range whitespaceChar {
		if c == rune(v) {
			b.exp = k
			return b
		}
	}
	b.exp = "#\\" + string(c)
	return b
}
func (self *Char) String() string {
	return self.exp
}
func (self *Char) isAtom() bool {
	return true
}
func (self *Char) clone() Expression {
	return NewChar(self.exp)
}
func (self *Char) equalValue(e Expression) bool {
	if v, ok := e.(*Char); ok {
		return self.Value == v.Value
	}
	return false
}
