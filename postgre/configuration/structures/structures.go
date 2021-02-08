package structures

//ContactDetails is ...
type ContactDetails struct {
	ID        string
	Age       string
	Email     string
	FirstName string
	LastName  string
	Image     string
	Img       []byte
}

//Pass is ...
type Pass struct {
	Data []ContactDetails
}
