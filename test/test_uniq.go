package main

import "polling_server/uniq"
import "fmt"

func main() {
	for i := 0; i < 100000; i++ {
		fmt.Println(uniq.Uniq())
	}
}
