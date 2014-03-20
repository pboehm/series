package main

import (
    "fmt"
    "github.com/pboehm/series/index"
    "github.com/pboehm/series/renamer"
)

func main() {
	fmt.Printf("Series ....\n")
	fmt.Printf("%s\n", index.GetIndex())
	fmt.Printf("%v\n", renamer.GetDirtyFiles())
}
