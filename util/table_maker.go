package util

import (
	"fmt"

	"github.com/carthewd/c3/internal/data"

	"github.com/alexeyco/simpletable"
)

// PrintTable will print data strcutures using struct fields as headers
// and slice data as rows
func PrintTable(t data.TableData) {
	table := simpletable.New()

	var headerCells []*simpletable.Cell

	for _, header := range t.GetHeaders() {
		newCell := simpletable.Cell{
			Align: simpletable.AlignCenter,
			Text:  header,
		}
		headerCells = append(headerCells, &newCell)
	}

	table.Header = &simpletable.Header{
		Cells: headerCells,
	}

	for _, rows := range t.GetRows() {
		var rowCells []*simpletable.Cell
		for _, item := range rows {
			newCell := simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  item,
			}
			rowCells = append(rowCells, &newCell)
		}
		table.Body.Cells = append(table.Body.Cells, rowCells)
	}

	fmt.Println(table.String())
}
