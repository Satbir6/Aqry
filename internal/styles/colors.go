package styles

type Theme struct {
	Primary string
	Muted   string
	Success string
	Warning string
	Danger  string
	Border  string
	Surface string
	Text    string
}

func DefaultTheme() Theme {
	return Theme{
		Primary: "#7C8CFF",
		Muted:   "#7B8496",
		Success: "#63D297",
		Warning: "#E7B75F",
		Danger:  "#F07178",
		Border:  "#495166",
		Surface: "#252A36",
		Text:    "#E6E9EF",
	}
}
