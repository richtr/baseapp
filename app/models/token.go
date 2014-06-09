package models

import (
	"fmt"
)

type Token struct {
	TokenId    int
	Email      string
	Type, Hash string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%s)", t.Hash)
}
