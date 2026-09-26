package memless

import "strings"

// Query runs sql against the instance and returns its rows. It returns a
// RefusalError when the query is refused, and a FaultError on a boundary or
// internal fault. The result is read and released before returning.
func (i *Instance) Query(sql string) (Rows, error) {
	if strings.IndexByte(sql, 0) >= 0 {
		return Rows{}, &FaultError{Status: statusInvalidArgument, Message: "sql contains a NUL byte"}
	}
	if err := ensureLoaded(); err != nil {
		return Rows{}, err
	}
	var result uint64
	var message *byte
	status := memlessQuery(i.handle, sql, &result, &message)
	return interpretQuery(status, result, message)
}

func interpretQuery(status int32, result uint64, message *byte) (Rows, error) {
	text := takeMessage(message)
	if status == statusOk {
		return readResult(result), nil
	}
	if status == statusRefused {
		return Rows{}, &RefusalError{Message: text}
	}
	return Rows{}, &FaultError{Status: status, Message: text}
}

func readResult(result uint64) Rows {
	columns := resultColumns(result)
	rows := resultRows(result, len(columns))
	memlessResultRelease(result)
	return Rows{Columns: columns, Rows: rows}
}

func resultColumns(result uint64) []string {
	count := int(memlessResultColumnCount(result))
	return fillColumns(result, newColumns(count))
}

func newColumns(count int) []string {
	return make([]string, count)
}

func fillColumns(result uint64, columns []string) []string {
	for index := range columns {
		columns[index] = memlessResultColumn(result, uint64(index))
	}
	return columns
}

func resultRows(result uint64, width int) [][]any {
	count := int(memlessResultRowCount(result))
	if count == 0 || width == 0 {
		return newHeader(count)
	}
	return sliceRows(newCells(result, count, width), newHeader(count), width)
}

func newHeader(count int) [][]any {
	return make([][]any, count)
}

func newCells(result uint64, count int, width int) []any {
	return fillCells(result, make([]any, count*width), width)
}

func fillCells(result uint64, backing []any, width int) []any {
	for index := range backing {
		backing[index] = cell(result, uint64(index/width), uint64(index%width))
	}
	return backing
}

func sliceRows(backing []any, rows [][]any, width int) [][]any {
	for row := range rows {
		high := (row + 1) * width
		rows[row] = backing[row*width : high : high]
	}
	return rows
}
