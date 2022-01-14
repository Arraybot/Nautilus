package requests

// PannelHealth is a wrapper for the panel health.
type PanelHealth struct {
	Connections int `json:"connected"`
}

// ListenerHealth is a wrapper for the listener health.
type ListenerHealth struct {
	Memory string `json:"memory"`
}

// UrbanResult is a wrapper for the urban dictionary result.
type UrbanResult struct {
	List []*UrbanDefinition `json:"list"`
}

// UrbanDefinition is a wrapper for urban dictionary definitions.
type UrbanDefinition struct {
	Description string `json:"definition"`
	Upvotes     int    `json:"thumbs_up"`
	Downvotes   int    `json:"thumbs_down"`
}

// Pet is an interface for pet links.
type Pet interface {
	link() string
	source() string
}

// Cat. Meow.
type Cat struct {
	File string `json:"file"`
}

// Dog. Woof.
type Dog struct {
	Url string `json:"url"`
}

// Implement Pet for Cat.
func (p *Cat) link() string {
	return p.File
}

// Implement Pet for Cat.
func (p *Cat) source() string {
	return "https://aws.random.cat/meow"
}

// Implement Pet for Dog.
func (p *Dog) link() string {
	return p.Url
}

// Implement Pet for Dog.
func (p *Dog) source() string {
	return "https://random.dog/woof.json"
}
