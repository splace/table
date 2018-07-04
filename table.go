package table

import "strings"
import "bufio"
import "fmt"
import "strconv"
import "io"
import "sort"

var Writer io.Writer
var SortColumn  int
var NumericNotAlphaSort bool
type codePoint []byte

// write a code point repeatedly
func (c codePoint) repeat(w int){
	for i:=0;i<w;i++{
		Writer.Write(c)
	}
}

var justifierPadding codePoint

type rowFormatting struct{
	left,padding,divider,right codePoint
}

// see Print
func Fprint(w io.Writer,tabulated string,headerRows int, justifiers ...func(string,int)) {
	Writer=w
	Print(tabulated,headerRows, justifiers...)
} 

// see Print
func Printf(format string, tabulated string,headerRows int, justifiers ...func(string,int)) {
	Format=format
	Print(tabulated,headerRows, justifiers...)
}

// see Print
func Fprintf(w io.Writer,format string, tabulated string,headerRows int, justifiers ...func(string,int)) {
	Writer=w
	Format=format
	Print(tabulated,headerRows, justifiers...)
}

// write string as text table, monospaced font assumed, rows from lines, columns from tab character.
// headerrows - specify how many rows, at the top, formatted specially.
// justifiers - use by columns, missing:use defaut, len=1:use for all cells, len=n:use n'th for n'th column  
// Not thread safe, uses globals, can be used multiple, fixed count, times using mulitple imports.
// Unicode supporting.
// many built-in table format styles, set global var `Format`
// output written to global var `Writer`
func Print(tabulated string,headerRows int, justifiers ...func(string,int)) {
	// find max rows/widths, record cell strings
	var columnMaxWidths []int 
	var cells [][]string
	lineScanner := bufio.NewScanner(strings.NewReader(tabulated))
	for lineScanner.Scan() {
		rowCells:=strings.Split(lineScanner.Text(), "\t")
		if needed:=len(rowCells)-len(columnMaxWidths); needed>0 {
			columnMaxWidths=append(columnMaxWidths,make([]int,needed)...)
		}
		for ci:=range(rowCells){
			if len(rowCells[ci])>columnMaxWidths[ci]{
				columnMaxWidths[ci]=len(rowCells[ci])
			} 
		}
		// order by function
		cells=append(cells,rowCells)
	}

	// oder sortColumn by NumericNotAlphaSort
	if SortColumn>0{
		if headerRows<len(cells){
			if headerRows<0 {
				if NumericNotAlphaSort{
					sort.Sort(byColumnNumeric{byColumn{SortColumn-1,cells}})
					}else{
					sort.Sort(byColumnAlpha{byColumn{SortColumn-1,cells}})
				}
			}else {	
				if NumericNotAlphaSort{
					sort.Sort(byColumnNumeric{byColumn{SortColumn-1,cells[headerRows:]}})
				}else{
					sort.Sort(byColumnAlpha{byColumn{SortColumn-1,cells[headerRows:]}})
				}
			}
		}
	}
	
	// determine justifier used for a column
	justifier := func(c int)func(string,int){
		if c<len(justifiers) {
			return justifiers[c]
		}
		if len(justifiers)==1{
			return justifiers[0]
		}
		return DefaultJustifier
	}
	
	// use a scanner to split format string into individual UTF8 code points
	runeScanner := bufio.NewScanner(strings.NewReader(Format))
	runeScanner.Split(bufio.ScanRunes)
	
	// scan four code points, for a row.
	scanRowFormatting:=func() *rowFormatting {
		rf:=new(rowFormatting)
		runeScanner.Scan()			
		rf.left=codePoint(runeScanner.Bytes())
		runeScanner.Scan()			
		rf.padding=codePoint(runeScanner.Bytes())
		runeScanner.Scan()			
		rf.divider=codePoint(runeScanner.Bytes())
		if !runeScanner.Scan() {return nil}
		rf.right=codePoint(runeScanner.Bytes())
		return rf
	}	

	// write a contentless row, if formatting presnt.
	writeRow:=func(rf *rowFormatting) {
		if rf==nil{return}
		Writer.Write(rf.left)
		for column,width:=range(columnMaxWidths){
			rf.padding.repeat(width)
			if column<len(columnMaxWidths)-1 {
				Writer.Write(rf.divider)
			}
		}
		Writer.Write(rf.right)
		fmt.Fprintln(Writer)
	}

	// parse row type formatting blocks from Format, use helpful assumptions when not all blocks present.
	var dividerRowFormatting,cellRowFormatting,topRowFormatting *rowFormatting
	firstRowFormatting:=scanRowFormatting()
	if firstRowFormatting==nil{
		fmt.Fprintf(Writer, "Format `%s` needs to have at least 4 characters.",Format)
		return 
	}
	secondRowFormatting:=scanRowFormatting()
	if secondRowFormatting==nil{
		secondRowFormatting=firstRowFormatting
	}
	thirdRowFormatting:=scanRowFormatting()
	if thirdRowFormatting==nil{
		dividerRowFormatting=firstRowFormatting
		cellRowFormatting=secondRowFormatting
		topRowFormatting=nil
	}else{
		dividerRowFormatting=secondRowFormatting
		cellRowFormatting=thirdRowFormatting
		topRowFormatting=firstRowFormatting
	}

	// write table
	writeRow(topRowFormatting)
	justifierPadding = cellRowFormatting.padding
	
	for row:=range cells{
		if row==headerRows{
			writeRow(dividerRowFormatting)
		}
		Writer.Write(cellRowFormatting.left)
		for column,cell:=range(cells[row]){
			justifier(column)(cell,columnMaxWidths[column])
			if column<len(columnMaxWidths)-1 {
				Writer.Write(cellRowFormatting.divider)
			}
		}
		Writer.Write(cellRowFormatting.right)
		fmt.Fprintln(Writer)
	}
	writeRow(scanRowFormatting())
}

