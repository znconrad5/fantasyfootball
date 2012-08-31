package fantasyfootball

import (
	"fmt"
	"testing"
)

type A struct {
	val int
}

func TestPushAndPop(t *testing.T) {
	stack := New()
	for i:=0; i<200; i++ {
		stack.Push(&A{i})
	}
	for i:=199; i>=0; i-- {
		v := stack.Pop().(*A)
		if (v.val != i) {
			t.Errorf("Expected %d, got %d", i, v.val)
		}
	}
}

func TestSlice(t *testing.T) {
	sli1 := []int{1,2,3,4,5}
	sli2 := sli1[:4]
	sli2 = append(sli2, 6)
	fmt.Println(sli1)
	fmt.Println(sli2)
}

func TestRemoveExists(t *testing.T) {
	stack := New()
	a := &A{0}
	b := &A{1}
	c := &A{2}
	d := &A{3}
	stack.Push(a)
	stack.Push(b)
	stack.Push(c)
	stack.Push(d)
	stack.Remove(a)
	stack.Remove(d)
	if stack.Pop() != c {
		t.Errorf("Expected %v", c)
	}
}

func TestRemoveNotExists(t *testing.T) {
	stack := New()
	a := &A{0}
	b := &A{1}
	c := &A{2}
	stack.Push(a)
	stack.Push(b)
	stack.Push(c)
	d := &A{3}
	defer func() {
	        if r := recover(); r != nil {
	            // expected
	        }
	    }()
	stack.Remove(d)
	t.Fail()
}
