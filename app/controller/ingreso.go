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
// IngrePerGET despliega formulario escoger periodo
func IngrePerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.PeriodsI()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "ingreso/ingresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// IngrePerPOST procesa la forma enviada con periodo
func IngrePerPOST(w http.ResponseWriter, r *http.Request) {
        var ingres model.IngresoN
        var period model.Periodo
        var  err  error
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            ingres.PeriodId,  _   =  atoi32(r.FormValue("periodId"))
            period.Id            =  ingres.PeriodId
            _                    =  (&period).PeriodById()
            ingres.Period         =  period.Inicio
            var lisTipo []model.Tipo
            var lisIngre []model.IngresoN
            lisTipo,  err        = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisIngre, _          = model.IngresPer(ingres.PeriodId)
	    v                  := view.New(r)
	    v.Name              = "ingreso/ingresoreg"
            v.Vars["token"]     = csrfbanana.Token(w, r, sess)
            v.Vars["Ingreso"]    = ingres
            v.Vars["LisTip"]    = lisTipo
            v.Vars["LisIngres"]  = lisIngre
            v.Vars["Level"]     =  sess.Values["level"]
            v.Render(w)
        }
	http.Redirect(w, r, "/ingreso/list", http.StatusFound)
 }
// ---------------------------------------------------
 func getIngreData(c *  model.IngresoN, r *http.Request)(err error){
           formato         := "2006/01/02"
           c.Period, _     = time.Parse(formato,r.FormValue("period"))
           c.TipoId, _     = atoi32(r.FormValue("tipId"))
           c.Fecha, _      =  time.Parse(layout,r.FormValue("fecha"))
	   var nro int64
           nro, err       = money2int64(r.FormValue("amount"))
           if err == nil {
		   c.Amount   =  nro
            }
           c.Descripcion   =  r.FormValue("descripcion")
       return
   }
// ---------------------------------------------------
// IngreRegPOST despliega formulario crear ingreso
func IngreRegPOST(w http.ResponseWriter, r *http.Request) {
        var ingres   model.IngresoN
        var period  model.Periodo
        var err  error
	sess   := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getIngreData(&ingres, r)
           period.Inicio       =  ingres.Period
           err                 =  (&period).PeriodByCode()
           ingres.PeriodId       =   period.Id
           err                 =  (&ingres).IngresCreate()
           if err != nil {  // uyy como fue esto ? 
               log.Println(err)
               fmt.Println(err)
               sess.AddFlash(view.Flash{"Error guardando Ingreso.", view.FlashError})
               return
           } else {  // todo bien
                sess.AddFlash(view.Flash{"Ingreso. creada: " , view.FlashSuccess})
           }

            var lisTipo []model.Tipo
            var lisIngre []model.IngresoN
            lisTipo, err  = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisIngre, _           = model.IngresPer(period.Id)
            v                   := view.New(r)
            v.Name               = "ingreso/ingresoreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
            v.Vars["Ingreso"]     = ingres
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisIngres"]   = lisIngre
            v.Vars["Level"]      =  sess.Values["level"]
	    v.Render(w)
        }
	http.Redirect(w, r, "/ingreso/list", http.StatusFound)
 }
// ---------------------------------------------------
// IngreUpGET despliega la pagina del usuario
func IngreUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var ingres model.IngresoN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	id,_   := atoi32(params.ByName("id"))
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/ingreso/list/%s", SPag)
        ingres.Id = id
	err := (&ingres).IngresById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "ingreso/ingreso"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Ingre"]       = ingres
        v.Vars["Level"]       =  sess.Values["level"]
        v.Render(w)
   }

// ---------------------------------------------------
 func   getIngreFormUp(r * http.Request)(st string){
        var sf string
        var nr  int64
        var sup []string
        if r.FormValue("ckingreso") == "true" {
	     nr, _  =  money2int64(  r.FormValue("ingreso") )
             sf  =  fmt.Sprintf( " ingreso = '%d' ", nr )
	     sup = append(sup, sf)
           }
        if r.FormValue("ckamount") == "true" {
	     nr, _  =  money2int64(  r.FormValue("amount") )
             sf  =  fmt.Sprintf( " amount = '%d' ", nr )
	     sup = append(sup, sf)
           }

         if len(sup) > 0 {
              st =  strings.Join(sup, ", ")
          }
         return
  }
// ---------------------------------------------------
// IngreUpPOST procesa la forma enviada con los datos
func IngreUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var ingres model.IngresoN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	SPag        := params.ByName("pg")
        Id,_        := atoi32(SId)
        ingres.Id      = Id
        path        :=  fmt.Sprintf("/ingreso/list/%s", SPag)
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sr          :=  fmt.Sprintf(" where ingresos.id = %s ", SId)
            sini        :=  "update ingresos set "
            st          :=  getIngreFormUp(r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st    = sini + st + sr
             err   =  ingres.IngresUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Ingreso actualizada exitosamente : " , view.FlashSuccess})
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
// IngreLis displays the ingres page
func IngreLis(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per model.Periodo
	sess            := model.Instance(r)
        lisPeriod,err    := model.PeriodsI()
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
// fmt.Println("IngreLis lon ", len(lisPeriod), " id ", Id, " ", r.Method)
	per.Id                 = Id
	err = (&per).PeriodById();
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
        lisIngre, err         := model.IngresLim(Id)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Ingresos.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "ingreso/ingresolis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["LisIngre"]    = lisIngre
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// UserDeleteGET handles the note deletion
 func IngreDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var ingres model.IngresoN
/*
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
*/
	id,_        := atoi32(r.FormValue("id"))
        path        :=  "/ingreso/list"
        ingres.Id   = id
	err         := (&ingres).IngresById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "ingreso/ingresodelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Ingre"]        = ingres
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
  }

// ---------------------------------------------------
// IngreUpPOST procesa la forma enviada con los datos
func IngreDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var ingres model.Ingreso
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        ingres.Id    = Id
        path        :=  "/ingreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = ingres.IngresDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Ingreso borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
