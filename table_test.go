package table

import "fmt"
import "testing"
import "strings"

var tableTests = []struct {
	text                                    string
	headerRows                              int
	justifiers                              []func(string, int)
	output                                  string
	outputACIIStyle                         string
	outputACIIStyleSortColumn1              string
	outputACIIStyleSortColumn1Column2ToLeft string
}{
	{
		"1\t2\t3",
		-1,
		nil,
		"|1|2|3|",
		"|1|2|3|",
		"|1|2|3|",
		"|2|1|3|",
	},
	{
		`A	B	C
1	2	3`,
		1,
		nil,
		`|A|B|C|
|-|-|-|
|1|2|3|`,
		`|A|B|C|
+-+-+-+
|1|2|3|`,
		`|A|B|C|
+-+-+-+
|1|2|3|`,
		`|B|A|C|
+-+-+-+
|2|1|3|`,
	},
	{
		`Name	Age	Height(m)
John Doe	47	1.89
Jane Roe	42	1.90
Alan Roe	42	1.90`,
		1,
		nil,
		`|  Name  |Age|Height(m)|
|--------|---|---------|
|John Doe| 47|   1.89  |
|Jane Roe| 42|   1.90  |
|Alan Roe| 42|   1.90  |`,
		`|  Name  |Age|Height(m)|
+--------+---+---------+
|John Doe| 47|   1.89  |
|Jane Roe| 42|   1.90  |
|Alan Roe| 42|   1.90  |`,
		`|  Name  |Age|Height(m)|
+--------+---+---------+
|Alan Roe| 42|   1.90  |
|Jane Roe| 42|   1.90  |
|John Doe| 47|   1.89  |`,
		`|Age|  Name  |Height(m)|
+---+--------+---------+
| 42|Alan Roe|   1.90  |
| 42|Jane Roe|   1.90  |
| 47|John Doe|   1.89  |`,
	},
	{
		`Name	Age	Height(m)
John Doe	47	1.89
Jane Roe	42	1.90
Alan Roe	42	1.90`,
		4,
		nil,
		`|  Name  |Age|Height(m)|
|John Doe| 47|   1.89  |
|Jane Roe| 42|   1.90  |
|Alan Roe| 42|   1.90  |`,
		`|  Name  |Age|Height(m)|
|John Doe| 47|   1.89  |
|Jane Roe| 42|   1.90  |
|Alan Roe| 42|   1.90  |`,
		`|  Name  |Age|Height(m)|
|John Doe| 47|   1.89  |
|Jane Roe| 42|   1.90  |
|Alan Roe| 42|   1.90  |`,
		`|Age|  Name  |Height(m)|
| 47|John Doe|   1.89  |
| 42|Jane Roe|   1.90  |
| 42|Alan Roe|   1.90  |`,
	},
	{
		`Name	Age	Height(m)
John Doe	47	1.89
Jane Roe	42	1.90
Alan Roe	42	1.90`,
		1,
		[]func(string, int){Centred, LeftJustified},
		`|  Name  |Age|Height(m)|
|--------|---|---------|
|John Doe|47 |   1.89  |
|Jane Roe|42 |   1.90  |
|Alan Roe|42 |   1.90  |`,
		`|  Name  |Age|Height(m)|
+--------+---+---------+
|John Doe|47 |   1.89  |
|Jane Roe|42 |   1.90  |
|Alan Roe|42 |   1.90  |`,
		`|  Name  |Age|Height(m)|
+--------+---+---------+
|Alan Roe|42 |   1.90  |
|Jane Roe|42 |   1.90  |
|John Doe|47 |   1.89  |`,
		`|Age|Name    |Height(m)|
+---+--------+---------+
| 42|Alan Roe|   1.90  |
| 42|Jane Roe|   1.90  |
| 47|John Doe|   1.89  |`,
	},
	{
		`Name	Age	Height(m)
John Doe	47	1.89
Jane Roe	42	1.90
Alan Roe	42	1.90`,
		1,
		[]func(string, int){Centred ,MinWidth(RightJustified, 5)},
		`|  Name  |  Age|Height(m)|
|--------|-----|---------|
|John Doe|   47|   1.89  |
|Jane Roe|   42|   1.90  |
|Alan Roe|   42|   1.90  |`,
		`|  Name  |  Age|Height(m)|
+--------+-----+---------+
|John Doe|   47|   1.89  |
|Jane Roe|   42|   1.90  |
|Alan Roe|   42|   1.90  |`,
		`|  Name  |  Age|Height(m)|
+--------+-----+---------+
|Alan Roe|   42|   1.90  |
|Jane Roe|   42|   1.90  |
|John Doe|   47|   1.89  |`,
		`|Age|    Name|Height(m)|
+---+--------+---------+
| 42|Alan Roe|   1.90  |
| 42|Jane Roe|   1.90  |
| 47|John Doe|   1.89  |`,
	},
}

func TestTable(t *testing.T) {
	var buf strings.Builder
	Writer = &buf
	for i, tt := range tableTests {
		HeaderRows = tt.headerRows
		Print(tt.text, tt.justifiers...)
		if buf.String() != fmt.Sprintln(tt.output) {
			t.Fatalf("#%v:found:\n%s wanted:\n%s", i, buf.String(), tt.output)
		}
		buf.Reset()
	}
	Style = ASCIIStyle
	for i, tt := range tableTests {
		HeaderRows = tt.headerRows
		Print(tt.text, tt.justifiers...)
		if buf.String() != fmt.Sprintln(tt.outputACIIStyle) {
			t.Fatalf("#%v:found:\n%s wanted:\n%s", i, buf.String(), tt.outputACIIStyle)
		}
		buf.Reset()
	}
	SortColumn = 1
	for i, tt := range tableTests {
		HeaderRows = tt.headerRows
		Print(tt.text, tt.justifiers...)
		if buf.String() != fmt.Sprintln(tt.outputACIIStyleSortColumn1) {
			t.Fatalf("#%v:found:\n%s wanted:\n%s", i, buf.String(), tt.outputACIIStyleSortColumn1)
		}
		buf.Reset()
	}

	ColumnMapper = MoveToLeftEdge(2)
	for i, tt := range tableTests {
		HeaderRows = tt.headerRows
		Print(tt.text, tt.justifiers...)
		if buf.String() != fmt.Sprintln(tt.outputACIIStyleSortColumn1Column2ToLeft) {
			t.Fatalf("#%v:found:\n%s wanted:\n%s", i, buf.String(), tt.outputACIIStyleSortColumn1Column2ToLeft)
		}
		buf.Reset()
	}
	HeaderRows = 1
	Style = MarkdownStyle
	SortColumn = 0
	ColumnMapper = nil

}
