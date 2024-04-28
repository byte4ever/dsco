package model

import (
	"github.com/byte4ever/dsco/internal"
)

type ExpandList []ExpandOp

func (s *ExpandList) Count() int {
	return len(*s)
}

func (s ExpandList) ApplyOn(g internal.StructExpander) error {
	var errs ApplyError

	for _, op := range s {
		err := op(g)

		if err != nil {
			errs.Add(err)
			continue
		}
	}

	if errs.None() {
		return nil
	}

	return errs
}

func (s *ExpandList) Push(o ExpandOp) {
	*s = append(*s, o)
}
