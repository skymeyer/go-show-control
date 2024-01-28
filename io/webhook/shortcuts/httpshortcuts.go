package shortcuts

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.skymyer.dev/show-control/common"
)

const (
	FILE_JSON = "shortcuts.json"
	FILE_ZIP  = "shortcuts.zip"
)

func ImportHandler(w http.ResponseWriter, r *http.Request) {

	export := Export{
		Version:       74,
		CompatVersion: 71,
		Variables: []Variable{
			{
				ID:    "c9481de9-79a8-48bd-8e82-605f36c60379",
				Key:   "host",
				Value: "192.168.86.63:8765",
			},
		},
		Categories: []Category{
			{
				ID:     "a4b3b259-92cb-428f-868f-ab7d4fc6c08e",
				Name:   "Control",
				Layout: "medium_grid",
				Shortcuts: []Shortcut{
					{
						ID:     "57e639d0-4419-4e7a-a0c7-87dc28369847",
						Name:   "PAGE1",
						Method: "POST",
						URL:    "http://{{c9481de9-79a8-48bd-8e82-605f36c60379}}/io/control/CTR_PAGE_1",
						Icon:   "flat_color_lightbulb",
						Response: Response{
							SuccessOutput: "none",
							UIType:        "toast",
						},
					},
					{
						ID:     "05b1c141-86a8-4a7c-8e52-7b275f7073fa",
						Name:   "BACK",
						Method: "POST",
						URL:    "http://{{c9481de9-79a8-48bd-8e82-605f36c60379}}/io/control/CTR_BACK",
						Icon:   "flat_color_lightbulb",
						Response: Response{
							SuccessOutput: "none",
							UIType:        "toast",
						},
					},
				},
			},
			{
				ID:     "5afee60c-86d8-4127-9ff5-e5179d223545",
				Name:   "Live Grid",
				Layout: "medium_grid",
				Shortcuts: []Shortcut{
					{
						ID:     "b9db6d98-3b18-4cfa-9ff2-f86e9de1f4ba",
						Name:   "1",
						Method: "POST",
						URL:    "http://{{c9481de9-79a8-48bd-8e82-605f36c60379}}/io/button/BTN_1_1",
						Icon:   "flat_color_lightbulb",
						Response: Response{
							SuccessOutput: "none",
							UIType:        "toast",
						},
					},
					{
						ID:     "754a0098-fec7-4b8a-8c85-8fde772b8b1b",
						Name:   "2",
						Method: "POST",
						URL:    "http://{{c9481de9-79a8-48bd-8e82-605f36c60379}}/io/button/BTN_1_2",
						Icon:   "flat_color_lightbulb",
						Response: Response{
							SuccessOutput: "none",
							UIType:        "toast",
						},
					},
					{
						ID:     "4e21a741-30e7-4687-a209-71e51b9a6d54",
						Name:   "3",
						Method: "POST",
						URL:    "http://{{c9481de9-79a8-48bd-8e82-605f36c60379}}/io/button/BTN_1_3",
						Icon:   "flat_color_lightbulb",
						Response: Response{
							SuccessOutput: "none",
							UIType:        "toast",
						},
					},
					{
						ID:     "c796dded-b23e-4f51-9086-cf248b4154d2",
						Name:   "4",
						Method: "POST",
						URL:    "http://{{c9481de9-79a8-48bd-8e82-605f36c60379}}/io/button/BTN_1_4",
						Icon:   "flat_color_lightbulb",
						Response: Response{
							SuccessOutput: "none",
							UIType:        "toast",
						},
					},
				},
			},
		},
	}

	var bufJson bytes.Buffer
	json.NewEncoder(&bufJson).Encode(export)

	w.Header().Set("Content-Type", "application/json")
	w.Write(bufJson.Bytes())
}

type Export struct {
	Categories    []Category `json:"categories"`
	Variables     []Variable `json:"variables"`
	Version       int        `json:"version"`
	CompatVersion int        `json:"compatibilityVersion"`
}

func NewVariable(key, value string) Variable {
	return Variable{
		ID:    common.StableUUIDString(key),
		Key:   key,
		Value: value,
	}
}

type Variable struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewCategory(name string) Category {
	return Category{
		ID:   common.StableUUIDString(name),
		Name: name,
	}
}

type Category struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Layout    string     `json:"layoutType"`
	Shortcuts []Shortcut `json:"shortcuts"`
}

func NewShortCut(id, name string) Shortcut {
	return Shortcut{
		ID:   common.StableUUIDString(id),
		Name: name,
	}
}

type Shortcut struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Method   string   `json:"method"`
	URL      string   `json:"url"`
	Icon     string   `json:"iconName"`
	Response Response `json:"responseHandling"`
}

type Response struct {
	SuccessOutput string `json:"successOutput"`
	UIType        string `json:"uiType"`
}
