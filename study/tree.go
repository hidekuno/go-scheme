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
	MaxLineSize   = 1024
	DelimiterChar = "."
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
	idx := strings.LastIndex(self.Name, DelimiterChar)
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

	var reader = bufio.NewReaderSize(stream, MaxLineSize)
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
		idx := strings.LastIndex(s, DelimiterChar)
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

	var reader = bufio.NewReaderSize(stream, MaxLineSize)
	cache := map[string]*Item{}
	var top string

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		itemName := ""
		for _, s := range strings.Split(string(line), DelimiterChar) {
			if itemName == "" {
				itemName = s
			} else {
				itemName = itemName + DelimiterChar + s
			}
			if _, ok := cache[itemName]; ok {
				continue
			}
			idx := strings.LastIndex(itemName, DelimiterChar)
			if idx == -1 {
				top = itemName
				cache[itemName] = NewItem(itemName, nil)

			} else {
				cache[itemName] = NewItem(s, cache[itemName[:idx]])
				cache[itemName[:idx]].Add(cache[itemName])
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
	createTree := CreateTree
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
		createTree = CreateTreeOrdered
	}
FINISH:
	return v, file, createTree, ret
}

func main() {
	v, fd, createTree, err := ArgInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fd.Close()

	fj := createTree(fd)
	fj.Accept(v)
}
