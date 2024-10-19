package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sjana7797/students/internal/config"
	"github.com/sjana7797/students/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS student (
	 id INTEGER PRIMARY KEY AUTOINCREMENT,
	 name TEXT,
	 age INTEGER,
	 email TEXT UNIQUE
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO student (name,email,age) VALUES(?,?,?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	rows, err := s.Db.Query("SELECT id,name,email,age FROM student")

	if err != nil {
		log.Fatalf("Failed to query students: %v", err)
	}
	defer rows.Close()

	var students []types.Student

	// Loop through and print each row
	for rows.Next() {

		var student types.Student

		err = rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			log.Fatal(err)
		}

		students = append(students, student)

	}

	return students, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare(`SELECT id,name,email,age FROM student WHERE id = ? LIMIT 1`)

	if err != nil {
		slog.Error(err.Error())
		return types.Student{}, nil
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		slog.Error(err.Error())
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error:%w", err)
	}

	return student, nil

}
