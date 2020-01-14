package controller

import (
	"log"
	"net/http"
        "time"
//        "strconv"
        "strings"
        "fmt"

	"app/model"
	"app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
const (
       layout = "2006-01-02"
 )

// ---------------------------------------------------

// PeriodGET despliega la pagina del periodo
func PeriodGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriods, _ := model.Periods()
	v := view.New(r)
	v.Name = "periodo/period"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriods"] =  lisPeriods
//      Refill any form fields
// view.Repopulate([]string{"name"}, r.Form, v.Vars)
	v.Render(w)
 }
// ---------------------------------------------------
// POST procesa la forma enviada con los datos
func PeriodPOST(w http.ResponseWriter, r *http.Request) {
        var period model.Periodo
	sess := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"inicio"}); !validate {
              sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
              sess.Save(r, w)
              PeriodGET(w, r)
              return
            }
            period.Inicio, _    = time.Parse(layout,r.FormValue("inicio"))
            period.Final, _     = time.Parse(layout,r.FormValue("final"))
            err := (&period).PeriodByCode()
            if err == model.ErrNoResult { // Exito:  no hay perioo creado aun 
                 ex := (&period).PeriodCreate()
                 if ex != nil {  // uyy como fue esto ? 
                     log.Println(ex)
                     sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                    return
                 } else {  // todo bien
                sess.AddFlash(view.Flash{"Periodo. creado: " + period.Inicio.Format(layout), view.FlashSuccess})
	         }
            }
         }
        sess.Save(r, w)
	http.Redirect(w, r, "/period/list/1", http.StatusFound)
 }

// ---------------------------------------------------
// PeriodUpGET despliega la pagina del usuario
func PeriodUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var period model.Periodo
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
        period.Id = id
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/period/list/%s", SPag)
        err    := (&period).PeriodById()
	if err != nil { // Si no existe el periodo
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No tenemos periodo.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v               := view.New(r)
	v.Name           = "periodo/periodupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["Period"] = period
//	view.Repopulate([]string{"name"}, r.Form, v.Vars)
        v.Render(w)
   }
// ---------------------------------------------------
 func   getPeriodFormUp(r * http.Request)(st string){
        var sf string
        var sup []string
        if r.FormValue("ckinicio") == "true" {
	     sf  =  fmt.Sprintf( " inicio = '%s' ", r.FormValue("inicio") )
	     sup = append(sup, sf)
        }
        if r.FormValue("ckfinal") == "true" {
	     sf  =  fmt.Sprintf( " final = '%s' ", r.FormValue("final") )
	     sup = append(sup, sf)
        }
        if len(sup) > 0 {
              st =  strings.Join(sup, ", ")
        }
        return
  }
// ---------------------------------------------------
// PeriodUpPOST procesa la forma enviada con los datos
func PeriodUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var period model.Periodo
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SPag        := params.ByName("pg")
	SId          := params.ByName("id")
	period.Id, _ = atoi32(SId)
        path        :=  fmt.Sprintf("/period/list/%s", SPag)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sr          :=  fmt.Sprintf(" where periods.id = %s ", SId)
            sini        :=  "update periods set "
            st          :=  getPeriodFormUp(r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st = sini + st + sr
             err =  period.PeriodUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Periodo actualizado exitosamente ", view.FlashSuccess})
             } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error ocurrio actualizando.", view.FlashError})
	     }
	sess.Save(r, w)
           }
	}
	http.Redirect(w, r, path, http.StatusFound)
     }
//------------------------------------------------
// PeriodLisGET displays the aparta page
func PeriodLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var params httprouter.Params
        params           = context.Get(r, "params").(httprouter.Params)
        SPg             := params.ByName("pg")
        pg,_            := atoi32(SPg)
        posact           = int(pg)
        offset           = posact  - 1
        offset           = offset * limit
        TotalCount       = model.PeriodCount()
        lisPeriods, err := model.PeriodLim(limit, offset)
        if err != nil {
           log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Periodos.", view.FlashError})
            sess.Save(r, w)
         }
	v                    := view.New(r)
	v.Name                = "periodo/periodlis"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        numberOfBtns         :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns            :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]       =  sliceBtns
        v.Vars["current"]     =  posact	
        v.Vars["LisPeriod"]   = lisPeriods
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
 }
//-----------------------------------------------------------------------
// PeriodDeleteGET handles the note deletion
 func PeriodDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var period model.Periodo
        var params httprouter.Params
        params     = context.Get(r, "params").(httprouter.Params)
        Id,_      := atoi32(params.ByName("id"))
        period.Id  = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/period/list/%s", SPag)
        err := (&period).PeriodById()
        if err != nil {
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro no tenemos periodo.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v                  := view.New(r)
	v.Name              = "periodo/perioddelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Period"]   =  period
        v.Vars["level"]     =  sess.Values["level"]
	v.Render(w)
  }
//-----------------------------------------------------------------------
//-----------------------------------------------------------------------
// PeriodDeletePOST handles the note deletion
 func PeriodDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        sess := model.Instance(r)
        var period model.Periodo
        var params httprouter.Params
        params     = context.Get(r, "params").(httprouter.Params)
//        Id,_      := atoi32(params.ByName("id"))
//        period.Id  = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/period/list/%s", SPag)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           err  = period.PeriodDelete()
           if err != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
           } else {
                sess.AddFlash(view.Flash{"Periodo. borrado!", view.FlashSuccess})
           }
                sess.Save(r, w)
       }
	http.Redirect(w, r, path, http.StatusFound)
  }
//-----------------------------------------------------------------------


