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
	v.Vars["Title"]     =  "Ingreso"
        v.Vars["Action"]    =  "/ingreso/periodo/register"
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
            var lisTipo []model.Tipo
            var lisEgre []model.EgresoN
            egres.PeriodId,  _   =  atoi32(r.FormValue("id"))
            period.Id            =  egres.PeriodId
            _                    =  (&period).PeriodById()
            egres.Period         =  period.Inicio
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
            v.Render(w)
        }
	http.Redirect(w, r, "/egreso/list", http.StatusFound)
 }
// ---------------------------------------------------
<<<<<<< HEAD
 func getFormEgre(e *  model.EgresoN, r *http.Request)(err error){
           formato         :=  "2006/01/02"
           e.PeriodId, _    =  atoi32(r.FormValue("id"))
           e.TipoId, _      =  atoi32(r.FormValue("tipoId"))
           e.Fecha, _       =  time.Parse(formato,r.FormValue("fecha"))
           e.Descripcion    =  r.FormValue("descripcion")
=======
 func getFormEgre(c *  model.EgresoN, r *http.Request)(err error){
           formato         := "2006/01/02"
           c.Period, _     = time.Parse(formato,r.FormValue("period"))
           c.TipoId, _     = atoi32(r.FormValue("tipId"))
           c.Fecha, _      =  time.Parse(layout,r.FormValue("fecha"))
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c
	   var nro int64
           nro, err         = money2int64(r.FormValue("amount"))
           if err == nil {
                 e.Amount   =  nro
            }
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
           getFormEgre(&egres, r)
<<<<<<< HEAD
           period.Id       =  egres.PeriodId
           err                 =  (&period).PeriodById()

=======
           period.Inicio       =  egres.Period
           err                 =  (&period).PeriodByCode()
           egres.PeriodId       =   period.Id
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c
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

	lisTipo,  err        := model.Tipos()
<<<<<<< HEAD
        if err != nil {
             sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
        }
=======
            if err != nil {
                 sess.AddFlash(view.Flash{"No hay tipos ", view.FlashError})
            }
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c

	err = (&egres).EgresById()
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
<<<<<<< HEAD
        v.Vars["Title"]     =  "Actualizar Ingreso"
        v.Vars["Action"]    =  "/ingreso/update"
        v.Vars["LisTip"]      = lisTipo
=======
       v.Vars["LisTip"]      = lisTipo
        v.Vars["Level"]       =  sess.Values["level"]
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c
        v.Render(w)
   }

// ---------------------------------------------------
 func   getEgreFormUp(e1, e2 model.EgresoN, r * http.Request)(stUp string){
        var sf string
        var sup []string
<<<<<<< HEAD

	if e1.PeriodId != e2.PeriodId {
             sf  =  fmt.Sprintf( " period_id = %d ", e2.PeriodId )
	     sup = append(sup, sf)
	}
	if e1.TipoId != e2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", e2.TipoId )
	     sup = append(sup, sf)
	}

	if e1.Amount  != e2.Amount {
             sf  =  fmt.Sprintf( " amount = %d ", e2.Amount )
	     sup = append(sup, sf)
	}
        if e1.Fecha != e2.Fecha {
             sf  =  fmt.Sprintf( " fecha = '%s' ", e2.Fecha.Format(layout) )
=======
        formato        :=  "2006-01-02"

	if e1.PeriodId != e2.PeriodId {
             sf  =  fmt.Sprintf( " period_id = %d ", e2.PeriodId )
	     sup = append(sup, sf)
	}
	if e1.TipoId != e2.TipoId {
             sf  =  fmt.Sprintf( " tipo_id = %d ", e2.TipoId )
	     sup = append(sup, sf)
	}

	if e1.Amount  != e2.Amount {
             sf  =  fmt.Sprintf( " amount = '%d' ", e2.Amount )
	     sup = append(sup, sf)
	}
        if e1.Fecha != e2.Fecha {
             sf  =  fmt.Sprintf( " fecha = '%s' ", e2.Fecha.Format(formato) )
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c
	     sup = append(sup, sf)
	}

	if e1.Descripcion != e2.Descripcion {
             sf  =  fmt.Sprintf( " descripcion = %s ", e2.Descripcion )
	     sup = append(sup, sf)
	}
        lon := len(sup)
        if lon  > 0 {
            sini :=  "update egresos set "
<<<<<<< HEAD
	    now         := time.Now()
	    sf           =  fmt.Sprintf( " updated_at = '%s' ", now.Format(layout) )
            stUp  =  strings.Join(sup, ", ")
            sr   :=  fmt.Sprintf(" where egresos.id = %d ", e1.Id)
            stUp = sini + stUp + sf + sr
       }

=======
            stUp  =  strings.Join(sup, ", ")
            sr   :=  fmt.Sprintf(" where egresos.id = %d ", e1.Id)
            stUp = sini + stUp + sr
       }

       if len(sup) > 0 {
              stUp =  strings.Join(sup, ", ")
          }
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c
         return
  }
