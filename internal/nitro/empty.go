package nitro

import "fmt"

func Empty(name string) []Command {
	fmt.Println("this is empty!", name)
	return nil
}
