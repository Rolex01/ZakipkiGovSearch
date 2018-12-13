package user

import (
	"database/sql"
	"github.com/jasonlvhit/gocron"
	"bitbucket.org/company-one/tender-one/mail"
)

func InitScheduler(db *sql.DB) {
	gocron.Every(1).Day().At("05:00").Do(CheckSubscriptionExpiration, db)
	gocron.Start()
}

func CheckSubscriptionExpiration(db *sql.DB) error {

	users, err := GetExpireIn3DaysUsers(db)

	if err != nil {
		return err
	}

	for _, user := range users {
		
		err = mail.SendSubscriptionExpirationNotification("y.babenko@gkcrm.ru", user.Email, user.Name)

		if err != nil {
			return err
		}
	}

	return err
}
