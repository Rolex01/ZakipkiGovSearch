package user

import (
	"database/sql"
	"time"
)

type User struct {
	Id       int
	Email    string
	Name     string
	Lifetime time.Time
}

func CreateFromRow(row *sql.Rows) (User, error) {

	var u User

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.Lifetime,
	)

	return u, err
}

func GetExpireIn3DaysUsers(db *sql.DB) ([]User, error) {

	var users []User

	rows, err := db.Query(`
		SELECT row, name, email, lifetime
		  FROM users
		 WHERE lifetime <= current_date + interval '3' day
		   AND lifetime > current_date + interval '2' day
		   AND permission >= 1
	`)

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {

		if u, err := CreateFromRow(rows); err != nil {
			return users, err
		} else {
			users = append(users, u)
		}
	}

	return users, err
}
