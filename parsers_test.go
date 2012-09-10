package fantasyfootball

import (
	"fmt"
	"testing"
)

var parserTestDir = "C:/Users/Dustin/Documents/golibs/src/github.com/znconrad5/fantasyfootball/parsed"

func TestDstFile(t *testing.T) {
	dstParser := newDstParser()
	for i := 1; i <= 14; i++ {
		dstParser.parseFile(fmt.Sprintf("%s/DEF-ST_%d.txt", parserTestDir, i), i)
	}
}

func TestQbFile(t *testing.T) {
	qbParser := newQbParser()
	for i := 1; i <= 3; i++ {
		qbParser.parseFile(fmt.Sprintf("%s/QB_%d.txt", parserTestDir, i), i)
	}
}

func TestRbFile(t *testing.T) {
	rbParser := newRbParser()
	for i := 1; i <= 3; i++ {
		rbParser.parseFile(fmt.Sprintf("%s/RB_%d.txt", parserTestDir, i), i)
	}
}

func TestWrFile(t *testing.T) {
	wrParser := newWrParser()
	for i := 1; i <= 3; i++ {
		wrParser.parseFile(fmt.Sprintf("%s/WR_%d.txt", parserTestDir, i), i)
	}
}

func TestTeFile(t *testing.T) {
	teParser := newTeParser()
	for i := 1; i <= 3; i++ {
		teParser.parseFile(fmt.Sprintf("%s/TE_%d.txt", parserTestDir, i), i)
	}
}

func TestKFile(t *testing.T) {
	kParser := newKParser()
	for i := 1; i <= 3; i++ {
		kParser.parseFile(fmt.Sprintf("%s/K_%d.txt", parserTestDir, i), i)
	}
}
