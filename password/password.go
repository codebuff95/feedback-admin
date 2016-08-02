package password

import(
  "crypto/rand"
  "encoding/base64"
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
