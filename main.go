package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
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
	app.Static("/js/", "./src/js/", fiber.Static{
		CacheDuration: 0,
		Browse:        true,
	})
	app.Static("/images/", "./src/images/")

	//Routs
	app.Get("/", GetHome)
	app.Get("/stats:id?", GetStats)
	app.Post("/stats", PostStats)
	app.Post("/sentence", PostSentence)

	app.Get("/account", GetAccount)
	app.Get("/profile:id?", GetProfile)
	app.Post("/profile", PostProfile)

	app.Get("/leaderboard", GetLeaderboard)
	app.Post("/leaderboard", PostLeaderboard)

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
			"Last":     true,
		}, "layouts/main")
	} else {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.SendStatus(400)
		}

		typerun, err := typeruns.FindById(i, db)
		if err != nil {
			return c.SendStatus(404)
		}

		return c.Render("stats", fiber.Map{
			"Title":    "type_club | Stats",
			"LoggedIn": loggedIn(c),
			"TypeRun":  typerun.Clean(),
			"Last":     false,
		}, "layouts/main")

	}
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
	} else {
		i, err := strconv.ParseInt(input.Id, 10, 64)
		if err != nil {
			return c.SendStatus(400)
		}

		typerun, err := typeruns.FindById(i, db)
		if err != nil {
			return c.SendStatus(404)
		}

		body, err := json.Marshal(typerun)

		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(body))
	}
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

func GetProfile(c *fiber.Ctx) error {
	id := int64(c.QueryInt("id", -1))
	if id == -1 {
		if loggedIn(c) {
			localUser, err := getUserFromSess(c)
			if err != nil {
				return c.SendStatus(404)
			}

			return c.Redirect(fmt.Sprintf("/profile?id=%d", localUser.Id))
		} else {
			return c.Redirect("/login")
		}
	}

	user, err := users.FindById(id, db)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.Render("profile", fiber.Map{
		"Title":    fmt.Sprintf("type_club | %v", user.Username),
		"LoggedIn": loggedIn(c),
		"User":     user,
	}, "layouts/main")
}

func PostProfile(c *fiber.Ctx) error {
	input := &struct {
		Id int64 `json:"id"`
	}{}
	err := c.BodyParser(input)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}

	fmt.Println(input.Id)
	user, err := users.FindById(input.Id, db)
	if err != nil {
		return c.SendStatus(404)
	}

	runs, err := typeruns.FindByOwner(user.Id, db)
	if err != nil {
		return c.SendStatus(404)
	}
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].Id > runs[j].Id
	})

	output := struct {
		User  users.User            `json:"user"`
		Runs  []typeruns.TypeRun    `json:"runs"`
		Stats typeruns.ProfileStats `json:"stats"`
	}{
		User:  *user,
		Runs:  runs,
		Stats: typeruns.CalculateStats(runs),
	}

	data, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(404)
	}

	return c.SendString(string(data))
}

func GetLeaderboard(c *fiber.Ctx) error {
	return c.Render("leaderboard", fiber.Map{
		"Title":    "type_club | Leaderboard",
		"LoggedIn": loggedIn(c),
	}, "layouts/main")
}

func PostLeaderboard(c *fiber.Ctx) error {
	runs, err := typeruns.FindBest(100, db)
	if err != nil {
		return c.SendStatus(404)
	}
	fmt.Println(runs)

	type UserRun struct {
		users.User       `json:"user"`
		typeruns.TypeRun `json:"run"`
	}

	userRuns := make([]UserRun, 0)

	for _, run := range runs {
		user, err := users.FindById(run.OwnerId, db)
		if err != nil {
			continue
		}

		userRuns = append(userRuns, UserRun{
			User:    *user,
			TypeRun: run,
		})
	}

	data, err := json.Marshal(userRuns)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.SendString(string(data))
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
	if !loggedIn(c) {
		return c.SendStatus(400)
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(400)
	}

	if sess.Get("last") == nil {
		return c.SendStatus(400)
	}

	user, err := getUserFromSess(c)
	if err != nil {
		return c.SendStatus(400)
	}

	body := []byte(fmt.Sprintf("%v", sess.Get("last")))
	typerun := typeruns.NewTypeRun()
	typerun.OwnerId = user.Id

	if err := json.Unmarshal(body, typerun); err != nil {
		return c.SendStatus(400)
	}

	err = typerun.AddToDb(db)
	if err != nil {
		return c.SendStatus(400)
	}

	return c.Redirect(fmt.Sprintf("/stats?id=%d", typerun.Id))
}
