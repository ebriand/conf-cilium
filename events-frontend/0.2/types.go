package main

import "github.com/google/uuid"

type Event struct {
	ID     uuid.UUID
	Name   string
	Heroes []string
}

type Hero struct {
	Name     string
	RealName string
}

type indexData struct {
	Events []Event
	Heroes []string
}

type detailData struct {
	Name   string
	Heroes []Hero
}
