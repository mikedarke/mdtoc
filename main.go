package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikedarke/mdtoc/toc"
)

const (
	DEFAULT_INDEX_FILE_NAME string = "index.md"
	DEFAULT_TITLE           string = "Table of contents"
)

func main() {
	var indexFilename string
	var title string
	flag.StringVar(&indexFilename, "o", DEFAULT_INDEX_FILE_NAME, "output file name to write the table of contents")
	flag.StringVar(&title, "t", DEFAULT_TITLE, "the title to display above the table of contents")
	flag.Usage = usage
	flag.Parse()
	dir := flag.Arg(0)
	opt := toc.Options{
		Directory:      dir,
		Outputfilename: indexFilename,
		Title:          title,
	}
	markdownFiles := toc.GetMarkdownDocs(dir, []string{indexFilename})
	documents, err := toc.ParseFiles(markdownFiles, dir)
	check(err)
	err = toc.WriteIndex(opt, documents)
	check(err)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  mdtoc options srcfolder\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  example...\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  mdtoc -o index.md ./example\n\n")
	flag.PrintDefaults()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
