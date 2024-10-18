package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
	app.Get("/", Home)
	app.Get("/stats/:id?", Stats)
	app.Get("/replay/:id?", Replay)
	app.Get("/random/:length?", Random)
	app.Post("/done", RunDone)

	log.Fatal(app.Listen(":3000"))
}

func Home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "type_club | Home",
	}, "layouts/main")
}

func Stats(c *fiber.Ctx) error {
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

func Replay(c *fiber.Ctx) error {
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

func Random(c *fiber.Ctx) error {
	param := c.Query("length")
	length, err := strconv.Atoi(param)
	if err != nil {
		return c.SendStatus(400)
	}

	output := textgen.GenerateSentence(modal, length, "")
	return c.SendString(output)
}

func RunDone(c *fiber.Ctx) error {
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
