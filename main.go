package main

import (
	"os"

	"github.com/mikedarke/mdtoc/mdtoc"
)

const (
	DEFAULT_INDEX_FILE_NAME string = "index.md"
)

func main() {
	argsWithoutProg := os.Args[1:]
	dir := argsWithoutProg[0]
	indexFilename := DEFAULT_INDEX_FILE_NAME
	markdownFiles := mdtoc.GetMarkdownDocs(dir, []string{indexFilename})
	documents, err := mdtoc.ParseFiles(markdownFiles, dir)
	check(err)
	err = mdtoc.WriteIndex(indexFilename, dir, documents)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
