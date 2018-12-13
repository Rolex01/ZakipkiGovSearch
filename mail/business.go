package mail

import (
	"fmt"
)

func SendSubscriptionExpirationNotification(to, email, name string) error {

	return SendMessage(
		to,
		fmt.Sprintf("У пользователя %s подписка истекает через 3 дня", name),
		fmt.Sprintf("У пользователя <a href=\"mailto:%s\">%s</a> подписка истекает через 3 дня", email, name),
	)
}
