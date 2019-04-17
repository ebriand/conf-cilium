package main

import (
	"fmt"
	"os"

	"github.com/dustinrc/marvel"
)

var publicKey, privateKey string

func init() {
	publicKey = os.Getenv("MARVEL_PUBLIC_KEY")
	if publicKey == "" {
		panic("env variable MARVEL_PUBLIC_KEY is missing")
	}
	privateKey = os.Getenv("MARVEL_PRIVATE_KEY")
	if privateKey == "" {
		panic("env variable MARVEL_PRIVATE_KEY is missing")
	}
}

func getHeroImage(name string) (string, error) {
	auth := marvel.NewServerSideAuth(publicKey, privateKey)
	client := marvel.NewClient(auth, nil)
	params := marvel.CharacterParams{
		Name: name,
	}
	characters, err := client.Characters.All(&params)
	if err != nil {
		return "", err
	}
	if len(characters) != 1 {
		return "", fmt.Errorf("No picture found")
	}
	t := characters[0].Thumbnail
	return t.Path + "/standard_xlarge." + t.Extension, nil
}
