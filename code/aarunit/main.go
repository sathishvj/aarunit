package aarinit

import (
	"appengine"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/post/add", addPostHandler)
	http.HandleFunc("/post/addF", addPostFHandler)

	http.HandleFunc("/user/add", addUserHandler)
	http.HandleFunc("/user/addF", addUserFHandler)
	http.HandleFunc("/user/login", loginUserHandler)
	http.HandleFunc("/user/loginF", loginFUserHandler)

	http.HandleFunc("/group/add", addGroupHandler)
	http.HandleFunc("/group/addF", addGroupFHandler)
	http.HandleFunc("/group/list", listGroupsHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	t, err := template.New("Main").ParseFiles("tmpl/index.tmpl")
	if err != nil {
		c.Errorf("main.go: rootHandler(): Error loading/parsing template: %s", err.Error())
		return
	}

	_, posts, _ := getPosts(r)

	type RetPost struct {
		PostNum     int
		UpVoteCount int
		Url         string
		Title       string
		Username    string
		TimeDiff    string
		Group       string
	}

	var retPosts []RetPost
	for i := 0; i < len(posts); i++ {
		dur := time.Now().Sub(posts[i].Timestamp)
		var timeDiff string
		if dur.Seconds() < 60 {
			timeDiff = strconv.Itoa(int(dur.Seconds())) + " secs "
		} else if dur.Seconds() >= 60 && dur.Minutes() < 60 {
			timeDiff = strconv.Itoa(int(dur.Minutes())) + " mins "
		} else if dur.Minutes() >= 60 && dur.Hours() < 24 {
			timeDiff = strconv.Itoa(int(dur.Hours())) + " hrs "
		} else if dur.Hours() > 24 {
			timeDiff = strconv.Itoa(int(dur.Hours()/24)) + " days "
		}

		var url string
		if posts[i].Kind == "url" {
			url = posts[i].Value
		} else {
			//TODO: link to details page
			url = "http://www.google.com"
		}

		newRetPost := RetPost{
			PostNum:     i + 1,
			UpVoteCount: 1234,
			Url:         url,
			Title:       posts[i].Title,
			Username:    "tbd",
			TimeDiff:    timeDiff,
			Group:       posts[i].Group,
		}
		retPosts = append(retPosts, newRetPost)
	}

	err = t.ExecuteTemplate(w, "Main", struct {
		Username string
		Posts    []RetPost
	}{
		"tbd",
		retPosts,
	})
	if err != nil {
		c.Errorf("main.go: rootFHandler(): Error executing template: %s", err.Error())
	}

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

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	t, err := template.New("Login").ParseFiles("tmpl/login.tmpl")
	if err != nil {
		c.Errorf("main.go: loginUserHandler(): Error loading/parsing template: %s", err.Error())
	}

	err = t.ExecuteTemplate(w, "Login", nil)
	if err != nil {
		c.Errorf("main.go: loginUserHandler(): Error executing template: %s", err.Error())
	}
}

func loginFUserHandler(w http.ResponseWriter, r *http.Request) {

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
	tags := r.FormValue("tags")

	//do data checking and cleansing here
	addGroup(r, name, tags)

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
