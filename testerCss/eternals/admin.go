package eternals

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func ControlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		spoints := r.FormValue("Points")
		point, err := strconv.Atoi(spoints)
		if err != nil {
			log.Fatal(err)
		}

		updatePoint(db, 5, point)
		fmt.Println("spoints:", spoints)
		// for id, pointStr := range spoints {
		// 	point, err := strconv.Atoi(pointStr)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	fmt.Println("id:", id+2)
		// 	fmt.Println("point:", point)
		// 	updatePoint(db, id, point)
		// }
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case "GET":
		tmpl, err := template.ParseFiles("templates/adminpage.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userID := id
		fmt.Println(id)
		row := db.QueryRow("SELECT ID, Password, Login, Username, Point FROM user WHERE ID = ?", userID)
		var Userinfo user
		err = row.Scan(&Userinfo.ID, &Userinfo.Password, &Userinfo.Login, &Userinfo.Username, &Userinfo.Point)
		if err != nil {
			fmt.Println("Error scanning row:", err)
		}
		fmt.Println(Userinfo.ID, Userinfo.Username, Userinfo.Point)
		allusers, err := getAllUsers()
		if err != nil {
			log.Fatal(err)
		}
		UserinfoWithAll := struct {
			User    user
			AllUser []user
		}{
			User:    Userinfo,
			AllUser: allusers,
		}
		fmt.Println(dayOfWeek.String())
		tmpl.Execute(w, UserinfoWithAll)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This function does not support " + r.Method + " method."))
	}
}

func updatePoint(db *sql.DB, userID, newPointValue int) error {
	// Подготовка SQL-запроса UPDATE
	stmt, err := db.Prepare("UPDATE user SET point = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса с передачей нового значения "point" и ID пользователя
	_, err = stmt.Exec(newPointValue, userID)
	return err
}

func convertStringSliceToIntSlice(strSlice []string) ([]int, error) {
	var intSlice []int

	for _, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to int: %v", str, err)
		}
		intSlice = append(intSlice, num)
	}

	return intSlice, nil
}
