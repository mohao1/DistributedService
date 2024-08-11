package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}

	return result / float32(len(s.Grades))
}

func (ss Students) GetByID(id int) (*Student, error) {
	for i, s := range ss {
		if s.ID == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("student with ID %d not found", id)
}

type Students []Student

var (
	students     Students
	studentMutex sync.RWMutex
)

type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
