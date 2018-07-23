// package table prettifies tab/line text tables to aligned/formatted/sorted/rearranged and justified text tables.
package table

import "strings"
import "bufio"
import "fmt"
import "strconv"
import "io"
import "sort"

// global options
var (
	Writer              io.Writer
	HeaderRows          = 1
	Style               = MarkdownStyle
	ColumnMapper        func(int) int // rearrange columns
	SortColumn          int
	NumericNotAlphaSort bool
	DefaultCellPrinter  = Centred
	DividerEvery        int
	FormfeedWithDivider bool
)

type codePoint []byte

// write a code point a number of times
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

// set the global var's 'Writer' and 'Style' then call Print.
func Fprintf(w io.Writer, s string, tabulated string, cellPrinters ...func(string, int)) {
	Writer = w
	Style = s
	Print(tabulated, cellPrinters...)
}

// Write 'tabulated' string as text table, rows coming from lines, columns separated by the tab character.
// Mono-spaced font required for alignment.
// cellPrinters - applied to columns:
// * missing - use default
// * len=1 - use for all cells
// * len=n - use n'th for n'th column, use default if column count>n
// Not thread safe, uses globals for options (see variables), however can be used multiple, fixed count, times by using multiple imports and different aliases.
// Unicode supporting.
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

	// order sortColumn
	if SortColumn > 0 {
		if HeaderRows < len(cells) {
			if HeaderRows < 0 {
				if NumericNotAlphaSort {
					sort.Sort(byColumnNumeric{byColumn{cells}})
				} else {
					sort.Sort(byColumnAlpha{byColumn{cells}})
				}
			} else {
				if NumericNotAlphaSort {
					sort.Sort(byColumnNumeric{byColumn{cells[HeaderRows:]}})
				} else {
					sort.Sort(byColumnAlpha{byColumn{cells[HeaderRows:]}})
				}
			}
		}
	}

	// the cellPrinter needed for a column
	cellPrinter := func(c int) func(string, int) {
		if len(cellPrinters) == 1 {
			return cellPrinters[0]
		}
		if c < len(cellPrinters) {
			return cellPrinters[c]
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

	// write a content-less row using a row style, do nothing if nil.
	// used for top/bottom border and divider rows
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
				cellPrinter(column)("", columnMaxWidths[c])
				if column < len(columnMaxWidths)-1 {
					Writer.Write(rf.divider)
				}
			}
		}
		Writer.Write(rf.right)
		fmt.Fprintln(Writer)
	}

	// scan and store row Stylings from Style string, use helpful assumptions when not all blocks present.
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
		if row-HeaderRows == 0 {
			writeRow(dividerRowStyling)
			cellPrinterPadding = cellRowStyling.padding
		}
		if DividerEvery > 0 && row-HeaderRows > DividerEvery && (row-HeaderRows)%DividerEvery == 0 {
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
				if column < len(cells[row])-1 {
					Writer.Write(cellRowStyling.divider)
				}
			}
		} else {
			for column := range cells[row] {
				c := ColumnMapper(column)
				cellPrinter(column)(cells[row][c], columnMaxWidths[c])
				if column < len(columnMaxWidths)-1 {
					Writer.Write(cellRowStyling.divider)
				}
			}
		}
		Writer.Write(cellRowStyling.right)
		fmt.Fprintln(Writer)
	}
	// scan remaining row styling, if any, from Style for bottom border row. 
	writeRow(scanRowStyling())
}

// #cellPrinters

// right justifier printer
func RightJustified(c string, w int) {
	cellPrinterPadding.repeat(w - len([]rune(c)))
	fmt.Fprint(Writer, c)
}

// left justifier printer
func LeftJustified(c string, w int) {
	fmt.Fprint(Writer, c)
	cellPrinterPadding.repeat(w - len([]rune(c)))
}

// centre printer
func Centred(c string, w int) {
	lc := len([]rune(c))
	offset := ((w - lc + 1) / 2)
	cellPrinterPadding.repeat(offset)
	fmt.Fprint(Writer, c)
	cellPrinterPadding.repeat(w - lc - offset)
}

// centre print if a boolean, right justify if a number, default otherwise.
func NumbersBoolJustified(c string, w int) {
	_, err := strconv.ParseBool(c)
	if err == nil {
		Centred(c, w)
		return
	}
	NumbersRightJustified(c, w)
}

// right justify if a number
func NumbersRightJustified(c string, w int) {
	_, err := strconv.ParseInt(c, 10, 64)
	if err == nil {
		RightJustified(c, w)
		return
	}
	DefaultCellPrinter(c, w)
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

// #sorters, implementing sort.Interface 

type byColumn struct {
	Rows        [][]string
}

func (a byColumn) Len() int      { return len(a.Rows) }
func (a byColumn) Swap(i, j int) { a.Rows[i], a.Rows[j] = a.Rows[j], a.Rows[i] }

type byColumnAlpha struct {
	byColumn
}

func (a byColumnAlpha) Less(i, j int) bool { return a.Rows[i][SortColumn-1] < a.Rows[j][SortColumn-1] }

type byColumnNumeric struct {
	byColumn
}

func (a byColumnNumeric) Less(i, j int) bool {
	v1, err1 := strconv.ParseFloat(a.Rows[i][SortColumn-1], 64)
	v2, err2 := strconv.ParseFloat(a.Rows[j][SortColumn-1], 64)
	if err1 == nil && err2 == nil {
		return v1 < v2
	}
	return err1 == nil
}

// #mappers

// returns a column mapper func, that puts a particular column first, (columns start from 1), otherwise preserves order.
func MoveToLeftEdge(column uint) func(int) int {
	c := int(column - 1)
	return func(n int) int {
		if n==0 {
			return c
		}
		if n <= c {
			return n - 1
		}
		return n
	}
}
