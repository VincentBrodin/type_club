package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"

	"github.com/VincentBrodin/type_club/pkgs/pswdhash"
	"github.com/VincentBrodin/type_club/pkgs/textgen"
	"github.com/VincentBrodin/type_club/pkgs/typeruns"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var modal map[string][]string
var store *session.Store

func main() {
	var err error
	modal, err = textgen.LoadModel("./pkgs/textgen/markov_model.json")
	if err != nil {
		panic(err)
	}

	_, err = sql.Open("sqlite3", "data.db")
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
	return c.Redirect("/login")
}

func GetLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "type_club | Login",
	}, "layouts/main")

}

func PostLogin(c *fiber.Ctx) error {
	return c.SendStatus(404)
}

func GetRegister(c *fiber.Ctx) error {
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
		return c.SendStatus(400)
	}

	fmt.Printf("%v %v %v %v\n", username, email, password, hashed)
	return c.SendStatus(404)
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
		return c.SendStatus(400)
	}

	return c.Redirect("/stats?id=last")
}
