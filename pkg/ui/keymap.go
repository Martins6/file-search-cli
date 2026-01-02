package ui

type KeyMap struct {
	Up         string
	Down       string
	Enter      string
	Left       string
	OpenEditor string
	Quit       string
	Backspace  string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:         "up",
		Down:       "down",
		Enter:      "enter",
		Left:       "left",
		OpenEditor: "e",
		Quit:       "q",
		Backspace:  "backspace",
	}
}
