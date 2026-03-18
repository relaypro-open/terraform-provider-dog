package dog

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(title string, incoming any) {
	d, _ := json.MarshalIndent(incoming, "", "  ")
	fmt.Println("=", title)
	fmt.Println(string(d))
	fmt.Println("=end", title)
}

func PrettyFmt(title string, incoming any) (str string) {
	d, _ := json.MarshalIndent(incoming, "", "  ")
	return fmt.Sprintf("=%s\n%s\n%s=", title, string(d), title)
}
