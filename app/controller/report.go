package controller

import (
//      "log"
	"net/http"
//        "fmt"
        "strings"
//        "time"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/view"
	"github.com/gatopardo/micondo/app/shared/email"

//        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
//	"github.com/julienschmidt/httprouter"
  )
// ---------------------------------------------------
// MailSendGet despliega formulario para enviar correo
func MailSendGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisApts, err := model.Apts()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay receptor", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "report/correo"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisApts"]    = lisApts
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// MailSendPOST procesa la forma enviada con contenido
func MailSendPOST(w http.ResponseWriter, r *http.Request) {
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            aptId,  _       := atoi32(r.FormValue("aptId"))
	    tema            := r.FormValue("tema")
	    content         := r.FormValue("content")
//	    fmt.Printf("%d %s %s\n",aptId, tema, content)
            person, err         :=  model.EmailByAptId(aptId)
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay correos ", view.FlashError})
            }else{
		    to := person.Email
//	    fmt.Printf("%d | %s | %s | %s | %s\n",aptId, to, person.Fname ,tema, content)
	        email.SendEmail(to, tema,content);
            }
        }
	http.Redirect(w, r, "/email", http.StatusFound)
 }
// ---------------------------------------------------
