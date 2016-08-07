package password

import(
  "crypto/rand"
  "encoding/base64"
  "log"
)

const(
  PASSLEN int = 8
)

func GenerateRandomPassword() string{
  pass := make([]byte,PASSLEN)
  rand.Read(pass)
  finalpass := base64.URLEncoding.EncodeToString(pass)
  return string(finalpass[0:PASSLEN])
}

func GenerateRandomPasswords(noOfPasswords int) []string{
  log.Println("*Generating random passwords*")
  myPasswords := make([]string,noOfPasswords)
  for i := 0; i < noOfPasswords; i++{
    myPasswords[i] = GenerateRandomPassword()
  }
  return myPasswords
}
