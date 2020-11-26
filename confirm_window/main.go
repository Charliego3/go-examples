package main

import (
	"fmt"
	"github.com/gen2brain/dlgs"
)

//// ConfirmDialog is like the standard Dialog but with an additional confirmation button
//type ConfirmDialog struct {
//	*fyne.dialog
//
//	confirm *widget.Button
//}
//
//// SetConfirmText allows custom text to be set in the confirmation button
//func (d *ConfirmDialog) SetConfirmText(label string) {
//	d.confirm.SetText(label)
//	widget.Refresh(d.win)
//}
//
//// NewConfirm creates a dialog over the specified window for user confirmation.
//// The title is used for the dialog window and message is the content.
//// The callback is executed when the user decides. After creation you should call Show().
//func NewConfirm(title, message string, callback func(bool), parent fyne.Window) *ConfirmDialog {
//	d := newDialog(title, message, theme.QuestionIcon(), callback, parent)
//
//	d.dismiss = &widget.Button{Text: "No", Icon: theme.CancelIcon(),
//		OnTapped: d.Hide,
//	}
//	confirm := &widget.Button{Text: "Yes", Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
//		OnTapped: func() {
//			d.hideWithResponse(true)
//		},
//	}
//	d.setButtons(newButtonList(d.dismiss, confirm))
//
//	return &ConfirmDialog{d, confirm}
//}
//
//// ShowConfirm shows a dialog over the specified window for a user
//// confirmation. The title is used for the dialog window and message is the content.
//// The callback is executed when the user decides.
//func ShowConfirm(title, message string, callback func(bool), parent fyne.Window) {
//	NewConfirm(title, message, callback, parent).Show()
//}

func main() {
	//ok := dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
	//println("OK:", ok)
	//directory, err := dialog.Directory().Title("Load images").Browse()
	//println(directory, err)

	//filename, err := dialog.File().Filter("XML files", "xml").Title("Export to XML").Save()
	//filename, err := dialog.File().Filter("Mp3 audio file", "csv").Load()
	//println(filename, err)

	//item, b, err := dlgs.List("List", "Select item from list:", []string{"Bug", "New Feature", "Improvement"})
	//println(item, b, err)
	//
	//passwd, b, err := dlgs.Password("Password", "Enter your API key:")
	//println(passwd, b, err)
	//
	//answer, err := dlgs.Question("Question", "Are you sure you want to format this media?", true)
	//println(answer, err)

	//file, b, err := dlgs.File("Select a csv", ".csv", false)
	//println(file, b, err)

	//file, b, err := dlgs.FileMulti("Select a csv", ".csv")
	//fmt.Printf("%v, %v, %v", file, b, err)

	//c, b, err := dlgs.Color("Select a color", "")
	//fmt.Printf("%v, %v, %v", c, b, err)

	//date, b, err := dlgs.Date("Select a Date", "Text", time.Now())
	//fmt.Printf("%v, %v, %v",date, b, err)

	//r, b, err := dlgs.Entry("Entry Title", "Entry Text", "Entry Default Text")
	//fmt.Printf("%v, %v, %v",r, b, err)

	//b, err := dlgs.Info("Info Title", "Info Text")
	//fmt.Printf("%v, %v", b, err)

	b, err := dlgs.Warning("Warning Title", "Warning Text")
	fmt.Printf("%v, %v", b, err)

	b, err = dlgs.Error("Error Title", "Error Text")
	fmt.Printf("%v, %v", b, err)
}
