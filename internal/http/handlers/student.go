package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sjana7797/students/internal/types"
	"github.com/sjana7797/students/internal/utils/response"
)

func New() http.HandlerFunc {
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

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
