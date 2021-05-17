 /*
This program is intented to 
*/


package main

import (
	"github.com/lxn/walk"
	."github.com/lxn/walk/declarative"
    "os/exec"
    "errors"
    "fmt"
    "log"
    "os"
)

func main() {
    mmw := new(MyMainWindow)
    barcodeLE := new(walk.LineEdit)
	err := MainWindow{
		AssignTo: &mmw.MainWindow,
		Title:    "CopyThat",
		Size:     Size{300, 100},
		MinSize:  Size{200, 100},
    Layout: Grid{},
		Children: []Widget{
            Label{
              Column: 0,
              Row: 0,
              Text: "Scan barcode: ",
            },
            LineEdit{
              Column: 1,
              Row: 0,
              AssignTo: &barcodeLE,
              // Allows scan return function to trigger this event
              OnKeyDown: func(key walk.Key) {
                if key == walk.KeyReturn {
                    barcodeText:= barcodeLE.Text()
                    if barcodeText == "" {
                        mmw.display("Warning", "The bardcode text box is emply.")
                        return
                    }
                    err := createZebraLabel(barcodeText)
                    if err != nil{
                        mmw.display("Error", fmt.Sprintf("There was an error while printing the label.\n%s",err.Error()))
                    }
                    // reset label LineEdit
                    barcodeLE.SetText("")
                  }
                },
              },
            PushButton{
              Column: 1,
              Row: 2,
              Text: "Generate",
              OnClicked: func(){
                barcodeText:= barcodeLE.Text()
                if barcodeText == "" {
                    mmw.display("Warning", "The bardcode text box is emply.")
                    return
                }
                err := createZebraLabel(barcodeText)
                if err != nil{
                    mmw.display("Error", fmt.Sprintf("There was an error while printing the label.\n%s",err.Error()))
                }
                // reset label LineEdit
                barcodeLE.SetText("")
              },
            },
        },
	}.Create()
	checkError("Error making window", err, true)
	mmw.Run()
}

/* Prints zebra label */
func createZebraLabel(label_val string) error {
  print_script := fmt.Sprintf("^XA^FO30,25^BXN,2,200,18,18,~,1^FD%s^FS"+ //barcode
        "^FO75,35^A0,25,25^FD%s^FS^XZ", // text
    label_val, label_val)
  return printLabel(print_script, "Zebra")
}


// 25AA040AT-I/OT print barcode 25AA040AT-I/OT{
func createBradyLabel(label_val string) error{
  print_script := fmt.Sprintf("<?xml version=\"1.0\"?>"+
        "<bpl-document xmlns=\"http://www.bradycorp.com/printers/bpl\">"+
        "<defaults>"+
            "<printer "+ 
                "match-media=\"B33-72-423\""+
                "match-ribbon=\"B30-R6000\""+
                "heat=\"3\"" +
                "tear-or-cut=\"between-labels\"/>"+
            "<document units=\"inches\"/>"+
        "</defaults>"+
        "<labels>"+
            "<label>"+
                "<barcode "+
                    "position-x=\"0.025\""+
                    "position-y=\"0.025\""+
                    "height=\".1\""+
                    "type=\"datamatrix\""+
                    "density=\"10\""+
                    "ratio=\"2:1\""+
                    "check-character=\"false\">"+
                    "<datasource>"+
                        "<static-text value=\"%s\"/>"+
                    "</datasource>"+
                "</barcode>"+
                    "<text "+
                        "position-x=\".3\""+
                        "position-y=\".08\""+
                        "align=\"right\""+
                        "font-name=\"Arial\""+
                        "font-size=\"4\""+
                        "bold=\"true\""+
                        "width=\".5\""+
                        "height=\".05\">"+
                        "<datasource>"+
                            "<static-text value=\"%s\"/>"+
                        "</datasource>"+
                    "</text>"+
            "</label>"+
        "</labels>"+
        "</bpl-document>", label_val,  label_val)
        return printLabel(print_script, "Brady")
}


/* Creates a print script and sends it via cmd line to PrintFile program
    PrintFile must be configured to have printer settings set up under that alias
    or it will print to default*/
func printLabel(print_script, printer_alias string) error {
  script_file := "C:\\CopyThat\\label_script.txt"
  // check that directory exists if not make it
  if _, err := os.Stat("C:\\CopyThat\\"); os.IsNotExist(err) {
        err := os.Mkdir("C:\\CopyThat\\", 0755)
        checkError("Error", err, false)
    }
  if fileExists(script_file){
    err := os.Remove(script_file)
    if err != nil {
      return errors.New("Error while removing print script file." + err.Error())
    }
  }
  f, err := os.Create(script_file)
  if err != nil {
    return errors.New("Error while creating print script file." + err.Error())
  }
  defer f.Close()
  _, err = f.WriteString(print_script)
  if err != nil {
    return errors.New("Error while writing to print script file." + err.Error())
  }
  f.Close()
  err = exec.Command("prfile32", "/q", "/n:"+printer_alias, script_file).Run()
  if err !=nil {
    return errors.New("Error while executing print command." + err.Error())
  }
  return nil
}

//checks if a file exists and isn't a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (mmw *MyMainWindow) display(title string, message string) {
	walk.MsgBox(mmw, title, message, walk.MsgBoxIconInformation)
}

type MyMainWindow struct {
	*walk.MainWindow
}

//checkin errers
func checkError(message string, err error, exitProgram bool) {
	if err != nil {
		f, _ := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		logger := log.New(f, "log: ", log.LstdFlags)
		logger.Println(message)
		logger.Println(err.Error() + "\n")
        fmt.Printf("Error: %s - %s\n", message, err.Error())
        if exitProgram {
          os.Exit(1)
        }
    }
}
