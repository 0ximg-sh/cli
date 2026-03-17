package models

type RenderRequest struct {
	Code             string `json:"code"`
	Language         string `json:"language,omitempty"`
	Title            string `json:"title,omitempty"`
	Theme            string `json:"theme,omitempty"`
	Background       string `json:"background,omitempty"`
	BackgroundImage  string `json:"backgroundImage,omitempty"`
	CodePadRight     *int   `json:"codePadRight,omitempty"`
	Font             string `json:"font,omitempty"`
	HighlightLines   string `json:"highlightLines,omitempty"`
	LineOffset       *int   `json:"lineOffset,omitempty"`
	LinePad          *int   `json:"linePad,omitempty"`
	NoLineNumber     bool   `json:"noLineNumber,omitempty"`
	NoRoundCorner    bool   `json:"noRoundCorner,omitempty"`
	NoWindowControls bool   `json:"noWindowControls,omitempty"`
	PadHoriz         *int   `json:"padHoriz,omitempty"`
	PadVert          *int   `json:"padVert,omitempty"`
	ShadowBlurRadius *int   `json:"shadowBlurRadius,omitempty"`
	ShadowColor      string `json:"shadowColor,omitempty"`
	ShadowOffsetX    *int   `json:"shadowOffsetX,omitempty"`
	ShadowOffsetY    *int   `json:"shadowOffsetY,omitempty"`
	TabWidth         *int   `json:"tabWidth,omitempty"`
	WindowTitle      string `json:"windowTitle,omitempty"`
}
