package controller

import (
	"log"
	"net/http"
         "strings"
         "fmt"
          "path/filepath"
          "os"
          "io"
//          "io/ioutil"
//          "encoding/base64"
          "mime/multipart"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/passhash"
	"github.com/gatopardo/micondo/app/shared/view"

        "github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
  )

 const maxUploadSize = 2 << 20 // 2 mb

// ---------------------------------------------------

  func saveFile( fileHeader *multipart.FileHeader , data []byte )(err error){
         fileName    := fileHeader.Filename

         newPath := filepath.Join("images", fileName)
         var newFile *os.File
         newFile, err = os.Create(newPath)
         if err != nil {
             fmt.Println(err)
             return
         }
         defer newFile.Close()
         if _, err = newFile.Write(data); err != nil {
              fmt.Println(err)
              return
         }
       return
    }

// ---------------------------------------------------

  func getImage(r *http.Request, sname string)( st string ,err error){
         var file multipart.File
         var fileHeader *multipart.FileHeader
//         var buff bytes.Buffer
         r.ParseMultipartForm(maxUploadSize)
         file, fileHeader, err = r.FormFile(sname)
         if err != nil {
             fmt.Println("Error Retrieving the File")
             fmt.Println(err)
             return
         }
         defer file.Close()
//         fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
//         fmt.Printf("File Size: %+v\n", fileHeader.Size)
//         fmt.Printf("MIME Header: %+v\n", fileHeader.Header)
         st = "./static/favicons/icon_a/"+fileHeader.Filename
         f, err := os.OpenFile(st, os.O_WRONLY|os.O_CREATE, 0666)
         if err != nil {
             fmt.Println(err)
             return
         }
         defer f.Close()
         io.Copy(f, file)
         return
   }

//----------------------------------------------------
 func DatFormPers(p *model.Person, r * http.Request){
	 p.ApartaId,_     = atoi32(r.FormValue("aptId"))
         p.Fname     =  r.FormValue("fname")
         p.Lname     =  r.FormValue("lname")
         p.Email     =  r.FormValue("email")
         p.Address   =  r.FormValue("address")
         p.Tele      =  r.FormValue("tele")
         p.Mobil     =  r.FormValue("mobil")
         p.Tipo      =  r.FormValue("tipo")
         rPhoto, _  :=  getImage(r, "photo")
         p.Photo     =  rPhoto
//         fmt.Println("Person ", *p )
   }

// --------------------------------------------------------

// RegisterGET despliega la pagina del usuario
 func RegisterGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess    := model.Instance(r)
	lsApts, err  := model.Apts()
	if err != nil {
	   fmt.Println(err)
	   log.Println(err)
	}
	// Display the view
	v := view.New(r)
        v.Vars["LisApts"] = lsApts
	v.Name = "register/register"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
 }
// ---------------------------------------------------
	// Refill any form fields
//	view.Repopulate([]string{"cuenta", "password", "level"}, r.Form, v.Vars)
// ---------------------------------------------------

// RegisterPOST procesa la forma enviada con los datos
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
        var user model.User
        var person model.Person
        var shad model.Shadow
	sess       := model.Instance(r)
        action     := r.FormValue("action")
