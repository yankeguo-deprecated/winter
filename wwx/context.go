package wwx

import (
	"context"
	"encoding/xml"
	"github.com/guoyk93/winter"
	"time"
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
	XMLName xml.Name `xml:"xml"`

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

type Response struct {
	XMLName xml.Name `xml:"xml"`

	ToUserName   CDATA `xml:"ToUserName"`
	FromUserName CDATA `xml:"FromUserName"`
	CreateTime   int64 `xml:"CreateTime"`
	MsgType      CDATA `xml:"MsgType"`
	Content      CDATA `xml:"Content,omitempty"`

	Empty bool `xml:"-"`
}

// Context WeChat handling context
type Context interface {
	context.Context

	// Req returns the decoded request
	Req() Request

	// Res returns the response
	Res() *Response

	// Text send text response
	Text(s string)

	// Empty send empty response
	Empty()
}

type wxContext struct {
	c   winter.Context
	req Request
	res *Response
}

func (w wxContext) Deadline() (deadline time.Time, ok bool) {
	return w.c.Deadline()
}

func (w wxContext) Done() <-chan struct{} {
	return w.c.Done()
}

func (w wxContext) Err() error {
	return w.c.Err()
}

func (w wxContext) Value(key any) any {
	return w.c.Value(key)
}

func (w wxContext) Req() Request {
	return w.req
}

func (w wxContext) Res() *Response {
	return w.res
}

func (w wxContext) Text(s string) {
}

func (w wxContext) Empty() {
}
