package lib

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	DEFAULT_TITLE          string = "Table of contents"
	SPACES_PER_LEVEL       int    = 4
	MAX_DISTANCE_TO_ANCHOR int    = 100
)

func WriteIndex(indexFilename string, dir string, files []Markdownfile) error {
	path := filepath.Join(dir, indexFilename)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	title := fmt.Sprintf("# %s  \n", DEFAULT_TITLE)
	w.WriteString(title)
	sortByOrder(files)
	for i := range files {
		for j := range files[i].root.subheadings {
			writeHeading(files[i].root.subheadings[j], files[i], w)
		}
	}
	footer := fmt.Sprintf("\n<!-- Generated at %s -->\n", time.Now().Format(time.RFC3339))
	w.WriteString(footer)
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func sortByOrder(files []Markdownfile) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].order < files[j].order
	})
}

func writeHeading(h *Heading, f Markdownfile, w *bufio.Writer) {
	space := getIndentSpaces(h.depth - 1)
	w.WriteString(space)
	url := f.name
	nearestAnchor, distance := findNearestSectionAnchor(h.linenumber, f.sections)
	if nearestAnchor != nil && distance < MAX_DISTANCE_TO_ANCHOR {
		url += fmt.Sprintf("#%s", nearestAnchor.id)
	}
	link := fmt.Sprintf("[%s](%s)", h.text, url)
	w.WriteString(fmt.Sprintf("* %s  \n", link))
	for i := range h.subheadings {
		writeHeading(h.subheadings[i], f, w)
	}
}

func getIndentSpaces(indent int) string {
	i := 0
	space := ""
	n := indent * SPACES_PER_LEVEL
	for i < n {
		space += " "
		i++
	}

	return space
}

func findNearestSectionAnchor(linenumber int, sections []SectionMarker) (nearestSection *SectionMarker, distance int) {
	distance = -1
	for i := range sections {
		d := getDistance(sections[i].linenumber, linenumber)
		if distance == -1 || d < distance {
			distance = d
			nearestSection = &sections[i]
		}
	}

	return
}

func getDistance(x int, y int) int {
	if x > y {
		return x - y
	}

	return y - x
}
