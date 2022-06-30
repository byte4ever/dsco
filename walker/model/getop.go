package model

import (
	"errors"

	"github.com/byte4ever/dsco/merror"
	"github.com/byte4ever/dsco/walker/fvalues"
	"github.com/byte4ever/dsco/walker/ifaces"
)

type GetList []GetOp

func (s *GetList) Count() int {
	return len(*s)
}

type ApplyError struct {
	merror.MError
}

var ErrApply = errors.New("")

func (m ApplyError) Is(err error) bool {
	return errors.Is(err, ErrApply)
}

func (s GetList) ApplyOn(g ifaces.Getter) (fvalues.FieldValues, error) {
	var errs ApplyError

	res := make(fvalues.FieldValues, len(s))

	for _, op := range s {
		uid, fieldValue, err := op(g)

		if err != nil {
			errs.Add(err)
			continue
		}

		if fieldValue != nil {
			res[uid] = fieldValue
		}
	}

	if errs.None() {
		return res, nil
	}

	return res, errs
}

func (s *GetList) Push(o GetOp) {
	*s = append(*s, o)
}

type GetOp func(g ifaces.Getter) (uid uint, fieldValue *fvalues.FieldValue,
	err error)
