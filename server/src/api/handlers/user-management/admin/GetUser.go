package admin

import (
	"net/http"
)

func AdminGetUsers(w http.ResponseWriter, req *http.Request) {

	// w.Header().Set("Content-Type", "application/json")

	// // Check User Group (Admin)
	// if !auth_management.RetrieveIssuer(w, req) {
	// 	return
	// }
	// if !utils.CheckUserGroup(w, w.Header().Get("username"), "Admin") {
	// 	return
	// }

	// var data []handlers.User
	// var users = make(map[int]handlers.User) // A map to store unique user records by UserId

	// // To handle nullable columns in a database table
	// var username, email, organisationName, userGroup, addedDate, updatedDate sql.NullString
	// var userId, isActive sql.NullInt16

	// rows, err := database.GetUsers()
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 	log.Println("Internal Server Error in GetUsers:", err)
	// 	return
	// }

	// defer rows.Close()

	// for rows.Next() {
	// 	err = rows.Scan(&userId, &username, &email, &isActive, &organisationName, &userGroup, &addedDate, &updatedDate)
	// 	if err != nil {
	// 		utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 		log.Println("Internal Server Error in Scanning GetUsers:", err)
	// 		return
	// 	}

	// 	// Check if user already exists in map
	// 	if user, ok := users[int(userId.Int16)]; ok {
	// 		// User already exists, append userGroup to UserGroup array
	// 		user.UserGroup = append(user.UserGroup, userGroup.String)
	// 		users[int(userId.Int16)] = user
	// 	} else {
	// 		// User does not exist in map, create a new User object
	// 		user := handlers.User{
	// 			UserId:           int(userId.Int16),
	// 			Username:         username.String,
	// 			Email:            email.String,
	// 			IsActive:         int(isActive.Int16),
	// 			OrganisationName: organisationName.String,
	// 			UserGroup:        []string{userGroup.String},
	// 			AddedDate:        addedDate.String,
	// 			UpdatedDate:      updatedDate.String,
	// 		}
	// 		users[int(userId.Int16)] = user
	// 	}
	// }

	// // Convert map to slice
	// for _, user := range users {
	// 	data = append(data, user)
	// }

	// // Sort users by UserId in ascending order
	// sort.Slice(data, func(i, j int) bool {
	// 	return data[i].UserId < data[j].UserId
	// })

	// jsonStatus := struct {
	// 	Code     int             `json:"code"`
	// 	Response []handlers.User `json:"response"`
	// }{
	// 	Response: data,
	// 	Code:     http.StatusOK,
	// }

	// bs, err := json.Marshal(jsonStatus)
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 	log.Println("Internal Server Error in Marshal JSON body in GetUsers:", err)
	// 	return
	// }

	// io.WriteString(w, string(bs))
}