package aarinit

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"net/http"
	"time"
)

type Post struct {
	Id        string
	Kind      string
	Title     string
	Value     string
	Group     string
	Timestamp time.Time
}

type Comment struct {
	Id        string
	PostId    string
	Value     string
	Username  string
	Timestamp time.Time
}

type User struct {
	Username  string
	Password  string
	Email     string
	Timestamp time.Time
}

type Group struct {
	Name      string
	Tags      string
	Timestamp time.Time
}

// the data given should be already cleansed.  No checks are done here.
func addPost(r *http.Request, kind, title, value, group string) (err error) {
	c := appengine.NewContext(r)

	k := datastore.NewIncompleteKey(c, "Post", nil)
	uuid, _ := getNewUuid()
	p := Post{
		uuid,
		kind,
		title,
		value,
		group,
		time.Now(),
	}

	_, err = datastore.Put(c, k, &p)
	if err != nil {
		c.Errorf("model.go: addPost(): Error adding post: %s", err.Error())
		return err
	}

	c.Infof("model.go: addPost(): Successfully added post")
	return err
}

func getPosts(r *http.Request) (pKeys []*datastore.Key, posts []Post, err error) {

	cnt := 10

	c := appengine.NewContext(r)
	q := datastore.NewQuery("Post").
		Limit(cnt).
		Order("-Timestamp")

	//posts = new([]Post)
	posts = make([]Post, 0)
	if pKeys, err = q.GetAll(c, &posts); err != nil {
		c.Errorf("model.go: getPosts(): Error getting posts: %s", err.Error())
		return nil, nil, err
	}

	return pKeys, posts, err
}

// the data given should be already cleansed.  No checks are done here.
func addUser(r *http.Request, username, password, email string) (err error) {
	c := appengine.NewContext(r)

	k := datastore.NewIncompleteKey(c, "User", nil)
	u := User{
		username,
		password,
		email,
		time.Now(),
	}

	_, err = datastore.Put(c, k, &u)
	if err != nil {
		c.Errorf("model.go: addUser(): Error adding user: %s", err.Error())
		return err
	}

	c.Infof("model.go: addUser(): Successfully added user")
	return err
}

func validateUser(r *http.Request, username, password string) (bool, error) {

	c := appengine.NewContext(r)

	q := datastore.NewQuery("User").
		Filter("Username =", username).
		Filter("Password =", password).
		KeysOnly()

	if cnt, err := q.Count(c); err != nil {
		c.Errorf("model.go: validateUser(): Error getting user count: %s", err.Error())
		return false, err
	} else if cnt > 1 {
		err = errors.New("More than one user with same name and password.")
		c.Errorf("model.go: validateUser(): %s", err)
		return false, err
	} else if cnt == 1 {
		return true, nil
	}

	return false, nil
}

func addGroup(r *http.Request, name, tags string) (err error) {
	c := appengine.NewContext(r)

	k := datastore.NewIncompleteKey(c, "Group", nil)
	g := Group{
		name,
		tags,
		time.Now(),
	}

	_, err = datastore.Put(c, k, &g)
	if err != nil {
		c.Errorf("model.go: addGroup(): Error adding group: %s", err.Error())
		return err
	}

	c.Infof("model.go: addGroup(): Successfully added group")
	return err
}

func getGroups(r *http.Request) (gKeys []*datastore.Key, groups *[]Group, err error) {

	c := appengine.NewContext(r)
	q := datastore.NewQuery("Group").
		Order("Name")

	groups = new([]Group)
	if gKeys, err = q.GetAll(c, groups); err != nil {
		c.Errorf("model.go: getGroups(): Error getting groups: %s", err.Error())
		return nil, nil, err
	}

	return gKeys, groups, err
}

// the data given should be already cleansed.  No checks are done here.
func addComment(r *http.Request, postId, value, username string) (err error) {
	c := appengine.NewContext(r)

	k := datastore.NewIncompleteKey(c, "Comment", nil)

	uuid, _ := getNewUuid()
	p := Comment{
		uuid,
		postId,
		value,
		username,
		time.Now(),
	}

	_, err = datastore.Put(c, k, &p)
	if err != nil {
		c.Errorf("model.go: addComment(): Error adding comment: %s", err.Error())
		return err
	}

	c.Infof("model.go: addComment(): Successfully added comment")
	return err
}

func getComments(r *http.Request, postId string) (cKeys []*datastore.Key, comments *[]Comment, err error) {

	c := appengine.NewContext(r)
	q := datastore.NewQuery("Comment").
		Filter("PostId=", postId).
		Order("-Timestamp")

	comments = new([]Comment)
	if cKeys, err = q.GetAll(c, comments); err != nil {
		c.Errorf("model.go: getComments(): Error getting posts: %s", err.Error())
		return nil, nil, err
	}

	return cKeys, comments, err
}
