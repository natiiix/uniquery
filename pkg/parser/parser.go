package parser

import (
	"log"
	"regexp"
	"strings"

	"github.com/natiiix/uniquery/pkg/filters"
)

type QueryPart struct {
	Specifier string
	Filters   []filters.Filter
}

const (
	escapeRune    = '\\'
	quoteRune     = '"'
	specifierRune = '.'
	equalityRune  = '='
	regexRune     = '~'
	invertRune    = '!'
)

const (
	filterEquality = iota
	filterRegex
)

func ParseSinglePart(query []rune) (string, int) {
	sb := strings.Builder{}
	escaped := false
	quoted := false

	for i, r := range query {
		if escaped {
			sb.WriteRune(r)
			escaped = false
		} else if quoted {
			if r == quoteRune {
				quoted = false
			} else {
				sb.WriteRune(r)
			}
		} else {
			switch r {
			case specifierRune, equalityRune, regexRune, invertRune:
				return sb.String(), i

			case escapeRune:
				escaped = true

			case quoteRune:
				quoted = true

			default:
				sb.WriteRune(r)
			}
		}
	}

	if escaped || quoted {
		log.Fatalln("Unexpected end of query - trailing escape or quote")
	}

	return sb.String(), len(query)
}

func ParseSingleFilter(query []rune) (filters.Filter, int) {
	if len(query) <= 0 {
		log.Fatalln("Unexpected end of filter")
		return nil, 0
	}

	switch query[0] {
	case specifierRune:
		return nil, 0

	case equalityRune:
		value, len := ParseSinglePart(query[1:])
		return filters.EqualityFilter{Value: value}, 1 + len

	case regexRune:
		regex, len := ParseSinglePart(query[1:])
		return filters.RegexFilter{Regex: regexp.MustCompile(regex)}, 1 + len

	case invertRune:
		inner, len := ParseSingleFilter(query[1:])
		return filters.InvertFilter{InnerFilter: inner}, 1 + len

	default:
		log.Fatalln("Unexpected filter prefix:", string(query[0]))
		return nil, 0
	}
}

func ParseQuery(query string) []QueryPart {
	queryRunes := []rune(query)
	parts := []QueryPart{}

	for i := 0; i < len(queryRunes); i++ {
		specifier, specifierLength := ParseSinglePart(queryRunes[i:])
		i += specifierLength

		filters := []filters.Filter{}

		for i < len(queryRunes) {
			if filter, filterLength := ParseSingleFilter(queryRunes[i:]); filter != nil {
				filters = append(filters, filter)
				i += filterLength
			} else {
				break
			}
		}

		parts = append(parts, QueryPart{Specifier: specifier, Filters: filters})
	}

	return parts
}
