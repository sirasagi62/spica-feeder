package main

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	r := executeSearch("react")
	fmt.Println(r)
	c := convertResult(r)
	fmt.Println(c)
}
