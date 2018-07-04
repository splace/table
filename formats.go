package table

var Format = MarkupFormat

// built-in formats
 const (
	ASCIIFormat string = "+-++| ||"
	ASCIIBoxFormat string = "+-+++-++| ||+-++"
	MarkupFormat string = "|-||| ||"
	SquareFormat string = "┌─┬┐├─┼┤│ ││└─┴┘"
	SquareThickHeaderDivideFormat string = "┌─┬┐┝━┿┥│ ││└─┴┘"
	SquareDoubleEdgeFormat string = "╔═╤╗╟─┼╢║ │║╚═╧╝"
	SquareSingleVerticalFormat string = "╔═╤╗╠═╪╣║ │║╚═╧╝"
	SquareDoubleEdgeHorizontalFormat string = "╔═╤╗╠═╪╣║ │║╚═╧╝"
	SquareDoubleTopBottomFormat string = "╒═╤╕├─┼┤│ ││╘═╧╛"
	SquareDoubleSidesFormat string = "╓─┬╖╟─┼╢║ │║╙─┴╜"
	SquareDoubleTopFormat string = "╒═╤╕├─┼┤│ ││└─┴┘"
	SquareDoubleDivideFormat string = "┌─┬┐╞═╪╡│ ││└─┴┘"
	SquareDoubleBottomFormat string = "┌─┬┐├─┼┤│ ││╘═╧╛"
	SquareDoubleRightFormat string = "┌─┬╖├─┼┤│ ││└─┴╜"
	SquareDoubleLeftFormat string = "╓─┬┐╟─┼┤║ ││╙─┴┘"
	SquareDoubleInsideFormat string = "┌─╥┐╞═╬╡│ ║│└─╨┘"
	SquareDoubleInsideVerticalFormat string = "┌─╥┐├─╫┤│ ║│└─╨┘"
	SquareDoubleInsideHorizontalFormat string = "┌─┬┐╞═╪╡│ ││└─┴┘"
	RoundedFormat string = "╭─┬╮├─┼┤│ ││╰─┴╯"
	RoundedDoubleInsideFormat string = "╭─╥╮╞═╬╡│ ║│╰─╨╯"
	RoundedDoubleInsideHorizontalFormat string = "╭─┬╮╞═╪╡│ ││╰─┴╯"
	RoundedDoubleInsideVerticalFormat string = "╭─╥╮├─╫┤│ ║│╰─╨╯"
	DoubleFormat string = "╔═╦╗╠═╬╣║ ║║╚═╩╝"
	DoubleVerticalFormat string = "╓─╥╖╟─╫╢║ ║║╙─╨╜"
	DoubleHorizontalFormat string = "╒═╤╕╞═╪╡│ ││╘═╧╛"
)
