package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/extensions/eval"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	ActivateTheme("azure light")
	// Proxying widgets requires that the eval extension be enabled.
	InitializeExtension("eval")

	var text TextWidgetProxy

	// textInsert will be called when text is about to be inserted.
	// The text might have been typed or pasted.
	textInsert := func(args []string) {
		fmt.Println("Insert args", args)

		newText := args[2]
		fmt.Printf("  text '%s'\n", newText)

		// Ignore text input that starts with a digit.
		if unicode.IsDigit(rune(newText[0])) {
			fmt.Println("  starts with a digit, so not inserted")
			return
		}

		// Transform new text to uppercase, and then insert it.
		newText = strings.ToUpper(newText)
		args[2] = newText
		text.EvalWrapped(args)
	}

	// textDelete will be called when text is about to be deleted.
	// It may be a single character, or it may be selected text.
	textDelete := func(args []string) {
		fmt.Println("Delete args :", args)

		index1 := args[1]
		index2 := args[2]
		fmt.Printf("  text '%s'\n", text.Get(index1, index2))

		// Only allow deletion of fewer than 5 characters
		count, err := strconv.Atoi(text.Count(Chars(), index1, index2)[0])
		if err != nil {
			fmt.Println(err)
		}
		if count < 5 {
			text.EvalWrapped(args)
		} else {
			fmt.Println("  too long to delete")
		}
	}

	label := Label(Txt(`Try entering and deleting text. Try cutting and pasting.

Entered text will be uppercase. Entered text that starts with a digit will not be inserted.

Selected text longer than 4 characters will not be deleted`),
		Justify("left"),
	)

	text = NewTextWidgetProxy(Text(
		Width(50),
		Height(5),
	))
	text.Insert(END, "Some text")
	text.Register("insert", textInsert)
	text.Register("delete", textDelete)

	opts := Opts{
		Padx("1m"), Pady("2m"),
		Ipadx("1m"), Ipady("1m"),
		Sticky("w"),
	}
	Grid(label, opts)
	Grid(text, opts)
	Grid(TExit(), Columnspan(2), opts)

	Focus(text)
	App.Wait()
}
