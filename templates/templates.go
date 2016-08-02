package templates

import(
  "html/template"
  "log"
)

var LoginFormTemplate *template.Template
var HomePageTemplate *template.Template
var BadPageTemplate *template.Template
var CoursePageTemplate *template.Template
var SubjectPageTemplate *template.Template
var FacultyPageTemplate *template.Template
var SectionPageTemplate *template.Template

func InitEssentialTemplates() error{
  var err error
  LoginFormTemplate,err = template.ParseFiles("feedbackadminres/login.html")
  if err != nil{
    log.Println("Error parsing LoginFormTemplate:",err)
    return err
  }
  BadPageTemplate,err = template.ParseFiles("feedbackadminres/badpage.html")
  if err != nil{
    log.Println("Error parsing BadPageTemplate:",err)
    return err
  }
  HomePageTemplate,err = template.ParseFiles("feedbackadminres/home.html")
  if err != nil{
    log.Println("Error parsing BadPageTemplate:",err)
    return err
  }
  CoursePageTemplate,err = template.ParseFiles("feedbackadminres/course.html")
  if err != nil{
    log.Println("Error parsing CoursePageTemplate:",err)
    return err
  }
  SubjectPageTemplate,err = template.ParseFiles("feedbackadminres/subject.html")
  if err != nil{
    log.Println("Error parsing SubjectPageTemplate:",err)
    return err
  }
  FacultyPageTemplate,err = template.ParseFiles("feedbackadminres/faculty.html")
  if err != nil{
    log.Println("Error parsing FacultyPageTemplate:",err)
    return err
  }
  SectionPageTemplate,err = template.ParseFiles("feedbackadminres/section.html")
  if err != nil{
    log.Println("Error parsing SectionPageTemplate:",err)
    return err
  }
  return nil
}
