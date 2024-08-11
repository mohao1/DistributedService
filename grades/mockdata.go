package grades

func init() {
	// 初始化模拟的数据
	students = []Student{
		{
			ID:        1,
			FirstName: "Nick",
			LastName:  "Carter",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 94,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 82,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Emma",
			LastName:  "Watson",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 78,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 88,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 91,
				},
			},
		},
		{
			ID:        3,
			FirstName: "Liam",
			LastName:  "Neeson",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 92,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 85,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 89,
				},
			},
		},
		{
			ID:        4,
			FirstName: "Sophia",
			LastName:  "Turner",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 80,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 92,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 84,
				},
			},
		},
		{
			ID:        5,
			FirstName: "Mason",
			LastName:  "Green",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 76,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 89,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 90,
				},
			},
		},
		{
			ID:        6,
			FirstName: "Isabella",
			LastName:  "Clark",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 95,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 88,
				},
			},
		},
		{
			ID:        7,
			FirstName: "Lucas",
			LastName:  "Brown",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 79,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 87,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 83,
				},
			},
		},
		{
			ID:        8,
			FirstName: "Ava",
			LastName:  "Williams",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 88,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 90,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 86,
				},
			},
		},
		{
			ID:        9,
			FirstName: "Ethan",
			LastName:  "Johnson",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 81,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 84,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 79,
				},
			},
		},
		{
			ID:        10,
			FirstName: "Mia",
			LastName:  "Smith",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 90,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 93,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 87,
				},
			},
		},
		{
			ID:        11,
			FirstName: "Oliver",
			LastName:  "Davis",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 86,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 89,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 84,
				},
			},
		},
	}

}
