package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gatopardo/micondo/app/model"
	"github.com/gatopardo/micondo/app/shared/passhash"
	"github.com/gatopardo/micondo/app/shared/view"

	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
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

       // Prevenir intentos login de fuerza bruta pretendiendo entrada invalida :-)
  if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 5 {
		log.Println("Intentos de Entrada Repetidos en Exceso")
		sess.AddFlash(view.Flash{"Favor, no mas intentos :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"cuenta", "password"}); !validate {
		sess.AddFlash(view.Flash{"Falta campo: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
	// Form values
	cuenta := r.FormValue("cuenta")
	password := r.FormValue("password")
	// Get database user
         var user  model.User
         user.Cuenta  = cuenta 
	 err := user.UserByCuenta()
	// Determine if user exists
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Clave incorrecta - Intento: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"Ocurrio un error. Favor probar mas tarde.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(user.Password, password) {
		if user.Nivel == 0 {
			// User inactive and display inactive message
			sess.AddFlash(view.Flash{"Cuenta inactiva entrada prohibida.", view.FlashNotice})
			sess.Save(r, w)
		} else {
			// Login successfully
			model.Empty(sess)
			sess.AddFlash(view.Flash{"Entrada exitosa!", view.FlashSuccess})
			sess.Values["id"] = user.Id
			sess.Values["cuenta"] = cuenta
			sess.Values["level"] = user.Nivel
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
