package fantasyfootball

import (
	"fmt"
	"os"
	"testing"
)

var parserTestDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
var parserTestStartWeek = 1
var parserTestEndWeek = 14

func TestDstFile(t *testing.T) {
	dstParser := newDstParser()
	for i := parserTestStartWeek; i <= parserTestEndWeek; i++ {
		dstParser.parseFile(fmt.Sprintf("%s/DEF-ST_%d.txt", parserTestDir, i), i)
	}
}

func TestQbFile(t *testing.T) {
	qbParser := newQbParser()
	for i := parserTestStartWeek; i <= parserTestEndWeek; i++ {
		qbParser.parseFile(fmt.Sprintf("%s/QB_%d.txt", parserTestDir, i), i)
	}
}

func TestRbFile(t *testing.T) {
	rbParser := newRbParser()
	for i := parserTestStartWeek; i <= parserTestEndWeek; i++ {
		rbParser.parseFile(fmt.Sprintf("%s/RB_%d.txt", parserTestDir, i), i)
	}
}

func TestWrFile(t *testing.T) {
	wrParser := newWrParser()
	for i := parserTestStartWeek; i <= parserTestEndWeek; i++ {
		wrParser.parseFile(fmt.Sprintf("%s/WR_%d.txt", parserTestDir, i), i)
	}
}

func TestTeFile(t *testing.T) {
	teParser := newTeParser()
	for i := parserTestStartWeek; i <= parserTestEndWeek; i++ {
		teParser.parseFile(fmt.Sprintf("%s/TE_%d.txt", parserTestDir, i), i)
	}
}

func TestKFile(t *testing.T) {
	kParser := newKParser()
	for i := parserTestStartWeek; i <= parserTestEndWeek; i++ {
		kParser.parseFile(fmt.Sprintf("%s/K_%d.txt", parserTestDir, i), i)
	}
}
