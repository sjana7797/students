package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sjana7797/students/internal/storage"
	"github.com/sjana7797/students/internal/types"
	"github.com/sjana7797/students/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			fmtErr := fmt.Errorf("empty body")
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmtErr))
			slog.Error("Student Failed to create", slog.String("error", fmtErr.Error()))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			slog.Error("Student Failed to create", slog.String("error", err.Error()))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidatorError(validateErr))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.Info("student create successfully with id", slog.String("id", string(id)))

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func GetStudents(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		studentsData, err := storage.GetStudents()

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.Info("students fetched successfully")

		students := struct {
			Students []types.Student `json:"students"`
			Total    int             `json:"total"`
		}{
			Students: studentsData,
			Total:    len(studentsData),
		}
		response.WriteJSON(w, http.StatusAccepted, students)
	}

}

func GetStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")

		slog.Info("Getting student with id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)

	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("Updating student with id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		studentBody := types.UpdateStudent{}

		// Read JSON body
		err = json.NewDecoder(r.Body).Decode(&studentBody)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		updateStudentId, err := storage.UpdateStudentById(intId, studentBody)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(updateStudentId)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusAccepted, student)
	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("Deleting student with id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		err = storage.DeleteStudentById(intId)

		if err != nil {
			slog.Error(err.Error())
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)

	}
}
