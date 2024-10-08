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
	b64 "encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	addons "github.com/technonotes/workspaceaddons"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getDataFromMail(
	event addons.EventObject,
	userToken string,
) (orderID string, maxAmount string, date string, err error) {
	message, err := getMail(event, userToken)
	if err != nil {
		slog.Error("Cloud not get mail", "error", err)
		return "", "", "", err
	}

	content, err := getMailContent(message)
	if err != nil {
		slog.Error("Cloud not get content from mail", "error", err)
		return "", "", "", err
	}

	maxAmount, err = getMaxAmount(content)
	if err != nil {
		return "", "", "", err
	}

	orderID, err = getOrderID(string(content))
	if err != nil {
		orderID = "Unknown"
	}
	date = getMailDate(message)
	return orderID, maxAmount, date, nil
}

func getMailDate(message *gmail.Message) string {
	t := time.UnixMilli(message.InternalDate)
	return t.Format("02.01.2006")
}

func getMail(e addons.EventObject, userToken string) (*gmail.Message, error) {
	ctx := context.Background()

	token := new(oauth2.Token)
	token.AccessToken = userToken

	var config oauth2.Config
	client := config.Client(ctx, token)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create service: %+v", err)
	}

	r := srv.Users.Messages.Get("me", *e.Gmail.MessageId)

	if r == nil {
		return nil, errors.New("unable to rettreive message")
	}
	r.Header().Add("X-Goog-Gmail-Access-Token", *e.Gmail.AccessToken)
	message, err := r.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve message: %v", err)
	}
	slog.Info("Message content", "info", message)
	return message, nil
}

func getMailContent(message *gmail.Message) (string, error) {
	b64Content := message.Payload.Parts[0].Body.Data

	content, err := b64.URLEncoding.DecodeString(b64Content)
	if err != nil {
		return "", fmt.Errorf("unable to b64 decode message: %+v", err)
	}
	return string(content), nil
}

func getOrderID(content string) (string, error) {
	idregexp := regexp.MustCompile(`Kvittering\ \-\ [a-z0-9]+`)
	idmatches := idregexp.FindAllString(content, -1)
	var id string

	if idmatches == nil {
		id = ""
		return "", errors.New("couldn't find id")
	} else {
		for _, match := range idmatches {
			id = strings.TrimSpace(strings.Split(match, "-")[1])
		}
	}
	return id, nil
}

func getMaxAmount(content string) (string, error) {
	re := regexp.MustCompile(`[\d]+\,\d\d`)
	matches := re.FindAllString(content, -1)
	if matches == nil {
		return "", errors.New("no match")
	}
	maxAmount := 0
	for _, match := range matches {
		sAmount := strings.Split(match, " ")[0]
		iAmount, err := strconv.ParseInt(strings.Replace(sAmount, ",", "", -1), 0, 0)
		if err != nil {
			return "", errors.New(fmt.Sprintf("not able to convert %s to int", sAmount))
		}
		if iAmount > int64(maxAmount) {
			maxAmount = int(iAmount)
		}
	}

	s := fmt.Sprint(maxAmount)
	if maxAmount < 10 {
		s = "00" + s
	} else if maxAmount < 100 {
		s = "0" + s
	}
	sp := len(s) - 2
	maxAmountStr := s[:sp] + "," + s[sp:]
	return maxAmountStr, nil
}
