package controller

import (
	"log"
	"net/http"
        "strings"
        "fmt"

	"app/model"
	"app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )

//      Refill any form fields
// view.Repopulate([]string{"name"}, r.Form, v.Vars)
// ---------------------------------------------------

// AptGET despliega la pagina del apto
func AptGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisApts, _ := model.Apts()
	v := view.New(r)
	v.Name = "aparta/apt"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["LisApts"] =  lisApts
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// POST procesa la forma enviada con los datos
func AptPOST(w http.ResponseWriter, r *http.Request) {
        var apt model.Aparta
	sess := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"codigo"}); !validate {
                 sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
                 sess.Save(r, w)
                 AptGET(w, r)
                 return
	      }
              apt.Codigo           = r.FormValue("codigo")
              apt.Descripcion      = r.FormValue("descripcion")
              err := (&apt).AptByCode()
              if err == model.ErrNoResult { // Exito: no hay apartamento creado aun 
                  ex := (&apt).AptCreate()
	          if ex != nil {  // uyy como fue esto ? 
                     log.Println(ex)
                     sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                     return
	         } else {  // todo bien
                    sess.AddFlash(view.Flash{"Apto. creado: " +apt.Codigo, view.FlashSuccess})
	         }
              }
          }
          sess.Save(r, w)
	  http.Redirect(w, r, "/apto/list/1", http.StatusFound)
  }

// ---------------------------------------------------
// AptUpGET despliega la pagina del usuario
func AptUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var apt model.Aparta
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/apto/list/%s", SPag)
        apt.Id = id
        err := (&apt).AptById()
	if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No tenemos apartamento.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "aparta/aptupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["Apto"] = apt
        v.Render(w)
   }
// ---------------------------------------------------
 func   getAptFormUp(r * http.Request)(st string){
        var sf string
        st = ""
        var sup []string
        if r.FormValue("ckdescrip") == "true" {
	     sf  =  fmt.Sprintf( " descrip = '%s' ", r.FormValue("descrip") )
	     sup = append(sup, sf)
        }
        if len(sup) > 0 {
              st =  strings.Join(sup, ", ")
        }
        return
  }
// ---------------------------------------------------
// AptUpPOST procesa la forma enviada con los datos
func AptUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var apt model.Aparta
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	SPag        := params.ByName("pg")
        Id,_        := atoi32(SId)
        apt.Id      = Id
        path        :=  fmt.Sprintf("/apto/list/%s", SPag)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sr          :=  fmt.Sprintf(" where apt.id = %s ", SId)
            sini        :=  "update apartas set "
            st          :=  getAptFormUp(r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st = sini + st + sr
             err =  apt.AptUpdate(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Apto actualizado exitosamente para: " +apt.Codigo, view.FlashSuccess})
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
// AptLisGET displays the aparta page
func AptLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
//        var params httprouter.Params
//        params     = context.Get(r, "params").(httprouter.Params)
//        SPg       := params.ByName("pg")
//        pg,_      := atoi32(SPg)
//        posact     = int(pg)
//        offset     = posact  - 1
//        offset     = offset * limit
//        TotalCount = model.AptCount()
        lisApts, err := model.Apts()
        if err != nil {
           log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Aptos.", view.FlashError})
            sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "aparta/aptlis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
 //       numberOfBtns      :=  getNumberOfButtonsForPagination(TotalCount, limit)
//        sliceBtns         :=  createSliceForBtns(numberOfBtns, posact)
//        v.Vars["slice"]    =  sliceBtns
        v.Vars["current"]  =  posact
        v.Vars["LisApt"]   = lisApts
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// AptDeleteGET handles the apto deletion
 func AptDeleteGET(w http.ResponseWriter, r *http.Request) {
        sess := model.Instance(r)
        var apt model.Aparta
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        apt.Id = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/apto/list/%s", SPag)
        err := (&apt).AptById()
        if err != nil {
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro no tenemos apto.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "aparta/aptdelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Apto"]      =  apt
        v.Vars["level"]     =  sess.Values["level"]
	v.Render(w)
  }
// ---------------------------------------------------
// AptDeletePOST handles the apto deletion
 func AptDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        sess := model.Instance(r)
        var apt model.Aparta
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        apt.Id = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/apto/list/%s", SPag)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = apt.AptDelete()
            if err != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
            } else {
                sess.AddFlash(view.Flash{"Apto. borrado!", view.FlashSuccess})
            }
            sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)

  }
// ---------------------------------------------------

