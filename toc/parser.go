package toc

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	MARKDOWN_FILE_SUFFIX       string = "md"
	MARKDOWN_HEADING_CHAR      string = "#"
	ORDINAL_COMMENT_EXPRESSION string = "<!-- index\\[([0-9]+)\\] -->"
	SECTION_DIV_EXPRESSION     string = "<div.+id=\"(.+?)\".*>"
	HEADING_EXPRESSION         string = "^(#+)\\s*(.*)"
)

type SectionMarker struct {
	id         string
	linenumber int
}

type Markdownfile struct {
	root     Heading
	order    int
	name     string
	filepath string
	sections []SectionMarker
	parseErr error
}

type Heading struct {
	text        string
	depth       int
	parent      *Heading
	subheadings []*Heading
	linenumber  int
}

func GetMarkdownDocs(dir string, excludedList []string) []fs.FileInfo {
	items, _ := ioutil.ReadDir(dir)
	markdownDocuments := make([]fs.FileInfo, 0)

	for i := range items {
		filename := items[i].Name()
		if strings.HasSuffix(filename, "."+MARKDOWN_FILE_SUFFIX) && !isExcludedFile(filename, excludedList) {
			markdownDocuments = append(markdownDocuments, items[i])
		}
	}

	return markdownDocuments
}

func ParseFiles(markdownDocuments []fs.FileInfo, dir string) ([]Markdownfile, error) {
	ordinalMatcher, err := regexp.Compile(ORDINAL_COMMENT_EXPRESSION)
	if err != nil {
		return []Markdownfile{}, err
	}
	sectionMatcher, err := regexp.Compile(SECTION_DIV_EXPRESSION)
	if err != nil {
		return []Markdownfile{}, err
	}

	c := make(chan Markdownfile)
	for i := range markdownDocuments {
		path := filepath.Join(dir, markdownDocuments[i].Name())
		go parseFile(path, ordinalMatcher, sectionMatcher, c)
	}

	mdfiles := make([]Markdownfile, len(markdownDocuments))
	for i := range mdfiles {
		mdfiles[i] = <-c
	}

	return mdfiles, nil
}

func parseFile(path string, ordinalMatcher *regexp.Regexp, sectionMatcher *regexp.Regexp, resultchan chan Markdownfile) {
	file, err := os.Open(path)
	if err != nil {
		resultchan <- Markdownfile{
			name:     filepath.Base(file.Name()),
			filepath: path,
			parseErr: err,
		}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var latest *Heading
	root := Heading{
		depth:       0,
		subheadings: make([]*Heading, 0),
	}
	latest = &root
	order := math.MaxInt32
	linenumber := 0
	sections := make([]SectionMarker, 0)
	for scanner.Scan() {
		line := scanner.Text()
		orderMatches := ordinalMatcher.FindStringSubmatch(line)
		if len(orderMatches) == 2 {
			// Found the comment that defines the order for this file
			i, err := strconv.Atoi(orderMatches[1])
			if err == nil {
				order = i
			}
		}

		sectionMatches := sectionMatcher.FindStringSubmatch(line)
		if len(sectionMatches) == 2 {
			sections = append(sections, SectionMarker{
				id:         sectionMatches[1],
				linenumber: linenumber,
			})
		}

		if strings.HasPrefix(line, MARKDOWN_HEADING_CHAR) {
			h := parseHeading(line)
			h.linenumber = linenumber
			parent := getParent(h, latest)
			h.parent = parent
			parent.subheadings = append(parent.subheadings, &h)
			latest = &h
		}

		linenumber++
	}

	resultchan <- Markdownfile{
		name:     filepath.Base(file.Name()),
		filepath: path,
		root:     root,
		sections: sections,
		order:    order,
	}
}

func isExcludedFile(filename string, excludedList []string) bool {
	for i := range excludedList {
		exp := fmt.Sprintf("%s$", excludedList[i])
		r, _ := regexp.Compile(exp)
		if r.MatchString(filename) {
			return true
		}
	}
	return false
}

func getParent(h Heading, mostRecent *Heading) *Heading {
	c := mostRecent
	for h.depth <= c.depth {
		c = c.parent
	}
	return c
}

func parseHeading(line string) Heading {
	r, _ := regexp.Compile(HEADING_EXPRESSION)
	matches := r.FindStringSubmatch(line)

	if len(matches) < 3 {
		// empty header so ignore
		return Heading{}
	}

	return Heading{
		text:        matches[2],
		depth:       len(matches[1]),
		subheadings: make([]*Heading, 0),
	}
}
