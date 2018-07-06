package table

import "os"

func ExampleFprint() {
	Fprint(os.Stdout, "**Firstname**\t**Lastname**\t**Age**\nJill\tSmith\t50\nEve\tJackson\t45")
	// Output:
	// |**Firstname**|**Lastname**|**Age**|
	// |-------------|------------|-------|
	// |     Jill    |    Smith   |   50  |
	// |     Eve     |   Jackson  |   45  |
}

func ExampleFprintf() {
	Fprintf(os.Stdout,BoxStyle ,"**Firstname**\t**Lastname**\t**Age**\nJill\tSmith\t50\nEve\tJackson\t45")
	// Output:
	//┌─────────────┬────────────┬───────┐
	//│**Firstname**│**Lastname**│**Age**│
	//├─────────────┼────────────┼───────┤
	//│     Jill    │    Smith   │   50  │
	//│     Eve     │   Jackson  │   45  │
	//└─────────────┴────────────┴───────┘
}

func ExamplePrint() {
	Writer=os.Stdout
	Style=RoundedDoubleInsideBoxStyle
	ColumnMapper=MoveToLeftEdge(3)
	SortColumn=1
	Print("**Firstname**\t**Lastname**\t**Age**\nJill\tSmith\t50\nEve\tJackson\t45")
	// Output:
	//╭───────╥─────────────╥────────────╮
	//│**Age**║**Firstname**║**Lastname**│
	//╞═══════╬═════════════╬════════════╡
	//│   45  ║     Eve     ║   Jackson  │
	//│   50  ║     Jill    ║    Smith   │
	//╰───────╨─────────────╨────────────╯
}

