package common

type CriterionOption[T any] struct {
	Include []T `json:"includes"`
	Exclude []T `json:"excludes"` // Operation is using OR for each member and AND operation with Include

	IsAndOperation bool `json:"is_and"` // Operation applied to Include member
}

func (c CriterionOption[T]) HasInclude() bool {
	return len(c.Include) != 0
}

func (c CriterionOption[T]) HasExclude() bool {
	return len(c.Exclude) != 0
}

type IncludeArray[T any] struct {
	Values    []T  `json:"values"`
	IsInclude bool `json:"is_include"`
}
