package walker

type GetList []GetOp

func (s GetList) ApplyOn(g Getter) (FieldValues, []error) {
	var errs []error

	res := make(FieldValues, len(s))

	for _, op := range s {
		uid, fieldValue, err := op(g)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if fieldValue != nil {
			res[uid] = fieldValue
		}
	}

	return res, errs
}

func (s *GetList) Push(o GetOp) {
	*s = append(*s, o)
}

type GetOp func(g Getter) (uid uint, fieldValue *FieldValue, err error)
