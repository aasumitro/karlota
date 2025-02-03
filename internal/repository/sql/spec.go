package sql

import (
	"fmt"
	"strings"
)

type Specification interface {
	GetQuery() string
	GetValues() []any
	GetJoin() (string, string)
}

type joinSpecification struct {
	specifications []Specification
	separator      string
}

func (s joinSpecification) GetQuery() string {
	queries := make([]string, 0, len(s.specifications))

	for _, spec := range s.specifications {
		queries = append(queries, spec.GetQuery())
	}

	return strings.Join(queries, fmt.Sprintf(" %s ", s.separator))
}

func (s joinSpecification) GetValues() []any {
	values := make([]any, 0)

	for _, spec := range s.specifications {
		values = append(values, spec.GetValues()...)
	}

	return values
}

func (s joinSpecification) GetJoin() (string, string) {
	return "", ""
}

func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

func Or(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "OR",
	}
}

type notSpecification struct {
	Specification
}

func (s notSpecification) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.Specification.GetQuery())
}

func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

type binaryOperatorSpecification[T any] struct {
	field    string
	operator string
	value    T
}

func (s binaryOperatorSpecification[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

func (s binaryOperatorSpecification[T]) GetValues() []any {
	return []any{s.value}
}

func (s binaryOperatorSpecification[T]) GetJoin() (string, string) {
	return "", ""
}

func Equal[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

func NotEqual[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<>",
		value:    value,
	}
}

func GreaterThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">",
		value:    value,
	}
}

func GreaterOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

func LessThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<",
		value:    value,
	}
}

func LessOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<=",
		value:    value,
	}
}

func In[T any](field string, value []T) Specification {
	return binaryOperatorSpecification[[]T]{
		field:    field,
		operator: "IN",
		value:    value,
	}
}

type stringSpecification string

func (s stringSpecification) GetQuery() string {
	return string(s)
}

func (s stringSpecification) GetValues() []any {
	return nil
}

func (s stringSpecification) GetJoin() (string, string) {
	return "", ""
}

func WithInSubquery(field string, subquery string) Specification {
	return stringSpecification(fmt.Sprintf("%s IN (%s)", field, subquery))
}

func IsNull(field string) Specification {
	return stringSpecification(fmt.Sprintf("%s IS NULL", field))
}

type relationSpecification struct {
	// specifications []Specification
	// separator      string
	joinType      string // e.g., "INNER JOIN", "LEFT JOIN"
	joinCondition string // e.g., "users ON users.id = posts.user_id"
}

func (s relationSpecification) GetQuery() string {
	// Return an empty string as this is a join specification,
	// and it doesn't contribute to the WHERE clause directly.
	return ""
}

func (s relationSpecification) GetValues() []any {
	// Return an empty slice as joins don't typically have bound values.
	return nil
}

func (s relationSpecification) GetJoin() (string, string) {
	if s.joinCondition != "" {
		return s.joinType, s.joinCondition
	}
	return "", ""
}

func Join(joinType, joinCondition string) Specification {
	return relationSpecification{
		joinType:      joinType,
		joinCondition: joinCondition,
	}
}
