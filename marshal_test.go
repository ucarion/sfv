package sfv_test

// import (
// 	"fmt"

// 	"github.com/ucarion/sfv"
// )

// func ExampleMarshal_raw_item() {
// 	item := sfv.Item{
// 		Value: sfv.Token("text/html"),
// 		Params: sfv.Params{
// 			Map: map[string]interface{}{
// 				"charset": sfv.Token("UTF-8"),
// 			},
// 			Keys: []string{"charset"},
// 		},
// 	}

// 	fmt.Println(sfv.Marshal(item))
// 	// Output: text/html;charset=UTF-8 <nil>
// }

// func ExampleMarshal_raw_list() {
// 	list := []sfv.Member{
// 		sfv.Member{
// 			IsItem: true,
// 			Item: sfv.Item{
// 				Value: sfv.Token("fr-CH"),
// 			},
// 		},
// 		sfv.Member{
// 			IsItem: true,
// 			Item: sfv.Item{
// 				Value: sfv.Token("fr"),
// 				Params: sfv.Params{
// 					Map: map[string]interface{}{
// 						"q": 0.9,
// 					},
// 					Keys: []string{"q"},
// 				},
// 			},
// 		},
// 		sfv.Member{
// 			IsItem: true,
// 			Item: sfv.Item{
// 				Value: sfv.Token("*"),
// 				Params: sfv.Params{
// 					Map: map[string]interface{}{
// 						"q": 0.5,
// 					},
// 					Keys: []string{"q"},
// 				},
// 			},
// 		},
// 	}

// 	fmt.Println(sfv.Marshal(list))
// 	// Output: fr-CH, fr;q=0.9, *;q=0.5 <nil>
// }

// func ExampleMarshal_raw_dict() {
// 	dict := sfv.Dictionary{
// 		Map: map[string]sfv.Member{
// 			"public": sfv.Member{
// 				IsItem: true,
// 				Item:   sfv.Item{Value: true},
// 			},
// 			"max-age": sfv.Member{
// 				IsItem: true,
// 				Item:   sfv.Item{Value: int64(604800)},
// 			},
// 			"immutable": sfv.Member{
// 				IsItem: true,
// 				Item:   sfv.Item{Value: true},
// 			},
// 		},
// 		Keys: []string{"public", "max-age", "immutable"},
// 	}

// 	fmt.Println(sfv.Marshal(dict))
// 	// Output: public, max-age=604800, immutable <nil>
// }