// ---------------------------------------------------
// EgreUpPOST procesa la forma enviada con los datos
func EgreUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var eg,egres model.EgresoN
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	SId         := params.ByName("id")
        Id,_        := atoi32(SId)
        egres.Id     = Id
        eg.Id     = Id
        path        :=  "/egreso/list"
        action      := r.FormValue("action")
        if ! (strings.Compare(action,"Cancelar") == 0) {
            err  = (&egres).EgresById()
	    if err != nil { // Si no existe cuota
<<<<<<< HEAD
                  sess.AddFlash(view.Flash{"Es raro. No esta egreso.", view.FlashError})
=======
                  sess.AddFlash(view.Flash{"Es raro. No esta ingreso.", view.FlashError})
>>>>>>> f465b1bf4aa3a90fd83cfd0ebb4b2020c1cba62c
            }
	    getFormEgre(&eg,r)

	    st          :=  getEgreFormUp(eg, egres, r)
            if len(st) == 0{
                 sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
            } else {
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
// ---------------------------------------------------
// EgreLisGET despliega formulario escoger periodo
func EgreLisGET(w http.ResponseWriter, r *http.Request) {
	sess := model.Instance(r)
        lisPeriod, err := model.Periods()
        if err != nil {
             sess.AddFlash(view.Flash{"No hay periodos ", view.FlashError})
         }
	v                  := view.New(r)
	v.Name              = "egreso/egresoper"
	v.Vars["token"]     = csrfbanana.Token(w, r, sess)
        v.Vars["LisPeriod"] = lisPeriod
        v.Vars["Title"]     =  "Listar"
        v.Vars["Action"]    =  "/egreso/list"
	v.Render(w)
 }

//------------------------------------------------
// EgreLis displays the egres page
func EgreLisPOST(w http.ResponseWriter, r *http.Request) {
        var Id  uint32
	var per  model.Periodo
	sess            := model.Instance(r)
        lisPeriod,err    := model.Periods()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Obteniendo Periodos.", view.FlashError})
            sess.Save(r, w)
         }
        Id,_             = atoi32(r.FormValue("id"))
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
	SId         := params.ByName("id")
	id,_        := atoi32(SId)
        path        :=  "/egreso/list"
        egres.Id   = id
	err         := (&egres).EgresById()
	if err != nil { // Si no existe ingreso
           log.Println(err)
           sess.AddFlash(view.Flash{"Es raro. No hay egreso.", view.FlashError})
           sess.Save(r, w)
           http.Redirect(w, r, path, http.StatusFound)
           return
	}
	v                    := view.New(r)
	v.Name                = "egreso/egresodelete"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
        v.Vars["Title"]     =  "Eliminar Egreso"
        v.Vars["Action"]    =  "/egreso/delete"
        v.Vars["Egre"]        = egres
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
