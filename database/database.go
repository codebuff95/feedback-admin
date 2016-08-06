package database

import(
  "gopkg.in/mgo.v2"
  "log"
  "os"
  "encoding/json"
)


var DatabaseName string
var DatabaseSession *mgo.Session

var AdminCollection *mgo.Collection
var CourseCollection *mgo.Collection
var SubjectCollection *mgo.Collection
var FacultyCollection *mgo.Collection
var SectionCollection *mgo.Collection
var FeedbackCollection *mgo.Collection

type DatabaseDetails struct{
  Url string `json:"url"`
  Dbname string `json:"dbname"`
}

type InitDatabaseError int

func (err InitDatabaseError) Error() string{
  if err == PROBLEMOPENINGFILE{
    return "Could not open file faconfig.json"
  }
  return "Miscellaneous Issues in database"
}

const(
  PROBLEMOPENINGFILE = 0
  PROBLEMDECODING = 1
)

func InitDatabaseSession() error{
  var err error

  myFile, err := os.Open("feedbackadminres/faconfig.json")
  defer myFile.Close()

  if err != nil{
    return InitDatabaseError(PROBLEMOPENINGFILE)
  }
  var myDBDetails DatabaseDetails
  err = json.NewDecoder(myFile).Decode(&myDBDetails)
  if err != nil{
    return InitDatabaseError(PROBLEMDECODING)
  }
  DatabaseSession,err = mgo.Dial(myDBDetails.Url)
  DatabaseName = myDBDetails.Dbname
  log.Println("feedback-admin: Initialised new database session to url:",myDBDetails.Url,"and Dbname:",myDBDetails.Dbname)
  return nil
}

func InitCollections(){
  //log.Println("**Initialising Essential Collections with Database",DatabaseName,"**")
  AdminCollection = DatabaseSession.DB(DatabaseName).C("admin")
  CourseCollection = DatabaseSession.DB(DatabaseName).C("course")
  SubjectCollection = DatabaseSession.DB(DatabaseName).C("subject")
  FacultyCollection = DatabaseSession.DB(DatabaseName).C("faculty")
  SectionCollection = DatabaseSession.DB(DatabaseName).C("section")
  FeedbackCollection = DatabaseSession.DB(DatabaseName).C("feedback")
}
