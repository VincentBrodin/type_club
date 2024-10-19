package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
	app.Static("/js/", "./src/js/")
	app.Static("/images/", "./src/images/")

	//Routs
	app.Get("/", GetHome)
	app.Get("/stats/:id?", GetStats)
	app.Get("/replay/:id?", GetReplay)
	app.Get("/random/:length?", GetRandom)

	app.Get("/account", GetAccount)

	app.Get("/login", GetLogin)
	app.Post("/login", PostLogin)

	app.Get("/register", GetRegister)
	app.Post("/register", PostRegister)

	app.Post("/validate", PostValidateAccount)

	app.Post("/done", PostDone)

	log.Fatal(app.Listen(":3000"))
}

func GetHome(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "type_club | Home",
	}, "layouts/main")
}

func GetStats(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Redirect("/")
	}

	id := c.Query("id")
	payload := typeruns.NewTypeRun()

	// if the last run was not to be saved
	if id == "last" {
		body := fmt.Sprintf("%v", sess.Get("last"))

		b := []byte(body)
		if err := json.Unmarshal(b, &payload); err != nil {
			return c.Redirect("/")
		}
		return c.Render("stats", fiber.Map{
			"Title":   "type_club | Stats",
			"TypeRun": typeruns.Clean(payload),
		}, "layouts/main")
	}
	return c.SendString(id)
}

func GetReplay(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(400)
	}

	id := c.Query("id")
	// if the last run was not to be saved
	if id == "last" {
		body := fmt.Sprintf("%v", sess.Get("last"))
		if body == "<nil>" {
			return c.SendStatus(404)
		}
		return c.SendString(body)
	}
	return c.SendString(id)
}

func GetRandom(c *fiber.Ctx) error {
	param := c.Query("length")
	length, err := strconv.Atoi(param)
	if err != nil {
		return c.SendStatus(400)
	}

	output := textgen.GenerateSentence(modal, length, "")
	return c.SendString(output)
}

func GetAccount(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(400)
	}

	body := fmt.Sprintf("%v", sess.Get("user"))
	b := []byte(body)
	user := &users.User{
		Id: -1,
	}
	err = json.Unmarshal(b, user)

	if err != nil || user.Id == -1 {
		return c.Redirect("/login")
	}

	return c.Render("account", fiber.Map{
		"Title": "type_club | Account",
	}, "layouts/main")
}

func loggedIn(c *fiber.Ctx) bool {
	sess, err := store.Get(c)
	if err != nil {
		return false
	}

	return sess.Get("user") != nil
}

func GetLogin(c *fiber.Ctx) error {
	if loggedIn(c) {
		return c.Redirect("/account")
	}
	return c.Render("login", fiber.Map{
		"Title": "type_club | Login",
	}, "layouts/main")

}

func PostLogin(c *fiber.Ctx) error {
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
		return c.Redirect("/")
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

func GetRegister(c *fiber.Ctx) error {
	if loggedIn(c) {
		return c.Redirect("/account")
	}
	return c.Render("register", fiber.Map{
		"Title": "type_club | Register",
	}, "layouts/main")

}

func PostRegister(c *fiber.Ctx) error {
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

	payload := typeruns.NewTypeRun()
	if err := c.BodyParser(&payload); err != nil {
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
