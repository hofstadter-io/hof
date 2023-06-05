{{ $ROOT := . }}
{{ $T := .TYPE }}

package types

import (
	"fmt"
)

// Represents a {{ .TYPE.Name }}
type {{ $T.Name }} struct {
	{{ range $T.OrderedFields }}
	{{ .Name }} {{ .Type }}
	{{- end }}

	{{ range $R := $T.Relations }}
	{{ $R.Name }} {{ $R.GoType }}
	{{ end }}
}

// A map type to store {{ .TYPE.Name }}
type {{ $T.Name }}Map map[string]*{{ $T.Name }}

// A var to work with
var {{ $T.Name }}Store {{ $T.Name }}Map

// Note, we are omitting locking and allowing concurrency issues

// initialize our storage
func init() {
	{{ $T.Name }}Store = make({{ $T.Name }}Map)
}

//
//// library funcs
//


func {{ $T.Name }}Create(in *{{ $T.Name }}) error {
	idx := in.{{ $T.Index }}

	// check if already exists
	if _, ok := {{ $T.Name }}Store[idx]; ok {
		return fmt.Errorf("Entry with %v already exists", idx)
	}

	// store the new value
	{{ $T.Name }}Store[idx] = in

	return nil
}

func {{ $T.Name }}Read(idx string) (*{{ $T.Name }}, error) {

	// return if exists
	if val, ok := {{ $T.Name }}Store[idx]; ok {
		return val, nil
	}

	// otherwise return error
	return nil, fmt.Errorf("Entry with %v does not exist", idx)
}

func {{ $T.Name }}List() ([]*{{ $T.Name }}, error) {
	ret := []*{{ $T.Name }}{}

	// return if exists
	for _, elem := range {{ $T.Name }}Store {
		ret = append(ret, elem)
	}

	return ret, nil
}

func {{ $T.Name }}Update(in *{{ $T.Name }}) error {
	idx := in.{{ $T.Index }}

	// replace if exists, note we are not dealing with partial updates here
	if _, ok := {{ $T.Name }}Store[idx]; ok {
		{{ $T.Name }}Store[idx] = in
		return nil
	}

	// otherwise return error
	return fmt.Errorf("Entry with %v does not exist", idx)
}

func {{ $T.Name }}Delete(idx string) error {

	// replace if exists, note we are not dealing with partial updates here
	if _, ok := {{ $T.Name }}Store[idx]; ok {
		delete({{ $T.Name }}Store, idx)
		return nil
	}

	// otherwise return error
	return fmt.Errorf("Entry with %v does not exist", idx)
}

{{ range $R := $T.Relations }}
{{/* 
	we need to look up the Model on the other side of the relation
	we use hof's dref custom template function
*/}}

{{- $M := (dref $R.Type $ROOT.DM.Models )}}
{{/* reverse lookup to find the relation which points back at our top-level TYPE for this template*/}}
{{ $D := ( printf "Relations.[:].[Type==%s]" $T.Name) }}
{{ $Reverse := (dref $D $M)}}
/*
{{ yaml $Reverse }}
*/

func {{ $T.Name }}ReadWith{{ $R.Name }}(idx string) (*{{ $T.Name }}, error) {

	val, ok := {{ $T.Name }}Store[idx]
	if !ok {
		return nil, fmt.Errorf("Entry with %v does not exist", idx)
	}

	// make copy, so we don't fill the relation in the store
	ret := *val

	for _, elem := range {{ $R.Type }}Store {
			// make copy, so we don't fill the relation in the store
		local := *elem
		{{ if (eq $R.Reln "HasMany" "ManyToMany") }}
		if local.{{ $Reverse.Name }}.{{ $T.Index }} == ret.{{ $T.Index }} {
			// avoid cyclic decoding by Echo return
			local.{{ $Reverse.Name }} = nil
			ret.{{ $R.Name }} = append(ret.{{ $R.Name }}, &local)	
		}
		{{ else if (eq $R.Reln "BelongsTo") }}
		if local.{{ $M.Index }} == ret.{{ $T.Index }} {
			ret.{{ $R.Name }} = &local
			break
		}
		{{ else }}
		if local.{{ $Reverse.Name }}.{{ $M.Index }} == ret.{{ $T.Index }} {
			ret.{{ $R.Name }} = &local
			break
		}
		{{ end }}
	}

	return &ret, nil
}
{{ end }}
