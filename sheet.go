/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Write the data to the given sheet
func updateSheet(
	ctx context.Context,
	client *http.Client,
	date string,
	description string,
	amount string,
	spreadsheetId string,
) error {
	ssrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return errors.New(fmt.Sprintf("unable to retieve sheets client: %v", err))
	}

	req := sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			InheritFromBefore: false,
			Range: &sheets.DimensionRange{
				Dimension:  "ROWS",
				StartIndex: 4,
				EndIndex:   4,
				SheetId:    0,
			},
		},
	}
	requests := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&req},
	}

	_, err = ssrv.Spreadsheets.BatchUpdate(spreadsheetId, requests).Context(ctx).Do()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to retrieve data from sheet: %v", err))
	}

	records := [][]interface{}{{date, description, amount}}
	valueInputOption := "USER_ENTERED"
	insertDataOption := "INSERT_ROWS"
	rb := &sheets.ValueRange{
		Values: records,
	}
	response2, err := ssrv.Spreadsheets.Values.Append(spreadsheetId, "Sheet1!A5:C5", rb).
		ValueInputOption(valueInputOption).
		InsertDataOption(insertDataOption).
		Context(ctx).
		Do()

	if err != nil || response2.HTTPStatusCode != 200 {
		return errors.New(fmt.Sprintf("error updating sheet: %+v", err))
	}
	return nil
}
