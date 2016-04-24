package cli

import (
	"fmt"
	"log"

	"github.com/ChrisMcKenzie/minesweeper/board"
	"github.com/jroimartin/gocui"
)

var b = board.NewBoard()

func flagItem(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		// bx, _ := b.Size()
		cx, cy := v.Cursor()

		v.Clear()
		b.Flag(cx/2, cy).Render(v)
	}
	return nil
}

func selectItem(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		// bx, _ := b.Size()
		cx, cy := v.Cursor()

		v.Clear()
		b.Select(cx/2, cy).Render(v)
	}
	return nil
}

func reset(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.Clear()
		b = board.NewBoard()
		b.Render(v)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, sy := b.Size()
		cx, cy := v.Cursor()
		if fy := cy + 1; fy < sy {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		sx, _ := b.Size()
		cx, cy := v.Cursor()
		if fx := cx + 2; fx-sx < sx {
			if err := v.SetCursor(cx+2, cy); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox+2, oy); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if fx := cx - 2; fx != -1 {
			if err := v.SetCursor(cx-2, cy); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox-2, oy); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, reset); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, selectItem); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, flagItem); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	height := ((maxY / 2) - 10) + 20
	width := ((maxX / 2) - 20) + 40
	mainX, mainY := (maxX/2)-20, (maxY/2)-10

	if v, err := g.SetView("legend", maxX-23, 0, maxX-1, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		legendLayout(v)
	}

	statusX, statusY := (maxX/2)-2, mainY-4
	if v, err := g.SetView("status", statusX, statusY, statusX+5, statusY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if !b.Over || b.Won {
			fmt.Fprintf(v, " 😀  ")
		}
	}

	if v, err := g.SetView("main", mainX, mainY, width, height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, "%s", b)

		v.Title = "Minesweeper"
		v.Editable = false
		v.Wrap = true
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}

		if b.Over {
			if status, err := g.View("status"); err != nil {
				status.Clear()
				fmt.Fprintf(status, " 🤕  ")
			}
		}

		if b.Won {
			if status, err := g.View("status"); err != nil {
				status.Clear()
				fmt.Fprintf(status, " 😎  ")
			}
		}

		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(cx+1, cy); err != nil && oy > 0 {
			if err := v.SetOrigin(ox+1, oy); err != nil {
				return err
			}
		}
	}
	return nil
}

func Start() {
	// return
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
