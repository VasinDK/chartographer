package entity

import "mime/multipart"

type Chart struct {
	X        int                   `json:"x,omitempty"`
	Y        int                   `json:"y,omitempty"`
	Width    int                   `json:"width,omitempty"`
	Height   int                   `json:"height,omitempty"`
	Id       string                `json:"id,omitempty"`
	IdParent string                `json:"idParent,omitempty"`
	File     multipart.File        `json:"file,omitempty"`
	Handle   *multipart.FileHeader `json:"handle,omitempty"`
}
