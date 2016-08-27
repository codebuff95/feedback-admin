package question

import (
	"errors"
	"feedback-admin/database"
	"log"
)

type Question struct {
	Questionid string `bson:"questionid"`
	Text       string `bson:"text"`
	Weightage  int    `bson:"weightage"`
}

type TextQuestion struct {
	Questionid string `bson:"questionid"`
	Text       string `bson:"text"`
}

type TextQuestionResponse struct {
	Questionid string `bson:"questionid"`
	Response   string `bson:"response"`
}

func GetAllQuestionsWeightage() (int, error) {
	log.Println("**Getting All Questions Weightage**")
	var myQuestionSlice []Question
	err := database.QuestionCollection.Find(nil).All(&myQuestionSlice)
	if err != nil || len(myQuestionSlice) == 0 {
		log.Println("Error finding questions:", err)
		return 0, err
	}
	log.Println("Success getting questionSlice of size:", len(myQuestionSlice))
	sum := 0
	for i := 0; i < len(myQuestionSlice); i++ {
		sum += myQuestionSlice[i].Weightage
	}
	return sum, nil
}

func GetQuestions() (*[]Question, error) {
	var myQuestions []Question
	err := database.QuestionCollection.Find(nil).All(&myQuestions)
	if len(myQuestions) == 0 || err != nil {
		return nil, errors.New("No questions found")
	}
	return &myQuestions, nil
}

func GetTextQuestions() (*[]TextQuestion, error) {
	var myTextQuestions []TextQuestion
	err := database.TextQuestionCollection.Find(nil).All(&myTextQuestions)
	if len(myTextQuestions) == 0 || err != nil {
		return nil, errors.New("No text questions found")
	}
	return &myTextQuestions, nil
}
