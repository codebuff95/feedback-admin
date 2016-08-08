package question

import(
  "feedback-admin/database"
  "log"
  "errors"
)

type Question struct{
  Questionid string `bson:"questionid"`
  Text string `bson:"text"`
  Weightage int `bson:"weightage"`
}

func GetAllQuestionsWeightage() (int,error){
  log.Println("**Getting All Questions Weightage**")
  var myQuestionSlice []Question
  err := database.QuestionCollection.Find(nil).All(&myQuestionSlice)
  if err != nil || len(myQuestionSlice) == 0{
    log.Println("Error finding questions:",err)
    return 0,err
  }
  log.Println("Success getting questionSlice of size:",len(myQuestionSlice))
  sum := 0
  for i := 0; i < len(myQuestionSlice); i++{
    sum += myQuestionSlice[i].Weightage
  }
  return sum,nil
}

func GetQuestions() (*[]Question,error){
  var myQuestions []Question
  err := database.QuestionCollection.Find(nil).All(&myQuestions)
  if len(myQuestions) == 0 || err != nil{
    return nil,errors.New("No questions found")
  }
  return &myQuestions,nil
}
