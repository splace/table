package table_test

import "../table"
import "os"

func ExampleTable(){
	table.Fprint(os.Stdout,"**Firstname**\t**Lastname**\t**Age**\nJill\tSmith\t50\nEve\tJackson\t45")
	// Output:
	// |**Firstname**|**Lastname**|**Age**|
	// |-------------|------------|-------|
	// |     Jill    |    Smith   |   50  |
	// |     Eve     |   Jackson  |   45  |
}


