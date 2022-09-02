package models

type Image struct {
	UserID  int
	Bucket  string
	Path    string
	Payload []byte
}
