/*
   Go lang 1st study program.
   This is thing for study data type and basic syntax.

   hidekuno@gmail.com
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	MAX_LINE_SIZE  = 1024
	DELIMITER_CHAR = "."
)

type Item struct {
	Name     string
	Parent   *Item
	Children []*Item
}

func NewItem(name string, parent *Item) *Item {
	item := new(Item)
	item.Name = name
	item.Parent = parent
	return item
}

func (self *Item) Add(c *Item) {
	self.Children = append(self.Children, c)
}

func (self *Item) LastName() string {
	str := self.Name
	idx := strings.LastIndex(self.Name, DELIMITER_CHAR)
	if idx > 0 {
		str = self.Name[idx+1:]
	}
	return str
}

func (self *Item) Accept(visitor Visitor) {
	visitor.Visit(self)
}

type Visitor interface {
	Visit(*Item)
}

type VisitorImpl struct {
	Visitor
	Level int
}

type VisitorLine struct{ Visitor }

func (self *VisitorImpl) Visit(item *Item) {
	var indent []byte
	for i := 0; i < self.Level*2; i++ {
		indent = append(indent, 0x20)
	}
	fmt.Print(string(indent))
	fmt.Println(item.LastName())
	for _, c := range item.Children {
		self.Level++
		c.Accept(self)
		self.Level--
	}
}

func (self *VisitorLine) Visit(item *Item) {
	var line []string

	line = append(line, item.LastName())
	if item.Parent != nil {
		p := item.Parent

		if p.Children[len(p.Children)-1] == item {
			line = append(line, "`--")
		} else {
			line = append(line, "|--")
		}
		for p.Parent != nil {
			if p.Parent.Children[len(p.Parent.Children)-1] == p {
				line = append(line, "   ")
			} else {
				line = append(line, "|  ")
			}
			p = p.Parent
		}
	}
	for i := len(line) - 1; i >= 0; i-- {
		fmt.Print(line[i])
	}
	fmt.Println("")
	for _, c := range item.Children {
		c.Accept(self)
	}
}

func NewDefaultVisitor() *VisitorImpl {
	v := new(VisitorImpl)
	v.Level = 0
	return v
}

func NewVisitorLine() *VisitorLine {
	v := new(VisitorLine)
	return v
}

func CreateTreeOrdered(stream *os.File) *Item {

	var reader = bufio.NewReaderSize(stream, MAX_LINE_SIZE)
	cache := map[string]*Item{}
	var top string

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		s := string(line)
		idx := strings.LastIndex(s, DELIMITER_CHAR)
		if idx == -1 {
			top = s
			cache[s] = NewItem(s, nil)
		} else {
			cache[s] = NewItem(s, cache[s[:idx]])
			cache[s[:idx]].Add(cache[s])
		}
	}
	return cache[top]
}

func CreateTree(stream *os.File) *Item {

	var reader = bufio.NewReaderSize(stream, MAX_LINE_SIZE)
	cache := map[string]*Item{}
	var top string

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		item_name := ""
		for _, s := range strings.Split(string(line), ".") {
			if item_name == "" {
				item_name = s
			} else {
				item_name = item_name + "." + s
			}
			if _, ok := cache[item_name]; ok {
				continue
			}
			idx := strings.LastIndex(item_name, DELIMITER_CHAR)
			if idx == -1 {
				top = item_name
				cache[item_name] = NewItem(item_name, nil)

			} else {
				cache[item_name] = NewItem(s, cache[item_name[:idx]])
				cache[item_name[:idx]].Add(cache[item_name])
			}
		}
	}
	return cache[top]
}

func ArgInit() (Visitor, *os.File, func(*os.File) *Item, error) {

	f := flag.String("f", "", "help fj data filename")
	l := flag.Bool("l", false, "help keisen line ")
	o := flag.Bool("o", false, "help sorted file ")
	flag.Parse()

	var (
		v   Visitor
		ret error
	)
	file := os.Stdin
	create_tree := CreateTree
	v = NewDefaultVisitor()
	ret = nil

	if *f != "" {
		fd, err := os.Open(*f)
		if os.IsNotExist(err) {
			ret = err
			goto FINISH
		} else if err != nil {
			panic(err)
		}
		file = fd
	}
	if *l {
		v = NewVisitorLine()
	}
	if *o {
		create_tree = CreateTreeOrdered
	}
FINISH:
	return v, file, create_tree, ret
}

func main() {
	v, fd, create_tree, err := ArgInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fd.Close()

	fj := create_tree(fd)
	fj.Accept(v)
}
