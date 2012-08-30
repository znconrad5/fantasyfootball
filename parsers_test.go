package fantasyfootball

import (
	"fmt"
	"testing"
)

func TestDstFile(t *testing.T) {
	dstParser := newDstParser()
	for i := 1; i <= 14; i++ {
		dstParser.parseFile(fmt.Sprintf("/Users/zachconrad/Documents/go/src/fantasyfootball/data/dst_%d.txt", i), i)
	}
}

func TestQbFile(t *testing.T) {
	qbParser := newQbParser()
	for i := 1; i <= 3; i++ {
		qbParser.parseFile(fmt.Sprintf("/Users/zachconrad/Documents/go/src/fantasyfootball/data/qb_%d.txt", i), i)
	}
}

func TestRbFile(t *testing.T) {
	rbParser := newRbParser()
	for i := 1; i <= 3; i++ {
		rbParser.parseFile(fmt.Sprintf("/Users/zachconrad/Documents/go/src/fantasyfootball/data/rb_%d.txt", i), i)
	}
}

func TestWrFile(t *testing.T) {
	wrParser := newWrParser()
	for i := 1; i <= 3; i++ {
		wrParser.parseFile(fmt.Sprintf("/Users/zachconrad/Documents/go/src/fantasyfootball/data/wr_%d.txt", i), i)
	}
}

func TestTeFile(t *testing.T) {
	teParser := newTeParser()
	for i := 1; i <= 3; i++ {
		teParser.parseFile(fmt.Sprintf("/Users/zachconrad/Documents/go/src/fantasyfootball/data/te_%d.txt", i), i)
	}
}

func TestKFile(t *testing.T) {
	kParser := newKParser()
	for i := 1; i <= 3; i++ {
		kParser.parseFile(fmt.Sprintf("/Users/zachconrad/Documents/go/src/fantasyfootball/data/k_%d.txt", i), i)
	}
}
