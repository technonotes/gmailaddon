package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// const spreadsheetId = "1Uskz7OqC5q8ldMlzHZTLiX2L18t7ISGTC5w9AP9fyxk"

/*
Reading the email
Finding the invoice id and amount
Create a form in case user wants to edit the data
*/
// func Invoice(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		slog.Error("Got GET request, expects POST", "error", "Invalid HTTP method")
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	case "POST":
// 		var e addons.EventObject
// 		var unmarshalErr *json.UnmarshalTypeError
// 		decoder := json.NewDecoder(r.Body)
// 		decoder.DisallowUnknownFields()
// 		err := decoder.Decode(&e)
// 		if err != nil {
// 			slog.Error(
// 				"JSON to struct error",
// 				"error",
// 				err.Error(),
// 				"unmarshal error",
// 				unmarshalErr.Error(),
// 			)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		slog.Debug("Received request", "event", e.Gmail)
// 		userToken := *e.AuthorizationEventObject.UserOAuthToken
//
// 		message, err := getMail(e, userToken)
// 		if err != nil {
// 			slog.Error("Cloud not get mail", "error", err)
// 			returnError(w)
// 			return
// 		}
//
// 		content, err := getMailContent(message)
// 		if err != nil {
// 			slog.Error("Cloud not get content from mail", "error", err)
// 			returnError(w)
// 			return
// 		}
//
// 		maxAmount, err := getMaxAmount(content)
// 		if err != nil {
// 			returnError(w)
// 		}
//
// 		orderID, err := getOrderID(string(content))
// 		if err != nil {
// 			orderID = "Unknown"
// 		}
// 		date := getMailDate(message)
//
// 		renderAction := getInvoiceCard(orderID, maxAmount, date, spreadsheetId)
//
// 		jsonData, err := json.Marshal(renderAction)
// 		if err != nil {
// 			slog.Debug("Marshal error", "error", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		slog.Debug("JSON data", "json", string(jsonData))
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(jsonData)
//
// 	default:
// 		w.WriteHeader(http.StatusNotImplemented)
// 		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
// 	}
// }

// Show when not in a mail
// func InvoiceMain(w http.ResponseWriter, _ *http.Request) {
// 	var renderAction addons.RenderAction
// 	var unmarshalErr *json.UnmarshalTypeError
// 	action := renderAction.CreateAction()
// 	navigation := action.AddNavigation()
// 	card := navigation.AddCard()
// 	card.AddHeader("CIRCLE")
// 	infoSection := card.AddSection("Invoice kvitteringer")
// 	infoSection.AddWidget().
// 		AddImage("Placeholder image", "https://storage.googleapis.com/innovatorshivepictures/gmailaddonlogo.png")
// 	infoSection.AddWidget().
// 		AddTextParagraph("Velg en Invoice kvitteringsmail for å trekke ut data som så kan lagres i Google Sheet")
//
// 	jsonData, err := json.Marshal(renderAction)
// 	if err != nil {
// 		slog.Error("Marshal error", "error", unmarshalErr.Error())
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonData)
// }

// Receive a submitted form, store data in sheet
// func Submit(w http.ResponseWriter, r *http.Request) {
// 	var e addons.EventObject
//
// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()
// 	err := decoder.Decode(&e)
// 	if err != nil {
// 		slog.Error("Error decoding Event object", "error", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	userToken := *e.AuthorizationEventObject.UserOAuthToken
//
// 	var unmarshalErr *json.UnmarshalTypeError
// 	var renderActionWrapper addons.RenderActionWrapper
// 	renderAction := renderActionWrapper.AddRenderAction()
//
// 	ctx := context.Background()
//
// 	token := new(oauth2.Token)
// 	token.AccessToken = userToken
//
// 	var config oauth2.Config
// 	client := config.Client(ctx, token)
// 	notificationStr := "Google sheet updated with invoice data"
// 	err = updateSheet(
// 		ctx,
// 		client,
// 		e.CommonEventObject.FormInputs["date"].StringInputs.Value[0],
// 		e.CommonEventObject.FormInputs["description"].StringInputs.Value[0],
// 		e.CommonEventObject.FormInputs["amount"].StringInputs.Value[0],
// 		e.CommonEventObject.FormInputs["sheetID"].StringInputs.Value[0],
// 	)
// 	if err != nil {
// 		slog.Error("updateSheet error", "error", err)
// 		notificationStr = "Error!!! Data not updated"
// 	}
//
// 	action := renderAction.CreateAction()
// 	action.AddNotification(notificationStr)
//
// 	jsonData, err := json.Marshal(renderActionWrapper)
// 	if err != nil {
// 		slog.Error("Marshal error", "error", unmarshalErr.Error())
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonData)
// }
//
// func getMailDate(message *gmail.Message) string {
// 	t := time.UnixMilli(message.InternalDate)
// 	return t.Format("02.01.2006")
// }

