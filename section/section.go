package section

import(
  //"github.com/codebuff95/uafm"
  //"github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-admin/user"
  "feedback-admin/course"
  "feedback-admin/database"
  "feedback-admin/faculty"
  "feedback-admin/password"
  "feedback-admin/subject"
  "feedback-admin/templates"
  "net/http"
  "strconv"
  "time"
  "log"
  "errors"
  "html/template"
  "gopkg.in/mgo.v2/bson"
)

type Teacher struct{
  Facultyid string `bson:"facultyid"`
  Subjectid string `bson:"subjectid"`
}

type DetailedTeacher struct{
  Facultyid string `bson:"facultyid"`
  Facultyname string `bson:"facultyname"`
  Subjectid string `bson:"subjectid"`
  Subjectname string `bson:"subjectname"`
}

type Section struct{
  Sectionid string `bson:"sectionid"`
  Sectionname string `bson:"sectionname"`
  Year int `bson:"year"`
  Session int `bson:"session"`
  Courseid string `bson:"courseid"`
  Teachers *[]Teacher `bson:"teachers,omitempty"`
  Password string `bson:"password"`
  Students int `bson:"students"`
  Addedon string  `bson:"addedon"`
}

type SectionPage struct{
  Sections *[]Section
  FormSid *string
  Session *string
}

func GetSections(mysession string) (*[]Section,error){
  log.Println("Getting Sections Slice for session:",mysession)
  mysessionint,err := strconv.Atoi(mysession)
  if err != nil{
    log.Println("Error converting section session to int:",err)
    return nil,err
  }
  var mySectionSlice []Section
  err = database.SectionCollection.Find(bson.M{"session":mysessionint}).Sort("courseid","year","sectionid").All(&mySectionSlice)
  if err != nil{
    log.Println("Error finding sections:",err)
    return nil,err
  }
  log.Println("Success getting sectionSlice of size:",len(mySectionSlice))
  if len(mySectionSlice) == 0{
    return nil,nil
  }
  for i := 0; i < len(mySectionSlice); i++{
    if len(*mySectionSlice[i].Teachers) == 0{
      mySectionSlice[i].Teachers = nil
    }
  }
  return &mySectionSlice,nil
}

func GetSection(sectionid string) (*Section,error){
  log.Println("Getting section with sectionid:",sectionid)
  var mySection *Section = &Section{}
  err := database.SectionCollection.Find(bson.M{"sectionid" : sectionid}).Limit(1).One(mySection)
  if err != nil{
    log.Println("Could not get section.")
    return nil,err
  }
  return mySection,nil
}

func displaySectionPage(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying section page to user.")
  r.ParseForm()
  mysession := r.Form.Get("sectionsession")
  //validate my session
  mySessionInt,err := strconv.Atoi(mysession)
  if err != nil || mySessionInt < 1990 || mySessionInt > 2030{
    displayBadPage(w,r,errors.New("Bad Session"))
    return
  }
  //
  var mySectionPage SectionPage
  formSid, err := formsession.CreateSession("0",time.Minute*10) //Form created will be valid for 10 minutes.
  if err != nil{
    log.Println("Error creating new session for Section form:",err)
    displayBadPage(w,r,err)
    return
  }
  mySectionPage.FormSid = formSid
  log.Println("Creating new Section page to client",r.RemoteAddr,"with formSid:",*mySectionPage.FormSid)  //Enter client ip address and new form SID.
  mySectionPage.Session = &mysession
  mySectionPage.Sections,_ = GetSections(mysession)
  templates.SectionPageTemplate.Execute(w,mySectionPage)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func SectionHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***SECTION HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  displaySectionPage(w,r)
}

func AddSectionHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***ADD SECTION HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating form and Parsing form for new Section data.")

  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please add new section properly"))
    return
  }
  r.ParseForm()

  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err,". Displaying badpage.")
    displayBadPage(w,r,err)
    return
  }

  enteredSectionId := template.HTMLEscapeString(r.Form.Get("sectionid"))
  enteredSectionName := template.HTMLEscapeString(r.Form.Get("sectionname"))
  enteredSectionYear,err := strconv.Atoi(template.HTMLEscapeString(r.Form.Get("sectionyear")))

  if err != nil{
    log.Println("Error validating entered form:",err,". Displaying badpage.")
    displayBadPage(w,r,errors.New("Entered SectionYear invalid."))
    return
  }

  enteredSectionSession,err := strconv.Atoi(template.HTMLEscapeString(r.Form.Get("sectionsession")))

  if err != nil{
    log.Println("Error validating entered form:",err,". Displaying badpage.")
    displayBadPage(w,r,errors.New("Entered SectionSession invalid."))
    return
  }

  enteredCourseId := template.HTMLEscapeString(r.Form.Get("sectioncourseid"))

  _,err = course.GetCourse(enteredCourseId)

  if err != nil{
    log.Println("Error validating entered form:",err,". Displaying badpage.")
    displayBadPage(w,r,errors.New("Entered CourseId invalid."))
    return
  }

  myPassword := password.GenerateRandomPassword()

  enteredStudents,err := strconv.Atoi(template.HTMLEscapeString(r.Form.Get("sectionstudents")))

  if err != nil{
    log.Println("Error validating entered form:",err,". Displaying badpage.")
    displayBadPage(w,r,errors.New("Entered SectionStudents invalid."))
    return
  }

  addedOn := time.Now().Format("2006-01-02 15:04:05")
  myNewSection := &Section{Sectionid : enteredSectionId, Sectionname : enteredSectionName, Year : enteredSectionYear, Session : enteredSectionSession, Courseid : enteredCourseId, Students : enteredStudents, Password : myPassword, Addedon : addedOn}
  log.Println("Adding new Section with details:",*myNewSection)
  err = AddSection(myNewSection)
  if err != nil{
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully added new Section with Id:",enteredSectionId,", name:",enteredSectionName)
  http.Redirect(w, r, "/home", http.StatusSeeOther)
}


