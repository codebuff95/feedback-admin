package logout

import(
  "feedback-admin/user"
  "github.com/codebuff95/uafm/usersession"
  "net/http"
  "log"
  //"html/template"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request){
  log.Println("***LOGOUT HANDLER***")
  log.Println("Serving client:",r.RemoteAddr)
  userRid, err := user.AuthenticateRequest(r)
  if err != nil{ //user session not authentic. No need to logout. Redirect to login page.
    log.Println("User session is not authentic, redirecting to login page.")
    http.Redirect(w, r, "/login", http.StatusSeeOther)
    return
  }
  log.Println("User Session is authentic. Deleting user session from the database.")
  userSidCookie,err := r.Cookie("usersid")
  userSid := userSidCookie.Value
  sessionsDeleted, err := usersession.DeleteSession(userSid)
  if err != nil{
    log.Println("Error in removing user session from database:",err)
  }else{
    log.Println("Deleted",sessionsDeleted,"usersessions from database for user:",userRid)
  }
  //Delete Cookie set at client.
  //...
  //No need.
  http.Redirect(w, r, "/login", http.StatusSeeOther)
  return
}
