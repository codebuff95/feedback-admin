package subject

import(
  //"github.com/codebuff95/uafm"
  //"github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-admin/user"
  "feedback-admin/college"
  "feedback-admin/database"
  "feedback-admin/templates"
  "net/http"
  "time"
  "log"
  "errors"
  "html/template"
  "gopkg.in/mgo.v2/bson"
)

type Subject struct{
  Subjectid string `bson:"subjectid"`
  Subjectname string `bson:"subjectname"`
  Addedon string  `bson:"addedon"`
}

type SubjectPage struct{
  Subjects *[]Subject
  FormSid *string
  Collegename *string
}

func GetSubjects() (*[]Subject,error){
  log.Println("Getting Subjects Slice")
  var mySubjectSlice []Subject
  err := database.SubjectCollection.Find(nil).All(&mySubjectSlice)
  if err != nil{
    log.Println("Error finding subjects:",err)
    return nil,err
  }
  log.Println("Success getting subjectSlice of size:",len(mySubjectSlice))
  if len(mySubjectSlice) == 0{
    return nil,nil
  }
  return &mySubjectSlice,nil
}

func GetSubject(sid string) (*Subject,error){
  log.Println("Getting subject with subjectid:",sid)
  var mySubject *Subject = &Subject{}
  err := database.SubjectCollection.Find(bson.M{"subjectid" : sid}).Limit(1).One(mySubject)
  if err != nil{
    log.Println("Could not get subject.")
    return nil,err
  }
  return mySubject,err
}

func displaySubjectPage(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying subject page to user.")
  var mySubjectPage SubjectPage
  formSid, err := formsession.CreateSession("0",time.Minute*10) //Form created will be valid for 10 minutes.
  if err != nil{
    log.Println("Error creating new session for Subject form:",err)
    displayBadPage(w,r,err)
    return
  }
  mySubjectPage.FormSid = formSid
  mySubjectPage.Collegename = &college.GlobalDetails.Collegename
  log.Println("Creating new Subject page to client",r.RemoteAddr,"with formSid:",*mySubjectPage.FormSid)  //Enter client ip address and new form SID.
  mySubjectPage.Subjects,_ = GetSubjects()
  templates.SubjectPageTemplate.Execute(w,mySubjectPage)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func SubjectHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***SUBJECT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  displaySubjectPage(w,r)
}

func AddSubjectHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***ADD SUBJECT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating form and Parsing form for new Subject data.")

  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please add new subject properly"))
    return
  }
  r.ParseForm()

  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err,". Displaying badpage.")
    displayBadPage(w,r,err)
    return
  }

  enteredSubjectId := template.HTMLEscapeString(r.Form.Get("subjectid"))
  enteredSubjectName := template.HTMLEscapeString(r.Form.Get("subjectname"))
  addedOn := time.Now().Format("2006-01-02 15:04:05")
  myNewSubject := &Subject{Subjectid : enteredSubjectId, Subjectname : enteredSubjectName, Addedon : addedOn}
  log.Println("Adding new Subject with details:",*myNewSubject)
  err = AddSubject(myNewSubject)
  if err != nil{
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully added new Subject with Id:",enteredSubjectId,", name:",enteredSubjectName)
  http.Redirect(w, r, "/subject", http.StatusSeeOther)
}


func RemoveSubjectHandler(w http.ResponseWriter,r *http.Request){
  log.Println("***REMOVE SUBJECT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating remove Subject request.")

  if r.Method != "GET"{
    displayBadPage(w,r,errors.New("Please remove subject properly"))
    return
  }

  mySubjectId := template.HTMLEscapeString(r.URL.Path[len("/removesubject/"):])
  _,err = RemoveSubject(mySubjectId)
  if err != nil{
    log.Println("Problem removing Subject:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Success removing Subject:",mySubjectId)
  http.Redirect(w, r, "/subject", http.StatusSeeOther)
}


func AddSubject(mysubject *Subject) error{
  return database.SubjectCollection.Insert(mysubject)
}

func RemoveSubject(id string) (int,error){
   changeinfo,err := database.SubjectCollection.RemoveAll(bson.M{"subjectid":id})
  if changeinfo == nil{
    return 0,err
  }
  return changeinfo.Removed,err
}