// func getMail(e addons.EventObject, userToken string) (*gmail.Message, error) {
// 	ctx := context.Background()
//
// 	token := new(oauth2.Token)
// 	token.AccessToken = userToken
//
// 	var config oauth2.Config
// 	client := config.Client(ctx, token)
//
// 	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to create service: %+v", err)
// 	}
//
// 	r := srv.Users.Messages.Get("me", *e.Gmail.MessageId)
//
// 	if r == nil {
// 		return nil, errors.New("unable to rettreive message")
// 	}
// 	r.Header().Add("X-Goog-Gmail-Access-Token", *e.Gmail.AccessToken)
// 	message, err := r.Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve message: %v", err)
// 	}
// 	return message, nil
// }
//
// func getMailContent(message *gmail.Message) (string, error) {
// 	b64Content := message.Payload.Parts[0].Body.Data
//
// 	content, err := b64.URLEncoding.DecodeString(b64Content)
// 	if err != nil {
// 		return "", fmt.Errorf("unable to b64 decode message: %+v", err)
// 	}
// 	return string(content), nil
// }
//
// func getOrderID(content string) (string, error) {
// 	idregexp := regexp.MustCompile(`Kvittering\ \-\ [a-z0-9]+`)
// 	idmatches := idregexp.FindAllString(content, -1)
// 	var id string
//
// 	if idmatches == nil {
// 		id = ""
// 		return "", errors.New("couldn't find id")
// 	} else {
// 		for _, match := range idmatches {
// 			id = strings.TrimSpace(strings.Split(match, "-")[1])
// 		}
// 	}
// 	return id, nil
// }

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

// func getMaxAmount(content string) (string, error) {
// 	re := regexp.MustCompile(`[\d]+\,\d\d`)
// 	matches := re.FindAllString(content, -1)
// 	if matches == nil {
// 		return "", errors.New("no match")
// 	}
// 	maxAmount := 0
// 	for _, match := range matches {
// 		sAmount := strings.Split(match, " ")[0]
// 		iAmount, err := strconv.ParseInt(strings.Replace(sAmount, ",", "", -1), 0, 0)
// 		if err != nil {
// 			return "", errors.New(fmt.Sprintf("not able to convert %s to int", sAmount))
// 		}
// 		if iAmount > int64(maxAmount) {
// 			maxAmount = int(iAmount)
// 		}
// 	}
//
// 	s := fmt.Sprint(maxAmount)
// 	if maxAmount < 10 {
// 		s = "00" + s
// 	} else if maxAmount < 100 {
// 		s = "0" + s
// 	}
// 	sp := len(s) - 2
// 	maxAmountStr := s[:sp] + "," + s[sp:]
// 	return maxAmountStr, nil
// }
//
// Create card to show form
// func getInvoiceCard(
// 	id string,
// 	maxAmountStr string,
// 	date string,
// 	sheetID string,
// ) *addons.RenderAction {
// 	var renderAction addons.RenderAction
// 	action := renderAction.CreateAction()
// 	navigation := action.AddNavigation()
// 	card := navigation.AddCard()
//
// 	card.AddHeader("CIRCLE")
// 	sectionDataFromMail := card.AddSection("Information")
// 	sectionDataFromMail.AddWidget().AddTextInput("date", "Date", date)
// 	sectionDataFromMail.AddWidget().AddTextInput("description", "Description", "Invoice - "+id)
// 	sectionDataFromMail.AddWidget().AddTextInput("amount", "Amount", maxAmountStr)
// 	sectionDataFromMail.AddWidget().AddTextInput("sheetID", "Sheet ID", sheetID)
//
// 	sectionButtons := card.AddSection("")
// 	sectionButtons.AddWidget().
// 		AddButtonList().
// 		AddSubmitButton("Submit", "https://gmailaddon-807377867216.europe-north1.run.app/submit")
// 	return &renderAction
// }
//
// func returnError(w http.ResponseWriter) {
// 	var unmarshalErr *json.UnmarshalTypeError
// 	var renderAction addons.RenderAction
//
// 	action := renderAction.CreateAction()
// 	navigation := action.AddNavigation()
// 	card := navigation.AddCard()
//
// 	card.AddHeader("CIRCLE")
// 	section := card.AddSection("Error")
// 	section.AddWidget().
// 		AddTextParagraph("Error in processing mail, is this an Invoice receipt mail?")
//
// 	jsonData, err := json.Marshal(renderAction)
// 	if err != nil {
// 		slog.Error("Marshal error", "error", unmarshalErr.Error())
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonData)
// }
//
// func main() {
// 	// Set up structured logging
// 	slog.SetDefault(slog.New(NewCloudLoggingHandler()))
//
// 	http.HandleFunc("/message", Invoice)
// 	http.HandleFunc("/main", InvoiceMain)
// 	http.HandleFunc("/submit", Submit)
// 	http.ListenAndServe(":8080", nil)
// }