//      fmt.Println(action)
        if ! (strings.Compare(action,"Cancelar") == 0) {
           if validate, missingField := view.Validate(r, []string{"cuenta", "level"}); !validate {
               sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
               sess.Save(r, w)
               RegisterGET(w, r)
               return
	   }
           stUuid := model.CreateUUID()
           // Get form values
           rPasswd         := r.FormValue("password")
	   user.Cuenta      = r.FormValue("cuenta")
           vPasswd         := r.FormValue("password_verify")
           user.Nivel, _    = atoi32( r.FormValue("level"))
           user.Uuid        = stUuid
           if strings.Compare(rPasswd, vPasswd) != 0{
		log.Println(rPasswd, " * ", vPasswd)
		sess.AddFlash(view.Flash{"Claves distintas no posible", view.FlashError})
		sess.Save(r, w)
                RegisterGET(w, r)
//		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
           }
            pass, errp := passhash.HashString(rPasswd)
	   // If password hashing failed
	   if errp != nil {
		log.Println(errp)
                sess.AddFlash(view.Flash{"Problema encriptando clave.", view.FlashError})
		sess.Save(r, w)
//              RegisterGET(w, r)
		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
	   }
           user.Password = pass
	   err := (&user).UserByCuenta()
           if err == model.ErrNoResult { // Exito:  no hay usuario creado aun 
                DatFormPers(&person, r)
// fmt.Println("Person ", person.Fname, person.Lname)
                err = (&person).PersonCreate()
                if err != nil {
                   log.Println(err)
 fmt.Println("Person ", person.Fname, person.Lname)
                   sess.AddFlash(view.Flash{"Error guardando Person.", view.FlashError})
                   sess.Save(r, w)
		  http.Redirect(w, r, "/user/register", http.StatusFound)
                   return
                }
                user.PersonId =  person.Id
                ex := (&user).UserCreate()
	        if ex != nil {  // uyy como fue esto ? 
                   log.Println(ex)
                   sess.AddFlash(view.Flash{"Error guardando User.", view.FlashError})
                   sess.Save(r, w)
		   http.Redirect(w, r, "/user/register", http.StatusFound)
                   return
	        } else {  // todo bien
                   shad.UserId    = user.Id
                   shad.Uuid      = stUuid
                   shad.Password  = rPasswd
                   if  err = (&shad).ShadCreate() ; err != nil{
                       sess.AddFlash(view.Flash{"Error guardando.", view.FlashError})
                       log.Println( err)
                       sess.Save(r, w)
		       http.Redirect(w, r, "/user/register", http.StatusFound)
                       return
                   }
                   sess.AddFlash(view.Flash{"Creando: " +user.Cuenta, view.FlashSuccess})
                   sess.Save(r, w)
	        }
           }
         }
	// Display list
	http.Redirect(w, r, "/user/list", http.StatusFound)
  }
// ---------------------------------------------------
// ---------------------------------------------------
// RegisUpGET despliega la pagina del usuario
func RegisUpGET(w http.ResponseWriter, r *http.Request) {
        var user model.User
        var pers model.Person
	sess := model.Instance(r)
        // necesitamos user id
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        id,_ := atoi32(params.ByName("id"))
        user.Id = id
        // Obtener usuario dado id
        err := (&user).UserById()
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Es raro. No usuario.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
// fmt.Println("User ", user.Cuenta )
        pers.Id = user.PersonId
        pers,err = model.PersonById(pers.Id)
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"No hay atributos.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
	// Display the view
	v := view.New(r)
	v.Name = "register/regisupdate"
	v.Vars["token"]  = csrfbanana.Token(w, r, sess)
        v.Vars["cuenta"]   = user.Cuenta
        v.Vars["pers"]     =  pers
        v.Vars["imgPhoto"] =  string(pers.Photo)
        v.Vars["Level"]  =  sess.Values["level"]

//    Refill any form fields
//	view.Repopulate([]string{"cuenta", "level"}, r.Form, v.Vars)
        v.Render(w)
   }

//----------------------------------------------------
//  gettin partial updates
     func getPartialUp(r *http.Request,fieldname,  dname string) ( sform string ){
         sform = ""
         stTrimf  :=  strings.Trim(r.FormValue(fieldname), " ")
         stTrimp  :=   strings.Trim(dname, " ")
         if  ( stTrimf != stTrimp) && (len(stTrimp) > 0){
             sform = fmt.Sprintf(" %s  = '%s' ",fieldname, stTrimf )
         }
        return
      }

