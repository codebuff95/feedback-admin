package user

import(
  "net/http"
  "github.com/codebuff95/uafm/usersession"
  "feedback-admin/database"
  "log"
  "html/template"
  //"gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Admin struct{
  Username *string `bson:"username" json:"username"`
  Password string `bson:"password" json:"password"`
}

func AuthenticateRequest(r *http.Request) (*string,error){
  log.Println("Authenticating incoming request from client:",r.RemoteAddr,"for authentic user")
  userSidCookie,err := r.Cookie("usersid")
  if err != nil{
    log.Println("User Authentication failed")
    return nil,err
  }
  userSid := userSidCookie.Value
  log.Println("Got Cookie Value")
  log.Println("Authenticating UserSession with Sid:\"",userSid,"\"")
  userRid, err := usersession.ValidateSession(userSid)
  if err != nil{
    log.Println("Error authenticating userSid:",err)
  }
  return userRid,err
}

func AuthenticateLoginAttempt(r *http.Request) (*Admin,error){
  log.Println("*Authenticating Login Attempt*")
  attemptUsername := template.HTMLEscapeString(r.Form.Get("username"))       //Escape special characters for security.
  attemptPassword := template.HTMLEscapeString(r.Form.Get("password"))       //Escape special characters for security.
  log.Println("Attempt username :", attemptUsername, ", Attempt Password:", attemptPassword)
  var myAdmin Admin = Admin{}
  err := database.AdminCollection.
            Find(bson.M{"username":attemptUsername, "password":attemptPassword}).
                Limit(1).One(&myAdmin)
  if err != nil{
    log.Println("Error finding requested admin in collection:",err)
    return nil,err
  }
  if myAdmin.Username == nil{
    log.Println("Could not find admin with supplied credentials. Returning nil Admin")
  }
  return &myAdmin,err
}
