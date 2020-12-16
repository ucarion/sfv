package sfv_test

import (
	"fmt"

	"github.com/ucarion/sfv"
)

func ExampleMarshal_raw_item() {
	item := sfv.Item{
		BareItem: sfv.BareItem{
			Type:  sfv.BareItemTypeToken,
			Token: "text/html",
		},
		Params: sfv.Params{
			Keys: []string{"charset"},
			Map: map[string]sfv.BareItem{
				"charset": sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "UTF-8",
				},
			},
		},
	}

	fmt.Println(sfv.Marshal(item))
	// Output: text/html;charset=UTF-8 <nil>
}

func ExampleMarshal_raw_list() {
	list := sfv.List{
		sfv.Member{
			IsItem: true,
			Item: sfv.Item{
				BareItem: sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "fr-CH",
				},
			},
		},
		sfv.Member{
			IsItem: true,
			Item: sfv.Item{
				BareItem: sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "fr",
				},
				Params: sfv.Params{
					Keys: []string{"q"},
					Map: map[string]sfv.BareItem{
						"q": sfv.BareItem{
							Type:    sfv.BareItemTypeDecimal,
							Decimal: 0.9,
						},
					},
				},
			},
		},
		sfv.Member{
			IsItem: true,
			Item: sfv.Item{
				BareItem: sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "*",
				},
				Params: sfv.Params{
					Keys: []string{"q"},
					Map: map[string]sfv.BareItem{
						"q": sfv.BareItem{
							Type:    sfv.BareItemTypeDecimal,
							Decimal: 0.5,
						},
					},
				},
			},
		},
	}

	fmt.Println(sfv.Marshal(list))
	// Output: fr-CH, fr;q=0.9, *;q=0.5 <nil>
}

func ExampleMarshal_raw_dict() {
	dict := sfv.Dictionary{
		Keys: []string{"public", "max-age", "immutable"},
		Map: map[string]sfv.Member{
			"public": sfv.Member{
				IsItem: true,
				Item: sfv.Item{
					BareItem: sfv.BareItem{
						Type:    sfv.BareItemTypeBoolean,
						Boolean: true,
					},
				},
			},
			"max-age": sfv.Member{
				IsItem: true,
				Item: sfv.Item{
					BareItem: sfv.BareItem{
						Type:    sfv.BareItemTypeInteger,
						Integer: 604800,
					},
				},
			},
			"immutable": sfv.Member{
				IsItem: true,
				Item: sfv.Item{
					BareItem: sfv.BareItem{
						Type:    sfv.BareItemTypeBoolean,
						Boolean: true,
					},
				},
			},
		},
	}

	fmt.Println(sfv.Marshal(dict))
	// Output: public, max-age=604800, immutable <nil>
}
