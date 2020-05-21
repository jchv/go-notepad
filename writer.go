// Package notepad implements support for logging to Notepad windows.
package notepad

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32            = windows.MustLoadDLL("user32")
	fnFindWindow      = user32.MustFindProc("FindWindowW")
	fnFindWindowEx    = user32.MustFindProc("FindWindowExW")
	fnSendMessage     = user32.MustFindProc("SendMessageW")
	notepadWinTitle   = windows.StringToUTF16Ptr("Untitled - Notepad")
	notepadWinTitle2  = windows.StringToUTF16Ptr("*Untitled - Notepad")
	notepadEditWindow = windows.StringToUTF16Ptr("EDIT")
)

const (
	msgEMREPLACESEL = 0x00C2
)

// Writer as an io.Writer that writes to Notepad.
type Writer struct {
	hwnd windows.Handle
}

// NewWriter creates a new Notepad writer.
func NewWriter() (*Writer, error) {
	notepadWin, _, err := fnFindWindow.Call(uintptr(0), uintptr(unsafe.Pointer(notepadWinTitle)))
	if notepadWin == 0 {
		notepadWin, _, err = fnFindWindow.Call(uintptr(0), uintptr(unsafe.Pointer(notepadWinTitle2)))
		if notepadWin == 0 {
			return nil, fmt.Errorf("untitled notepad window not found: %w", err)
		}
	}

	editWindow, _, err := fnFindWindowEx.Call(notepadWin, uintptr(0), uintptr(unsafe.Pointer(notepadEditWindow)), uintptr(0))
	if editWindow == 0 {
		return nil, fmt.Errorf("notepad editbox not found: %w", err)
	}

	return &Writer{hwnd: windows.Handle(editWindow)}, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	buf := windows.StringToUTF16Ptr(string(p))
	fnSendMessage.Call(uintptr(w.hwnd), msgEMREPLACESEL, 1, uintptr(unsafe.Pointer(buf)))
	return len(p), nil
}
