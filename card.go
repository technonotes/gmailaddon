package main

import (
	addons "github.com/technonotes/workspaceaddons"
)

// Show when not in a mail
func placeholderCard() addons.RenderAction {
	var renderAction addons.RenderAction
	action := renderAction.CreateAction()
	navigation := action.AddNavigation()
	card := navigation.AddCard()
	card.AddHeader("CIRCLE")
	infoSection := card.AddSection("Invoice kvitteringer")
	infoSection.AddWidget().
		AddImage("GMail add-on logo", "https://storage.googleapis.com/innovatorshivepictures/gmailaddonlogo.png")
	infoSection.AddWidget().
		AddTextParagraph("Select an Invoice mail to extract data to a Google Sheet")
	return renderAction
}

func submittedCard(notificationStr string) addons.RenderActionWrapper {
	var renderActionWrapper addons.RenderActionWrapper
	renderAction := renderActionWrapper.AddRenderAction()
	action := renderAction.CreateAction()
	action.AddNotification(notificationStr)
	return renderActionWrapper
}

// Create card to show form
func getInvoiceCard(
	id string,
	maxAmountStr string,
	date string,
	sheetID string,
) *addons.RenderAction {
	var renderAction addons.RenderAction
	action := renderAction.CreateAction()
	navigation := action.AddNavigation()
	card := navigation.AddCard()

	card.AddHeader("CIRCLE")
	sectionDataFromMail := card.AddSection("Information")
	sectionDataFromMail.AddWidget().AddTextInput("date", "Date", date)
	sectionDataFromMail.AddWidget().AddTextInput("description", "Description", "Invoice - "+id)
	sectionDataFromMail.AddWidget().AddTextInput("amount", "Amount", maxAmountStr)
	sectionDataFromMail.AddWidget().AddTextInput("sheetID", "Sheet ID", sheetID)

	sectionButtons := card.AddSection("")
	sectionButtons.AddWidget().
		AddButtonList().
		AddSubmitButton("Submit", "https://gmailaddon-807377867216.europe-north1.run.app/submit")
	return &renderAction
}

func getErrorCard(errorMessage string) *addons.RenderAction {
	var renderAction addons.RenderAction

	action := renderAction.CreateAction()
	navigation := action.AddNavigation()
	card := navigation.AddCard()

	card.AddHeader("CIRCLE")
	section := card.AddSection("Error")
	section.AddWidget().
		AddImage("GMaill add-on logo", "https://storage.googleapis.com/innovatorshivepictures/gmailaddonlogo.png")
	section.AddWidget().
		AddTextParagraph(errorMessage)
	return &renderAction
}
