package course

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

type Course struct{
  Courseid string `bson:"courseid"`
  Coursename string `bson:"coursename"`
  Addedon string  `bson:"addedon"`
}

type CoursePage struct{
  Courses *[]Course
  FormSid *string
}

func GetCourses() (*[]Course,error){
  log.Println("Getting Courses Slice")
  var myCourseSlice []Course
  err := database.CourseCollection.Find(nil).All(&myCourseSlice)
  if err != nil{
    log.Println("Error finding courses:",err)
    return nil,err
  }
  log.Println("Success getting courseSlice of size:",len(myCourseSlice))
  if len(myCourseSlice) == 0{
    return nil,nil
  }
  return &myCourseSlice,nil
}

func GetCourse(cid string) (*Course,error){
  log.Println("Getting course with courseid:",cid)
  var myCourse *Course = &Course{}
  err := database.CourseCollection.Find(bson.M{"courseid" : cid}).Limit(1).One(myCourse)
  if err != nil{
    log.Println("Could not get course.")
    return nil,err
  }
  return myCourse,err
}

func displayCoursePage(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying course page to user.")
  var myCoursePage CoursePage
  formSid, err := formsession.CreateSession("0",time.Minute*10) //Form created will be valid for 10 minutes.
  if err != nil{
    log.Println("Error creating new session for course form:",err)
    displayBadPage(w,r,err)
    return
  }
  myCoursePage.FormSid = formSid
  log.Println("Creating new course page to client",r.RemoteAddr,"with formSid:",*myCoursePage.FormSid)  //Enter client ip address and new form SID.
  myCoursePage.Courses,_ = GetCourses()
  templates.CoursePageTemplate.Execute(w,myCoursePage)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func CourseHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***COURSE HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  displayCoursePage(w,r)
}

func AddCourseHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***ADD COURSE HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating form and Parsing form for new course data.")

  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please add new course properly"))
    return
  }
  r.ParseForm()

  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err,". Displaying badpage.")
    displayBadPage(w,r,err)
    return
  }

  enteredCourseId := template.HTMLEscapeString(r.Form.Get("courseid"))
  enteredCourseName := template.HTMLEscapeString(r.Form.Get("coursename"))
  addedOn := time.Now().Format("2006-01-02 15:04:05")
  myNewCourse := &Course{Courseid : enteredCourseId, Coursename : enteredCourseName, Addedon : addedOn}
  log.Println("Adding new course with details:",*myNewCourse)
  err = AddCourse(myNewCourse)
  if err != nil{
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully added new course with Id:",enteredCourseId,", name:",enteredCourseName)
  http.Redirect(w, r, "/course", http.StatusSeeOther)
}


func RemoveCourseHandler(w http.ResponseWriter,r *http.Request){
  log.Println("***REMOVE COURSE HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating remove course request.")

  if r.Method != "GET"{
    displayBadPage(w,r,errors.New("Please remove course properly"))
    return
  }

  myCourseId := template.HTMLEscapeString(r.URL.Path[len("/removecourse/"):])
  _,err = RemoveCourse(myCourseId)
  if err != nil{
    log.Println("Problem removing course:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Success removing course:",myCourseId)
  http.Redirect(w, r, "/course", http.StatusSeeOther)
}


func AddCourse(mycourse *Course) error{
  return database.CourseCollection.Insert(mycourse)
}

func RemoveCourse(id string) (int,error){
   changeinfo,err := database.CourseCollection.RemoveAll(bson.M{"courseid":id})
  if changeinfo == nil{
    return 0,err
  }
  return changeinfo.Removed,err
}
