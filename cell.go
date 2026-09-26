package memless

const (
	kindText    int32 = 1
	kindInteger int32 = 2
	kindDecimal int32 = 3
	kindBoolean int32 = 4
)

func cell(result uint64, row uint64, column uint64) any {
	var integer int64
	var decimal float64
	var boolean int32
	var text *byte
	kind := memlessResultCell(result, row, column, &integer, &decimal, &boolean, &text)
	return cellValue(kind, integer, decimal, boolean, text)
}

func cellValue(kind int32, integer int64, decimal float64, boolean int32, text *byte) any {
	if kind == kindText {
		return cString(text)
	}
	if kind == kindInteger {
		return integer
	}
	if kind == kindDecimal {
		return decimal
	}
	if kind == kindBoolean {
		return boolean != 0
	}
	return nil
}