//----------------------------------------------------
    func getPersFormUp(p model.Person,r *http.Request)(stUp string){
      var sform string
      var  sArrSup []string
      sform =  getPartialUp(r, "fname", p.Fname)
      if len(sform) > 0 {
          sArrSup = append(sArrSup, sform)
       }
      sform =  getPartialUp(r, "lname", p.Lname)
      if len(sform) > 0 {
          sArrSup = append(sArrSup, sform)
       }
      sform =  getPartialUp(r, "address", p.Address)
      if len(sform) > 0 {
          sArrSup = append(sArrSup, sform)
       }
      sform =  getPartialUp(r, "tele", string(p.Tele))
      if len(sform) > 0 {
          sArrSup = append(sArrSup, sform)
       }
      sform =  getPartialUp(r, "mobil", p.Mobil)
      if len(sform) > 0 {
          sArrSup = append(sArrSup, sform)
       }
        rPhoto, _ := getImage(r, "photo")
        if len(rPhoto)  > 0  { // != len(p.Photo){
            sform = fmt.Sprintf(" photo  = '%s' ", string(rPhoto) )
//      sform =  getPartialUp(r,"photo", string( p.Photo))
             if len(sform) > 0 {
                  sArrSup = append(sArrSup, sform)
             }
       }

//      fmt.Println("getPersFormUp photo ", len(p.Photo))
        lon := len(sArrSup)
       if lon  > 0 {
            sini        :=  "update persons set "
            stUp =  strings.Join(sArrSup, ", ")
            sr          :=  fmt.Sprintf(" where persons.id = %d ", p.Id)
             stUp = sini + stUp + sr
// fmt.Println(sArrSup[:lon - 1])
       }
         return
    }
//----------------------------------------------------
// RegisUpPOST procesa la forma enviada con los datos
func RegisUpPOST(w http.ResponseWriter, r *http.Request) {
        var err error
        var user model.User
        var pers model.Person
	sess := model.Instance(r)
        var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
        action        := r.FormValue("action")
        path :=  "/user/list/"
//      fmt.Println(action)
        if (strings.Compare(action,"Cancelar") == 0) {
               sess.Save(r, w)
               http.Redirect(w, r, path, http.StatusFound)
               return
        }
	user.Id, _ = atoi32(params.ByName("id"))
        // Obtener usuario dado id
        err  = (&user).UserById()
        if err != nil { // Si no existe el usuario
            log.Println(err)
            sess.AddFlash(view.Flash{"Raro update. No usuario.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
        pers.Id = user.PersonId
        pers, err = model.PersonById(pers.Id)
        if err != nil { // Si no existe persona
            log.Println(err)
            sess.AddFlash(view.Flash{"No hay atributos.", view.FlashError})
            sess.Save(r, w)
            http.Redirect(w, r, "/user/list", http.StatusFound)
            return
        }
	user.Cuenta     = r.FormValue("cuenta")
        st          :=  getPersFormUp(pers,r)
        if len(st) == 0{
            sess.AddFlash(view.Flash{"No actualizacion solicitada", view.FlashSuccess})
        } else {
            err =  pers.Update(st)
            if err == nil{
                 sess.AddFlash(view.Flash{"Persona actualizado: " +user.Cuenta, view.FlashSuccess})
            } else       {
		log.Println(err)
		sess.AddFlash(view.Flash{"Un error actualizando.", view.FlashError})
	    }
          }
		sess.Save(r, w)
	http.Redirect(w, r, path, http.StatusFound)
     }

// ---------------------------------------------------
// RegisSearchPOST procesa la forma enviada con los datos
func RegisSearchPOST(w http.ResponseWriter, r *http.Request) {
	sess       := model.Instance(r)
        rSearch    := r.FormValue("search")
        if rSearch == ""{
           fmt.Println("Nada a buscar")
           return
         }
        lisUsers, err := model.SUsers(rSearch)
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Usuarios.", view.FlashError})
            sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "register/regislis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["LisRegis"] = lisUsers
	// Refill any form fields
	v.Render(w)
     }
// ---------------------------------------------------
//------------------------------------------------
// RegisLisGET displays the register page
func RegisLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisUsers, err := model.Users()
        if err != nil {
            log.Println(err)
	    sess.AddFlash(view.Flash{"Error Listando Usuarios.", view.FlashError})
            sess.Save(r, w)
         }
//	 fmt.Println(lisUsers)
	// Display the view
	v := view.New(r)
	v.Name = "register/regislis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["Level"]    =  sess.Values["level"]
        v.Vars["LisRegis"] = lisUsers
	v.Render(w)
}

// UserDelGET handles the user deletion
 func RegisDelGET(w http.ResponseWriter, r *http.Request) {
	// Get session
        sess := model.Instance(r)
        v    := view.New(r)
        v.Name = "register/regisdelete"
        var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        Id,_   := atoi32(params.ByName("id"))
	SPag        := params.ByName("pagi")
        path :=  fmt.Sprintf("/user/list/%s", SPag)
// fmt.Println(Id,SPag)
        var user model.User
        user.Id = Id
        err := (&user).UserById()
        if err != nil {
            log.Println(err)
            fmt.Println(err)
            sess.AddFlash(view.Flash{"Error Usuario no hallado.", view.FlashError})
                http.Redirect(w, r, path, http.StatusFound)
        }else{
//           v.Vars["cuenta"]  = user.Cuenta
//           v.Vars["level"]    = user.Level
              v.Vars["Level"]    =  sess.Values["level"]

      }
// fmt.Println(path)
	   v.Vars["token"]  = csrfbanana.Token(w, r, sess)
           v.Render(w)

  }

// PersDelPOST handles the user deletion
 func RegisDelPOST(w http.ResponseWriter, r *http.Request) {

 }
// PersRegisGET displays the register page
func PersRegisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "person/persreg"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"cuenta", "lname", "fname","email","address","tele", "mobil"  }, r.Form, v.Vars)
	v.Render(w)
}

