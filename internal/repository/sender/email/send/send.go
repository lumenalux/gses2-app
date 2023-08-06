package send

import (
	"errors"
	"io"
)

var (
	errNoRecipients = errors.New("no recipients")
)

type SenderSMTPClient interface {
	Mail(string) error
	Rcpt(string) error
	Data() (io.WriteCloser, error)
	Quit() error
}

func setMail(client SenderSMTPClient, from string) error {
	return client.Mail(from)
}

func setRecipients(client SenderSMTPClient, to []string) error {
	if len(to) == 0 {
		return errNoRecipients
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	return nil
}

func writeAndClose(client SenderSMTPClient, message []byte) error {
	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(message)
	if err != nil {
		return err
	}

	return writer.Close()
}

func SendEmail(client SenderSMTPClient, email *EmailMessage) error {
	err := setMail(client, email.From)
	if err != nil {
		return err
	}

	err = setRecipients(client, email.To)
	if errors.Is(err, errNoRecipients) {
		return nil
	}

	if err != nil {
		return err
	}

	emailMessage, err := email.Prepare()
	if err != nil {
		return err
	}

	return writeAndClose(client, emailMessage)
}
