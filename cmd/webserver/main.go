package main

import (
	"encoding/base64"
	"log"
	"subgen/internal/config"
	"subgen/internal/db"
	"subgen/internal/userlink"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db.InitDB()
	db.Migrate()

	engine := html.New("./web/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./web/static")

	// Add authentication check middleware
	authRequired := func(c *fiber.Ctx) error {
		if c.Cookies("admin_auth") != "1" {
			return c.Redirect("/login")
		}
		return c.Next()
	}

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{"Error": ""})
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		var adminRec struct{ Username, Password string }
		row := db.DB.Raw("SELECT username, password FROM admins WHERE username = ?", username).Row()
		_ = row.Scan(&adminRec.Username, &adminRec.Password)
		if adminRec.Username == username && bcrypt.CompareHashAndPassword([]byte(adminRec.Password), []byte(password)) == nil {
			c.Cookie(&fiber.Cookie{Name: "admin_auth", Value: "1", Path: "/"})
			return c.Redirect("/dashboard")
		}
		return c.Render("login", fiber.Map{"Error": "Invalid credentials"})
	})
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "admin_auth", Value: "", Path: "/", MaxAge: -1})
		return c.Redirect("/login")
	})

	app.Get("/dashboard", authRequired, func(c *fiber.Ctx) error {
		return c.Render("dashboard", fiber.Map{})
	})
	app.Get("/config/add", authRequired, func(c *fiber.Ctx) error {
		return c.Render("add_config", fiber.Map{})
	})
	app.Post("/config/add", authRequired, func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		content := c.FormValue("content")
		if name == "" || content == "" {
			return c.Render("add_config", fiber.Map{"Error": "Name and content are required."})
		}
		cfg := config.Config{Name: name, Content: content}
		if err := db.DB.Create(&cfg).Error; err != nil {
			return c.Render("add_config", fiber.Map{"Error": "Failed to save config."})
		}
		return c.Redirect("/config/list")
	})
	app.Get("/config/list", authRequired, func(c *fiber.Ctx) error {
		var configs []config.Config
		if err := db.DB.Find(&configs).Error; err != nil {
			return c.Render("list_configs", fiber.Map{"Configs": []interface{}{}, "Error": "Failed to fetch configs."})
		}
		return c.Render("list_configs", fiber.Map{"Configs": configs})
	})
	app.Get("/config/edit/:id", authRequired, func(c *fiber.Ctx) error {
		id := c.Params("id")
		var cfg config.Config
		if err := db.DB.First(&cfg, id).Error; err != nil {
			return c.Render("edit_config", fiber.Map{"ID": id, "Name": "", "Content": "", "Error": "Config not found."})
		}
		return c.Render("edit_config", fiber.Map{"ID": cfg.ID, "Name": cfg.Name, "Content": cfg.Content})
	})
	app.Post("/config/edit/:id", authRequired, func(c *fiber.Ctx) error {
		id := c.Params("id")
		name := c.FormValue("name")
		content := c.FormValue("content")
		if name == "" || content == "" {
			return c.Render("edit_config", fiber.Map{"ID": id, "Name": name, "Content": content, "Error": "Name and content are required."})
		}
		var cfg config.Config
		if err := db.DB.First(&cfg, id).Error; err != nil {
			return c.Render("edit_config", fiber.Map{"ID": id, "Name": name, "Content": content, "Error": "Config not found."})
		}
		cfg.Name = name
		cfg.Content = content
		if err := db.DB.Save(&cfg).Error; err != nil {
			return c.Render("edit_config", fiber.Map{"ID": id, "Name": name, "Content": content, "Error": "Failed to update config."})
		}
		return c.Redirect("/config/list")
	})
	app.Get("/userlink", authRequired, func(c *fiber.Ctx) error {
		var uuidRec userlink.UUID
		if err := db.DB.First(&uuidRec).Error; err != nil {
			return c.Render("userlink", fiber.Map{"UserLink": "UUID not set"})
		}
		baseURL := c.BaseURL()
		userLink := baseURL + "/" + uuidRec.UUID
		return c.Render("userlink", fiber.Map{"UserLink": userLink})
	})
	// User config endpoint (public, returns base64 config by UUID)
	app.Get("/:uuid", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")

		var uuidRec userlink.UUID
		if err := db.DB.Where("uuid = ?", uuid).First(&uuidRec).Error; err != nil {
			return c.Status(404).SendString("UUID not found")
		}
		var configs []config.Config
		if err := db.DB.Find(&configs).Error; err != nil {
			return c.Status(500).SendString("Failed to fetch configs")
		}

		var combined string
		for _, cfg := range configs {
			combined += cfg.Content + "\n"
		}
		b64 := base64.StdEncoding.EncodeToString([]byte(combined))
		return c.SendString(b64)
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	log.Println("Fiber webserver running on :8095")
	log.Fatal(app.Listen(":8095"))
}
