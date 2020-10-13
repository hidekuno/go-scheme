/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
package scheme

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
func (self *Char) String() string {
	return self.exp
}
func (self *Char) isAtom() bool {
	return true
}
func (self *Char) clone() Expression {
	return NewChar(self.exp)
}
