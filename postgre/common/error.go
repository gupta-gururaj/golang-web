package common

import "fmt"

//HandleError is ...
func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
