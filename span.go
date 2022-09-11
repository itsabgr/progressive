package span

import (
	"errors"
	"fmt"
)

type number interface {
	~float32 | ~float64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type span[T number] struct {
	Start, End T
}

func (x span[T]) IndexIn(list []span[T]) int {

	for i, span := range list {
		if span.Equal(x) {
			return i
		}
	}
	return -1
}

func (x span[T]) Len() T {
	l := x.End - x.Start
	if l < 0 {
		return -l
	}
	return l
}

func (x span[T]) Equal(y span[T]) bool {
	r := x.End == y.End && x.Start == y.Start
	return r
}
func (x span[T]) String() string {
	return fmt.Sprintf("[%+v, %+v)", x.Start, x.End)
}
func (x span[T]) relationTo(y span[T]) relation {
	lxly := compare(x.Start, y.Start)
	lxgy := compare(x.Start, y.End)
	gxly := compare(x.End, y.Start)
	gxgy := compare(x.End, y.End)
	switch {
	case lxly == 0 && gxgy == 0:
		return relationEqual
	case gxly < 0:
		return relationBefore
	case lxly < 0 && gxly == 0 && gxgy < 0:
		return relationMeets
	case gxly == 0:
		return relationOverlaps
	case lxly > 0 && lxgy == 0 && gxgy > 0:
		return relationMetBy
	case lxgy == 0:
		return relationOverlappedBy
	case lxgy > 0:
		return relationAfter
	case lxly < 0 && gxgy < 0:
		return relationOverlaps
	case lxly < 0 && gxgy == 0:
		return relationFinishedBy
	case lxly < 0 && gxgy > 0:
		return relationContains
	case lxly == 0 && gxgy < 0:
		return relationStarts
	case lxly == 0 && gxgy > 0:
		return relationStartedBy
	case lxly > 0 && gxgy < 0:
		return relationDuring
	case lxly > 0 && gxgy == 0:
		return relationFinishes
	case lxly > 0 && gxgy > 0:
		return relationOverlappedBy
	default:
		return relationUnknown
	}
}
func (x span[T]) SubFrom(list []span[T]) []span[T] {
	list2 := make([]span[T], 0, len(list)+1)
	for _, y := range list {
		for _, s := range y.sub(x) {
			if s.Len() != 0 && s.IndexIn(list2) < 0 {
				list2 = append(list2, s)
			}
		}

	}
	return list2
}
func (x span[T]) Contains(i T) bool {
	return x.Start <= i && x.End > i
}

func (x span[T]) AddTo(list []span[T]) []span[T] {
	list2 := make([]span[T], 0, len(list)+1)
	for _, y := range list {
		switch x.relationTo(y) {
		case relationBefore,
			relationAfter:
			list2 = append(list2, y)
		case relationMeets,
			relationOverlaps,
			relationFinishedBy:
			x.End = y.End
		case relationStarts,
			relationDuring,
			relationFinishes:
			x = y
		case relationOverlappedBy,
			relationMetBy:
			x.Start = y.Start
		}
	}
	list2 = append(list2, x)
	return list2
}

func (x span[T]) sub(y span[T]) []span[T] {
	switch x.relationTo(y) {
	case relationStarts,
		relationEqual,
		relationDuring,
		relationFinishes:
		return nil
	case relationBefore,
		relationMeets,
		relationMetBy,
		relationAfter:
		return []span[T]{x}
	case relationOverlaps,
		relationFinishedBy:
		return []span[T]{{x.Start, y.Start}}
	case relationContains:
		return []span[T]{{x.Start, y.Start}, {y.End, x.End}}
	case relationStartedBy,
		relationOverlappedBy:
		return []span[T]{{y.End, x.End}}
	}
	panic(errors.New("span: unreachable reached"))
}
