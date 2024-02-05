package server

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PathTestCase struct {
	Origin string
	Expect string
	Desc   string
}

func TestServerPath(t *testing.T) {
	pathTest := []PathTestCase{
		{
			Origin: ".",
			Expect: ".",
		},
		{
			Origin: "",
			Expect: ".",
		},
		{
			Origin: "/hello",
			Expect: "/hello",
		},
	}

	for _, tc := range pathTest {
		res := path.Base(tc.Origin)
		assert.Equal(t, tc.Expect, res)
	}
}
