package feedback

import (
	"feedback-admin/database"
	"feedback-admin/question"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Feedback struct {
	Sectionid             string                          `bson:"sectionid"`
	Ratings               []Rating                        `bson:"ratings"`
	Textquestionresponses []question.TextQuestionResponse `bson:"textquestionresponses"`
	Addedon               string                          `bson:"addedon"`
}

type Rating struct {
	Facultyid string  `bson:"facultyid"`
	Subjectid string  `bson:"subjectid"`
	Points    []Point `bson:"points"`
}

type Point struct {
	Questionid string `bson:"questionid"`
	Marks      int    `bson:"marks"`
}

func GetFeedbacks(sectionId string) (*[]Feedback, error) {
	var myFeedbackSlice []Feedback
	err := database.FeedbackCollection.Find(bson.M{"sectionid": sectionId}).All(&myFeedbackSlice)
	if err != nil {
		return nil, err
	}
	if len(myFeedbackSlice) == 0 {
		return nil, nil
	}
	return &myFeedbackSlice, nil
}
