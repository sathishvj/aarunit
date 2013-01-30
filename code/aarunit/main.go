package aarunit

import (
	"appengine"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/post/add", addPostHandler)
	http.HandleFunc("/post/addF", addPostFHandler)

	http.HandleFunc("/user/add", addUserHandler)
	http.HandleFunc("/user/addF", addUserFHandler)
	http.HandleFunc("/user/validate", validateUserHandler)

	http.HandleFunc("/group/add", addGroupHandler)
	http.HandleFunc("/group/addF", addGroupFHandler)
	http.HandleFunc("/group/list", listGroupsHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	_, posts, _ := getPosts(r)
	b, _ := json.Marshal(posts)
	fmt.Fprintf(w, "%s", string(b))
}

func addPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	t, err := template.New("AddPost").ParseFiles("tmpl/addPost.tmpl")
	if err != nil {
		c.Errorf("main.go: addPostHandler(): Error loading/parsing template: %s", err.Error())
	}

	_, groups, _ := getGroups(r)
	err = t.ExecuteTemplate(w, "AddPost", groups)
	if err != nil {
		c.Errorf("main.go: addPostHandler(): Error executing template: %s", err.Error())
	}
}

func addPostFHandler(w http.ResponseWriter, r *http.Request) {
	kind := r.FormValue("kind")
	title := r.FormValue("title")
	value := r.FormValue("value")
	group := r.FormValue("group")

	//do data checking and cleansing here
	addPost(r, kind, title, value, group)

	http.Redirect(w, r, "/", http.StatusFound)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	t, err := template.New("AddUser").ParseFiles("tmpl/addUser.tmpl")
	if err != nil {
		c.Errorf("main.go: addUserHandler(): Error loading/parsing template: %s", err.Error())
	}

	err = t.ExecuteTemplate(w, "AddUser", nil)
	if err != nil {
		c.Errorf("main.go: addUserHandler(): Error executing template: %s", err.Error())
	}
}

func addUserFHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	//do data checking and cleansing here
	addUser(r, username, password, email)

	http.Redirect(w, r, "/", http.StatusFound)
}

func validateUserHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	username := r.FormValue("username")
	password := r.FormValue("password")

	//do data checking and cleansing here
	ok, err := validateUser(r, username, password)
	if err != nil || !ok {
		c.Errorf("main.go: validateUserHandler(): Error validating user: %s", err)

		fmt.Fprintf(w, "%s", getSrvRetErrStr(err))

		return
	}

	//http.Redirect(w, r, "/", http.StatusFound)
	fmt.Fprintf(w, "%s", getSrvRetSuccessStr("User "+username+" successfully validated."))
}

func addGroupHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	t, err := template.New("AddGroup").ParseFiles("tmpl/addGroup.tmpl")
	if err != nil {
		c.Errorf("main.go: addGroupFHandler(): Error loading/parsing template: %s", err.Error())
	}

	err = t.ExecuteTemplate(w, "AddGroup", nil)
	if err != nil {
		c.Errorf("main.go: addGroupFHandler(): Error executing template: %s", err.Error())
	}
}

func addGroupFHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	//do data checking and cleansing here
	addGroup(r, name)

	http.Redirect(w, r, "/", http.StatusFound)
}

func listGroupsHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	_, groups, err := getGroups(r)

	if err != nil {
		c.Errorf("main.go: listGroupsHandler(): Error calling getGroups: %s", err)

		fmt.Fprintf(w, "%s", getSrvRetErrStr(err))

		return
	}

	//checking return value
	s := getSrvRetSuccessStr(groups)
	c.Infof("main.go: listGroupsHandler(): List of groups: %s", s)

	fmt.Fprintf(w, "%s", getSrvRetSuccessStr(groups))
}
