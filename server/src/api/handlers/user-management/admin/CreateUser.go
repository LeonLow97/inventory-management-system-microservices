package admin

import (
	"net/http"
)

func AdminCreateUser(w http.ResponseWriter, req *http.Request) {
	// // Set Headers
	// w.Header().Set("Content-Type", "application/json")
	// var adminNewUser auth_management.AdminUserMgmtJson

	// // Reading the request body and UnMarshal the body to the AdminUserMgmt struct
	// bs, _ := io.ReadAll(req.Body)
	// if err := json.Unmarshal(bs, &adminNewUser); err != nil {
	// 	utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 	log.Println("Internal Server Error in UnMarshal JSON body in AdminCreateUser route:", err)
	// 	return
	// }

	// // Check User Group Admin
	// if !auth_management.RetrieveIssuer(w, req) {
	// 	return
	// }
	// if !utils.CheckUserGroup(w, w.Header().Get("username"), "Admin") {
	// 	return
	// }

	// // Trim white spaces (username, password, email, company name)
	// adminNewUser = adminNewUser.AdminUserMgmtFieldsTrimSpaces()

	// // Validate form inputs
	// if !auth_management.AdminUserMgmtFormValidation(w, adminNewUser, "CREATE_USER") {
	// 	return
	// }

	// hashedPassword := utils.GenerateHash(adminNewUser.Password)

	// // Check if username already exists in database (duplicates not allowed)
	// isExistingUsername := database.GetUsername(adminNewUser.Username)
	// if isExistingUsername {
	// 	utils.WriteJSON(w, http.StatusBadRequest, "Username has already been taken. Please try again.")
	// 	return
	// }

	// // Check if email already exists in database (duplicates not allowed)
	// isExistingEmail := database.GetEmail(adminNewUser.Email)
	// if isExistingEmail {
	// 	utils.WriteJSON(w, http.StatusBadRequest, "Email address has already been taken. Please try again.")
	// 	return
	// }

	// // Check if organisation exists in database
	// isExistingOrganisation := database.GetOrganisationName(adminNewUser.OrganisationName)
	// if !isExistingOrganisation {
	// 	utils.WriteJSON(w, http.StatusNotFound, "Organisation name cannot be found. Please try again.")
	// 	return
	// }

	// // Check if user group is valid and trim user group
	// isValidUserGroup := auth_management.UserGroupFormValidation(w, adminNewUser.UserGroup)
	// if !isValidUserGroup {
	// 	return
	// }

	// // Insert users table
	// userId, err := database.InsertNewUser(adminNewUser.Username, hashedPassword, adminNewUser.Email, adminNewUser.IsActive)
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 	log.Println("Error inserting new user via admin route:", err)
	// 	return
	// }

	// // Insert user_organisation_mapping table
	// err = database.InsertIntoUserOrganisationMapping(userId, adminNewUser.OrganisationName)
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 	log.Println("Error inserting into user_organisation_mapping table:", err)
	// 	return
	// }

	// // Insert multiple user groups to user_group_mapping table
	// for _, ug := range adminNewUser.UserGroup {
	// 	err = database.InsertIntoUserGroupMapping(userId, ug)
	// 	if err != nil {
	// 		utils.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
	// 		log.Println("Error inserting into user_group_mapping table:", err)
	// 		return
	// 	}
	// }

	// utils.WriteJSON(w, http.StatusOK, "Admin Successfully Created User!")
}