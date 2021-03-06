package composition

import (
	"testing"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
)

func stylesheetAttrs(href string) []html.Attribute {
	commonAttr := []html.Attribute{{Key: "rel", Val: "stylesheet"}, {Key: "type", Val: "text/css"}}
	return append(commonAttr, html.Attribute{Key: "href", Val: href})
}

func Test_IdentityDeduplicationStrategy(t *testing.T) {
	a := assert.New(t)
	stylesheets := [][]html.Attribute{stylesheetAttrs("/a"), stylesheetAttrs("/b")}
	deduper := new(IdentityDeduplicationStrategy)
	result := deduper.Deduplicate(stylesheets)
	a.EqualValues(stylesheets, result)
}

func Test_SimpleDeduplicationStrategy(t *testing.T) {
	a := assert.New(t)
	stylesheets := [][]html.Attribute{
		stylesheetAttrs("/a"),
		stylesheetAttrs("/b"),
		stylesheetAttrs("/a"),
		stylesheetAttrs("/b"),
		stylesheetAttrs("/c"),
		stylesheetAttrs("/a"),
	}
	expected := [][]html.Attribute{
		stylesheetAttrs("/a"),
		stylesheetAttrs("/b"),
		stylesheetAttrs("/c"),
	}
	deduper := new(SimpleDeduplicationStrategy)
	result := deduper.Deduplicate(stylesheets)
	a.EqualValues(expected, result)
}

// Tests for setting an own deduplication strategy
type Strategy struct {
}

func (strategy *Strategy) Deduplicate(stylesheets [][]html.Attribute) (result [][]html.Attribute) {
	for i, stylesheetAttrs := range stylesheets {
		if i%2 == 0 {
			result = append(result, stylesheetAttrs)
		}
	}
	return result
}
