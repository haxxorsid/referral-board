package app

import (
	"bytes"
	"os"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/haxxorsid/referralboard-private/server/config"
	"github.com/haxxorsid/referralboard-private/server/models"
)

var a = &App{}

func TestMain(m *testing.M) {
	config := config.GetConfig()
	a.Initialize(config)
	a.SetUpDB()
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestAddUser(t *testing.T) {

	testUsers := &models.User{
		FirstName:           "Mallikaaa",
		LastName:            "Kumar",
		CurrentLocation:     "Mumbai",
		CurrentCompanyName:  "Company ABCD",
		CurrentPosition:     "Software Engineer",
		School:              "University of Mumbai",
		YearsOfExperienceId: 1,
		Email:               "mailaddress4@asd.com",
		Password:            "root",
	}
	userformValue, _ := json.Marshal(testUsers)

	req, err := http.NewRequest("POST", "/api/users/newuser", bytes.NewBuffer(userformValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["lastName"] != "Kumar" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

// invalid email
func TestInvalidEmail(t *testing.T) {
	testUsers := &models.User{
		FirstName:           "Mallikaaa",
		LastName:            "Kumar",
		CurrentLocation:     "Mumbai",
		CurrentCompanyName:  "Company ABCD",
		CurrentPosition:     "Software Engineer",
		School:              "University of Mumbai",
		YearsOfExperienceId: 1,
		Email:               "mailaddress4",
		Password:            "root",
	}
	userformValue, _ := json.Marshal(testUsers)

	req, err := http.NewRequest("POST", "/api/users/newuser", bytes.NewBuffer(userformValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// email is comapanyA -- success if companyId == 1
func TestUserRegistrationWithExistingCompany(t *testing.T) {
	testUsers := &models.User{
		FirstName:           "Mallikaaa",
		LastName:            "Kumar",
		CurrentLocation:     "Mumbai",
		CurrentCompanyName:  "Company A",
		CurrentPosition:     "Software Engineer",
		School:              "University of Mumbai",
		YearsOfExperienceId: 1,
		Email:               "mailaddress2a@companya.com",
		Password:            "root",
	}
	userformValue, _ := json.Marshal(testUsers)

	req, err := http.NewRequest("POST", "/api/users/newuser", bytes.NewBuffer(userformValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	// fmt.Println("type: ", reflect.ValueOf(m["currentCompanyId"]).Kind())
	if m["currentCompanyId"] != float64(1) {
		t.Errorf("Expected current company ID expected to be <nil>. Got '%v'", m["currentCompanyId"])
	}
}

// email is comapanyRandom -- success if companyId == null
func TestUserRegistrationWithRandomCompany(t *testing.T) {
	testUsers := &models.User{
		FirstName:           "Mallikaaa",
		LastName:            "Kumar",
		CurrentLocation:     "Mumbai",
		CurrentCompanyName:  "Company ABCD",
		CurrentPosition:     "Software Engineer",
		School:              "University of Mumbai",
		YearsOfExperienceId: 1,
		Email:               "mailaddress2a@asd.com",
		Password:            "root",
	}
	userformValue, _ := json.Marshal(testUsers)

	req, err := http.NewRequest("POST", "/api/users/newuser", bytes.NewBuffer(userformValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["currentCompanyId"] != float64(0) {
		t.Errorf("Expected current comapny ID to be 0. Got '%v'", m["currentCompanyId"])
	}
}

// not unique email format
func TestNotUniqueEmail(t *testing.T) {
	testUsers := &models.User{
		FirstName:           "Mallikaaa",
		LastName:            "Kumar",
		CurrentLocation:     "Mumbai",
		CurrentCompanyName:  "Company A",
		CurrentPosition:     "Software Engineer",
		School:              "University of Mumbai",
		YearsOfExperienceId: 1,
		Email:               "mailaddress2@companya.com",
		Password:            "root",
	}
	userformValue, _ := json.Marshal(testUsers)

	req, err := http.NewRequest("POST", "/api/users/newuser", bytes.NewBuffer(userformValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// check company name is in sync with company table -- success if companyName != provided company anem
func TestUserRegistrationWithExistingCompanyName(t *testing.T) {
	testUsers := &models.User{
		FirstName:           "Mallikaaa",
		LastName:            "Kumar",
		CurrentLocation:     "Mumbai",
		CurrentCompanyName:  "Company A",
		CurrentPosition:     "Software Engineer",
		School:              "University of Mumbai",
		YearsOfExperienceId: 1,
		Email:               "mailaddress2ab@companya.com",
		Password:            "root",
	}
	userformValue, _ := json.Marshal(testUsers)

	req, err := http.NewRequest("POST", "/api/users/newuser", bytes.NewBuffer(userformValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["currentCompanyName"] != "Company A" {
		t.Errorf("Expected current comapny ID to be 0. Got '%v'", m["currentCompanyId"])
	}
}
