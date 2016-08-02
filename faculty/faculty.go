package faculty

import(
  //"github.com/codebuff95/uafm"
  //"github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-admin/user"
  "feedback-admin/database"
  "feedback-admin/templates"
  "net/http"
  "time"
  "log"
  "errors"
  "html/template"
  "gopkg.in/mgo.v2/bson"
)

type Faculty struct{
  Facultyid string `bson:"facultyid"`
  Facultyname string `bson:"facultyname"`
  Addedon string  `bson:"addedon"`
}

type FacultyPage struct{
  Facultys *[]Faculty
  FormSid *string
}

func GetFacultys() (*[]Faculty,error){
  log.Println("Getting Facultys Slice")
  var myFacultySlice []Faculty
  err := database.FacultyCollection.Find(nil).All(&myFacultySlice)
  if err != nil{
    log.Println("Error finding facultys:",err)
    return nil,err
  }
  log.Println("Success getting facultySlice of size:",len(myFacultySlice))
  if len(myFacultySlice) == 0{
    return nil,nil
  }
  return &myFacultySlice,nil
}

func GetFaculty(fid string) (*Faculty,error){
  log.Println("Getting faculty with facultyid:",fid)
  var myFaculty *Faculty = &Faculty{}
  err := database.FacultyCollection.Find(bson.M{"facultyid" : fid}).Limit(1).One(myFaculty)
  if err != nil{
    log.Println("Could not get faculty.")
    return nil,err
  }
  return myFaculty,err
}

func displayFacultyPage(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying faculty page to user.")
  var myFacultyPage FacultyPage
  formSid, err := formsession.CreateSession("0",time.Minute*10) //Form created will be valid for 10 minutes.
  if err != nil{
    log.Println("Error creating new session for Faculty form:",err)
    displayBadPage(w,r,err)
    return
  }
  myFacultyPage.FormSid = formSid
  log.Println("Creating new Faculty page to client",r.RemoteAddr,"with formSid:",*myFacultyPage.FormSid)  //Enter client ip address and new form SID.
  myFacultyPage.Facultys,_ = GetFacultys()
  templates.FacultyPageTemplate.Execute(w,myFacultyPage)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func FacultyHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***FACULTY HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  displayFacultyPage(w,r)
}

func AddFacultyHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***ADD FACULTY HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating form and Parsing form for new Faculty data.")

  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please add new faculty properly"))
    return
  }
  r.ParseForm()

  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err,". Displaying badpage.")
    displayBadPage(w,r,err)
    return
  }

  enteredFacultyId := template.HTMLEscapeString(r.Form.Get("facultyid"))
  enteredFacultyName := template.HTMLEscapeString(r.Form.Get("facultyname"))
  addedOn := time.Now().Format("2006-01-02 15:04:05")
  myNewFaculty := &Faculty{Facultyid : enteredFacultyId, Facultyname : enteredFacultyName, Addedon : addedOn}
  log.Println("Adding new Faculty with details:",*myNewFaculty)
  err = AddFaculty(myNewFaculty)
  if err != nil{
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully added new Faculty with Id:",enteredFacultyId,", name:",enteredFacultyName)
  http.Redirect(w, r, "/faculty", http.StatusSeeOther)
}


func RemoveFacultyHandler(w http.ResponseWriter,r *http.Request){
  log.Println("***REMOVE FACULTY HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating remove Faculty request.")

  if r.Method != "GET"{
    displayBadPage(w,r,errors.New("Please remove faculty properly"))
    return
  }

  myFacultyId := template.HTMLEscapeString(r.URL.Path[len("/removefaculty/"):])
  _,err = RemoveFaculty(myFacultyId)
  if err != nil{
    log.Println("Problem removing Faculty:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Success removing Faculty:",myFacultyId)
  http.Redirect(w, r, "/faculty", http.StatusSeeOther)
}


func AddFaculty(myfaculty *Faculty) error{
  return database.FacultyCollection.Insert(myfaculty)
}

func RemoveFaculty(id string) (int,error){
   changeinfo,err := database.FacultyCollection.RemoveAll(bson.M{"facultyid":id})
  if changeinfo == nil{
    return 0,err
  }
  return changeinfo.Removed,err
}
