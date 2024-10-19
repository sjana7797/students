package storage

import (
	"github.com/sjana7797/students/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudents() ([]types.Student, error)
	GetStudentById(id int64) (types.Student, error)
}
