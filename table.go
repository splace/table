package table

import "strings"
import "bufio"
import "fmt"
import "strconv"
import "io"
import "sort"

// global options
var (
	Writer io.Writer
	HeaderRows int = 1
	Style = MarkupStyle
	ColumnMapper func(int) int // rearrange columns
	SortColumn int
	NumericNotAlphaSort bool
	DefaultCellPrinter = Centred
	DividerEvery int
	FormfeedWithDivider bool 
)

type codePoint []byte

// write a code point repeatedly
func (c codePoint) repeat(w int) {
	for i := 0; i < w; i++ {
		Writer.Write(c)
	}
}

var cellPrinterPadding codePoint

type rowStyling struct {
	left, padding, divider, right codePoint
}

// set global var 'Writer' then call Print.
func Fprint(w io.Writer, tabulated string, cellPrinters ...func(string, int)) {
	Writer = w
	Print(tabulated, cellPrinters...)
}

// set global var 'Style' then call Print.
func Printf(s string, tabulated string, cellPrinters ...func(string, int)) {
	Style = s
	Print(tabulated, cellPrinters...)
}

// set global var's 'Writer' and 'Style' then call Print.
func Fprintf(w io.Writer, s string, tabulated string, cellPrinters ...func(string, int)) {
	Writer = w
	Style = s
	Print(tabulated, cellPrinters...)
}

// write string as text table, mono-spaced font assumed, rows from lines, columns from tab character.
// cellPrinters - use by columns, missing:-use default, len=1:-use for all cells, len=n:-use n'th for n'th column
// Not thread safe, uses globals, but can be used multiple, fixed count, times using multiple imports.
// Unicode supporting.
// many built-in table styles, set global var `Style`
// output written to global var `Writer`
// number of header row set by  var `HeaderRows`
func Print(tabulated string, cellPrinters ...func(string, int)) {
	// find max rows/widths, record cell strings
	var columnMaxWidths []int
	var cells [][]string
	lineScanner := bufio.NewScanner(strings.NewReader(tabulated))
	for lineScanner.Scan() {
		rowCells := strings.Split(lineScanner.Text(), "\t")
		if needed := len(rowCells) - len(columnMaxWidths); needed > 0 {
			columnMaxWidths = append(columnMaxWidths, make([]int, needed)...)
		}
		for ci := range rowCells {
			if len(rowCells[ci]) > columnMaxWidths[ci] {
				columnMaxWidths[ci] = len(rowCells[ci])
			}
		}
		// order by function
		cells = append(cells, rowCells)
	}

	// order sortColumn by NumericNotAlphaSort
	if SortColumn > 0 {
		if HeaderRows < len(cells) {
			if HeaderRows < 0 {
				if NumericNotAlphaSort {
					sort.Sort(byColumnNumeric{byColumn{SortColumn - 1, cells}})
				} else {
					sort.Sort(byColumnAlpha{byColumn{SortColumn - 1, cells}})
				}
			} else {
				if NumericNotAlphaSort {
					sort.Sort(byColumnNumeric{byColumn{SortColumn - 1, cells[HeaderRows:]}})
				} else {
					sort.Sort(byColumnAlpha{byColumn{SortColumn - 1, cells[HeaderRows:]}})
				}
			}
		}
	}

	// find the cellPrinter needed for a column
	cellPrinter := func(c int) func(string, int) {
		if c < len(cellPrinters) {
			return cellPrinters[c]
		}
		if len(cellPrinters) == 1 {
			return cellPrinters[0]
		}
		return DefaultCellPrinter
	}

	// use a scanner to split Style string into individual UTF8 code points
	runeScanner := bufio.NewScanner(strings.NewReader(Style))
	runeScanner.Split(bufio.ScanRunes)

	// scan a row style, 4 code points.
	scanRowStyling := func() *rowStyling {
		rf := new(rowStyling)
		runeScanner.Scan()
		rf.left = codePoint(runeScanner.Bytes())
		runeScanner.Scan()
		rf.padding = codePoint(runeScanner.Bytes())
		runeScanner.Scan()
		rf.divider = codePoint(runeScanner.Bytes())
		if !runeScanner.Scan() {
			return nil
		}
		rf.right = codePoint(runeScanner.Bytes())
		return rf
	}

	// write a content-less row, if Styling present.
	writeRow := func(rf *rowStyling) {
		if rf == nil {
			return
		}
		Writer.Write(rf.left)
		cellPrinterPadding = rf.padding
		if ColumnMapper == nil {
			for column, width := range columnMaxWidths {
				cellPrinter(column)("", width)
				if column < len(columnMaxWidths)-1 {
					Writer.Write(rf.divider)
				}
			}
		} else {
			for column := range columnMaxWidths {
				c := ColumnMapper(column)
				cellPrinter(c)("", columnMaxWidths[c])
				if column < len(columnMaxWidths)-1 {
					Writer.Write(rf.divider)
				}
			}
		}
		Writer.Write(rf.right)
		fmt.Fprintln(Writer)
	}

	// parse row type Styleting blocks from Style, use helpful assumptions when not all blocks present.
	var dividerRowStyling, cellRowStyling, topRowStyling *rowStyling
	firstRowStyling := scanRowStyling()
	if firstRowStyling == nil {
		fmt.Fprintf(Writer, "Style `%s` needs to have at least 4 characters.", Style)
		return
	}
	secondRowStyling := scanRowStyling()
	if secondRowStyling == nil {
		secondRowStyling = firstRowStyling
	}
	thirdRowStyling := scanRowStyling()
	if thirdRowStyling == nil {
		dividerRowStyling = firstRowStyling
		cellRowStyling = secondRowStyling
		topRowStyling = nil
	} else {
		dividerRowStyling = secondRowStyling
		cellRowStyling = thirdRowStyling
		topRowStyling = firstRowStyling
	}

	// write table
	writeRow(topRowStyling)
	cellPrinterPadding = cellRowStyling.padding
	for row := range cells {
		if row-HeaderRows ==0 {
			writeRow(dividerRowStyling)
			cellPrinterPadding = cellRowStyling.padding
		}
		if DividerEvery>0 && row-HeaderRows > DividerEvery && (row - HeaderRows) % DividerEvery ==0{
			writeRow(dividerRowStyling)
			if FormfeedWithDivider {
				Writer.Write([]byte("\f"))
				writeRow(dividerRowStyling)
			}
			cellPrinterPadding = cellRowStyling.padding
		}
		Writer.Write(cellRowStyling.left)
		if ColumnMapper == nil {
			for column, cell := range cells[row] {
				cellPrinter(column)(cell, columnMaxWidths[column])
				if column < len(columnMaxWidths)-1 {
					Writer.Write(cellRowStyling.divider)
				}
			}
		} else {
			for column := range cells[row] {
				c := ColumnMapper(column)
				cellPrinter(c)(cells[row][c], columnMaxWidths[c])
				if column < len(columnMaxWidths)-1 {
					Writer.Write(cellRowStyling.divider)
				}
			}
		}
		Writer.Write(cellRowStyling.right)
		fmt.Fprintln(Writer)
	}
	writeRow(scanRowStyling())
}

