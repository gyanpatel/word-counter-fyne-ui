package main

import (
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jaytaylor/html2text"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var data = [][]string{{"Word", "Count"}}
var inputVal []string
var skipWordVal []string

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Word Counter")
	saveFile := &fyne.MenuItem{
		Label: "Save result as csv",
		Action: func() {
			saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}

				if write == nil {
					// user canceled
					return
				}

				// save file
				for _, v := range data {

					write.Write([]byte(strings.Join(v, ",") + "\n"))
				}
				defer write.Close()

				myWindow.SetTitle(myWindow.Title() + " - " + write.URI().Name())

			}, myWindow)
			saveDialog.SetFileName("wordcount_result.csv")
			saveDialog.Show()
		},
	}
	about := &fyne.MenuItem{
		Label: "About",
		Action: func() {
			dialog.ShowInformation("About this app", "Counts all the words for the given text \nand prints the occurrence for each word", myWindow)
		},
	}
	menuItems := []*fyne.MenuItem{saveFile, about}
	menu := &fyne.Menu{
		Label: "File",
		Items: menuItems,
	}
	mainMenu := fyne.NewMainMenu(menu)
	myWindow.SetMainMenu(mainMenu)
	myWindow.Resize(fyne.Size{Width: 400, Height: 400})
	url := widget.NewEntry()
	url.SetPlaceHolder("Or Alternatively enter a URL...")
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter text...")
	skipWord := widget.NewMultiLineEntry()
	skipWord.SetPlaceHolder("Enter words to skip separated by space...")
	outLabel := widget.NewLabel("")

	// table
	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})
	table.SetColumnWidth(0, 250)
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Widget: input},
			{Widget: url},
			{Widget: skipWord}},
		OnSubmit: func() {
			if len(url.Text) > 0 {
				urlData, err := http.Get((url.Text))

				if err != nil {
					dialog.ShowInformation("Error", "Unabe to read "+url.Text, myWindow)
					return
				}
				bodyBytes, err := io.ReadAll(urlData.Body)
				if err != nil {
					dialog.ShowInformation("Error", "Unabe to read "+url.Text, myWindow)
					return
				}
				err = urlData.Body.Close()
				if err != nil {
					dialog.ShowInformation("Error", "Unabe to read "+url.Text, myWindow)
					return
				}
				text, err := html2text.FromString(string(bodyBytes))
				if err != nil {
					dialog.ShowInformation("Error", "Unabe to read "+url.Text, myWindow)
					return
				}
				input.Text = text
			}
			data = [][]string{{"Word", "Count"}}
			processWords(outLabel, skipWord.Text, input.Text)

		},
		OnCancel: func() {
			input.SetText("")
			skipWord.SetText("")
			url.SetText("")
			data = [][]string{{"Word", "Count"}}
			table.Refresh()

		},
		SubmitText: "Count",
		CancelText: "Reset",
	}

	// log.Println("tapped", input.Text, skipWord.Text)
	content := container.New(layout.NewGridLayout(1), form, table) //, list, text, check)

	myWindow.SetContent(content)
	myWindow.Show()
	myApp.Run()
}

func processWords(lbl *widget.Label, skipWordStr, inputStr string) {
	// log.Println("Content was:", inputStr, skipWordStr)
	inputVal = strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ToLower(inputStr), "\n", " ")), " ")
	skipWordVal = strings.Split(strings.TrimSpace(strings.ReplaceAll(strings.ToLower(skipWordStr), "\n", " ")), " ")
	countData := countWords()
	// log.Println("countData", countData)
	for _, rec := range countData {
		data = append(data, []string{rec.word, strconv.Itoa(rec.ct)})

	}
	// log.Println(inputVal, skipWordVal, data)
	//str.Set(inputVal)

}

type wordCt struct {
	word string
	ct   int
}

func countWords() (wordDet []wordCt) {
	for _, wp := range inputVal {
		ct := 0
		if skipYn(wp) {
			continue
		}
		for _, wc := range inputVal {
			if strings.Compare(wc, wp) == 0 {
				ct = ct + 1
			}
		}
		skipWordVal = append(skipWordVal, wp)
		wordDet = append(wordDet, wordCt{word: wp, ct: ct})

	}
	sort.Slice(wordDet, func(i, j int) bool {
		return wordDet[i].ct > wordDet[j].ct

	})
	return
}
func skipYn(wp string) bool {
	for _, sk := range skipWordVal {
		// log.Println("sk", wp, sk)

		if strings.Compare(wp, sk) == 0 {
			// log.Println("sk1", wp, sk)
			return true
		}

	}
	return false
}
