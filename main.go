package main

import(
  "feedback-admin/course"
  "feedback-admin/login"
  "feedback-admin/logout"
  "feedback-admin/database"
  "feedback-admin/subject"
  "feedback-admin/faculty"
  "feedback-admin/section"
  "feedback-admin/templates"
  "feedback-admin/home"
  "github.com/codebuff95/uafm"
  "log"
  "net/http"
)

func handlefatalerror(err error){
  if err != nil{
    log.Fatal("*_*_* Fatal Error:",err,"*_*_*")
  }
}

func main(){

  err := uafm.Init("feedbackadminres")
  handlefatalerror(err)

  err = database.InitDatabaseSession()
  handlefatalerror(err)

  database.InitCollections()

  err = templates.InitEssentialTemplates()
  handlefatalerror(err)

  http.HandleFunc("/login",login.LoginHandler)
  http.HandleFunc("/home",home.HomeHandler)
  http.HandleFunc("/logout",logout.LogoutHandler)

  http.HandleFunc("/course",course.CourseHandler)
  http.HandleFunc("/addcourse",course.AddCourseHandler)
  http.HandleFunc("/removecourse/",course.RemoveCourseHandler)

  http.HandleFunc("/subject",subject.SubjectHandler)
  http.HandleFunc("/addsubject",subject.AddSubjectHandler)
  http.HandleFunc("/removesubject/",subject.RemoveSubjectHandler)

  http.HandleFunc("/faculty",faculty.FacultyHandler)
  http.HandleFunc("/addfaculty",faculty.AddFacultyHandler)
  http.HandleFunc("/removefaculty/",faculty.RemoveFacultyHandler)

  http.HandleFunc("/section",section.SectionHandler)
  http.HandleFunc("/addsection",section.AddSectionHandler)
  http.HandleFunc("/removesection/",section.RemoveSectionHandler)
  http.HandleFunc("/addteacher",section.AddTeacherHandler)
  http.HandleFunc("/removeteacher",section.RemoveTeacherHandler)

  //Serving static files: only files in directory feedbackadminres/publicres are being made public (for security purposes)
  http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("feedbackadminres/publicres"))))

  http.ListenAndServe(":8080",nil)
}
