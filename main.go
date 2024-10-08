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
	"encoding/json"
	"log/slog"
	"net/http"

	addons "github.com/technonotes/workspaceaddons"
	"golang.org/x/oauth2"
)

const spreadsheetId = "1cxwE4vtDeLtnVBbiN8pRz_VqQ5cxTj9px8xvmVkP0HA"

/*
Reading the email
Finding the invoice id and amount
Create a form in case user wants to edit the data
*/
func Invoice(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		slog.Error("Got GET request, expects POST", "error", "Invalid HTTP method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "POST":
		var event addons.EventObject
		var unmarshalErr *json.UnmarshalTypeError
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&event)
		if err != nil {
			slog.Error(
				"JSON to struct error",
				"error",
				err.Error(),
				"unmarshal error",
				unmarshalErr.Error(),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Debug(
			"Received request",
			"event",
			event.Gmail,
		)
		userToken := *event.AuthorizationEventObject.UserOAuthToken
		orderID, maxAmount, date, err := getDataFromMail(event, userToken)
		if err != nil {
			returnError(w)
			return
		}
		renderAction := getInvoiceCard(orderID, maxAmount, date, spreadsheetId)

		jsonResponse, err := json.Marshal(renderAction)
		if err != nil {
			slog.Debug("Marshal error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Debug("JSON data", "json", string(jsonResponse))
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

// Show when not in a mail
func InvoiceMain(w http.ResponseWriter, _ *http.Request) {
	var unmarshalErr *json.UnmarshalTypeError
	renderAction := placeholderCard()
	jsonResponse, err := json.Marshal(renderAction)
	if err != nil {
		slog.Error("Marshal error", "error", unmarshalErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Receive a submitted form, store data in sheet
func Submit(w http.ResponseWriter, r *http.Request) {
	var e addons.EventObject

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		slog.Error("Error decoding Event object", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userToken := *e.AuthorizationEventObject.UserOAuthToken

	var unmarshalErr *json.UnmarshalTypeError

	ctx := context.Background()

	token := new(oauth2.Token)
	token.AccessToken = userToken

	var config oauth2.Config
	client := config.Client(ctx, token)
	notificationStr := "Google sheet updated with invoice data"
	err = updateSheet(
		ctx,
		client,
		e.CommonEventObject.FormInputs["date"].StringInputs.Value[0],
		e.CommonEventObject.FormInputs["description"].StringInputs.Value[0],
		e.CommonEventObject.FormInputs["amount"].StringInputs.Value[0],
		e.CommonEventObject.FormInputs["sheetID"].StringInputs.Value[0],
	)
	if err != nil {
		slog.Error("updateSheet error", "error", err)
		notificationStr = "Error!!! Data not updated"
	}

	renderActionWrapper := submittedCard(notificationStr)
	jsonData, err := json.Marshal(renderActionWrapper)
	if err != nil {
		slog.Error("Marshal error", "error", unmarshalErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func returnError(w http.ResponseWriter) {
	var unmarshalErr *json.UnmarshalTypeError
	renderAction := getErrorCard("Error in processing mail, is this an Invoice mail?")
	jsonResponse, err := json.Marshal(renderAction)
	if err != nil {
		slog.Error("Marshal error", "error", unmarshalErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	// Set up structured logging
	slog.SetDefault(slog.New(NewCloudLoggingHandler()))

	http.HandleFunc("/message", Invoice)
	http.HandleFunc("/main", InvoiceMain)
	http.HandleFunc("/submit", Submit)
	http.ListenAndServe(":8080", nil)
}
