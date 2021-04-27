package main

type Viewer struct {
	Username      string
	Messages      []string
	MessagesCount int
}

type ByMessagesCount []Viewer

func (v ByMessagesCount) Len() int           { return len(v) }
func (v ByMessagesCount) Less(i, j int) bool { return v[i].MessagesCount > v[j].MessagesCount }
func (v ByMessagesCount) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

type Streamer struct {
	Username string
	Viewers  map[string]Viewer
}
