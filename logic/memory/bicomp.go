package memory

import "github.com/cheekybits/genny/generic"

//go:generate genny -pkg=gen -in=$GOFILE -out=gen/$GOFILE gen "Number=NUMBERS"

type Number generic.Number

func NumberBiComp(a, b Number) uint8 {
	if a > b {
		return 1
	} else if a < b {
		return 2
	}

	return 0
}
