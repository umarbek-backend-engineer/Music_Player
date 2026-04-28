package modules

// this struct will be used to request and response
type Register struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password,omitempty"`
}

type RestPassowrd struct {
	Id               string `json:"id"`
	Current_Password string `json:"current_password"`
	New_Password     string `json:"new_password"`
}


