package structural

import (
	"cuelang.org/go/cue"
)

func CueQuery(squery, sdata string) (string, error) {
	out := NewpvList()

	vqueryi, err := r.Compile("", squery)
	if err != nil {
		return "", err
	}
	vquery := vqueryi.Value()
	if vquery.Err() != nil {
		return "", vquery.Err()
	}
	vdatai, err := r.Compile("", sdata)
	if err != nil {
		return "", err
	}
	vdata := vdatai.Value()
	if vdata.Err() != nil {
		return "", vdata.Err()
	}

	err = cueQuery(out, vquery, vdata)
	if err != nil {
		return "", err
	}

	return out.ToString()
}

func cueQuery(out *pvList, vquery, vdata cue.Value) error {
	// Loop over the query args in vquery
	queryListIter, err := vquery.List()
	if err != nil {
		return err
	}
	var list = []cue.Value{vdata}
	for queryListIter.Next() {
		queryVal := queryListIter.Value()
		list, err = queryStep(queryVal, list)
		if err != nil {
			return err
		}
	}
	for _, v := range list {
		out.Append(*ExprFromValue(v))
	}

	return nil
}

func queryStep(query cue.Value, list []cue.Value) ([]cue.Value, error) {
	out := make([]cue.Value, 0)

	for _, v := range list {
		if isStruct(v) {
			viter, err := v.Fields()
			if err != nil {
				return out, err
			}
			for viter.Next() {
				val := viter.Value()
				if query.Unify(val).Kind() != cue.BottomKind {
					out = append(out, val)
				}
			}
		}
	}

	return out, nil
}
