package controller

import (
      "log"
	"net/http"
//        "fmt"
        "strings"
//        "time"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/view"
	"github.com/gatopardo/micondo/app/shared/email"

//        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
//       "github.com/julienschmidt/httprouter"
  )
// ---------------------------------------------------
// MailSendGet despliega formulario para enviar correo
func MailSendGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisApts, err := model.Apts()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay aptos", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "report/rptmail"
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
// RptAptGet reporte estado de apto
func RptAptGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos", view.FlashError})
         }
	v                     := view.New(r)
	v.Name                 = "report/rptper"
	v.Vars["token"]        = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"]   = lisPeriods
        v.Vars["Level"]        =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// RptAptPOST reporte estado de apto
func RptAptPOST(w http.ResponseWriter, r *http.Request) {
	var pers model.Person
	var apt model.Aparta
	var prid model.Periodo
	var err error
	sess := model.Instance(r)
        uid, ok       := sess.Values["id"].(uint32)
	if ! ok {
             log.Println("No uint32 value in session")
	}
        sPerid    :=  r.FormValue("id")
	perid,_   := atoi32(sPerid)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             pers, apt, err = model.ApartaByUserId(uid)
	     if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No apto", view.FlashError})
	     }
            prid.Id = perid
            err := (&prid).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No hay periodo", view.FlashError})
             }
	     lisPaym, _            := model.Payments(apt.Id, prid.Inicio)
             v                     := view.New(r)
             v.Name                 = "report/rptapto"
	     v.Vars["token"]        = csrfbanana.Token(w, r, sess)
             v.Vars["Apt"]          = apt
             v.Vars["Pers"]         = pers
             v.Vars["LisPaym"]      = lisPaym
             v.Vars["Level"]        =  sess.Values["level"]
	     v.Render(w)
         }else{
	  http.Redirect(w, r, "/cuota/list", http.StatusFound)
	 }
 }
// ---------------------------------------------------
// RptCondGet reporte estado de apto
func RptCondGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos", view.FlashError})
         }
	v                     := view.New(r)
	v.Name                 = "report/condper"
	v.Vars["token"]        = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"]   = lisPeriods
        v.Vars["Level"]        =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// RptCondPOST reporte estado de apto
func RptCondPOST(w http.ResponseWriter, r *http.Request) {
	var periodo model.Periodo
	sess := model.Instance(r)
        sPerid    :=  r.FormValue("id")
	perid,_   := atoi32(sPerid)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            periodo.Id = perid
            err := (&periodo).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No hay periodo", view.FlashError})
             }
	     lisAmt, _            := model.Amounts( periodo.Final)
             v                     := view.New(r)
             v.Name                 = "report/condapto"
	     v.Vars["token"]        =  csrfbanana.Token(w, r, sess)
             v.Vars["LisAmt"]       =  lisAmt
             v.Vars["Per"]          =  periodo
             v.Vars["Level"]        =  sess.Values["level"]
	     v.Render(w)
         }else{
	  http.Redirect(w, r, "/cuota/list", http.StatusFound)
	 }
 }
// ---------------------------------------------------
