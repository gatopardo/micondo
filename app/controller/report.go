package controller

import (
      "log"
	"net/http"
        "fmt"
        "strings"
        "time"

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
	v                  := view.New(r)
	v.Name              = "report/rptmail"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// MailSendPOST procesa la forma enviada con contenido
func MailSendPOST(w http.ResponseWriter, r *http.Request) {
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            dateLayout := "2006-01-02"
            timeLayout := "15:04:05"
            tim        := time.Now()
            fec        := tim.Format(dateLayout)
            hour       := tim.Format(timeLayout)
            stm        := "Fecha "+fec +" Hora : " + hour 
fmt.Println("MailSendPost ", stm)
	    tema            := r.FormValue("tema")
	    content         := stm +"\n" + r.FormValue("content")
            lisPers, err    :=  model.Persons()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay usuarios ", view.FlashError})
		   fmt.Println("Error correo ", err)
	       http.Redirect(w, r, "/email", http.StatusFound)
            }
fmt.Println("MailSendPost lon  ", len(lisPers))
             for _,person := range lisPers{
                 to := person.Email
    fmt.Printf(" %s | %s\n", person.Fname ,to)
	        err = email.SendEmail(to, tema,content);
                if err != nil {
                   sess.AddFlash(view.Flash{"Error enviando ", view.FlashError})
		   fmt.Println("Error Enviando", err)
                }
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
	var peridi model.Periodo
	var peridf model.Periodo
	var err error
	sess := model.Instance(r)
        uid, ok       := sess.Values["id"].(uint32)
	if ! ok {
             log.Println("No uint32 value in session")
	}
        sPeridf    :=  r.FormValue("idf")
	fperid,_   := atoi32(sPeridf)
        sPeridi    :=  r.FormValue("idi")
	iperid,_   := atoi32(sPeridi)
	action    := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             pers, apt, err = model.ApartaByUserId(uid)
	     if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No apto", view.FlashError})
	     }
            peridf.Id = fperid
            err := (&peridf).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No hay periodo", view.FlashError})
             }
            peridi.Id = iperid
            err = (&peridi).PeriodById()
            if err != nil {
	        log.Println(err)
                sess.AddFlash(view.Flash{"No hay periodo", view.FlashError})
             }
	     lisPaym, _            := model.Payments(apt.Id, peridf.Inicio, peridi.Inicio)
	     value                 := lisPaym[len(lisPaym) - 1].Balance
             v                     := view.New(r)
             v.Name                 = "report/rptapto"
	     v.Vars["token"]        = csrfbanana.Token(w, r, sess)
             v.Vars["Apt"]          = apt
             v.Vars["Pers"]         = pers
             v.Vars["Perid"]        = peridf
	     v.Vars["Valor"]        = value
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
