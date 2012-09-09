package fantasyfootball

import (
	"fmt"
	"testing"
)

var parserTestDir = "C:/Users/Dustin/Documents/golibs/src/github.com/znconrad5/fantasyfootball/data"

func TestDstFile(t *testing.T) {
	dstParser := newDstParser()
	for i := 1; i <= 14; i++ {
		dstParser.parseFile(fmt.Sprintf("%s/def-st_%d.txt", parserTestDir, i), i)
	}
}

func TestQbFile(t *testing.T) {
	qbParser := newQbParser()
	for i := 1; i <= 3; i++ {
		qbParser.parseFile(fmt.Sprintf("%s/qb_%d.txt", parserTestDir, i), i)
	}
}

func TestRbFile(t *testing.T) {
	rbParser := newRbParser()
	for i := 1; i <= 3; i++ {
		rbParser.parseFile(fmt.Sprintf("%s/rb_%d.txt", parserTestDir, i), i)
	}
}

func TestWrFile(t *testing.T) {
	wrParser := newWrParser()
	for i := 1; i <= 3; i++ {
		wrParser.parseFile(fmt.Sprintf("%s/wr_%d.txt", parserTestDir, i), i)
	}
}

func TestTeFile(t *testing.T) {
	teParser := newTeParser()
	for i := 1; i <= 3; i++ {
		teParser.parseFile(fmt.Sprintf("%s/te_%d.txt", parserTestDir, i), i)
	}
}

func TestKFile(t *testing.T) {
	kParser := newKParser()
	for i := 1; i <= 3; i++ {
		kParser.parseFile(fmt.Sprintf("%s/k_%d.txt", parserTestDir, i), i)
	}
}
