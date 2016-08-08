package report

import(
  //"github.com/codebuff95/uafm"
  "feedback-admin/user"
  "feedback-admin/question"
  "feedback-admin/section"
  "feedback-admin/templates"
  "feedback-admin/feedback"
  "feedback-admin/college"
  "strings"
  "os"
  "html/template"
  "net/http"
  "log"
  "errors"
  //"gopkg.in/mgo.v2"
  //"html/template"
)

type SectionWiseReportPage struct{
  Collegename *string
  Sectionname *string
  Sectionid *string
  Customfeedbacks *[]SectionWiseCustomFeedback
  Detailedteachers *[]section.DetailedTeacher
}

type SectionWiseCustomFeedback struct{
  Index int
  Total []int
}

type SubjectWiseReportPage struct{
  Collegename *string
  Sectionname *string
  Sectionid *string
  Customfeedbacks *[]SubjectWiseCustomFeedback
  Outof *int
}

type SubjectWiseCustomFeedback struct{
  Facultyname string
  Facultyid string
  Subjectid string
  Totalfeedbacks int
  Totalscore int
  Average float64
  Percentage float64
}

type PointWiseReportPage struct{
  Collegename *string
  Sectionname *string
  Sectionid *string
  PointWiseTeachers *[]PointWiseTeacher
  Students *[]int
}

type PointWiseTeacher struct{
  Questions []PointWiseQuestion
  Facultyname string
  Facultyid string
}

type PointWiseQuestion struct{
  Text string
  Ratings []int
}

func displayReportOptionsPage(w http.ResponseWriter, r *http.Request){
  templates.ReportOptionsPageTemplate.Execute(w,nil)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func ReportHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***REPORT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  displayReportOptionsPage(w,r)
}

func SectionWiseReportHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***SECTION WISE REPORT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please generate report properly"))
    return
  }
  r.ParseForm()
  var myReportPage SectionWiseReportPage
  mySection,err := section.GetSection(template.HTMLEscapeString(r.Form.Get("sectionid")))
  if err != nil{
    displayBadPage(w,r,errors.New("Error finding given Section ID"))
    return
  }
  myReportPage.Collegename = &college.GlobalDetails.Collegename
  myReportPage.Sectionname = &mySection.Sectionname
  myReportPage.Sectionid = &mySection.Sectionid
  myReportPage.Detailedteachers,err = section.GetDetailedTeachers(mySection.Teachers)
  log.Println("MyReportPage.Sectionname:",*myReportPage.Sectionname)

  if err != nil{
    displayBadPage(w,r,errors.New("Error finding detailed teachers of section."))
    return
  }

  myReportPage.Customfeedbacks,err = getSectionWiseCustomFeedbacks(mySection.Sectionid,mySection.Teachers)
  if err != nil || myReportPage.Customfeedbacks == nil{
    log.Println("Problem getting custom feedbacks.")
    displayBadPage(w,r,errors.New("Error finding feedbacks for this section"))
    return
  }
  log.Println("Successfully got Customfeedbacks:",*myReportPage.Customfeedbacks,". Executing Template.")
  filewriter, err := os.Create("sectionwise-"+mySection.Sectionid+".csv")
  defer filewriter.Close()
  if err != nil{
    log.Println("Problem opening file for writing")
    displayBadPage(w,r,errors.New("Problem opening file for writing"))
    return
  }
  templates.SectionWiseReportTemplate.Execute(filewriter,myReportPage)
  http.Redirect(w, r, "/report", http.StatusSeeOther)
  return
}

func getSectionWiseCustomFeedbacks(sectionId string, myTeachers *[]section.Teacher) (*[]SectionWiseCustomFeedback,error){
  log.Println("**Getting Section Wise Custom Feedbacks**")
  myFeedbackSlice, err := feedback.GetFeedbacks(sectionId)
  if err != nil || myFeedbackSlice == nil{
    log.Println("Error finding Feedbacks:",err)
    return nil,err
  }
  log.Println("Success getting feedbackSlice of size:",len(*myFeedbackSlice))

  var mySectionWiseCustomFeedbackSlice []SectionWiseCustomFeedback = make([]SectionWiseCustomFeedback,len(*myFeedbackSlice))
  for i,myFeedback := range *myFeedbackSlice{
    for _,myTeacher := range *myTeachers{
      //var r int
      var found bool = false
      var myRating feedback.Rating
      for _,myRating = range myFeedback.Ratings{
        if strings.Compare(myRating.Facultyid,myTeacher.Facultyid) == 0 && strings.Compare(myRating.Subjectid,myTeacher.Subjectid) == 0{
          found = true
          break
        }
      }
      var sum int = 0
      if found == false{
        log.Println("Problem Finding feedback record for faculty ID:",myTeacher.Facultyid,"in a feedback that was added on:",myFeedback.Addedon)
        return nil,errors.New("Could not find feedback record for teacher with facultyid: " + myTeacher.Facultyid +" in feedback which was addedon: " + myFeedback.Addedon)
      }
      for _,points := range myRating.Points{
        sum += points.Marks
      }
      mySectionWiseCustomFeedbackSlice[i].Index = i+1
      mySectionWiseCustomFeedbackSlice[i].Total = append(mySectionWiseCustomFeedbackSlice[i].Total,sum)
    }
  }
  log.Println("Success getting Section Wise Custom Feedback Slice")
  return &mySectionWiseCustomFeedbackSlice,nil
}

func SubjectWiseReportHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***SUBJECT WISE REPORT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please generate report properly"))
    return
  }
  r.ParseForm()
  var myReportPage SubjectWiseReportPage
  mySection,err := section.GetSection(template.HTMLEscapeString(r.Form.Get("sectionid")))
  if err != nil{
    displayBadPage(w,r,errors.New("Error finding given Section ID"))
    return
  }
  myReportPage.Collegename = &college.GlobalDetails.Collegename
  myReportPage.Sectionname = &mySection.Sectionname
  myReportPage.Sectionid = &mySection.Sectionid
  myDetailedTeachers,err := section.GetDetailedTeachers(mySection.Teachers)

  if err != nil{
    log.Println("Error finding detailed teachers of section:",mySection.Sectionid)
    displayBadPage(w,r,errors.New("Error finding teachers of section."))
    return
  }
  log.Println("Success finding detailed teachers of section:",mySection.Sectionid)

  myOutof,err := question.GetAllQuestionsWeightage()

  if err != nil{
    log.Println("Error finding all questions weightage")
    displayBadPage(w,r,errors.New("Error finding all questions' weightage"))
  }

  myReportPage.Outof = &myOutof

  myReportPage.Customfeedbacks,err = getSubjectWiseCustomFeedbacks(mySection.Sectionid,myDetailedTeachers, myOutof)
  if err != nil || myReportPage.Customfeedbacks == nil{
    log.Println("Problem getting custom feedbacks.")
    displayBadPage(w,r,errors.New("Error finding feedbacks for this section"))
    return
  }
  log.Println("Successfully got Customfeedbacks:",*myReportPage.Customfeedbacks,". Executing Template.")
  filewriter, err := os.Create("subjectwise-"+mySection.Sectionid+".csv")
  defer filewriter.Close()
  if err != nil{
    log.Println("Problem opening file for writing")
    displayBadPage(w,r,errors.New("Problem opening file for creating report"))
    return
  }
  templates.SubjectWiseReportTemplate.Execute(filewriter,myReportPage)
  http.Redirect(w, r, "/report", http.StatusSeeOther)
  return
}

func getSubjectWiseCustomFeedbacks(sectionId string, myDetailedTeachers *[]section.DetailedTeacher, Outof int) (*[]SubjectWiseCustomFeedback,error){
  log.Println("**Getting Subject Wise Custom Feedbacks**")
  myFeedbackSlice, err := feedback.GetFeedbacks(sectionId)
  if err != nil || myFeedbackSlice == nil{
    log.Println("Error finding Feedbacks:",err)
    return nil,err
  }
  log.Println("Success getting feedbackSlice of size:",len(*myFeedbackSlice))

  var myMap map[string]*SubjectWiseCustomFeedback = make(map[string]*SubjectWiseCustomFeedback)
  var mySubjectWiseCustomFeedbackSlice []SubjectWiseCustomFeedback = make([]SubjectWiseCustomFeedback,len(*myDetailedTeachers))
  for i,myDetailedTeacher := range *myDetailedTeachers{
    mySubjectWiseCustomFeedbackSlice[i].Facultyname = myDetailedTeacher.Facultyname
    mySubjectWiseCustomFeedbackSlice[i].Facultyid = myDetailedTeacher.Facultyid
    mySubjectWiseCustomFeedbackSlice[i].Subjectid = myDetailedTeacher.Subjectid
    myMap[myDetailedTeacher.Facultyid+"~"+myDetailedTeacher.Subjectid] = &mySubjectWiseCustomFeedbackSlice[i]
  }
  for _,myFeedback := range *myFeedbackSlice{
    for _,myRating := range myFeedback.Ratings{
      totalPoints := 0
      for _,myPoint := range myRating.Points{
        totalPoints += myPoint.Marks
      }
      myMap[myRating.Facultyid+"~"+myRating.Subjectid].Totalscore += totalPoints
      myMap[myRating.Facultyid+"~"+myRating.Subjectid].Totalfeedbacks++
    }
  }

  for i,_ := range mySubjectWiseCustomFeedbackSlice{
    if mySubjectWiseCustomFeedbackSlice[i].Totalfeedbacks != 0{
      mySubjectWiseCustomFeedbackSlice[i].Average = float64(mySubjectWiseCustomFeedbackSlice[i].Totalscore) / float64(mySubjectWiseCustomFeedbackSlice[i].Totalfeedbacks)
      mySubjectWiseCustomFeedbackSlice[i].Percentage = (mySubjectWiseCustomFeedbackSlice[i].Average / float64(Outof)) * 100
    }
  }

  log.Println("Success getting Subject Wise Custom Feedback Slice")
  return &mySubjectWiseCustomFeedbackSlice,nil
}

func PointWiseReportHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***POINT WISE REPORT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage.")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic")
  if r.Method != "POST"{
    displayBadPage(w,r,errors.New("Please generate report properly"))
    return
  }
  r.ParseForm()
  var myReportPage PointWiseReportPage
  mySection,err := section.GetSection(template.HTMLEscapeString(r.Form.Get("sectionid")))
  if err != nil{
    displayBadPage(w,r,errors.New("Error finding given Section ID"))
    return
  }
  myReportPage.Collegename = &college.GlobalDetails.Collegename
  myReportPage.Sectionname = &mySection.Sectionname
  myReportPage.Sectionid = &mySection.Sectionid

  myStudents := make([]int,mySection.Students)

  for i := range myStudents{
    myStudents[i] = i+1
  }

  myReportPage.Students = &myStudents

  myReportPage.PointWiseTeachers,err = getPointWiseTeachers(mySection.Sectionid)
  if err != nil || myReportPage.PointWiseTeachers == nil{
    log.Println("Error getting Pointwise Teachers:",err)
    displayBadPage(w,r,errors.New("Error generating pointwise teachers for this section: "+err.Error()))
    return
  }
  log.Println("Successfully got PointWiseTeachers:",*myReportPage.PointWiseTeachers,". Executing Template.")
  filewriter, err := os.Create("pointwise-"+mySection.Sectionid+".csv")
  defer filewriter.Close()
  if err != nil{
    log.Println("Problem opening file for writing")
    displayBadPage(w,r,errors.New("Problem opening file for creating report"))
    return
  }
  templates.PointWiseReportTemplate.Execute(filewriter,myReportPage)
  http.Redirect(w, r, "/report", http.StatusSeeOther)
  return
}

func getPointWiseTeachers(sectionId string) (*[]PointWiseTeacher,error){
  log.Println("**Getting Point Wise Teachers**")
  myFeedbacks, err := feedback.GetFeedbacks(sectionId)
  if err != nil || myFeedbacks == nil{
    log.Println("Error finding Feedbacks:",err)
    return nil,err
  }
  log.Println("Success getting feedbacks of size:",len(*myFeedbacks))

  mySection,_ := section.GetSection(sectionId)
  myDetailedTeachers, err := section.GetDetailedTeachers(mySection.Teachers)
  if myDetailedTeachers == nil || err != nil {
    log.Println("No teachers found for this section.")
    return nil,errors.New("No teachers found for this section")
  }
  myPointWiseTeachers := make([]PointWiseTeacher,len(*myDetailedTeachers))
  myQuestions,err := question.GetQuestions()
  if err != nil || myQuestions == nil{
    log.Println("No questions found")
    return nil,errors.New("No questions found in database")
  }
  var myMap map[string]*PointWiseQuestion = make(map[string]*PointWiseQuestion)
  for i,myDetailedTeacher := range *myDetailedTeachers {
    myPointWiseTeachers[i].Facultyname = myDetailedTeacher.Facultyname
    myPointWiseTeachers[i].Facultyid = myDetailedTeacher.Facultyid
    myPointWiseTeachers[i].Questions = make([]PointWiseQuestion,len(*myQuestions))
    for q,myQuestion := range *myQuestions{
      myPointWiseTeachers[i].Questions[q].Text = myQuestion.Text
      myMap[myQuestion.Questionid] = &myPointWiseTeachers[i].Questions[q]
    }
    for _,myFeedback := range *myFeedbacks{
      var myRating *feedback.Rating = nil
      for r,myFeedbackRating := range myFeedback.Ratings{
        if myFeedbackRating.Facultyid == myPointWiseTeachers[i].Facultyid{
          myRating = &myFeedback.Ratings[r]
        }
      }
      if myRating == nil{
        log.Println("A feedback submitted on:",myFeedback.Addedon,"does not have a rating for faculty ID:",myPointWiseTeachers[i].Facultyid)
        return nil, errors.New("A feedback submitted on: "+myFeedback.Addedon+" does not have a rating for faculty ID: "+myPointWiseTeachers[i].Facultyid)
      }
      for _,myPoint := range myRating.Points{
        myMap[myPoint.Questionid].Ratings = append(myMap[myPoint.Questionid].Ratings,myPoint.Marks)
      }
    }
  }
  return &myPointWiseTeachers,nil
}
