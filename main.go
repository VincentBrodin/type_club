package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	textgenerator "github.com/VincentBrodin/type_club/pkgs/text_generator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var modal map[string][]string
var store *session.Store

type RunInputs struct {
	Value string  `json:"value"`
	Time  float64 `json:"time"`
}

type TypeRun struct {
	Target   string      `json:"target"`
	Html     string      `json:"html"`
	Accuracy float64     `json:"accuracy"`
	Wpm      float64     `json:"wpm"`
	Awpm     float64     `json:"awpm"`
	Time     float64     `json:"time"`
	Inputs   []RunInputs `json:"inputs"`
}

func Clean(dirty TypeRun) TypeRun {
	clean := dirty

	clean.Accuracy = math.Round(dirty.Accuracy*1000) / 10
	clean.Wpm = math.Round(dirty.Wpm*100) / 100
	clean.Awpm = math.Round(dirty.Awpm*100) / 100
	clean.Time = math.Round(dirty.Time*100) / 100

	return clean
}

func main() {
	var err error
	modal, err = textgenerator.LoadModel("./pkgs/text_generator/markov_model.json")
	if err != nil {
		panic(err)
	}

	store = session.New()

	engine := html.New("./views", ".html")
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
		fmt.Println(err)
		return c.SendStatus(400)
	}

	id := c.Query("id")
	payload := TypeRun{}

	// if the last run was not to be saved
	if id == "last" {
		body := fmt.Sprintf("%v", sess.Get("last"))

		b := []byte(body)
		if err := json.Unmarshal(b, &payload); err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}
		return c.Render("stats", fiber.Map{
			"Title":   "type_club | Stats",
			"TypeRun": Clean(payload),
		}, "layouts/main")
	}
	return c.SendString(id)
}

func Replay(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}

	id := c.Query("id")
	payload := TypeRun{}

	// if the last run was not to be saved
	if id == "last" {
		body := fmt.Sprintf("%v", sess.Get("last"))

		b := []byte(body)
		if err := json.Unmarshal(b, &payload); err != nil {
			return c.SendStatus(400)
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return c.SendStatus(400)
		}
		return c.SendString(string(data))
	}
	return c.SendString(id)
}

func Random(c *fiber.Ctx) error {
	param := c.Query("length")
	length, err := strconv.Atoi(param)
	if err != nil {
		return c.SendStatus(400)
	}

	output := textgenerator.GenerateSentence(modal, length, "")
	return c.SendString(output)
}

func RunDone(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}

	payload := TypeRun{}
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
