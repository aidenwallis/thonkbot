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
