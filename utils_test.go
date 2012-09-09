package fantasyfootball

import (
	"fmt"
	"testing"
)

type A struct {
	val int
}

func TestChan(t *testing.T) {
	bChan := make(chan bool)
	close(bChan)
	select {
	case <- bChan:
		fmt.Println(<-bChan)
	default:
		fmt.Println("default")
	}
}

func TestPushAndPop(t *testing.T) {
	stack := NewStack()
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

func TestRemoveFirst(t *testing.T) {
	stack := NewStack()
	a := &A{0}
	b := &A{1}
	c := &A{2}
	stack.Push(a)
	stack.Push(b)
	stack.Push(c)
	stack.Remove(a)
	if stack.Pop() != c {
		t.Errorf("Expected %v", c)
	}
	if stack.Pop() != b {
		t.Errorf("Expected %v", b)
	}
	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}
}

func testRemoveMid(t *testing.T) {
	stack := NewStack()
	a := &A{0}
	b := &A{1}
	c := &A{2}
	stack.Push(a)
	stack.Push(b)
	stack.Push(c)
	stack.Remove(b)
	if stack.Pop() != c {
		t.Errorf("Expected %v", c)
	}
	if stack.Pop() != a {
		t.Errorf("Expected %v", b)
	}
	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}
}

func testRemoveLast(t *testing.T) {
	stack := NewStack()
	a := &A{0}
	b := &A{1}
	c := &A{2}
	stack.Push(a)
	stack.Push(b)
	stack.Push(c)
	stack.Remove(c)
	if stack.Pop() != b {
		t.Errorf("Expected %v", c)
	}
	if stack.Pop() != a {
		t.Errorf("Expected %v", b)
	}
	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}
}

func TestRemoveNotExists(t *testing.T) {
	stack := NewStack()
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
