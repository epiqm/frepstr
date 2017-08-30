// frepstr
// utility for replacing text in specified folders
// inluding recursive search in sub-folders
// written by @epiqmax <epiqmax@gmail.com>
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"log"
)

const DIR_SEP = string(os.PathSeparator)

var dirsAmount int
var filesAmount int
var foundFilesAmount int
var replacedFilesAmount int

var dirTree []string
var fileTree []string

func parseDirs(dir string, find string, replace string) {
	if dir[len(dir)-1 : len(dir)] != DIR_SEP {
		dir = fmt.Sprintf("%s%s", dir, DIR_SEP)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, f := range files {
		var ename string
		if f.IsDir() {
			dirsAmount++
			ename = fmt.Sprintf("%s%s%s", dir, f.Name(), DIR_SEP)
			fmt.Println(ename)
			dirTree = append(dirTree, ename)
			parseDirs(ename, find, replace)
		} else {
			filesAmount++
			ename = fmt.Sprintf("%s%s", dir, f.Name())
			fmt.Println(ename)

			if len(find) > 0 {
				input, err := ioutil.ReadFile(ename)
				if err != nil {
					fmt.Println("error:", err.Error())
				}

				var fileText string = string(input)
				if strings.Contains(fileText, find) {
					foundFilesAmount++
					fmt.Println("+", ename)
					if len(replace) > 0 {
						output := strings.Replace(fileText, find, replace, -1)
						err = ioutil.WriteFile(ename, []byte(output), 0)
						if err != nil {
							fmt.Println("error:", err.Error())
						} else {
							replacedFilesAmount++
							fmt.Println("R", ename)
						}
					}
				}
			}

			fileTree = append(fileTree, ename)
		}
	}
}

func main() {
	var	dirPtr = flag.String("dir", "./", "working directory")
	var	findPtr = flag.String("search", "", "text to search in the directory and sub directories")
	var	replacePtr = flag.String("replace", "", "replace search value")

	flag.Parse()

	_, err := os.Stat(*dirPtr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if os.IsNotExist(err) {
		log.Fatalln(err.Error())
	}

	if len(*replacePtr) > 0 {
		fmt.Println("Replace: ", *replacePtr)
	}

	parseDirs(string(*dirPtr), string(*findPtr), string(*replacePtr))

	fmt.Print(fmt.Sprintf("Checked %d dirs, %d files.\n", dirsAmount, filesAmount))
	if foundFilesAmount > 0 {
		fmt.Println(fmt.Sprintf("Files contain search text: %d.", foundFilesAmount))
	}
	if replacedFilesAmount > 0 {
		fmt.Println(fmt.Sprintf("Replaced: %d files.", replacedFilesAmount))
	}
}
