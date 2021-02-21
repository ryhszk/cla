package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	a := tview.NewTextView()
	a.SetText("textarea(a)")

	b := tview.NewTextView()
	b.SetTitle("title(b)").
		SetBorder(true)
	b.SetText("bbbbbbbbbb")

	c := tview.NewTextView()
	c.SetTitle("title(c)").
		SetTitleAlign(tview.AlignRight).
		SetBorder(true)
	c.SetText("ccccccccccc")

	flex := tview.NewFlex().
		AddItem(a, 0, 1, false).
		AddItem(b, 0, 1, false).
		AddItem(c, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
