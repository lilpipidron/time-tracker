package requests

type ChangeUserInfo struct {
	Name           string `json:"name,omitempty"`
	Surname        string `json:"surname,omitempty"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address,omitempty"`
	PassportNumber string `json:"passport_number"`
}
