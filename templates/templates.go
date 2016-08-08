package templates

import(
  "html/template"
  texttemplate "text/template"
  "log"
)

var LoginFormTemplate *template.Template
var HomePageTemplate *template.Template
var BadPageTemplate *template.Template
var CoursePageTemplate *template.Template
var SubjectPageTemplate *template.Template
var FacultyPageTemplate *template.Template
var SectionPageTemplate *template.Template
var ReportOptionsPageTemplate *template.Template
var SectionWiseReportTemplate *texttemplate.Template
var SubjectWiseReportTemplate *texttemplate.Template
var PointWiseReportTemplate *texttemplate.Template
var PasswordsTemplate *texttemplate.Template

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
  ReportOptionsPageTemplate,err = template.ParseFiles("feedbackadminres/reportoptions.html")
  if err != nil{
    log.Println("Error parsing ReportOptionsPageTemplate:",err)
    return err
  }
  SectionWiseReportTemplate,err = texttemplate.ParseFiles("feedbackadminres/reporttemplates/sectionwise")
  if err != nil{
    log.Println("Error parsing SectionWiseReportTemplate:",err)
    return err
  }
  SubjectWiseReportTemplate,err = texttemplate.ParseFiles("feedbackadminres/reporttemplates/subjectwise")
  if err != nil{
    log.Println("Error parsing SubjectWiseReportTemplate:",err)
    return err
  }
  PointWiseReportTemplate,err = texttemplate.ParseFiles("feedbackadminres/reporttemplates/pointwise")
  if err != nil{
    log.Println("Error parsing PointWiseReportTemplate:",err)
    return err
  }
  PasswordsTemplate,err = texttemplate.ParseFiles("feedbackadminres/reporttemplates/passwords")
  if err != nil{
    log.Println("Error parsing PasswordsTemplate:",err)
    return err
  }

  return nil
}
