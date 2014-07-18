package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go"
	"log"
	"net/http"
)

const (
	mgDomain      = "sandbox27534.mailgun.org"
	mgApiKey      = "key-7mfaud20r5uyzxhbonjvci335qhojeb3"
	mgFromUser    = "FCE <fce@sandbox27534.mailgun.org>"
	mgSubject     = "How was class/gym tonight?"
	mgText        = "You will need HTML e-mail support to use this application.\n"
	webhookDomain = "localhost:8081"
)

type Query struct {
	Id string
	Attendance
}

type User struct {
	Id          string
	First, Last string
	Email       string
}

type Attendance struct {
	UserId string
	Gym    string
	When   string
}

var gyms []string
var users []User
var attendance []Attendance
var queries []Query
var nextId int

var ok = []byte(`{"Tag": "OK"}`)

func GetGymsHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(struct {
		Tag  string
		Gyms []string
	}{
		"gyms", gyms,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Tag": "error", "Msg": "Unable to marshal gym data"}`))
		return
	}
	w.Write(j)
}

func PostGymHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	g := vs["gym"]
	gyms = append(gyms, g)
	w.Write(ok)
}

func DelGymHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	gym := vs["gym"]
	// There's totally a better way to do this.
	newGyms := make([]string, 0)
	for _, g := range gyms {
		if g != gym {
			newGyms = append(newGyms, g)
		}
	}
	gyms = newGyms
	w.WriteHeader(http.StatusNoContent)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(struct {
		Tag   string
		Users []User
	}{
		"users", users,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Tag": "error", "Msg": "Unable to marshal user data"}`))
		return
	}
	w.Write(j)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	id := vs["id"]
	fn := vs["first"]
	ln := vs["last"]
	em := vs["email"]
	u := User{
		Id:    id,
		First: fn,
		Last:  ln,
		Email: em,
	}
	// TODO(sfalvo): sanity checking here...
	users = append(users, u)
	w.Write(ok)
}

func DelUserHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	id := vs["id"]
	// There's totally a better way to do this.
	newUsers := make([]User, 0)
	for _, u := range users {
		if u.Id != id {
			newUsers = append(newUsers, u)
		}
	}
	users = newUsers
	w.WriteHeader(http.StatusNoContent)
}

func GetAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(struct {
		Tag        string
		Attendance []Attendance
	}{
		"attendance", attendance,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Tag": "error", "Msg": "Unable to marshal attendance data"}`))
		return
	}
	w.Write(j)
}

func recordAttendance(a Attendance) {
	attendance = append(attendance, a)
}

func PostAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	id := vs["id"]
	gym := vs["gym"]
	when := vs["when"]
	a := Attendance{
		UserId: id,
		Gym:    gym,
		When:   when,
	}
	// TODO(sfalvo): sanity checking here...
	recordAttendance(a)
	w.Write(ok)
}

func UserById(id string) User {
	for _, u := range users {
		if u.Id == id {
			return u
		}
	}
	log.Println("User ID(" + id + ") not found.  Returning empty user.")
	return User{}
}

func firstNameForUser(id string) string {
	return UserById(id).First
}

func emailFor(id string) string {
	return UserById(id).Email
}

func AskUserHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	nextId++

	id := vs["id"]
	gym := vs["gym"]
	when := vs["when"]

	fname := firstNameForUser(id)

	msg := `<html>
<head></head>
<body>Hey, ` + fname + `!  How was your ` + gym + ` class tonight?<br /><br />
<a href="` + fmt.Sprintf("http://%s/webhooks/ack/%d", webhookDomain, nextId) + `">It was awesome!  I'm so happy and so thoroughly satisfied that I went tonight!</a>
</form>
Otherwise, ignore this message, and your absence will be recorded automatically.  Thanks!
</html>
`
	mg := mailgun.NewMailgun(mgDomain, mgApiKey, "")
	m := mg.NewMessage(mgFromUser, mgSubject, mgText, emailFor(id))
	m.SetHtml(msg)
	_, _, err := mg.Send(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to send e-mail: " + err.Error()))
		return
	}

	queries = append(queries, Query{
		fmt.Sprintf("%d", nextId),
		Attendance{
			UserId: id,
			Gym:    gym,
			When:   when,
		},
	})
	w.Write(ok)
}

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	n := vs["n"]
	newQueries := make([]Query, 0)
	for _, q := range queries {
		if q.Id == n {
			recordAttendance(q.Attendance)
		} else {
			newQueries = append(newQueries, q)
		}
	}
	queries = newQueries
	w.Write(ok)
}

func main() {
	gyms = make([]string, 0)
	users = make([]User, 0)
	attendance = make([]Attendance, 0)
	queries = make([]Query, 0)

	r := mux.NewRouter()
	// Resources: gyms, users, attendance records.
	r.Methods("GET").Path("/gyms").HandlerFunc(GetGymsHandler)
	r.Methods("POST").Path("/gyms/{gym}").HandlerFunc(PostGymHandler)
	r.Methods("DELETE").Path("/gyms/{gym}").HandlerFunc(DelGymHandler)
	r.Methods("GET").Path("/users").HandlerFunc(GetUsersHandler)
	r.Methods("POST").Path("/users/{id}/{first}/{last}/{email}").HandlerFunc(PostUserHandler)
	r.Methods("DELETE").Path("/users/{id}").HandlerFunc(DelUserHandler)
	r.Methods("GET").Path("/attendance").HandlerFunc(GetAttendanceHandler)
	r.Methods("POST").Path("/attendance/{id}/{gym}/{when}").HandlerFunc(PostAttendanceHandler)
	// Actions: email user
	r.Methods("POST").Path("/webhooks/ask/{id}/{gym}/{when}").HandlerFunc(AskUserHandler)
	r.Methods("GET").Path("/webhooks/ack/{n}").HandlerFunc(ClickHandler)

	log.Fatal(http.ListenAndServe(webhookDomain, r))
}
