package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/sirupsen/logrus"
)

func DecodeJsonOrBadRequest(r io.Reader, w http.ResponseWriter, v interface{}, l logrus.FieldLogger) bool {

	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		SendInternalError(err, "Error decoding json", w, l)
		return false
	}
	return true
}

func VerifyUserOrNotFound(db database.AppDatabase, uid string, w http.ResponseWriter) bool {

	user_exists, err := db.UserExists(uid)

	if err != nil {
		SendInternalError(err, "Error verifying user existence", w, nil)
		return false
	}

	if !user_exists {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	return true
}

func SendInternalError(err error, description string, w http.ResponseWriter, l logrus.FieldLogger) {
	w.WriteHeader(http.StatusInternalServerError)
	l.WithError(err).Error(description)
	w.Write([]byte(description)) // todo: maybe in json. But it's not important to send the full error to the client
}
