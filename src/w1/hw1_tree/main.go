package main

import (
	"fmt"
	"io"
	"os"
)

func dirTreeString(path string, printFiles bool, pref string) (res string, err error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return
	}
	var list []os.DirEntry
	for _, entry := range dirEntry {
		if entry.IsDir() || printFiles {
			list = append(list, entry)
		}
	}
	for i, item := range list {
		curPref, subPref := func(last bool) (string, string) {
			if last {
				return "└───", "\t"
			} else {
				return "├───", "│\t"
			}
		}(i == len(list)-1)
		if item.IsDir() {
			res += pref + curPref + item.Name() + "\n"
			subPath := path + string(os.PathSeparator) + item.Name()
			var subTree string
			subTree, err = dirTreeString(subPath, printFiles, pref+subPref)
			if err != nil {
				return
			}
			res += subTree
		} else if printFiles {
			var info os.FileInfo
			info, err = item.Info()
			res += fmt.Sprintf("%v%v%v", pref, curPref, item.Name())
			if info.Size() == 0 {
				res += " (empty)\n"
			} else {
				res += fmt.Sprintf(" (%vb)\n", info.Size())
			}
		}
	}
	return
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	tree, err := dirTreeString(path, printFiles, "")
	_, err = fmt.Fprint(out, tree)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
