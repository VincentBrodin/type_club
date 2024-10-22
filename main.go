package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/VincentBrodin/type_club/pkgs/pswdhash"
	"github.com/VincentBrodin/type_club/pkgs/textgen"
	"github.com/VincentBrodin/type_club/pkgs/typeruns"
	"github.com/VincentBrodin/type_club/pkgs/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var modal map[string][]string
var store *session.Store
var db *sql.DB

func main() {
	var err error
	modal, err = textgen.LoadModel("./pkgs/textgen/markov_model.json")
	if err != nil {
		panic(err)
	}

	db, err = sql.Open("sqlite3", "db/base.db")
	if err != nil {
		panic(err)
	}

	store = session.New()

	engine := html.New("./views", ".html")
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//Static
	app.Static("/js/", "./src/js/", fiber.Static{
		CacheDuration: 0,
		Browse:        true,
	})
	app.Static("/images/", "./src/images/")

	//Routs
	app.Get("/", GetHome)
	app.Get("/stats/:id?", GetStats)
	app.Post("/stats", PostStats)
	app.Post("/sentence", PostSentence)

	app.Get("/account", GetAccount)

	app.Get("/login", GetLogin)
	app.Post("/login", PostLogin)

	app.Get("/register", GetRegister)
	app.Post("/register", PostRegister)

	app.Get("/logout", GetLogout)

	app.Post("/validate", PostValidateAccount)
	app.Post("/update", PostUpdate)
	app.Post("/check", PostCheck)

	app.Post("/done", PostDone)
	app.Post("/save", PostSave)

	log.Fatal(app.Listen(":3000"))
}

func GetHome(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":    "type_club | Home",
		"LoggedIn": loggedIn(c),
	}, "layouts/main")
}

func GetStats(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Redirect("/")
	}

	id := c.Query("id")
	typerun := typeruns.NewTypeRun()

	// if the last run was not to be saved
	if id == "last" {
		body := fmt.Sprintf("%v", sess.Get("last"))

		b := []byte(body)
		if err := json.Unmarshal(b, &typerun); err != nil {
			return c.Redirect("/")
		}
		return c.Render("stats", fiber.Map{
			"Title":    "type_club | Stats",
			"LoggedIn": loggedIn(c),
			"TypeRun":  typerun.Clean(),
		}, "layouts/main")
	}
	return c.SendString(id)
}

func PostStats(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(400)
	}

	input := &struct {
		Id string `json:"id"`
	}{}

	err = c.BodyParser(input)
	if err != nil {
		return c.SendStatus(400)
	}

	if input.Id == "last" {
		body := fmt.Sprintf("%v", sess.Get("last"))
		return c.SendString(body)
	}
	return c.SendString(input.Id)
}

func PostSentence(c *fiber.Ctx) error {
	input := &struct {
		Start  string `json:"start"`
		Length int    `json:"length"`
	}{}

	err := c.BodyParser(input)
	if err != nil {
		return c.SendStatus(400)
	}

	output := textgen.GenerateSentence(modal, input.Length, input.Start)
	return c.SendString(output)
}

func GetAccount(c *fiber.Ctx) error {
	if !loggedIn(c) {
		return c.Redirect("/login")
	}

	user, err := getUserFromSess(c)
	if err != nil {
		return c.Redirect("/login")
	}

	return c.Render("account", fiber.Map{
		"Title":    "type_club | Account",
		"LoggedIn": loggedIn(c),
		"User":     *user,
	}, "layouts/main")
}

func loggedIn(c *fiber.Ctx) bool {
	sess, err := store.Get(c)
	if err != nil {
		return false
	}
	return !(sess.Get("user") == nil)
}

func GetLogin(c *fiber.Ctx) error {
	if loggedIn(c) {
		return c.Redirect("/account")
	}
	return c.Render("login", fiber.Map{
		"Title":    "type_club | Login",
		"LoggedIn": loggedIn(c),
	}, "layouts/main")

}

func PostLogin(c *fiber.Ctx) error {
	if loggedIn(c) {
		return c.Redirect("/account")
	}
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := users.FindByUsername(username, db)
	if err != nil {
		return c.SendStatus(400)
	}

	if pswdhash.VerifyPassword(password, user.Password) {
		err = addUserToSess(user, c)
		if err != nil {
			return c.SendStatus(404)
		}
		return c.Redirect("/account")
	}
	return c.Redirect("/login")
}

func addUserToSess(user *users.User, c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	sess.Set("user", string(data))
	err = sess.Save()
	if err != nil {
		return err
	}
	return nil
}

