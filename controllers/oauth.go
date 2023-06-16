// Handler function for routes
package controllers

import (
	"encoding/json"
	"fmt"
	"kwoc-backend/utils"
	"net/http"

	"github.com/kossiitkgp/kwoc-db-models/models"
	"github.com/rs/zerolog/log"
)

type OAuthReqBodyFields struct {
	// Code generated by Github OAuth
	Code string `json:"code"`
	// `mentor` or `student`
	Type string `json:"type"`
}

type OAuthResBodyFields struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	// `mentor` or `student`
	Type string `json:"type"`
	// Whether the user has newly registered or was registered before
	IsNewUser bool   `json:"isNewUser"`
	Jwt       string `json:"jwt"`
}

func OAuth(w http.ResponseWriter, r *http.Request) {
	var reqFields = OAuthReqBodyFields{}
	err := json.NewDecoder(r.Body).Decode(&reqFields)

	if err != nil {
		log.Err(err).Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			"Error parsing body parameters.",
		)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error parsing body parameters.")
		return
	}

	if reqFields.Code == "" || reqFields.Type == "" {
		log.Warn().Msgf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			"Empty body parameters.",
		)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Empty body parameters.")
		return
	}

	// Get a Github OAuth access token
	accessToken, err := utils.GetOauthAccessToken(reqFields.Code)
	if err != nil {
		log.Err(err).Msg("Error getting OAuth access token.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error getting OAuth access token.")
		return
	}

	// Get the user's information from the Github API
	userInfo, err := utils.GetOauthUserInfo(accessToken)
	if err != nil {
		log.Err(err).Msg("Error getting OAuth user info.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error getting OAuth user info.")
		return
	}

	db, err := utils.GetDB()
	if err != nil {
		log.Err(err).Msg("Error connecting to the database.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error connecting to the database.")
		return
	}

	// Check if the user has already registered
	var isNewUser bool = false

	if reqFields.Type == "student" {
		student := models.Student{}
		db.
			Table("students").
			Where("username = ?", userInfo.Username).
			First(&student)

		isNewUser = student.Username != userInfo.Username
	} else if reqFields.Type == "mentor" {
		mentor := models.Mentor{}
		db.
			Table("mentors").
			Where("username = ?", userInfo.Username).
			First(&mentor)

		isNewUser = mentor.Username != userInfo.Username
	}

	// Generate a JWT string for the user
	jwtString, err := utils.GenerateLoginJwtString(utils.LoginJwtFields{
		Username: userInfo.Username,
	})
	if err != nil {
		log.Err(err).Msg("Error generating a JWT string.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error generating a JWT string.")
		return
	}

	resFields := OAuthResBodyFields{
		Username:  userInfo.Username,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Type:      reqFields.Type,
		IsNewUser: isNewUser,
		Jwt:       jwtString,
	}

	resJson, err := json.Marshal(resFields)
	if err != nil {
		log.Err(err).Msg("Error generating response JSON.")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error generating response JSON.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
