package filters

import "regexp"

type EqualityFilter struct {
	Value string
}

func (f EqualityFilter) IsMatch(value interface{}) bool {
	if valueStr, ok := value.(string); ok {
		return valueStr == f.Value
	}

	return false
}

type RegexFilter struct {
	Regex *regexp.Regexp
}

func (f RegexFilter) IsMatch(value interface{}) bool {
	if valueStr, ok := value.(string); ok {
		return f.Regex.MatchString(valueStr)
	}

	return false
}

type InvertFilter struct {
	InnerFilter Filter
}

func (f InvertFilter) IsMatch(value interface{}) bool {
	return !f.InnerFilter.IsMatch(value)
}

type Filter interface {
	IsMatch(value interface{}) bool
}
