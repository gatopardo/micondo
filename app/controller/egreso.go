package controller

import (
	"log"
	"net/http"
        "fmt"
        "strings"
        "time"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
// ---------------------------------------------------
// EgrePerGET despliega formulario escoger periodo
func EgrePerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "egreso/egresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// EgrePerPOST procesa la forma enviada con periodo
func EgrePerPOST(w http.ResponseWriter, r *http.Request) {
        var egres model.EgresoN
        var period model.Periodo
        var  err  error
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            egres.PeriodId,  _   =  atoi32(r.FormValue("periodId"))
            period.Id            =  egres.PeriodId
            _                    =  (&period).PeriodById()
            egres.Period         =  period.Inicio
            var lisTipo []model.Tipo
            var lisEgre []model.EgresoN
            lisTipo,  err        = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisEgre, _          = (&egres).EgresPer()
	    v                  := view.New(r)
	    v.Name              = "egreso/egresoreg"
            v.Vars["token"]     = csrfbanana.Token(w, r, sess)
            v.Vars["Egreso"]    = egres
            v.Vars["LisTip"]    = lisTipo
            v.Vars["LisEgres"]  = lisEgre
            v.Vars["Level"]     =  sess.Values["level"]
            v.Render(w)
        }
	http.Redirect(w, r, "/egreso/list", http.StatusFound)
 }
// ---------------------------------------------------
 func getEgreData(c *  model.EgresoN, r *http.Request)(err error){
           formato         := "2006/01/02"
           c.Period, _     = time.Parse(formato,r.FormValue("period"))
           c.TipoId, _     = atoi32(r.FormValue("tipId"))
           c.Fecha, _      =  time.Parse(layout,r.FormValue("fecha"))
	   var nro int64
           nro, err        = money2int64(r.FormValue("amount"))
           if err == nil {
                 c.Amount   =  nro
            }
           c.Descripcion   =  r.FormValue("descripcion")
       return
   }
// ---------------------------------------------------
// EgreRegPOST despliega formulario crear egreso
func EgreRegPOST(w http.ResponseWriter, r *http.Request) {
        var egres   model.EgresoN
        var period  model.Periodo
        var err  error
	sess   := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getEgreData(&egres, r)
           period.Inicio       =  egres.Period
           err                 =  (&period).PeriodByCode()
           egres.PeriodId       =   period.Id
           err                 =  (&egres).EgresCreate()
           if err != nil {  // uyy como fue esto ? 
               log.Println(err)
               fmt.Println(err)
               sess.AddFlash(view.Flash{"Error guardando Egreso.", view.FlashError})
               return
           } else {  // todo bien
                sess.AddFlash(view.Flash{"Egreso. creada: " , view.FlashSuccess})
           }

            var lisTipo []model.Tipo
            var lisEgre []model.EgresoN
            lisTipo, err  = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisEgre,err           = (&egres).EgresPer()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay egresos ", view.FlashError})
            }
            v                   := view.New(r)
            v.Name               = "egreso/egresoreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
            v.Vars["Egreso"]     = egres
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisEgres"]   = lisEgre
            v.Vars["Level"]      =  sess.Values["level"]
	    v.Render(w)
        }
	http.Redirect(w, r, "/egreso/list", http.StatusFound)
 }
// ---------------------------------------------------
// EgreUpGET despliega la pagina del usuario
func EgreUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var egres model.EgresoN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	Sid         := params.ByName("id")
	id,_        := atoi32(Sid)
        path        := "/egreso/list"
        egres.Id = id
	err := (&egres).EgresById()
	if err != nil { // Si no existe Egreso
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta egreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "egreso/egresodelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Egre"]       = egres
        v.Vars["Level"]       =  sess.Values["level"]
        v.Render(w)
   }

// ---------------------------------------------------
 func   getEgreFormUp(r * http.Request)(st string){
        var sf string
        var nr  int64
        var sup []string
        if r.FormValue("ckegreso") == "true" {
	     nr, _  =  money2int64(  r.FormValue("egreso") )
             sf     =  fmt.Sprintf( " egreso = '%d' ", nr )
	     sup    = append(sup, sf)
           }
        if r.FormValue("ckamount") == "true" {
             nr, _  =  money2int64(  r.FormValue("amount") )
             sf     =  fmt.Sprintf( " amount = '%d' ", nr )
	     sup    = append(sup, sf)
           }

         if len(sup) > 0 {
              st =  strings.Join(sup, ", ")
          }
         return
  }
// ---------------------------------------------------
// EgreUpPOST procesa la forma enviada con los datos
func EgreUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var egres model.EgresoN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        egres.Id      = Id
        path        :=  "/egreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sr          :=  fmt.Sprintf(" where egresos.id = %s ", SId)
            sini        :=  "update egresos set "
            st          :=  getEgreFormUp(r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st    = sini + st + sr
             err   =  egres.EgresUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Egreso actualizada exitosamente : " , view.FlashSuccess})
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
// EgreLis displays the egres page
func EgreLis(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per  model.Periodo
	sess            := model.Instance(r)
        lisPeriod,err    := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
        if r.Method == "GET" {
            Id = lisPeriod[0].Id
        }else{
            Id,_             = atoi32(r.FormValue("id"))
        }
	per.Id               = Id
	err  = (&per).PeriodById()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error con Periodo.", view.FlashError})
            sess.Save(r, w)
         }

// fmt.Println("List Egreso ", Id)
        lisEgre, err         := model.EgresLim(Id)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Egresos.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "egreso/egresolis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisEgre"]    = lisEgre
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// EgreDeleteGET handles the note deletion
 func EgreDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var egres model.EgresoN
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
	Sid         := params.ByName("id")
	id,_        := atoi32(Sid)
        path        :=  "/egreso/list"
        egres.Id     = id
	err         := (&egres).EgresById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No hay egreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "egreso/egresodelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Egre"]        = egres
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
  }

// ---------------------------------------------------
// EgreDeletePOST procesa la forma enviada con los datos
func EgreDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var egres model.Egreso
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        egres.Id      = Id
        path        :=  "/egreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = egres.EgresDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Egreso borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
