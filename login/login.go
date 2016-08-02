package login

import(
  //"github.com/codebuff95/uafm"
  "github.com/codebuff95/uafm/usersession"
  "github.com/codebuff95/uafm/formsession"
  "feedback-admin/user"
  "feedback-admin/templates"
  "net/http"
  "time"
  "log"
  "errors"
  //"html/template"
)

func displayLoginForm(w http.ResponseWriter, r *http.Request){
  log.Println("Displaying login form to user.")
  formSid, err := formsession.CreateSession("0",time.Minute*10) //Form created will be valid for 10 minutes.
  if err != nil{
    log.Println("Error creating new session for login form:",err)
    displayBadPage(w,r,err)
    return
  }
  log.Println("Creating new login form to client",r.RemoteAddr,"with fomSid:",*formSid)  //Enter client ip address and new form SID.
  templates.LoginFormTemplate.Execute(w,formSid)
}

func displayBadPage(w http.ResponseWriter, r *http.Request, err error){
  templates.BadPageTemplate.Execute(w,err.Error())
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***LOGIN HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  userId, err := user.AuthenticateRequest(r)
  if err == nil{ //user session authentic. Redirect to home page.
    log.Println("User session is authentic with UserId:",*userId,", redirecting to homepage.")
    http.Redirect(w, r, "/home", http.StatusSeeOther)
    return
  }
  log.Println("User Session not authentic:",err)
  if r.Method == "GET"{
    log.Println("Request method is GET")
    displayLoginForm(w,r)
    return
  }
  //method = POST
  log.Println("Request method is POST")
  r.ParseForm()

  //Validate entered form session.
  _,err = formsession.ValidateSession(r.Form.Get("formsid"))
  if err != nil{
    log.Println("Error validating formSid:",err)
    displayBadPage(w,r,err)
    return
  }

  //Form is valid. Continue with deleting the form session.
  formSidDeleted,err := formsession.DeleteSession(r.Form.Get("formsid"))
  log.Println("Deleted",formSidDeleted,"form sessions for sid",r.Form.Get("formsid"),".")

  //Login Attempt Authentication begin.
  attemptAdmin,err := user.AuthenticateLoginAttempt(r)
  if err != nil || attemptAdmin == nil{
    displayBadPage(w,r,errors.New("User does not exist"))
    return
  }

  //Authentic Login Attempt. Set cookie on client, and redirect to homepage.
  newUserSid,err := usersession.CreateSession(*attemptAdmin.Username,time.Hour*24*3)
  if err != nil{
    log.Println("Error creating new UserSid:",err)
    displayBadPage(w,r,err)
    return
  }
  myCookie := &http.Cookie{Name:"usersid",Value:*newUserSid}  //Cookie is not persistent for admin rights security.
  http.SetCookie(w, myCookie)
  log.Println("usersid Cookie successfully set on client. Redirecting to homepage.")
  http.Redirect(w, r, "/home", http.StatusSeeOther)
  return
}
