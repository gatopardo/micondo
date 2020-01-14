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

// TipoGET despliega la pagina del tipo
func TipoGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisTipos, _ := model.Tipos()
	v := view.New(r)
	v.Name = "tipo/tipo"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["LisTipos"] =  lisTipos
	v.Render(w)
 }
// ---------------------------------------------------
// ---------------------------------------------------
// POST procesa la forma enviada con los datos
func TipoPOST(w http.ResponseWriter, r *http.Request) {
        var tipo model.Tipo
	sess := model.Instance(r)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
	    if validate, missingField := view.Validate(r, []string{"codigo"}); !validate {
                 sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
                 sess.Save(r, w)
                 TipoGET(w, r)
                 return
	      }
              tipo.Codigo           = r.FormValue("codigo")
              tipo.Descripcion      = r.FormValue("descripcion")
              err := (&tipo).TipoByCode()
              if err == model.ErrNoResult { // Exito: no hay tipo creado aun 
                  ex := (&tipo).TipoCreate()
	          if ex != nil {  // uyy como fue esto ? 
                     log.Println(ex)
                     sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                     return
	         } else {  // todo bien
                    sess.AddFlash(view.Flash{"Tipoo. creado: " +tipo.Codigo, view.FlashSuccess})
	         }
              }
          }
          sess.Save(r, w)
	  http.Redirect(w, r, "/categoria/list/1", http.StatusFound)
  }

// ---------------------------------------------------
// TipoUpGET despliega la pagina del usuario
func TipoUpGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        var tipo model.Tipo
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	id,_ := atoi32(params.ByName("id"))
	SPag   := params.ByName("pg")
        path   :=  fmt.Sprintf("/categoria/list/%s", SPag)
        tipo.Id = id
	 err := (&tipo).TipoById()
	if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No tenemos tipo.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "tipo/tipoupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["Tipo"] = tipo
        v.Render(w)
   }
   // ---------------------------------------------------
 func   getTipoFormUp(r * http.Request)(st string){
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

// ---------------------------------------------------
// TipoUpPOST procesa la forma enviada con los datos
func TipoUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var tipo model.Tipo
	// Get session
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
	SPag        := params.ByName("pg")
        Id,_        := atoi32(SId)
        tipo.Id      = Id
        path        :=  fmt.Sprintf("/categoria/list/%s", SPag)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            sr          :=  fmt.Sprintf(" where tipo.id = %s ", SId)
            sini        :=  "update tipos set "
            st          :=  getTipoFormUp(r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No hay actualizacion solicitada", view.FlashSuccess})
            } else {
             st = sini + st + sr
             err =  tipo.Update(st)
             if err == nil{
                 sess.AddFlash(view.Flash{"Tipo actualizado exitosamente para: " +tipo.Codigo, view.FlashSuccess})
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
// TipoLisGET displays the tipo page
func TipoLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
//        var params httprouter.Params
//        params = context.Get(r, "params").(httprouter.Params)
//        SPg   := params.ByName("pg")
//        pg,_ := atoi32(SPg)
//        posact = int(pg)
//        offset = posact  - 1
//        offset = offset * limit
//        TotalCount = model.TipoCount()
        lisTipos, err := model.Tipos()
        if err != nil {
           log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Tipos.", view.FlashError})
            sess.Save(r, w)
         }
	v                 := view.New(r)
	v.Name             = "tipo/tipolis"
	v.Vars["token"]    = csrfbanana.Token(w, r, sess)
//        numberOfBtns      :=  getNumberOfButtonsForPagination(TotalCount, limit)
//        sliceBtns         :=  createSliceForBtns(numberOfBtns, posact)
//        v.Vars["slice"]    =  sliceBtns
        v.Vars["current"]  =  posact
        v.Vars["LisTipo"]   = lisTipos
        v.Vars["Level"]    =  sess.Values["level"]
	v.Render(w)
 }
// ---------------------------------------------------
// TipoDeleteGET handles the tipo deletion
 func TipoDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        var tipo model.Tipo
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        tipo.Id = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/categoria/list/%s", SPag)
        err := (&tipo).TipoById()
        if err != nil {
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro no tenemos tipo.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, path, http.StatusFound)
            return
        }
	v := view.New(r)
	v.Name = "tipo/tipodelete"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["Tipo"]      =  tipo
        v.Vars["level"]     =  sess.Values["level"]
	v.Render(w)
  }
// ---------------------------------------------------
// TipoDeletePOST handles the tipo deletion
 func TipoDeletePOST(w http.ResponseWriter, r *http.Request) {
        var err error
        sess := model.Instance(r)
        var tipo model.Tipo
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_ := atoi32(params.ByName("id"))
        tipo.Id = Id
	SPag        := params.ByName("pg")
        path        :=  fmt.Sprintf("/categoria/list/%s", SPag)
        action        := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = tipo.Delete()
            if err != nil {
                log.Println(err)
                sess.AddFlash(view.Flash{"Error no posible. Auxilio.", view.FlashError})
            } else {
                sess.AddFlash(view.Flash{"Tipo. borrado!", view.FlashSuccess})
            }
            sess.Save(r, w)
        }
	http.Redirect(w, r, path, http.StatusFound)

  }
// ---------------------------------------------------

