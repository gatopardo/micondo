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
// CuotPerGET despliega formulario escoger periodo
func CuotPerGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "cuota/cuotper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
 //----------------------------------------------------
// CuotPerPOST procesa la forma enviada con periodo
func CuotPerPOST(w http.ResponseWriter, r *http.Request) {
        var cuot model.CuotaN
        var period model.Periodo
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            var lisTipo []model.Tipo
            var lisCuot []model.CuotaN
            cuot.PeriodId,  _   =  atoi32(r.FormValue("id"))
            period.Id           =  cuot.PeriodId
            _                   =  (&period).PeriodById()
            cuot.Period         =  period.Inicio
            lisApts, err       :=  model.Apts()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay aptos ", view.FlashError})
            }
            lisTipo,  err        = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
             }
             lisCuot, _          = (&cuot).CuotsPer()

	     v                  := view.New(r)
	     v.Name              = "cuota/cuotreg"
             v.Vars["token"]     = csrfbanana.Token(w, r, sess)
             v.Vars["Cuota"]     = cuot
             v.Vars["LisApt"]    = lisApts
             v.Vars["LisTip"]    = lisTipo
             v.Vars["LisCuots"]   = lisCuot
             v.Vars["Level"]     =  sess.Values["level"]
             v.Render(w)
        }
	http.Redirect(w, r, "/cuota/list", http.StatusFound)
 }
// ---------------------------------------------------
 func getFormCuot(c *  model.CuotaN, r *http.Request)(err error){
           formato        :=  "2006/01/02"
           c.PeriodId, _   =  atoi32(r.FormValue("periodId"))
           c.ApartaId, _   =  atoi32(r.FormValue("aptId"))
           c.TipoId, _     =  atoi32(r.FormValue("tipId"))
           c.Fecha, _      =  time.Parse(formato,r.FormValue("fecha"))
           unr, err       :=  money2int64(r.FormValue("amount"))
           if err == nil {
                 c.Amount   =  unr
            }
       return
   }
// ---------------------------------------------------
// CuotRegPOST despliega formulario crear cuota
func CuotRegPOST(w http.ResponseWriter, r *http.Request) {
        var cuot   model.CuotaN
        var period  model.Periodo
        var err  error
//fmt.Println("CuotRegPost 1")
	sess   := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           getFormCuot(&cuot, r)
           period.Inicio       =  cuot.Period
           err                 =  (&period).PeriodByCode()
           cuot.PeriodId       =   period.Id
           err                 =  (&cuot).CuotCreate()
           if err != nil {  // uyy como fue esto ? 
               log.Println(err)
               sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
               return
           } else {  // todo bien
                sess.AddFlash(view.Flash{"Cuota. creada: " , view.FlashSuccess})
           }

            var lisApto []model.Aparta
            var lisTipo []model.Tipo
            var lisCuot []model.CuotaN
            lisApto, err  = model.Apts()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay apartas ", view.FlashError})
            }
            lisTipo, err  = model.Tipos()
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
            lisCuot, _           = (&cuot).CuotsPer()

            v                   := view.New(r)
            v.Name               = "cuota/cuotreg"
            v.Vars["token"]      = csrfbanana.Token(w, r, sess)
            v.Vars["Cuota"]      = cuot
            v.Vars["LisApt"]     = lisApto
            v.Vars["LisTip"]     = lisTipo
            v.Vars["LisCuots"]   = lisCuot
            v.Vars["Level"]      =  sess.Values["level"]
	    v.Render(w)
        }
	http.Redirect(w, r, "/cuota/list", http.StatusFound)
 }
// ---------------------------------------------------
// CuotUpGET despliega la pagina del usuario
func CuotUpGET(w http.ResponseWriter, r *http.Request) {
        var lisTipo []model.Tipo
	sess := model.Instance(r)
        var cuot model.CuotaN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	id,_   := atoi32(params.ByName("id"))
        cuot.Id = id
        path   :=  "/cuota/list"
        lisApts, err       :=  model.Apts()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay aptos ", view.FlashError})
        }
        lisTipo,  err        = model.Tipos()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
        }
	err = (&cuot).CuotById()
	if err != nil { // Si no existe cuota
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta cuota.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "cuota/cuotupdate"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Cuot"]       = cuot
        v.Vars["LisApt"]    = lisApts
        v.Vars["LisTip"]    = lisTipo
        v.Vars["Level"]       =  sess.Values["level"]
        v.Render(w)
   }