// #cellPrinters


func RightJustified(c string, w int) {
	cellPrinterPadding.repeat(w - len([]rune(c)))
	fmt.Fprint(Writer, c)
}

func LeftJustified(c string, w int) {
	fmt.Fprint(Writer, c)
	cellPrinterPadding.repeat(w - len([]rune(c)))
}

func Centred(c string, w int) {
	lc := len([]rune(c))
	offset := ((w - lc + 1) / 2)
	cellPrinterPadding.repeat(offset)
	fmt.Fprint(Writer, c)
	cellPrinterPadding.repeat(w - lc - offset)
}

func NumberBoolJustified(c string, w int) {
	_, err := strconv.ParseBool(c)
	if err == nil {
		Centred(c, w)
		return
	}
	NumbersRightJustified(c, w)
}

func NumbersRightJustified(c string, w int) {
	_, err := strconv.ParseInt(c, 10, 64)
	if err == nil {
		RightJustified(c, w)
		return
	}
	LeftJustified(c, w)
}

// modify a cellPrinter to have a minimum width
func MinWidth(form func(string, int), min uint) func(string, int) {
	m := int(min)
	return func(s string, w int) {
		if w < m {
			form(s, m)
			return
		}
		form(s, w)
	}
}

// #sorters

type byColumn struct {
	ColumnIndex int
	Rows        [][]string
}

func (a byColumn) Len() int      { return len(a.Rows) }
func (a byColumn) Swap(i, j int) { a.Rows[i], a.Rows[j] = a.Rows[j], a.Rows[i] }

type byColumnAlpha struct {
	byColumn
}

func (a byColumnAlpha) Less(i, j int) bool { return a.Rows[i][a.ColumnIndex] < a.Rows[j][a.ColumnIndex] }

type byColumnNumeric struct {
	byColumn
}

func (a byColumnNumeric) Less(i, j int) bool {
	v1, err1 := strconv.ParseFloat(a.Rows[i][a.ColumnIndex], 64)
	v2, err2 := strconv.ParseFloat(a.Rows[j][a.ColumnIndex], 64)
	if err1 == nil && err2 == nil {
		return v1 < v2
	}
	return err1 == nil
}

// #mappers

// column mapper to put a particular column first, columns are numbered from 1, otherwise preserves order.
func MoveToLeftEdge(column uint) func(int) int {
	c := int(column - 1)
	return func(n int) int {
		if n < c {
			return n + 1
		}
		if n == c {
			return 0
		}
		return n
	}
}
