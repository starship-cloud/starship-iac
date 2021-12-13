package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenUserId() string{
	id, err := gonanoid.Generate("abcdefgefghijklmnopqrstuvwxyz", 10)
	if err != nil {
		panic(err)
	}
	return "user-" + id
}

func GenProjectId() string{
	id, err := gonanoid.Generate("abcdefgefghijklmnopqrstuvwxyz", 8)
	if err != nil {
		panic(err)
	}
	return "project-" + id
}
