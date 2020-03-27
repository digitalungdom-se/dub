package internal

import (
	"context"
	"os"

	"github.com/cbroglie/mustache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer          *gomail.Dialer
	emailCollection *mongo.Collection
}

func NewMailer(collection *mongo.Collection) Mailer {
	var mailer Mailer

	mailer.dialer = gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("NOREPLY_EMAIL"), os.Getenv("NOREPLY_PASSWORD"))
	mailer.emailCollection = collection

	return mailer
}

func (mailer *Mailer) newEmail(toEmail string, emailTemplateID string, subject string, mustacheData map[string]string) error {
	filter := bson.M{"type": emailTemplateID}
	var emailTemplate bson.M

	err := mailer.emailCollection.FindOne(
		context.Background(),
		filter).Decode(&emailTemplate)
	if err != nil {
		return err
	}

	emailTemplateSTR := emailTemplate["email"].(string)

	data, err := mustache.Render(emailTemplateSTR, mustacheData)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(os.Getenv("NOREPLY_EMAIL"), os.Getenv("NOREPLY_NAME")))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", data)

	err = mailer.dialer.DialAndSend(m)

	return err
}

func (mailer *Mailer) SendVerifyDiscord(toEmail string, token string) error {
	data := map[string]string{"token": token}
	return mailer.newEmail(toEmail, "discordVerification", "Koppla ditt Discord konto till Digital Ungdom", data)
}

func (mailer *Mailer) SendVerifyDiscordConfirmation(toEmail string, name string) error {
	data := map[string]string{"name": name}
	return mailer.newEmail(toEmail, "discordVerificationConfirmation", "Grattis ditt Discord är nu kopplat!", data)
}
