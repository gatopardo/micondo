package controller

import (
	"fmt"
	"log"
	"net/http"
        "encoding/json"
        "encoding/base64"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/passhash"
	"github.com/gatopardo/micondo/app/shared/view"

	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"

        "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

  )

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt = "login_attempt"
  )

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
  }

 func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
  }

 // jdecrypt private function to "decrypt" password
  func jdecrypt( stCuen string , stPass string)(stRes string,err error){
           var stEnc  []byte
	   stEnc, err = base64.StdEncoding.DecodeString(stPass)
	   if err != nil {
//                fmt.Println("jdecrypt ", err )
                log.Println("jdecrypt ", err )
           } else{
		   lon    :=  len(stEnc)
		   lan    :=  len(stCuen)
		   if lon >  lan {
                           stCuen = PadRight(stCuen, " ", lon)
		   }
		   rCuen  :=  []byte(stCuen)
		   rEnc   :=  []byte(stEnc)
		   rRes   :=  make([]byte, lon)
		   for i  := 0; i < lon; i++ {
                        rRes[i] = rCuen[i] ^ rEnc[i]
		   }
//                   fmt.Println("jdecrypt  ",rRes)
		   stRes = string(rRes)
	   }
	   return
  }

// JLoginGET service to return persons data
 func JLoginGET(w http.ResponseWriter, r *http.Request) {
      // Get session
        var params httprouter.Params
	sess           := model.Instance(r)

	v := view.New(r)
	v.Vars["token"] = csrfbanana.Token(w, r, sess)

        params          = context.Get(r, "params").(httprouter.Params)
        cuenta         := params.ByName("cuenta")
        password       := params.ByName("password")
	password, _      =  jdecrypt(cuenta, password)
        var jpers  model.Jperson
        jpers.Cuenta  = cuenta
	pass, err    := (&jpers).JPersByCuenta()
	// Determine if user exists
	if err == model.ErrNoResult {
             loginAttempt(sess)
	} else {
		b:= passhash.MatchString(pass, password)
                if b && jpers.Nivel > 0{ var js []byte
		   js, err =  json.Marshal(jpers)
                   if err == nil{
//		        fmt.Println(" json " + string(js))
			model.Empty(sess)
			fmt.Println("JLoginGET 1 sess.Values ", sess.Values )
			sess.Values["id"] = jpers.Id
                        sess.Save(r, w)
			fmt.Println("JLoginGET 2 sess.Values ", sess.Values )
                        w.Header().Set("Content-Type", "application/json")
                        w.Write(js)
			return
                    }
	        }
            }
               log.Println(err)
//	       http.Error(w, err.Error(), http.StatusBadRequest)
               http.Error(w, err.Error(), http.StatusInternalServerError)
		return
      }

// LoginGET displays the login page
func LoginGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"cuenta","password"}, r.Form, v.Vars)
	v.Render(w)
  }

// LoginPOST handles the login form submission
 func LoginPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"cuenta", "pass"}); !validate {
		sess.AddFlash(view.Flash{"Falta campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
       // Prevenir intentos login de fuerza bruta pretendiendo entrada invalida :-)
        if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 3 {
		log.Println("Intentos de Entrada Repetidos en Exceso")
		sess.AddFlash(view.Flash{"No mas intentos :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
	// Form values
	cuenta := r.FormValue("cuenta")
	password := r.FormValue("pass")
	// Get database user
         var user  model.User
         user.Cuenta  = cuenta
	 err := (&user).UserByCuenta()
	// Determine if user exists
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Cuenta incorrecta - Intento: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
                log.Println("Cuenta Incorreta ", err)
	} else if err != nil {
		// Display error message
		log.Println(" Error busqueda ",err)
		sess.AddFlash(view.Flash{"Un error. Favor probar mas tarde.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(user.Password, password) {
             if user.Nivel == 0 {
                 // User inactive and display inactive message
                 sess.AddFlash(view.Flash{"Cuenta inactiva entrada prohibida.", view.FlashNotice})
//			sess.Save(r, w)
		} else { // Login successfully
			model.Empty(sess)
			fmt.Println("LoginPost 1 sess.Values ", sess.Values )
			sess.AddFlash(view.Flash{"Entrada exitosa!", view.FlashSuccess})
			sess.Values["id"] = user.Id
			sess.Values["cuenta"] = cuenta
			sess.Values["level"] = user.Nivel
			fmt.Println("LoginPost 2 sess.Values ", sess.Values )
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Clave incorrecta - Intento: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	}
	// Show the login page again
	LoginGET(w, r)
  }

// JLoginPOST service to return persons data
 func JLoginPOST(w http.ResponseWriter, r *http.Request) {
      // Get session
	sess           := model.Instance(r)
	v := view.New(r)
	v.Vars["token"] = csrfbanana.Token(w, r, sess)

	cuenta := r.FormValue("Cuenta")
	password := r.FormValue("Pass")
	password, _      =  jdecrypt(cuenta, password)
	fmt.Println("JLoginPOST cuenta ", cuenta, " Pass ", password)
        var jpers  model.Jperson
        jpers.Cuenta  = cuenta
	pass, err    := (&jpers).JPersByCuenta()
	// Determine if user exists
	if err == model.ErrNoResult {
             loginAttempt(sess)
	} else {
		b:= passhash.MatchString(pass, password)
                if b && jpers.Nivel > 0{ var js []byte
		   js, err =  json.Marshal(jpers)
                   if err == nil{
		        fmt.Println(" json " + string(js))
			sess.Values["id"] = jpers.Id
                        sess.Save(r, w)
                        w.Header().Set("Content-Type", "application/json")
                        w.Write(js)
			return
                    }
	        }
            }
               log.Println(err)
//	       http.Error(w, err.Error(), http.StatusBadRequest)
               http.Error(w, err.Error(), http.StatusInternalServerError)
		return
      }


// LogoutGET clears the session and logs the user out
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := model.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
	        model.Empty(sess)
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
