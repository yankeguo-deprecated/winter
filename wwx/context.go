package wwx

import (
	"context"
)

const (
	TypeText     = "text"
	TypeImage    = "image"
	TypeVoice    = "voice"
	TypeVideo    = "video"
	TypeLocation = "location"
	TypeLink     = "link"

	TypeEvent = "event"

	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "SCAN"
	EventLocation    = "LOCATION"
	EventClick       = "CLICK"
	EventView        = "VIEW"
)

type Request struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`

	MsgID   int64  `xml:"MsgId"`
	MsgType string `xml:"MsgType"`

	MsgDataID string `xml:"MsgDataId"`
	Idx       int64  `xml:"Idx"`

	// text
	Content string `xml:"Content"`

	// image
	PicURL  string `xml:"PicUrl"`
	MediaID string `xml:"MediaId"`

	// voice
	//MediaID string `xml:"MediaId"`
	Format      string `xml:"Format"`
	Recognition string `xml:"Recognition"`

	// video
	//MediaID string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`

	// location
	LocationX float64 `xml:"Location_X"`
	LocationY float64 `xml:"Location_Y"`
	Scale     int     `xml:"Scale"`
	Label     string  `xml:"Label"`

	// link
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	URL         string `xml:"Url"`

	// Event
	Event     string  `xml:"Event"`
	EventKey  string  `xml:"EventKey"`
	Ticket    string  `xml:"Ticket"`
	Latitude  float64 `xml:"Latitude"`
	Longitude float64 `xml:"Longitude"`
	Precision float64 `xml:"Precision"`
}

type TextResponse struct {
	ToUserName   string `xml:"ToUserName,cdata"`
	FromUserName string `xml:"FromUserName,cdata"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType,cdata"`
	Content      string `xml:"Content,cdata,omitempty"`
}

// Context WeChat handling context
type Context interface {
	context.Context

	// Req returns the decoded request
	Req() Request

	// AppID returns the appid
	AppID() string

	// Empty respond with empty
	Empty() string

	// Text respond with text
	Text(s string) string
}
