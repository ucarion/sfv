package sfv

type Dictionary struct {
	Map  map[string]Member
	Keys []string
}

type List = []Member

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
	BareItem BareItem
	Params   Params
}

type Params struct {
	Map  map[string]BareItem
	Keys []string
}

type BareItem struct {
	Type    BareItemType
	Integer int64
	Decimal float64
	String  string
	Token   string
	Binary  []byte
	Boolean bool
}

func (b BareItem) isBoolTrue() bool {
	return b.Type == BareItemTypeBoolean && b.Boolean == true
}

type BareItemType int

func (t BareItemType) String() string {
	switch t {
	case BareItemTypeInteger:
		return "integer"
	case BareItemTypeDecimal:
		return "decimal"
	case BareItemTypeString:
		return "string"
	case BareItemTypeToken:
		return "token"
	case BareItemTypeBinary:
		return "binary"
	case BareItemTypeBoolean:
		return "boolean"
	default:
		return "invalid type"
	}
}

const (
	BareItemTypeInteger BareItemType = iota + 1
	BareItemTypeDecimal
	BareItemTypeString
	BareItemTypeToken
	BareItemTypeBinary
	BareItemTypeBoolean
)
