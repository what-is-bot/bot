package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/what-is-bot/bot/internal"
)

type answerProvider struct {
	db *sql.DB
}

func (a *answerProvider) Ask(internal.Question) ([]internal.Answer, error) {
	return nil, nil
}

func NewAnswerProvider() internal.AnswerProvider {
	viper.SetDefault("postgres_port", 5432)
	viper.SetDefault("postgres_sslmode", "disable")

	user := viper.GetString("postgres_user")
	pass := viper.GetString("postgres_password")
	database := viper.GetString("postgres_database")
	host := viper.GetString("postgres_host")
	port := viper.GetInt("postgres_port")
	sslmode := viper.GetString("postgres_sslmode")

	connString := fmt.Sprintf("user=%s password= %s dbname=%s host=%s port=%d sslmode=%s",
		user, pass, database, host, port, sslmode)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	return &answerProvider{db: db}
}
