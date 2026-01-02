package ui

type KeyMap struct {
	Up         string
	Down       string
	Enter      string
	Left       string
	OpenEditor string
	PrintPath  string
	Quit       string
	Backspace  string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:         "up",
		Down:       "down",
		Enter:      "enter",
		Left:       "left",
		OpenEditor: "/e",
		PrintPath:  "/p",
		Quit:       "/q",
		Backspace:  "backspace",
	}
}
