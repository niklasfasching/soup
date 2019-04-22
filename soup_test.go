package soup

import (
	"strings"
	"testing"
)

func TestSoup(t *testing.T) {
	d := MustParse(strings.NewReader(`<ul><li>foo</li><li>bar</li></ul>`))
	if actual := d.All("li").Text("\n"); actual != "foo\nbar" {
		t.Errorf("Got %s, expected foo\\nbar", actual)
	}

}
