package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// initialize fiber app
	app := fiber.New()

	// register the GET '/' endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 80")
	if err := app.Listen(":80"); err != nil {
		log.Panic(err)
	}
}

func render(c *fiber.Ctx, tmplName string) error {
	// define the partial templates
	const templatesPath = "../../internal/front-end/templates"

	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", templatesPath),
		fmt.Sprintf("%s/header.partial.gohtml", templatesPath),
		fmt.Sprintf("%s/footer.partial.gohtml", templatesPath),
	}

	// create a slice with the main template and then append the partials
	templateFiles := make([]string, len(partials)+1)
	templateFiles[0] = fmt.Sprintf("%s/%s", templatesPath, tmplName)
	for i, v := range partials {
		templateFiles[i+1] = v
	}

	// parse all template files
	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// execute teh template into a buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// set the content-type header to text/html and send the result
	c.Type("html")
	return c.Send(buf.Bytes())
}
