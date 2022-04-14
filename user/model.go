package user

import "time"

/*<<<<---User--->>>>
Contains fields which discuss about a employee or a person in a organization. ///Checking
*/
type User struct {
	Id                   int
	Name                 string
	Department           string
	SocailSecurityNumber int
	DOB                  time.Time
	City                 string
	State                string
	Address              string
	JobTitle             string
	IsPermanent          bool
	Gender               string
	PAN                  string
	Passport             string
	ReportingManager     int
}
