package app

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var memoryCreds = make(map[string]string)

func init() {
	memoryCreds["gleb"] = "password"
}