package lexers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/styles"
)

/*
	WARNING: Running tests will leave html formatted files in testdata/aql for local testing purpose. Delete them after usage.
*/

func TestInsert(t *testing.T) {
	testFile, err := ioutil.ReadFile("testdata/aql/test_insert.actual")
	code := string(testFile)

	err = write("insert", code)
	if err != nil {
		log.Print(err)
	}
}

func TestMatchRange(t *testing.T) {
	testFile, err := ioutil.ReadFile("testdata/aql/test_match_range.actual")
	code := string(testFile)

	err = write("match_range", code)
	if err != nil {
		log.Print(err)
	}
}

func TestCreateEdges(t *testing.T) {
	testFile, err := ioutil.ReadFile("testdata/aql/test_create_edges.actual")
	code := string(testFile)

	err = write("create_edges", code)
	if err != nil {
		log.Print(err)
	}
}

func TestMergeTraits(t *testing.T) {
	testFile, err := ioutil.ReadFile("testdata/aql/test_merge_traits.actual")
	code := string(testFile)

	err = write("merge_traits", code)
	if err != nil {
		log.Print(err)
	}
}

func TestSortMultipleAttrs(t *testing.T) {
	testFile, err := ioutil.ReadFile("testdata/aql/test_sort_multiple_attrs.actual")
	code := string(testFile)

	err = write("sort_multiple_attrs", code)
	if err != nil {
		log.Print(err)
	}
}

func write(testName, code string) error {
	l := Get("aql")
	if l == nil {
		return errors.New("cannot find aql lexer")
	}

	l = chroma.Coalesce(l)

	// Determine formatter and style.
	f, s := formatters.Get("html"), styles.Get("monokai")
	if f == nil {
		f = formatters.Fallback
	}

	it, err := l.Tokenise(nil, code)
	if err != nil {
		log.Print(err)
	}

	// Prepare output file
	os.Truncate("testdata/aql/%s.html", 100)
	file, err := os.Create(fmt.Sprintf("testdata/aql/%s.html", testName))
	err = f.Format(file, s, it)

	return err
}