// PersRegisPOST maneja registro personas  con  form submission
func  PersRegisPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
        var user model.User
        var person model.Person
	sess := model.Instance(r)
	// Se previene ataque de fuerza bruta con entradas invalidas múltiples :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		log.Println("Detenido intento en persona entrada repetida multiples veces")
		http.Redirect(w, r, "/regispers", http.StatusFound)
		sess.AddFlash(view.Flash{"Faltan Campos", view.FlashError})
		return
	}
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"cuenta", "lname", "fname", "email", "address", "tele", "mobil" }); !validate {
		sess.AddFlash(view.Flash{"Falta Campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		RegisterGET(w, r)
		return
	}
	// Get form values
	person.Fname   =  r.FormValue("fname")
	person.Lname   =  r.FormValue("lname")
	person.Email   =  r.FormValue("email")
	person.Address =  r.FormValue("address")
	person.Tele    =  r.FormValue("tele")
	person.Mobil   =  r.FormValue("mobil")
	// Busquemos datos de cuenta
	 err := user.UserByCuenta()
	if err == nil { // Exito  usuario con esa cuenta
	    ex := (&person).PersonCreate() // Si hay error es por el query
            if ex == nil { // todo bien creando persona
		sess.AddFlash(view.Flash{"Persona creada exitosamente para: " +user.Cuenta, view.FlashSuccess})
                sess.Save(r, w)
                http.Redirect(w, r, "/regispers", http.StatusFound)
		return
             }else{ // persona ya existía
                    sess.AddFlash(view.Flash{"La cuenta ya existe para: " + user.Cuenta, view.FlashError})
                sess.Save(r, w)
             }
        } else { // no hay usuario para esta persona    
	           log.Println(err)
		   sess.AddFlash(view.Flash{"Vamos a Crear Usuario.", view.FlashError})
		   sess.Save(r, w)
	}
//	if err == model.ErrNoResult  // If success (no usuario con esa cuenta)
	// Display the page
	PersRegisGET(w, r)
}

// PersLisGET displays the register page
func PersLisGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
        lisPers, err := model.Persons()
        if err != nil {
	    log.Println(err)
	    sess.AddFlash(view.Flash{"Error ocurrio Listando Personas.", view.FlashError})
             sess.Save(r, w)
         }
	// Display the view
	v := view.New(r)
	v.Name = "person/perslis"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
        v.Vars["lisPers"] = lisPers
	v.Render(w)
}



