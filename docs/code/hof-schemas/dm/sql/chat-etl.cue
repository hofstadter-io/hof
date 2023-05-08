package sql

import (
	"strings"

	"github.com/hofstadter-io/hof/schema/dm/fields"
)

ChatETL: {
	// input  models  model     fields    value    nested map (relns)
	Original: Models: [string]: [string]: string | {[string]: string}
	Datamodel: Models: {
		for m, M in Original.Models {
			(m): {
				Name: m
				Fields: {
					for f, F in M {
						// regular field
						if !strings.HasPrefix(f, "$") {
							(f): {
								Name: f
								[
									if F == "string" {Varchar},
									if F == "int" {fields.Int},
									if F == "bool" {fields.Bool},
									if F == "float" {fields.Float},
									if F == "uuid" {fields.UUID},
									if F == "datetime" {fields.Datetime},
									if F == "email" {fields.Email},
									if F == "password" {fields.Password},
									if F == "url" {Varchar},
									"UNKNOWN TYPE: \(F)" & false,
								][0]
							}
						}

						// $relations
						if f == "$relations" {
							for f2, F2 in F {
								(f2): {
									fields.UUID
									Name: f2
									Relation: {
										Name:  F2.name
										Type:  F2.type
										Other: F2.model
									}
								}
							}
						}

						// special fields
						if f == "id" {
							(f): SQL: PrimaryKey: true
						}
					}
				}
			}
		}
	}
}
