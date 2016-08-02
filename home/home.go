package home

import(
  "feedback-admin/user"
  "feedback-admin/templates"
  "net/http"
  "log"
  //"html/template"
)

func displayHomePage(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying home page to user.")
  templates.HomePageTemplate.Execute(w,nil) //any admin info?
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***HOME HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  _, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. Redirect to login page.
    log.Println("User session is not authentic, displaying badpage")
    //http.Redirect(w, r, "/login", http.StatusSeeOther)
    displayBadPage(w,r,err)
    return
  }
  log.Println("User Session is authentic.")
  displayHomePage(w,r)
  return
}