// used to justify when no specific justifiers provided
var DefaultJustifier = Centred

func RightJustified(c string,w int){
	justifierPadding.repeat(w-len([]rune(c)))
	fmt.Fprint(Writer,c)
}

func LeftJustified(c string,w int){
	fmt.Fprint(Writer,c)
	justifierPadding.repeat(w-len([]rune(c)))
}

func Centred(c string,w int){
	lc:=len([]rune(c))
	offset:=((w-lc+1)/2)
	justifierPadding.repeat(offset)
	fmt.Fprint(Writer,c)
	justifierPadding.repeat(w-lc-offset)
}

func NumberBoolJustified(c string,w int){	
	_, err := strconv.ParseBool(c)
	if err==nil{
		Centred(c,w)
	}
	NumbersRightJustified(c,w)
}


func NumbersRightJustified(c string,w int){	
	_, err := strconv.ParseInt(c, 10, 64)
	if err==nil{
		RightJustified(c,w)
	}
	LeftJustified(c,w)
}

type byColumn  struct{
	ColumnIndex int
	Rows [][]string
}

func (a byColumn) Len() int           { return len(a.Rows) }
func (a byColumn) Swap(i, j int)      { a.Rows[i], a.Rows[j] = a.Rows[j], a.Rows[i] }

type byColumnAlpha  struct{
	byColumn
}

func (a byColumnAlpha) Less(i, j int) bool { return a.Rows[i][a.ColumnIndex] < a.Rows[j][a.ColumnIndex] }


type byColumnNumeric  struct{
	byColumn
}

func (a byColumnNumeric) Less(i, j int) bool { 
	v1,err1:= strconv.ParseFloat(a.Rows[i][a.ColumnIndex],64)
	v2,err2:= strconv.ParseFloat(a.Rows[j][a.ColumnIndex],64)
	if err1==nil && err2==nil {
		return v1 < v2
	}
	return err1==nil
}

