package schemas

import "github.com/MarcusSanchez/go-z"

var ListTitleSchema = z.String().
	NotEmpty("title must be set").
	Min(1, "title must be at least 1 character long").
	Max(100, "title must be at most 100 characters long")
