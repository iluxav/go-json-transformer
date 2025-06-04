package jsontransform

// DSL represents the transformation rules provided by the udser.....
type DSL struct {
	Transformations []Transformation `json:"transformations"`
	Directive       string           `json:"directive,omitempty"`
	Fields          []string         `json:"fields,omitempty"`
	Sources         []string         `json:"sources,omitempty"`
}

// Transformation represents an individual transformation rule.
type Transformation struct {
	Source      interface{} `json:"source"`
	Destination string      `json:"destination"`
	IsStatic    bool        `json:"isStatic,omitempty"`
}
