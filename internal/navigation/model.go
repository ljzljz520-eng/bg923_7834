package navigation

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type Link struct {
	ID          string `json:"id"`
	Group       string `json:"group"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}

type ImportLink struct {
	ID          string `json:"id"`
	Group       string `json:"group"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}

type ImportResult struct {
	Imported int      `json:"imported"`
	Rejected int      `json:"rejected"`
	Errors   []string `json:"errors,omitempty"`
}

type ValidationResult struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Valid  bool   `json:"valid"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type HomeSection struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type HomePage struct {
	Title          string        `json:"title"`
	Primary        []HomeSection `json:"primary"`
	RightPanel     []HomeSection `json:"rightPanel"`
	ResourceGroups []string      `json:"resourceGroups"`
}
