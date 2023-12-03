package eternals

func getAllUsers() ([]user, error) {
	rows, err := db.Query("SELECT id, password, login, username, point FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var user user
		err := rows.Scan(&user.ID, &user.Password, &user.Login, &user.Username, &user.Point)
		if err != nil {
			return nil, err
		}
		if user.ID == 1 {
			continue
		} else {
			users = append(users, user)
		}

	}

	return users, nil
}
