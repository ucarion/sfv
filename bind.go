package sfv

import (
	"fmt"
)

func bindItemToSink(item Item, v interface{}) error {
	switch v := v.(type) {
	case *bool:
		if item.BareItem.Type == BareItemTypeBoolean {
			*v = item.BareItem.Boolean
			return nil
		}
	case *int:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = int(item.BareItem.Integer)
			return nil
		}
	case *int8:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = int8(item.BareItem.Integer)
			return nil
		}
	case *int16:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = int16(item.BareItem.Integer)
			return nil
		}
	case *int32:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = int32(item.BareItem.Integer)
			return nil
		}
	case *int64:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = int64(item.BareItem.Integer)
			return nil
		}
	case *uint:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = uint(item.BareItem.Integer)
			return nil
		}
	case *uint8:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = uint8(item.BareItem.Integer)
			return nil
		}
	case *uint16:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = uint16(item.BareItem.Integer)
			return nil
		}
	case *uint32:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = uint32(item.BareItem.Integer)
			return nil
		}
	case *uint64:
		if item.BareItem.Type == BareItemTypeInteger {
			*v = uint64(item.BareItem.Integer)
			return nil
		}
	case *float32:
		if item.BareItem.Type == BareItemTypeDecimal {
			*v = float32(item.BareItem.Decimal)
			return nil
		}
	case *float64:
		if item.BareItem.Type == BareItemTypeDecimal {
			*v = item.BareItem.Decimal
			return nil
		}
	case *string:
		if item.BareItem.Type == BareItemTypeString {
			*v = item.BareItem.String
			return nil
		}

		if item.BareItem.Type == BareItemTypeToken {
			*v = item.BareItem.Token
			return nil
		}
	case *[]byte:
		if item.BareItem.Type == BareItemTypeBinary {
			*v = item.BareItem.Binary
			return nil
		}
	}

	return fmt.Errorf("cannot marshal to %T from %s", v, item.BareItem.Type)
}