// ---------------------------------------------------
 func   getCuotFormUp(c1,c2 model.CuotaN, r * http.Request)(stUp string){
        var sf string
	var sup  []string
        formato        :=  "2006-01-02"

	if c1.PeriodId != c2.PeriodId {
             sf  =  fmt.Sprintf( " period_id = %d ", c2.PeriodId )
	     sup = append(sup, sf)
	}
	if c1.ApartaId != c2.ApartaId {
             sf  =  fmt.Sprintf( " aparta_id = %d ", c2.ApartaId )
	     sup = append(sup, sf)
	}
	if c1.TipoId != c2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", c2.TipoId )
	     sup = append(sup, sf)
	}
	if c1.Fecha != c2.Fecha {
             sf  =  fmt.Sprintf( " fecha = '%s' ", c2.Fecha.Format(formato) )
	     sup = append(sup, sf)
	}
	if c1.Amount != c2.Amount {
             sf  =  fmt.Sprintf( " amount = %d ", c2.Amount )
	     sup = append(sup, sf)
	}
       lon := len(sup)
       if lon  > 0 {
            sini        :=  "update cuotas set "
            stUp =  strings.Join(sup, ", ")
            sr          :=  fmt.Sprintf(" where cuotas.id = %d ", c1.Id)
             stUp = sini + stUp + sr
       }

         return
  }
// ---------------------------------------------------
// CuotUpPOST procesa la forma enviada con los datos
func CuotUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var c, cuot model.CuotaN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        cuot.Id      = Id
        path        :=  "/cuota/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    err = (&cuot).CuotById()
	    if err != nil { // Si no existe cuota
                  sess.AddFlash(view.Flash{"Es raro. No esta cuota.", view.FlashError})
            }
	    getFormCuot(&c,r)
            st          :=  getCuotFormUp(cuot, c, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
             err   =  cuot.CuotUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Cuota actualizada exitosamente : " , view.FlashSuccess})
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
// CuotLisGET displays the cuot page
func CuotLisGET(w http.ResponseWriter, r *http.Request) {
	sess            := model.Instance(r)
        lisPeriod,err   := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "cuota/cuotper"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"]  = lisPeriod
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// CuotLis displays the cuot page
func CuotLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per model.Periodo
	sess            := model.Instance(r)
      lisPeriod,err    := model.PeriodsI()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
            Id,_             = atoi32(r.FormValue("id"))
	    per.Id           = Id
	    (&per).PeriodById()
        lisCuot, err         := model.CuotLim(Id)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Cuotas.", view.FlashError})
            sess.Save(r, w)
         }
	v                   := view.New(r)
	v.Name               = "cuota/cuotlis"
	v.Vars["token"]      = csrfbanana.Token(w, r, sess)
	v.Vars["Per"]        = per
        v.Vars["LisPeriod"]    = lisPeriod
        v.Vars["LisCuot"]    = lisCuot
        v.Vars["Level"]      =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// CuotDeleteGET handles the note deletion
 func CuotDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var cuot model.CuotaN
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	id,_        := atoi32(SId)
        path        :=  "/cuota/list"
        cuot.Id      = id
	err         := (&cuot).CuotById()
	if err != nil { // Si no existe cuota
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta cuota.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "cuota/cuotdelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Cuot"]        = cuot
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
  }

// ---------------------------------------------------
// CuotDeletePOST procesa la forma enviada con los datos
func CuotDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var cuot model.Cuota
	sess := model.Instance(r)
        var params httprouter.Params
        params       = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        cuot.Id      = Id
        path        :=  "/cuota/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = cuot.CuotDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Cuota borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
