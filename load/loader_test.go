package load

import (
	"path"
	"testing"

	"github.com/jakubDoka/goml/core"
)

func TestAppDataDir(t *testing.T) {
	dir, err := AppDataDir()
	if err != nil {
		t.Error(err)
		return
	}

	if dir != "C:\\Users\\jakub\\AppData\\Roaming" {
		t.Error(dir)
	}
}

func TestUtil(t *testing.T) {
	loader := OS

	root := "C:/Users/jakub/Documents/programming/golang/src/github.com/jakubDoka/mlok/load/test_data"

	testCases := []struct {
		desc string
		rec  bool
		ext  string
		out  map[string]struct{}
	}{
		{
			desc: "default args",
			out: map[string]struct{}{
				path.Join(root, "a"):     {},
				path.Join(root, "b.txt"): {},
				path.Join(root, "c.ttf"): {},
			},
		},
		{
			desc: "all",
			rec:  true,
			out: map[string]struct{}{
				path.Join(root, "a"):              {},
				path.Join(root, "b.txt"):          {},
				path.Join(root, "c.ttf"):          {},
				path.Join(root, "inner", "e.h"):   {},
				path.Join(root, "inner", "hello"): {},
				path.Join(root, "inner", "k.c"):   {},
			},
		},
		{
			desc: "txt",
			rec:  true,
			ext:  "txt",
			out: map[string]struct{}{
				path.Join(root, "b.txt"): {},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			list, err := loader.List(
				root,
				nil, tC.rec, tC.ext,
			)

			if err != nil {
				t.Error(err)
				return
			}

			mp := map[string]struct{}{}
			for _, e := range list {
				mp[e] = struct{}{}
			}

			core.TestEqual(t, mp, tC.out)
		})
	}
}
