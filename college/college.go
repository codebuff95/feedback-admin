package college

import(
  "os"
  "encoding/json"
)

type CollegeDetails struct{
  Collegename string  `bson:"collegename"`
}

var GlobalDetails CollegeDetails

func InitCollegeDetails() error{
  myFile, err := os.Open("feedbackadminres/faconfig.json")
  defer myFile.Close()

  if err != nil{
    return err
  }
  err = json.NewDecoder(myFile).Decode(&GlobalDetails)
  if err != nil{
    return err
  }
  return nil
}
