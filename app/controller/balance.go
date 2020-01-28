package controller

import (
	"log"
	"net/http"
        "fmt"
        "strings"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )
// ---------------------------------------------------
// ---------------------------------------------------
// BalanGET despliega formulario crear balances
func BalanGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
        var lisBalans []model.BalanceN
        lisBalans, err      =  model.Balans()
        if err != nil {
             sess.AddFlash(view.Flash{"No Balances ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "balance/balan"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisBalans"] = lisBalans
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Level"]     =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
 func getBalanData(b *  model.Balance, r *http.Request)(err error){
	   var nro int64
           b.PeriodId, _   = atoi32(r.FormValue("periodId"))
           nro, err        = money2int64(r.FormValue("amount"))
           if err == nil {
                 b.Amount  =  nro
                 nro,err   = money2int64(r.FormValue("cuota"))
            }
           if err == nil {
                 b.Cuota   =  nro
            }
       return
   }
// ---------------------------------------------------
// BalanPOST procesa la forma enviada con los datos
func BalanPOST(w http.ResponseWriter, r *http.Request) {
        var balan model.Balance
	sess          := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
           if validate, missingField := view.Validate(r, []string{"amount",  "cuota"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
               BalanGET(w, r)
               return
	     }
           getBalanData(&balan, r)
           err := (&balan).BalanByPeriod()
           if err == model.ErrNoResult { // Exito:  no hay usuario creado aun 
               ex := (&balan).BalanCreate()
               log.Println("Creating balance")
	       if ex != nil {  // uyy como fue esto ? 
                   log.Println(ex)
//   fmt.Println(ex)
                   sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                   sess.Save(r, w)
                   return
	       }
               sess.AddFlash(view.Flash{"Balance creado: " , view.FlashSuccess})
               sess.Save(r, w)
	   }
      }
	http.Redirect(w, r, "/balance/list/1", http.StatusFound)
 }

// ---------------------------------------------------
// BalanUpGET despliega la pagina del usuario
func BalanUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var balan model.BalanceN
	var params httprouter.Params
	params  = context.Get(r, "params").(httprouter.Params)
	id,_   := atoi32(params.ByName("id"))
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/balance/list/%s", SPag)
        balan.Id = id
	err := (&balan).BalanById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta balance.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "balance/balanupdate"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Balan"]       = balan
        v.Vars["Level"]       =  sess.Values["level"]
        v.Render(w)
   }

// ---------------------------------------------------
 func   getBalanFormUp(r * http.Request)(st string){
        var sf string
        var nro  int64
        var sup []string
        if r.FormValue("ckcuota") == "true" {
	     nro, _  =  money2int64(  r.FormValue("cuota") )
             sf  =  fmt.Sprintf( " cuota = '%d' ", nro )
	     sup = append(sup, sf)
           }
        if r.FormValue("ckamount") == "true" {
	     nro, _  =  money2int64(  r.FormValue("amount") )
             sf  =  fmt.Sprintf( " amount = '%d' ", nro )
	     sup = append(sup, sf)
           }

         if len(sup) > 0 {
              st =  strings.Join(sup, ", ")
          }
         return
  }
// ---------------------------------------------------
// BalanUpPOST procesa la forma enviada con los datos
func BalanUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var balan model.BalanceN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	SPag        := params.ByName("pg")
        Id,_        := atoi32(SId)
        balan.Id      = Id
        path        :=  fmt.Sprintf("/balance/list/%s", SPag)
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sr          :=  fmt.Sprintf(" where balances.id = %s ", SId)
            sini        :=  "update balances set "
            st          :=  getBalanFormUp(r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st    = sini + st + sr
             err   =  balan.BalanUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Balance actualizada exitosamente : " , view.FlashSuccess})
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
// BalanLisGET displays the balance page
func BalanLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
/*	
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)

        SPg := params.ByName("pg")
        Pg,_ := atoi32(SPg)
        posact = int(Pg)
        offset =  posact   -  1
        offset = offset * limit
        TotalCount =  model.BalansCount()
*/
        lisBalan, err := model.Balans()
        if err != nil {
//            fmt.Println(err)
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Usuarios.", view.FlashError})
            sess.Save(r, w)
         }
	v := view.New(r)
	v.Name             = "balance/balanlis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
/*
        numberOfBtns      :=  getNumberOfButtonsForPagination(TotalCount, limit)
        sliceBtns         :=  createSliceForBtns(numberOfBtns, posact)
        v.Vars["slice"]    =  sliceBtns
        v.Vars["current"]  =  posact
*/
        v.Vars["LisBalan"] = lisBalan
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }

//------------------------------------------------
// UserDeleteGET handles the note deletion
 func BalanDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var balan model.BalanceN
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
	id,_        := atoi32(params.ByName("id"))
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/balance/list/%s", SPag)
        balan.Id      = id
	err         := (&balan).BalanById()
	if err != nil { // Si no existe el usuario
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No esta balance.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "balance/balandelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Balan"]        = balan
        v.Vars["Level"]       =  sess.Values["level"]
	v.Render(w)
  }

// ---------------------------------------------------
// BalanUpPOST procesa la forma enviada con los datos
func BalanDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var balan model.Balance
	sess := model.Instance(r)
/*
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        balan.Id      = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/balance/list", SPag)
*/
        path        :=  "/balance/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
             err = balan.BalanDelete()
             if err != nil {
                 log.Println(err)
                 sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
              } else {
                  sess.AddFlash(view.Flash{"Balance borrado!", view.FlashSuccess})
              }
              sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)
 }
// ---------------------------------------------------
