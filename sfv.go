package sfv

type Dictionary struct {
	Map  map[string]Member
	Keys []string
}

type Member struct {
	IsItem    bool
	Item      Item
	InnerList InnerList
}

type InnerList struct {
	Items  []Item
	Params Params
}

type Item struct {
	Value  interface{}
	Params Params
}

type Params struct {
	Map  map[string]interface{}
	Keys []string
}
