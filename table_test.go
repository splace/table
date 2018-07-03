package table

import "fmt"
import "testing"
import "strings"

var tableTests = []struct {
	text       string
	headerRows int        
	justifiers []func(string,int)
	output  string
}{
	{
		"1\t2\t3",
		-1,
		nil,
		"|1|2|3|",
	},
	{
`A	B	C
1	2	3`,
		1,
		nil,
`|A|B|C|
|-|-|-|
|1|2|3|`,
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
|Alan Roe| 42|   1.90  |
|Jane Roe| 42|   1.90  |
|John Doe| 47|   1.89  |`,
	},
}


func TestTable(t *testing.T) {
	var buf strings.Builder
	Writer=&buf
	for i, tt := range tableTests {
		Print(tt.text,tt.headerRows,tt.justifiers...)
		if buf.String()!=fmt.Sprintln(tt.output){
			t.Fatalf("#%v:found:%s wanted:%s",i,buf.String(),tt.output)
		}
		buf.Reset()
	}
}