func RemoveSectionHandler(w http.ResponseWriter,r *http.Request){
  log.Println("***REMOVE SECTION HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating remove Section request.")

  if r.Method != "GET"{
    displayBadPage(w,r,errors.New("Please remove section properly"))
    return
  }

  mySectionId := template.HTMLEscapeString(r.URL.Path[len("/removesection/"):])
  _,err = RemoveSection(mySectionId)
  if err != nil{
    log.Println("Problem removing Section:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Success removing Section:",mySectionId)
  http.Redirect(w, r, "/home", http.StatusSeeOther)
}


func AddSection(mysection *Section) error{
  return database.SectionCollection.Insert(mysection)
}

func RemoveSection(id string) (int,error){
   changeinfo,err := database.SectionCollection.RemoveAll(bson.M{"sectionid":id})
  if changeinfo == nil{
    return 0,err
  }
  return changeinfo.Removed,err
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***ADD TEACHER HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating form and Parsing form for new teacher data.")

  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please add new teacher properly"))
    return
  }
  r.ParseForm()

  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err,". Displaying badpage.")
    displayBadPage(w,r,err)
    return
  }

  enteredSectionId := template.HTMLEscapeString(r.Form.Get("sectionid"))

  _,err = GetSection(enteredSectionId)

  if err != nil{
    log.Println("Bad Sectionid:",err)
    displayBadPage(w,r,errors.New("Bad Section ID"))
    return
  }

  enteredFacultyId := template.HTMLEscapeString(r.Form.Get("facultyid"))

  _,err = faculty.GetFaculty(enteredFacultyId)

  if err != nil{
    log.Println("Bad Facultyid:",err)
    displayBadPage(w,r,errors.New("Bad Faculty ID"))
    return
  }

  enteredSubjectId := template.HTMLEscapeString(r.Form.Get("subjectid"))

  _,err = subject.GetSubject(enteredSubjectId)

  if err != nil{
    log.Println("Bad Subjectid:",err)
    displayBadPage(w,r,errors.New("Bad Subject ID"))
    return
  }

  myNewTeacher := &Teacher{Facultyid : enteredFacultyId, Subjectid : enteredSubjectId}
  log.Println("Adding new teacher to sectionid:",enteredSectionId," with details:",*myNewTeacher)
  err = AddTeacher(enteredSectionId,myNewTeacher)
  if err != nil{
    log.Println("Error adding teacher to sectionid:",enteredSectionId,":",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully added new teacher to sectionid:",enteredSectionId," with FacultyId:",enteredFacultyId,", SubjectId:",enteredSubjectId)
  http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func AddTeacher(sectionId string, myTeacher *Teacher) error{
  _,err := GetSection(sectionId)
  if err != nil{
    return err
  }
  return database.SectionCollection.Update(bson.M{"sectionid" : sectionId},bson.M{"$push":bson.M{"teachers" : myTeacher}})
}

func RemoveTeacherHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***REMOVE TEACHER HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic. Validating form and Parsing form for remove teacher data.")

  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please remove teacher properly"))
    return
  }
  r.ParseForm()

  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err,". Displaying badpage.")
    displayBadPage(w,r,err)
    return
  }

  enteredSectionId := template.HTMLEscapeString(r.Form.Get("sectionid"))

  _,err = GetSection(enteredSectionId)

  if err != nil{
    log.Println("Bad Sectionid:",err)
    displayBadPage(w,r,errors.New("Bad Section ID"))
    return
  }

  enteredFacultyId := template.HTMLEscapeString(r.Form.Get("facultyid"))

  _,err = faculty.GetFaculty(enteredFacultyId)

  if err != nil{
    log.Println("Bad Facultyid:",err)
    displayBadPage(w,r,errors.New("Bad Faculty ID"))
    return
  }

  enteredSubjectId := template.HTMLEscapeString(r.Form.Get("subjectid"))

  _,err = subject.GetSubject(enteredSubjectId)

  if err != nil{
    log.Println("Bad Subjectid:",err)
    displayBadPage(w,r,errors.New("Bad Subject ID"))
    return
  }

  myTeacher := &Teacher{Facultyid : enteredFacultyId, Subjectid : enteredSubjectId}
  log.Println("Removing teacher from sectionid:",enteredSectionId," with details:",*myTeacher)
  err = RemoveTeacher(enteredSectionId,myTeacher)
  if err != nil{
    log.Println("Error removing teacher from sectionid:",enteredSectionId,":",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Successfully removed teacher from sectionid:",enteredSectionId," with FacultyId:",enteredFacultyId,", SubjectId:",enteredSubjectId)
  http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func RemoveTeacher(sectionId string, myTeacher *Teacher) error{
  _,err := GetSection(sectionId)
  if err != nil{
    return err
  }
  err = database.SectionCollection.Update(bson.M{"sectionid" : sectionId},bson.M{"$pull":bson.M{"teachers" : myTeacher}})
  if err != nil{
    return err
  }
  return err
  //check if teacher array of section with sectionid = sectionId is empty. if yes, $unset teacher array.
}

func GetDetailedTeachers(myTeachers *[]Teacher) (*[]DetailedTeacher,error){
  log.Println("*Getting Detailed Teachers*")
  if myTeachers == nil || len(*myTeachers) == 0{
    log.Println("myTeachers slice passed to GetDetailedTeachers is nil. Returning nil DetailedTeacher slice.")
    return nil,errors.New("No Teachers")
  }
  myDetailedTeachers := make([]DetailedTeacher,len(*myTeachers))
  for i,myTeacher := range *myTeachers{
    myDetailedTeachers[i].Facultyid = myTeacher.Facultyid
    myDetailedTeachers[i].Subjectid = myTeacher.Subjectid
    myFaculty,err := faculty.GetFaculty(myTeacher.Facultyid)
    if err != nil{
      return nil,errors.New("Could not find faculty with facultyid:" + myTeacher.Facultyid)
    }
    myDetailedTeachers[i].Facultyname = myFaculty.Facultyname

    mySubject,err := subject.GetSubject(myTeacher.Subjectid)
    if err != nil{
      return nil,errors.New("Could not find subject with subjectid:" + mySubject.Subjectid)
    }
    myDetailedTeachers[i].Subjectname = mySubject.Subjectname
  }
  log.Println("Success getting DetailedTeachers from Teachers. Returning slice of DetailedTeachers")
  return &myDetailedTeachers,nil
}
