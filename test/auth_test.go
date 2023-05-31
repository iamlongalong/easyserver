package test

import (
	"path/filepath"
	"testing"

	"github.com/magiconair/properties/assert"
)

type PathTestItem struct {
	Path string
	Root string

	Expect string
}

func TestPathClean(t *testing.T) {
	paths := []PathTestItem{
		{
			Path:   "../hello",
			Root:   "/",
			Expect: "/hello",
		},
		{
			Path:   "../../hello",
			Root:   "/",
			Expect: "/hello",
		},
		{
			Path:   "/../hello",
			Root:   "/",
			Expect: "/hello",
		},
		{
			Path:   "/../hello",
			Root:   "/hi",
			Expect: "/hi/hello",
		},
		{
			Path:   "/../../hello",
			Root:   "/hi/hey",
			Expect: "/hi/hey/hello",
		},
		{
			Path:   "/../..",
			Root:   "/hi",
			Expect: "/hi",
		},
		{
			Path:   "../../../hey",
			Root:   "/hi",
			Expect: "/hi/hey",
		},
	}

	for _, item := range paths {
		relativePath := filepath.Join("/", filepath.Clean(item.Path))
		real := filepath.Join(item.Root, relativePath)
		assert.Equal(t, real, item.Expect)
	}
}
