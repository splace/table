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
}


func TestTable(t *testing.T) {
	var buf strings.Builder
	Writer=&buf
	for i, tt := range tableTests {
		Print(tt.text,tt.headerRows,tt.justifiers...)
		if buf.String()!=fmt.Sprintln(tt.output){
			t.Fatalf("#%v:found:%s wanted:%s",i,buf.String(),tt.output)
		}
	}
}
