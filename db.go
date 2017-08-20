package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type CurrencyModel struct {
	ID        int
	Code      string
	Rate      float64
	UpdatedAt string
}

type UserModel struct {
	ID        int
	Username  string
	CreatedAt string
}

type SubscriptionModel struct {
	ID       int
	User     *UserModel
	Currency *CurrencyModel
}

var db *sql.DB

func dbConnect(name string, user string, password string, host string, port uint) error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, name))
	return err
}

func setUser(id int, username string) *UserModel {
	var err error

	user := getUser(id)
	if user != nil {
		_, err = db.Exec("UPDATE users SET username = $2 WHERE id = $1", id, username)
	} else {
		_, err = db.Exec("INSERT INTO users VALUES($1, $2)", id, username)
	}

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return getUser(id)
}

func getUser(id int) *UserModel {
	rows, err := db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	user := new(UserModel)
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		return user
	}
	return nil
}

func getSubscriptionsByUserID(id int) []SubscriptionModel {
	rows, err := db.Query(`
		SELECT 
			s.id,
			u.id,
			u.username,
			c.id,
			c.code,
			c.rate,
			c.updated_at
		FROM subscriptions s
		INNER JOIN users u ON u.id = s.user_id
		INNER JOIN currencies c ON c.id = s.currency_id
		WHERE u.id = $1`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var subscriptions []SubscriptionModel

	for rows.Next() {
		subscription := SubscriptionModel{
			User:     &UserModel{},
			Currency: &CurrencyModel{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.User.ID,
			&subscription.User.Username,
			&subscription.Currency.ID,
			&subscription.Currency.Code,
			&subscription.Currency.Rate,
			&subscription.Currency.UpdatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}

		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions
}

func getCurrency(code string) *CurrencyModel {
	rows, err := db.Query("SELECT * FROM currencies WHERE code = $1 LIMIT 1", code)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		currency := new(CurrencyModel)
		err := rows.Scan(&currency.ID, &currency.Code, &currency.Rate, &currency.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}

		return currency
	}
	return nil
}

func setCurrencyRate(code string, rate float64) *CurrencyModel {
	var err error
	currency := getCurrency(code)

	if currency != nil {
		_, err = db.Exec("UPDATE currencies SET rate = $1, updated_at = $2 WHERE code = $3", rate, time.Now(), code)
	} else {
		_, err = db.Exec("INSERT INTO currencies (code, rate, updated_at) VALUES($1, $2, $3)", code, rate, time.Now())
	}

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return getCurrency(code)
}

func subscribe(userID int, code string) bool {
	subscription := getSubscription(userID, code)

	if subscription == nil {
		currency := getCurrency(code)
		ret, err := db.Exec(`
			INSERT INTO subscriptions (user_id, currency_id) VALUES($1,$2)`,
			userID, currency.ID)

		if err != nil {
			log.Fatal(err)
		}
		lastInsertID, _ := ret.LastInsertId()
		return lastInsertID > 0
	}

	return true
}

func unsubscribe(userID int, code string) bool {
	ret, err := db.Exec(`
		DELETE FROM subscriptions s
		WHERE s.user_id = $1 
		AND s.currency_id IN 
			(SELECT id FROM currencies c WHERE c.code = $2)`, userID, code)

	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := ret.RowsAffected()
	return rowsAffected > 0
}

func getSubscription(userID int, code string) *SubscriptionModel {
	rows, err := db.Query(`
		SELECT 
			s.id,
			u.id,
			u.username,
			c.id,
			c.code,
			c.rate,
			c.updated_at
		FROM subscriptions s
		INNER JOIN users u ON u.id = s.user_id
		INNER JOIN currencies c ON c.id = s.currency_id
		WHERE u.id = $1 AND c.code = $2`, userID, code)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		subscription := SubscriptionModel{
			User:     &UserModel{},
			Currency: &CurrencyModel{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.User.ID,
			&subscription.User.Username,
			&subscription.Currency.ID,
			&subscription.Currency.Code,
			&subscription.Currency.Rate,
			&subscription.Currency.UpdatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}

		return &subscription
	}
	return nil
}
