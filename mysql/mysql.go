package mysql

import (
	"database/sql"
	"strings"

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

func FetchRandomQuote(username string, channel string) (*common.Quote, error) {
	stmt, err := db.Prepare(`
		SELECT username, message, created_at
			FROM messages AS r1 JOIN
				(SELECT CEIL(RAND() * (SELECT MAX(id) FROM messages WHERE username = ? AND channel_name = ?)) AS id) AS r2
			WHERE r1.id >= r2.id AND username = ? AND channel_name = ?
			ORDER BY r1.id ASC
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}
	quote := common.Quote{}
	row := stmt.QueryRow(username, channel, username, channel)
	err = row.Scan(&quote.Username, &quote.Message, &quote.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
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
