package mysql

import "database/sql"

func IsAdmin(username string) (bool, error) {
	stmt, err := db.Prepare("SELECT userlevel FROM bot_admins WHERE username = ? LIMIT 1")
	if err != nil {
		return false, err
	}

	var userlevel int
	row := stmt.QueryRow(username)
	err = row.Scan(&userlevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
