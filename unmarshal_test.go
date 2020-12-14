package sfv_test

import (
	"fmt"

	"github.com/ucarion/sfv"
)

func ExampleUnmarshal_raw_item() {
	item := sfv.Item{
		Value: sfv.Token("text/html"),
		Params: sfv.Params{
			Map: map[string]interface{}{
				"charset": sfv.Token("UTF-8"),
			},
			Keys: []string{"charset"},
		},
	}

	fmt.Println(sfv.Marshal(item))
	// Output: text/html;charset=UTF-8 <nil>
}
