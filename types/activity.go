package types

import "time"

type Activity struct {
	Id        int64  `storm:"id,increment"`
	Binding   string `storm:"index"`
	CreatedAt time.Time
	Request   *Request
	Response  *Response
	Error     error
}

func NewActivity() *Activity {
	return &Activity{
		CreatedAt: time.Now(),
	}
}

//func (a *Activity) SetBindings  {
//
//}
