package cli

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func legendLayout(v *gocui.View) {
	v.Title = "Keybindings"
	fmt.Fprintln(v, "^c: Exit")
	fmt.Fprintln(v, "^r: Reset")
	fmt.Fprintln(v, "Space: Flag")
	fmt.Fprintln(v, "Enter: Select")
	fmt.Fprintln(v, "<up>: Move Up")
	fmt.Fprintln(v, "<down>: Move Down")
	fmt.Fprintln(v, "<left>: Move Left")
	fmt.Fprintln(v, "<right>: Move Right")
}
