package render

import (
	"log"
	"net/http"
)

func (t *TemplatesHTML) Render(w http.ResponseWriter, r *http.Request, name string, data *PageData) {
	tmlp, ok := (*t)[name]
	if !ok {

		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmlp.Execute(w, data)
	if err != nil {
		log.Println("Error executing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