func getUserFromSess(c *fiber.Ctx) (*users.User, error) {
	sess, err := store.Get(c)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Id: -1,
	}
	data := []byte(fmt.Sprintf("%v", sess.Get("user")))
	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}

	if user.Id == -1 {
		return nil, fmt.Errorf("Could not get user")
	}

	return user, nil
}

func GetRegister(c *fiber.Ctx) error {
	if loggedIn(c) {
		return c.Redirect("/account")
	}
	return c.Render("register", fiber.Map{
		"Title":    "type_club | Register",
		"LoggedIn": loggedIn(c),
	}, "layouts/main")

}

func PostRegister(c *fiber.Ctx) error {
	if loggedIn(c) {
		return c.Redirect("/account")
	}
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	hashed, err := pswdhash.HashPassword(password)

	if err != nil {
		fmt.Println(err)
		return c.Redirect("/register")
	}

	if !validUsername(username) || !validEmail(email) {
		return c.Redirect("/register")

	}

	user := users.New(username, email, hashed)
	err = user.AddToDb(db)
	if err != nil {
		fmt.Println(err)
		return c.Redirect("/register")
	}

	err = addUserToSess(user, c)
	if err != nil {
		fmt.Println(err)
		return c.Redirect("/register")
	}

	return c.Redirect("/account")
}

func GetLogout(c *fiber.Ctx) error {
	if !loggedIn(c) {
		return c.Redirect("/login")
	}

	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.Delete("user")
	err = sess.Save()
	if err != nil {
		return err
	}

	return c.Redirect("/login")
}

func PostValidateAccount(c *fiber.Ctx) error {
	input := &struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{}

	err := c.BodyParser(input)
	if err != nil {
		return c.SendStatus(400)
	}

	output := struct {
		Username bool `json:"username"`
		Email    bool `json:"email"`
	}{
		Username: validUsername(input.Username),
		Email:    validEmail(input.Email),
	}

	data, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(404)
	}

	return c.SendString(string(data))
}

func PostUpdate(c *fiber.Ctx) error {
	input := &struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.BodyParser(input)
	if err != nil {
		return c.SendStatus(400)
	}

	user, err := getUserFromSess(c)
	if err != nil {
		return c.SendStatus(400)
	}

	if input.Username != "" {
		if !validUsername(input.Username) {
			return c.SendStatus(400)
		}
		user.Username = input.Username
	}
	if input.Email != "" {
		if !validEmail(input.Email) {
			return c.SendStatus(400)
		}
		user.Email = input.Email
	}
	if input.Password != "" {
		hash, err := pswdhash.HashPassword(input.Password)
		if err != nil {
			return c.SendStatus(400)
		}
		user.Password = hash
	}

	err = user.Update(db)
	if err != nil {
		return c.SendStatus(400)
	}

	err = addUserToSess(user, c)
	if err != nil {
		return c.SendStatus(400)
	}

	return c.SendStatus(200)
}

func PostCheck(c *fiber.Ctx) error {
	input := &struct {
		Password string `json:"password"`
	}{}

	err := c.BodyParser(input)
	if err != nil {
		return c.SendStatus(400)
	}

	user, err := getUserFromSess(c)
	if err != nil {
		return c.SendStatus(400)
	}

	valid := pswdhash.VerifyPassword(input.Password, user.Password)

	output := struct {
		Valid bool `json:"valid"`
	}{
		Valid: valid,
	}

	data, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(404)
	}

	return c.SendString(string(data))

}

func validEmail(email string) bool {
	user, err := users.FindByEmail(email, db)
	if err != nil {
		return true
	}
	if user.Id == -1 {
		return true
	}
	return false
}

func validUsername(username string) bool {
	user, err := users.FindByUsername(username, db)
	if err != nil {
		return true
	}
	if user.Id == -1 {
		return true
	}
	return false
}

func PostDone(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}

	typerun := typeruns.NewTypeRun()
	if err := c.BodyParser(typerun); err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}

	// Store the last run as a cookie
	body := string(c.Body())
	sess.Set("last", body)

	if err := sess.Save(); err != nil {
		fmt.Println(err)
		return c.Redirect("/")

	}

	return c.Redirect("/stats?id=last")
}

func PostSave(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(400)
	}

	if sess.Get("last") == nil {
		return c.SendStatus(400)
	}

	body := []byte(fmt.Sprintf("%v", sess.Get("last")))

	typerun := typeruns.NewTypeRun()

	if err := json.Unmarshal(body, typerun); err != nil {
		return c.SendStatus(400)
	}

	err = typerun.AddToDb(db)
	if err != nil {
		return c.SendStatus(400)
	}

	return c.SendStatus(200)
}
