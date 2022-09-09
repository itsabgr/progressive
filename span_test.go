package span

import (
	"fmt"
	"testing"
)

func TestSpan(t *testing.T) {
	var list []Span[int]
	list = Span[int]{50, 100}.AddTo(list)
	list = Span[int]{30, 110}.AddTo(list)
	list = Span[int]{200, 500}.AddTo(list)
	list = Span[int]{200, 500}.AddTo(list)
	list = Span[int]{1, 600}.AddTo(list)
	fmt.Println(list)

}
