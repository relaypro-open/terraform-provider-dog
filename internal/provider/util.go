package dog

import (
	"fmt"
	"encoding/json"
)


func PrettyPrint(title string, incoming interface{}) {
	d, _ := json.MarshalIndent(incoming, "", "  ")
	fmt.Println("=", title)
	fmt.Println(string(d))
	fmt.Println("=end", title)
}

func PrettyFmt(title string, incoming interface{}) (str string) {
	d, _ := json.MarshalIndent(incoming, "", "  ")
	return fmt.Sprintf("=%s\n%s\n%s=", title, string(d), title)
}
