package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/notherealmarco/WASAPhoto/service/structures"
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

func VerifyUserOrNotFound(db database.AppDatabase, uid string, w http.ResponseWriter, l logrus.FieldLogger) bool {

	user_exists, err := db.UserExists(uid)

	if err != nil {
		SendInternalError(err, "Error verifying user existence", w, l)
		return false
	}

	if !user_exists {
		SendNotFound(w, "User not found", l)
		return false
	}
	return true
}

func SendStatus(httpStatus int, w http.ResponseWriter, description string, l logrus.FieldLogger) {
	w.WriteHeader(httpStatus)
	err := json.NewEncoder(w).Encode(structures.GenericResponse{Status: description})
	if err != nil {
		l.WithError(err).Error("Error encoding json")
	}
}

func SendNotFound(w http.ResponseWriter, description string, l logrus.FieldLogger) {
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(structures.GenericResponse{Status: description})
	if err != nil {
		l.WithError(err).Error("Error encoding json")
	}
}

func SendBadRequest(w http.ResponseWriter, description string, l logrus.FieldLogger) {
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(structures.GenericResponse{Status: description})
	if err != nil {
		l.WithError(err).Error("Error encoding json")
	}
}

func SendBadRequestError(err error, description string, w http.ResponseWriter, l logrus.FieldLogger) {
	w.WriteHeader(http.StatusBadRequest)
	l.WithError(err).Error(description)
	err = json.NewEncoder(w).Encode(structures.GenericResponse{Status: description})
	if err != nil {
		l.WithError(err).Error("Error encoding json")
	}
}

func SendInternalError(err error, description string, w http.ResponseWriter, l logrus.FieldLogger) {
	w.WriteHeader(http.StatusInternalServerError)
	l.WithError(err).Error(description)
	err = json.NewEncoder(w).Encode(structures.GenericResponse{Status: description})
	if err != nil {
		l.WithError(err).Error("Error encoding json")
	}
}

func RollbackOrLogError(tx database.DBTransaction, l logrus.FieldLogger) {
	err := tx.Rollback()
	if err != nil {
		l.WithError(err).Error("Error rolling back transaction")
	}
}

func SendNotFoundIfBanned(db database.AppDatabase, uid string, banner string, w http.ResponseWriter, l logrus.FieldLogger) bool {
	banned, err := db.IsBanned(uid, banner)
	if err != nil {
		SendInternalError(err, "Database error: IsBanned", w, l)
		return false
	}
	if banned {
		SendNotFound(w, "User not found", l)
		return false
	}
	return true
}
