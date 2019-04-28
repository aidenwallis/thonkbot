package mysql

import (
	"database/sql"

	"github.com/aidenwallis/thonkbot/common"
)

func FetchChannels() ([]common.Channel, error) {
	channels := []common.Channel{}
	rows, err := db.Query("SELECT id, name FROM channels")
	if err != nil {
		return channels, err
	}

	for rows.Next() {
		channel := common.Channel{}
		err = rows.Scan(&channel.ID, &channel.Name)
		if err != nil {
			return channels, err
		}
		channels = append(channels, channel)
	}

	return channels, err
}

func CheckJoined(channel string) (bool, error) {
	stmt, err := db.Prepare("SELECT id FROM channels WHERE name = ? LIMIT 1")
	if err != nil {
		return false, err
	}

	var id int
	row := stmt.QueryRow(channel)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func AddChannel(channelName string, joinedBy string) (common.Channel, error) {
	channel := common.Channel{}
	stmt, err := db.Prepare("INSERT INTO channels (name, added_by) VALUES (?, ?)")
	if err != nil {
		return channel, err
	}

	res, err := stmt.Exec(channelName, joinedBy)
	if err != nil {
		return channel, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return channel, err
	}

	stmt, err = db.Prepare("SELECT id, name FROM channels WHERE id = ? LIMIT 1")
	if err != nil {
		return channel, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(&channel.ID, &channel.Name)
	return channel, err
}

func DeleteChannel(id int, name string) error {
	stmt, err := db.Prepare("DELETE FROM channels WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("DELETE FROM messages WHERE channel_name = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(name)
	return err
}

func GetUserLogs(channel string, username string, limit int) ([]common.Quote, error) {
	quotes := []common.Quote{}
	stmt, err := db.Prepare("SELECT username, message, created_at FROM messages WHERE channel_name = ? AND username = ? ORDER BY created_at DESC LIMIT ?")
	if err != nil {
		return quotes, err
	}

	rows, err := stmt.Query(channel, username, limit)
	if err != nil {
		return quotes, err
	}

	for rows.Next() {
		quote := common.Quote{}
		err = rows.Scan(&quote.Username, &quote.Message, &quote.CreatedAt)
		if err != nil {
			continue
		}
		quotes = append([]common.Quote{quote}, quotes...)
	}

	return quotes, err
}

func GetAllLogs(channel string) ([]common.Quote, error) {
	quotes := []common.Quote{}
	stmt, err := db.Prepare("SELECT username, message, created_at FROM messages WHERE channel_name = ?")
	if err != nil {
		return quotes, err
	}

	rows, err := stmt.Query(channel)
	if err != nil {
		return quotes, err
	}

	for rows.Next() {
		quote := common.Quote{}
		err = rows.Scan(&quote.Username, &quote.Message, &quote.CreatedAt)
		if err != nil {
			continue
		}
		quotes = append([]common.Quote{quote}, quotes...)
	}

	return quotes, err
}

func LogSearch(channel string, query string, mins int) ([]common.Quote, error) {
	quotes := []common.Quote{}
	stmt, err := db.Prepare(`
		SELECT username, message, created_at FROM messages WHERE channel_name = ? AND
		(created_at >= CURRENT_TIMESTAMP - INTERVAL ? MINUTE) AND (message LIKE ?)
	`)
	if err != nil {
		return quotes, err
	}

	nameMap := map[string]bool{}
	rows, err := stmt.Query(channel, mins, "%"+query+"%")
	if err != nil {
		return quotes, err
	}

	for rows.Next() {
		quote := common.Quote{}
		err = rows.Scan(&quote.Username, &quote.Message, &quote.CreatedAt)
		if err != nil {
			continue
		}
		_, isKnown := nameMap[quote.Username]
		if !isKnown {
			nameMap[quote.Username] = true
			quotes = append(quotes, quote)
		}
	}

	return quotes, err
}

func UpdateUsersTable(username string, channel string, msg string) error {
	discriminator := username + "::_::" + channel
	stmt, err := db.Prepare(`
		INSERT INTO users (username, channel, discriminator, first_message, last_message)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE message_count = message_count + 1, last_message = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, channel, discriminator, msg, msg, msg)
	return err
}

func IncrementChannel(channel string) error {
	stmt, err := db.Prepare(`UPDATE channels SET message_count = message_count + 1 WHERE name = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(channel)
	return err
}
