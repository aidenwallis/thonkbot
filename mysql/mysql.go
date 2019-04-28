package mysql

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/aidenwallis/thonkbot/common"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect(dsn string) error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db = conn
	return nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func FetchLineCount(username, channel string) (int, error) {
	discriminator := username + "::_::" + channel
	stmt, err := db.Prepare("SELECT message_count FROM users WHERE discriminator = ? LIMIT 1")
	if err != nil {
		return 0, err
	}
	msgCount := 0
	row := stmt.QueryRow(discriminator)
	err = row.Scan(&msgCount)
	if err == sql.ErrNoRows {
		return msgCount, nil
	}
	return msgCount, err
}

func FetchRandomQuote(username string, channel string) (*common.Quote, error) {
	discriminator := username + "::_::" + channel
	stmt, err := db.Prepare("SELECT message_count FROM users WHERE discriminator = ? LIMIT 1")
	if err != nil {
		return nil, err
	}
	msgCount := 0
	row := stmt.QueryRow(discriminator)
	err = row.Scan(&msgCount)
	if err == sql.ErrNoRows || msgCount == 0 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	quote := common.Quote{}
	rand.Seed(time.Now().Unix())
	offset := rand.Int() % msgCount

	stmt, err = db.Prepare(`SELECT username, message, created_at FROM messages WHERE username = ? AND channel_name = ? LIMIT 1 OFFSET ?`)
	if err != nil {
		return nil, err
	}

	row = stmt.QueryRow(username, channel, offset)
	err = row.Scan(&quote.Username, &quote.Message, &quote.CreatedAt)
	return &quote, err
}

func LogMessage(username string, message string, channel string) error {
	stmt, err := db.Prepare("INSERT INTO messages (username, message, channel_name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, message, channel)
	return err
}

func ScanMessages(channel string, username string, query string) (int, error) {
	query = strings.ToLower(query)
	stmt, err := db.Prepare("SELECT message FROM messages WHERE channel_name = ? AND username = ? AND message LIKE ?")
	if err != nil {
		return 0, err
	}

	rows, err := stmt.Query(channel, username, "%"+query+"%")
	if err != nil {
		return 0, err
	}

	count := 0
	for rows.Next() {
		var msg string
		err = rows.Scan(&msg)
		if err != nil {
			continue
		}
		count += strings.Count(strings.ToLower(msg), query)
	}
	return count, nil
}

func GlobalScanMessages(channel string, query string) (int, error) {
	query = strings.ToLower(query)
	stmt, err := db.Prepare("SELECT message FROM messages WHERE channel_name = ? AND message LIKE ?")
	if err != nil {
		return 0, err
	}

	rows, err := stmt.Query(channel, "%"+query+"%")
	if err != nil {
		return 0, err
	}

	count := 0
	for rows.Next() {
		var msg string
		err = rows.Scan(&msg)
		if err != nil {
			continue
		}
		count += strings.Count(strings.ToLower(msg), query)
	}
	return count, nil
}
