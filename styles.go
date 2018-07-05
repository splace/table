package table

var Style = MarkupStyle

// built-in Styles
 const (
	ASCIIStyle string = "+-++| ||"
	ASCIIBoxStyle string = "+-+++-++| ||+-++"
	MarkupStyle string = "|-||| ||"
	SquareStyle string = "┌─┬┐├─┼┤│ ││└─┴┘"
	SquareThickHeaderDivideStyle string = "┌─┬┐┝━┿┥│ ││└─┴┘"
	SquareDoubleEdgeStyle string = "╔═╤╗╟─┼╢║ │║╚═╧╝"
	SquareSingleVerticalStyle string = "╔═╤╗╠═╪╣║ │║╚═╧╝"
	SquareDoubleEdgeHorizontalStyle string = "╔═╤╗╠═╪╣║ │║╚═╧╝"
	SquareDoubleTopBottomStyle string = "╒═╤╕├─┼┤│ ││╘═╧╛"
	SquareDoubleSidesStyle string = "╓─┬╖╟─┼╢║ │║╙─┴╜"
	SquareDoubleTopStyle string = "╒═╤╕├─┼┤│ ││└─┴┘"
	SquareDoubleDivideStyle string = "┌─┬┐╞═╪╡│ ││└─┴┘"
	SquareDoubleBottomStyle string = "┌─┬┐├─┼┤│ ││╘═╧╛"
	SquareDoubleRightStyle string = "┌─┬╖├─┼┤│ ││└─┴╜"
	SquareDoubleLeftStyle string = "╓─┬┐╟─┼┤║ ││╙─┴┘"
	SquareDoubleInsideStyle string = "┌─╥┐╞═╬╡│ ║│└─╨┘"
	SquareDoubleInsideVerticalStyle string = "┌─╥┐├─╫┤│ ║│└─╨┘"
	SquareDoubleInsideHorizontalStyle string = "┌─┬┐╞═╪╡│ ││└─┴┘"
	RoundedStyle string = "╭─┬╮├─┼┤│ ││╰─┴╯"
	RoundedDoubleInsideStyle string = "╭─╥╮╞═╬╡│ ║│╰─╨╯"
	RoundedDoubleInsideHorizontalStyle string = "╭─┬╮╞═╪╡│ ││╰─┴╯"
	RoundedDoubleInsideVerticalStyle string = "╭─╥╮├─╫┤│ ║│╰─╨╯"
	DoubleStyle string = "╔═╦╗╠═╬╣║ ║║╚═╩╝"
	DoubleVerticalStyle string = "╓─╥╖╟─╫╢║ ║║╙─╨╜"
	DoubleHorizontalStyle string = "╒═╤╕╞═╪╡│ ││╘═╧╛"
)
