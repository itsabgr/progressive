package span

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"sort"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Span[T Number] struct {
	Start, End T
}

func (x Span[T]) IndexIn(list []Span[T]) int {

	for i, span := range list {
		if span.Equal(x) {
			return i
		}
	}
	return -1
}

func Sort[T Number](list []Span[T]) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Start < list[j].Start
	})
}

func (x Span[T]) Len() uint {
	l := x.End - x.Start
	if l < 0 {
		return uint(-l)
	}
	return uint(l)
}

func (x Span[T]) Equal(y Span[T]) bool {
	r := x.End == y.End && x.Start == y.Start
	return r
}
func (x Span[T]) String() string {
	return fmt.Sprintf("[%+v, %+v)", x.Start, x.End)
}
func (x Span[T]) RelationTo(y Span[T]) Relation {
	lxly := compare(x.Start, y.Start)
	lxgy := compare(x.Start, y.End)
	gxly := compare(x.End, y.Start)
	gxgy := compare(x.End, y.End)
	switch {
	case lxly == 0 && gxgy == 0:
		return RelationEqual
	case gxly < 0:
		return RelationBefore
	case lxly < 0 && gxly == 0 && gxgy < 0:
		return RelationMeets
	case gxly == 0:
		return RelationOverlaps
	case lxly > 0 && lxgy == 0 && gxgy > 0:
		return RelationMetBy
	case lxgy == 0:
		return RelationOverlappedBy
	case lxgy > 0:
		return RelationAfter
	case lxly < 0 && gxgy < 0:
		return RelationOverlaps
	case lxly < 0 && gxgy == 0:
		return RelationFinishedBy
	case lxly < 0 && gxgy > 0:
		return RelationContains
	case lxly == 0 && gxgy < 0:
		return RelationStarts
	case lxly == 0 && gxgy > 0:
		return RelationStartedBy
	case lxly > 0 && gxgy < 0:
		return RelationDuring
	case lxly > 0 && gxgy == 0:
		return RelationFinishes
	case lxly > 0 && gxgy > 0:
		return RelationOverlappedBy
	default:
		return RelationUnknown
	}
}
func (x Span[T]) SubFrom(list []Span[T]) []Span[T] {
	list2 := make([]Span[T], 0, len(list)+1)
	for _, y := range list {
		for _, s := range y.sub(x) {
			if s.Len() != 0 && s.IndexIn(list2) < 0 {
				list2 = append(list2, s)
			}
		}

	}

	return list2
}
func (x Span[T]) Contains(i T) bool {
	return x.Start <= i && x.End > i
}

func (x Span[T]) AddTo(list []Span[T]) []Span[T] {
	if len(list) == 0 {
		return []Span[T]{x}
	}
	list2 := make([]Span[T], 0, len(list)+1)
	tobe := x
	for _, y := range list {
		if result := y.add(tobe); len(result) == 1 {
			tobe = result[0]
		} else {
			list2 = append(list2, y)
		}
	}
	list2 = append(list2, tobe)
	return list2
}

func (x Span[T]) add(y Span[T]) []Span[T] {
	switch x.RelationTo(y) {
	case RelationBefore,
		RelationAfter:
		return []Span[T]{x, y}
	case RelationMeets,
		RelationOverlaps,
		RelationFinishedBy:
		return []Span[T]{{x.Start, y.End}}
	case RelationContains,
		RelationEqual,
		RelationStartedBy:
		return []Span[T]{x}
	case RelationStarts,
		RelationDuring,
		RelationFinishes:
		return []Span[T]{y}
	case RelationOverlappedBy,
		RelationMetBy:
		return []Span[T]{{y.Start, x.End}}
	}
	panic(errors.New("span: unreachable reached"))
}

func (x Span[T]) sub(y Span[T]) []Span[T] {
	switch x.RelationTo(y) {
	case RelationStarts,
		RelationEqual,
		RelationDuring,
		RelationFinishes:
		return nil
	case RelationBefore,
		RelationMeets,
		RelationMetBy,
		RelationAfter:
		return []Span[T]{x}
	case RelationOverlaps,
		RelationFinishedBy:
		return []Span[T]{{x.Start, y.Start}}
	case RelationContains:
		return []Span[T]{{x.Start, y.Start}, {y.End, x.End}}
	case RelationStartedBy,
		RelationOverlappedBy:
		return []Span[T]{{y.End, x.End}}
	}
	panic(errors.New("span: unreachable reached"))
}
