package fmt_io

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Tree map[string]string

func (data Tree) Children(parent string) (children []string) {
	for k, v := range data {
		if v == parent {
			children = append(children, k)
		}
	}
	return
}

func (data Tree) print(root string, cb func(string) string, file *os.File) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
		}
	}()
	//fmt_io := make([]string, 0, 40960)
	var f func(_, _, _, _ string)
	f = func(parent, prefix, add, node string) {
		var children = data.Children(parent)
		var extra string
		if cb != nil {
			extra = cb(parent)
			if extra != "" {
				if len(children) == 0 {
					extra = ": " + strings.ReplaceAll(extra, "\n", fmt.Sprintf("\n%s%s  %s", prefix, add, strings.Repeat(" ", len(parent))))
				} else {
					extra = ": " + strings.ReplaceAll(extra, "\n", fmt.Sprintf("\n%s%s│ %s", prefix, add, strings.Repeat(" ", len(parent))))
				}
			}
		}
		//fmt_io = append(fmt_io, fmt.Sprintf("%s%s%s%s", prefix, node, parent, extra))
		if _, err := fmt.Fprintf(file, "%s%s%s%s\n", prefix, node, parent, extra); err != nil {
			panic(err)
		}
		sort.Slice(children, func(i, j int) bool {
			return children[i] < children[j]
		})
		for i, child := range children {
			if i == len(children)-1 {
				f(child, prefix+add, "    ", "└── ")
			} else {
				f(child, prefix+add, "│   ", "├── ")
			}
		}
	}
	f(root, "", "", "")
	//for _, line := range fmt_io {
	//	fmt.Println(line)
	//}
	return
}

func (data Tree) Print(root string, cb func(string) string) error {
	return data.print(root, cb, os.Stdout)
}

func (data Tree) Dump(root string, cb func(string) string, file string) error {
	f, err := os.Create(file)
	defer f.Close()
	if err == nil {
		err = data.print(root, cb, f)
	}
	return err
}
