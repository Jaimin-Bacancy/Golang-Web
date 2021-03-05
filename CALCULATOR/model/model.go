package model

type Calculator struct {
	Input1    string `json:"input1"`
	Input2    string `json:"input2"`
	Operation string `json:"operator"`
	Result    string `json:"result"`
}
