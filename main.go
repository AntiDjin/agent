package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 1 - бот принимает город пользователя
	// 2 - бот сохраняет город пользователя и его город в бд
	// 3 - бот узнает координаты этого города
	// 4 - бот по команде /weather узнает погоду в этом городе и отправляет пользователю
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI("Тут должен быть токен")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Println("Бот запущен")

	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "/weather" {
				row := db.QueryRow("SELECT city FROM users WHERE id = $1", update.Message.From.ID)
				log.Panicln(row.Err())
				// получить из бд город этого пользователя
				// получить координаты этого города
				// получить погоду по этим координатам
				// скинуть погоду пользователю
			} else {
				row := db.QueryRow("SELECT id FROM users WHERE id = $1", update.Message.From.ID)
				var id int
				err = row.Scan(&id)
				if err != nil {
					if err == sql.ErrNoRows {
						db.Exec("INSERT INTO users (id) VALUES ($1)", update.Message.From.ID)
					} else {
						log.Println("Произошла ошибка при поиске пользователся с id ", id, err)
						continue
					}
				}
				_, err := db.Exec("UPDATE users SET city = $1 WHERE id = $2", update.Message.Text, update.Message.From.ID)
				if err != nil {
					log.Println("Возникла ошибка про обновлении города", err)
				}
				msg := tgbotapi.NewMessage(update.Message.From.ID, "Ваш город сохранен!")
				bot.Send(msg)
			}
		}
	}
}
