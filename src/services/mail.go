package services

import (
	"fmt"
	"socious-id/src/apps/utils"
	"socious-id/src/lib"
)

var EmailChannel = CategorizeChannel("email")

type EmailApproachType string

const (
	EmailApproachTemplate EmailApproachType = "TEMPLATE"
	EmailApproachDirect   EmailApproachType = "DIRECT"
)

type EmailConfig struct {
	Approach    EmailApproachType
	Destination string
	Title       string
	Template    string
	Args        map[string]string
}

func SendEmail(emailConfig EmailConfig) {
	Mq.sendJson(EmailChannel, emailConfig)
}

func EmailWorker(message interface{}) {
	emailConfig := new(EmailConfig)
	utils.Copy(message, emailConfig)

	var (
		destination = emailConfig.Destination
		title       = emailConfig.Title
		template    = lib.SendGridTemplates[emailConfig.Template]
		args        = emailConfig.Args
	)

	if emailConfig.Approach == EmailApproachTemplate {
		//Sending email with template
		err := lib.SendGridClient.SendWithTemplate(destination, title, template, args)
		if err != nil {
			fmt.Println("Coudn't Send Email, Error: ", err.Error())
		}
	}
}
