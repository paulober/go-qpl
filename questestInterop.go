package qpl

import (
	"encoding/xml"
	"os"
)

// enum matTextType can be text/plain or text/xhtml
const (
	TextPlain = "text/plain"
	TextXhtml = "text/xhtml"
)

type MatText struct {
	TextType string `xml:"texttype,attr"`
	Text     string `xml:",chardata"`
}

type MatImage struct {
	Label string `xml:"label"`
	Uri   string `xml:"uri"`
}

type Material struct {
	MatText  MatText  `xml:"mattext"`
	MatImage MatImage `xml:"matimage"`
}

type FlowMat struct {
	Material Material `xml:"material"`
}

type QTIMetadataField struct {
	FieldLabel string `xml:"fieldlabel"`
	FieldEntry string `xml:"fieldentry"`
}

type QTIMetadata struct {
	QTIMetadataFields []QTIMetadataField `xml:"qtimetadatafield"`
}

type SolutionHint struct {
	Index  int `xml:"index"`
	Points int `xml:"points"`
}

type ItemFeedback struct {
	Ident   string  `xml:"ident"`
	View    string  `xml:"view"`
	FlowMat FlowMat `xml:"flow_mat"`
}

type Outcomes struct {
	Decvar string `xml:"decvar"`
}

type VarEqual struct {
	Respident string `xml:"respident,attr"`
	Value     int    `xml:",chardata"`
}

type Not struct {
	VarEqual VarEqual `xml:"varequal"`
}

type ConditionVar struct {
	VarEqual VarEqual `xml:"varequal"`
	Not      *Not     `xml:"not"`
}

type SetVar struct {
	Action string `xml:"action,attr"`
	Value  string `xml:",chardata"`
}

type DisplayFeedback struct {
	FeedbackType string `xml:"feedbacktype"`
	LinkRefID    string `xml:"linkrefid"`
}

type RespCondition struct {
	DoContinue      string          `xml:"continue,attr"`
	ConditionVar    ConditionVar    `xml:"conditionvar"`
	SetVar          SetVar          `xml:"setvar"`
	DisplayFeedback DisplayFeedback `xml:"displayfeedback"`
}

type ResponseLabel struct {
	Ident    int      `xml:"ident,attr"`
	Material Material `xml:"material"`
}

type RenderChoice struct {
	Shuffle        string          `xml:"shuffle,attr"`
	ResponseLabels []ResponseLabel `xml:"response_label"`
}

type ResponseLid struct {
	Ident        string       `xml:"ident,attr"`
	RCardinality string       `xml:"rcardinality,attr"`
	RenderChoice RenderChoice `xml:"render_choice"`
}

type ResProcessing struct {
	Outcomes       Outcomes        `xml:"outcomes"`
	RespConditions []RespCondition `xml:"respcondition"`
}

type Flow struct {
	Material    Material    `xml:"material"`
	ResponseLid ResponseLid `xml:"response_lid"`
}

type Presentation struct {
	Label string `xml:"label,attr"`
	Flow  Flow   `xml:"flow"`
}

type ItemMetadata struct {
	QTIMetadata QTIMetadata `xml:"qtimetadata"`
}

type Item struct {
	Ident         string         `xml:"ident,attr"`
	Title         string         `xml:"title,attr"`
	MaxAttempts   int            `xml:"maxattempts,attr"`
	QTIComment    string         `xml:"qticomment"`
	Duration      string         `xml:"duration"`
	ItemMetadata  ItemMetadata   `xml:"itemmetadata"`
	Presentation  Presentation   `xml:"presentation"`
	ResProcessing ResProcessing  `xml:"resprocessing"`
	ItemFeedbacks []ItemFeedback `xml:"itemfeedback"`
	SolutionHint  SolutionHint   `xml:"solutionhint"`
}

type QuestestInterop struct {
	Items []Item `xml:"item"`
}

func ReadQTIXMLFile(filePath string) (QuestestInterop, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return QuestestInterop{}, err
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)

	var qplContentObject QuestestInterop
	err = decoder.Decode(&qplContentObject)
	if err != nil {
		return QuestestInterop{}, err
	}

	return qplContentObject, nil
}
