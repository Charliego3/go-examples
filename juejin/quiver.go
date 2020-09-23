package main

type Meta struct {
	Children []Meta `json:"children,omitempty"`
	UUID string `json:"uuid"`
}

type QSection struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}