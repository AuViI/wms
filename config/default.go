package config

// DefaultConfig to use if no file found
var DefaultConfig = Configuration{
	Title: "Weather Monitoring System",
	ExampleCities: []string{
		"Berlin",
		"Braunschweig",
		"Frankfurt",
		"Hamburg",
		"Holbæk",
		"Kühlungsborn",
		"New York",
		"Oslo",
		"Rostock",
		"Tokio",
	},
	ExampleModi: []string{
		"txt",
		"forecast",
		"list",
		"csv",
		"dtage",
		"view",
		"normlist",
	},
	Rendering: RenderConfig{
		Cities: []string{
			"Frankfurt",
			"Köln",
			"Kühlungsborn",
			"Rostock",
			"Warnemünde",
		},
		Interval: 12,
	},
	DTageLinks: []string{"1/aktuell", "3/meteo", "5/meteo"},

	read: "default",
}
