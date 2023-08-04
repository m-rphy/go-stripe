package main

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

type templateData struct {
    StringMap map[string]string
    IntMap map[string]int
    Float map[string]float32
    Data map[string]interface{}
    CSRToken string
    Flash string
    Warning string
    Error string
    IsAuthenticated int
    API string
    CSSVersion string
}


var functions = template.FuncMap {

}

// You get to compile the app and all the associated templates into
// a single binary

//go:embed templates
var templateFS embed.FS



func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
    return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
    var t *template.Template
    var err error

    templateToRender := fmt.Sprintf("templates/%s.page.tmpl", page)

    _, templateInMap := app.templateCache[templateToRender]

    // To only build templates in production mode
    if app.config.env == "production" && templateInMap {
        t = app.templateCache[templateToRender]
    } else {
        t, err = app.parseTemplate(partials, page, templateToRender)
        if err != nil {
            app.errorLog.Println(err)
            return err
        }
    }

    if td == nil {
        td = &templateData{}
    }

    td = app.addDefaultData(td, r)

    err = t.Execute(w, td)
    if err != nil {
        app.errorLog.Println(err)
        return err
    }

    return nil
}

func (app *application) parseTemplate(partials []string, page string, templateToRender string) (*template.Template, error) {
    
    var t *template.Template
    var err error

    // build partials
    if len(partials) > 0 {
        for i, x := range partials {
            partials[i] = fmt.Sprintf("templates/%s.partial.tmpl", x)
        }
    }
    
    if len(partials) > 0 {
        t, err = template.New(fmt.Sprint("%s.page.tmpl", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.tmpl", strings.Join(partials, ","), templateToRender)
    } else {
        t, err = template.New(fmt.Sprint("%s.page.tmpl", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.tmpl", templateToRender)
    }
    
    if err != nil {
        app.errorLog.Println(err)
        return nil, err
    }

    app.templateCache[templateToRender] = t
    return t, nil
}
