package model

import (
        "database/sql"
	"fmt"
	"time"
	"strings"
	"crypto/rand"
        "log"
        "net/http"
  )

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
   type User struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
        PersonId     uint32        `db:"person_id"  bson: "person_id"`
	Uuid        string        `db:"uuid" bson:"uuid,omitempty"`
	Cuenta      string        `db:"cuenta" bson:"cuenta"`
	Password    string        `db:"password" bson:"password"`
	Nivel       uint32        `db:"nivel" bson:"nivel"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"created_at"`
}

// Person table contains user profile
  type Person struct {
	Id          uint32        `db:"id" bson:"id,omitempty"`
        ApartaId     uint32        `db:"aparta_id"  bson: "aparta_id"`
	Fname       string        `db:"first_name" bson:"first_name"`
	Lname       string        `db:"last_name" bson:"last_name"`
	Email       string        `db:"email" bson:"email"`
	Address     string        `db:"address" bson:"address"`
	Tele        string        `db:"tele" bson:"tele"`
	Mobil       string        `db:"mobil" bson:"mobil"`
        Tipo        string        `db:"type"  bson: "type"`
        Photo       string        `db:"photo"  bson: "photo"`
	CreatedAt   time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" bson:"updated_at"`
 }

 type Individ struct{
           Usu User
           Pers Person
	   Apto Aparta
   }

// Shadow table contains text password
  type Shadow struct {
        Id          uint32        `db:"id" bson:"id,omitempty"`
	UserId      uint32        `db:"user_id" bson:"id,omitempty"`
	Uuid        string        `db:"uuid" bson:"uuid,omitempty"`
	Password    string        `db:"password" bson:"password"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"created_at"`
  }
// --------------------------------------------------------


 func (u *  User)DatFormUser(r * http.Request){
           u.Cuenta   =   r.FormValue("cuenta")
   }

// --------------------------------------------------------
// Crear una nueva sesion para un usuario
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid,  user_id, created_at, updated_at) values ($1, $2, $3, $4) returning id, uuid,  user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// usar QueryRow para retornar una fila y buscar el id para struct Session 
	err = stmt.QueryRow(CreateUUID(),  user.Id, time.Now()).Scan(session.Id, session.Uuid,  session.UserId, session.CreatedAt )
	return
}

// --------------------------------------------------------
// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid,  user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		Scan(&session.Id, &session.Uuid,  &session.UserId, &session.CreatedAt)
	return
}

// --------------------------------------------------------
// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid,  &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// --------------------------------------------------------
// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

// --------------------------------------------------------
// Obtener usuario desde la la session
func (session *Session) User() (user User, err error) {
        user = User{}
        err = Db.QueryRow("SELECT id, uuid, cuenta,password nivel  created_at FROM users WHERE id = $1", session.UserId).
        Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Nivel, &user.CreatedAt)
	return
}

// --------------------------------------------------------
// Delete all sessions from database
func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

// -----------------------------------------
// ShadCreate creates shadow
func (sd *Shadow)ShadCreate() error {
	var err error
        var stmt  *sql.Stmt
        stq := "INSERT INTO shadows (user_id, uuid,  password, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5) returning id"
	now  := time.Now()
        if stmt, err = Db.Prepare(stq ); err != nil  {
	      return standardizeError(err)
        }
        defer stmt.Close()
        var id uint32
        err = stmt.QueryRow(&sd.UserId, &sd.Uuid, &sd.Password,  now, now ).Scan(&id)
        if err == nil {
              sd.Id = id
        }
	return standardizeError(err)
}
// -----------------------------------------------------
// ===========================================================//

    func  (p *Person)person_Trim(){
              p.Lname = strings.Trim(p.Lname, " ")
              p.Fname = strings.Trim(p.Fname, " ")
              p.Email = strings.Trim(p.Email, " ")
              p.Address = strings.Trim(p.Address, " ")
              p.Tele = strings.Trim(p.Tele, " ")
              p.Mobil = strings.Trim(p.Mobil, " ")
              p.Photo = strings.Trim(p.Photo, " ")
         }

// -----------------------------------------------------
// PersonById obtenemos la persona dado id
func PersonById(id uint32) (p Person,err error) {
         stq    :=   "SELECT fname,lname, email, address, tele, mobil, tipo,photo, created_at, updated_at FROM  persons   WHERE  id = $1"
         err = Db.QueryRow(stq, id).Scan( &p.Fname, &p.Lname, &p.Email, &p.Address, &p.Tele, &p.Mobil, &p.Tipo, &p.Photo, &p.CreatedAt, &p.UpdatedAt )
//         return  standardizeError(err)
         return
 }

// -----------------------------------------------------
// PersonById obtenemos la persona dado id
func EmailByAptId(id uint32) (p Person,err error) {
         stq    :=   "SELECT fname, lname, email from persons p where p.aparta_id = $1"
         err = Db.QueryRow(stq, id).Scan(&p.Fname, &p.Lname, &p.Email  )
//         return  standardizeError(err)
          if err != nil{
		  return
	  }
         return

 }

// -----------------------------------------------------

// PersonCreate creates person
func (p * Person)PersonCreate()( err error) {
        var id uint32
        var stmt  *sql.Stmt
        stq := "INSERT INTO persons (aparta_id, fname, lname, email, address, tele, mobil,tipo,photo, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6, $7, $8, $9, $10, $11) returning id"
	now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	       return standardizeError(err)
         }
         defer stmt.Close()
        err = stmt.QueryRow(p.ApartaId,p.Fname, p.Lname, p.Email, p.Address, p.Tele, p.Mobil, p.Tipo, p.Photo, now, now ).Scan(&id)
        if err == nil {
             p.Id = id
        }
	return standardizeError(err)
  }

// -----------------------------------------------------

// Update person information in the database
func (p * Person) Update(stq string ) (err error) {
            _, err = Db.Exec(stq )
            if err != nil{
              fmt.Println("Person Update ", err)
            }
	    return standardizeError(err)
}

// -----------------------------------------------------

// Delete all users from database
func PersDeleteAll() (err error) {
	statement := "delete from persons"
	_, err = Db.Exec(statement)
	return
}
// -----------------------------------------------------

// Get all users in the database and returns it
  func Persons() (persons []Person, err error) {
        stq :=   "SELECT p.id,   fname, lname, email, address, tele, mobil,p.photo, p.created_at, p.updated_at FROM persons p, users u where p.user_id = u.id  order by cuenta"
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		person := Person{}
		if err = rows.Scan(&person.Id, &person.Fname, &person.Lname, &person.Email, &person.Address, &person.Tele, &person.Mobil, &person.Photo ,&person.CreatedAt, &person.UpdatedAt); err != nil {

			return
		}
	}
	return
}
// -----------------------------------------------------

// ===========================================================//

// UserById tenemos el usuario dado id
func (user * User)UserById() (err error) {
        stq  :=   "SELECT person_id, uuid, cuenta, password,nivel, created_at, updated_at FROM users WHERE id=$1"
	err = Db.QueryRow(stq, &user.Id).Scan( &user.PersonId, &user.Uuid,  &user.Cuenta, &user.Password, &user.Nivel, &user.CreatedAt, &user.UpdatedAt)
	return  standardizeError(err)
}

// ----------------------------------------
// UserByCuenta gets user information from cuenta
func (user *User)UserByCuenta() ( error) {
	var err error
        stq  :=   "SELECT id, uuid, cuenta, password,nivel, created_at, updated_at FROM users WHERE cuenta=$1"
         err = Db.QueryRow(stq, &user.Cuenta).Scan(&user.Id, &user.Uuid, &user.Cuenta, &user.Password, &user.Nivel, &user.CreatedAt, &user.UpdatedAt)
	return   standardizeError(err)
}

// -----------------------------------------------------
// UserCreate crear usuario
func (u *User)UserCreate() (err error) {
	var stmt  *sql.Stmt
         stq := "INSERT INTO users (person_id, uuid, cuenta, password, nivel, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6, $7) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	       return standardizeError(err)
         }
         defer stmt.Close()
         var id uint32
         err = stmt.QueryRow( &u.PersonId, &u.Uuid, &u.Cuenta, &u.Password, &u.Nivel, now , now).Scan(&id)
         if err == nil {
              u.Id = id
         }
	 return standardizeError(err)
 }

// -----------------------------------------------------
// delete user by id
 func (user * User) UserDeleteById() ( err error){
         stqd :=  "DELETE FROM users where id = $1"
           _, err = Db.Exec(stqd, user.Id)
         return
       }

// -----------------------------------------------------
// Delete user from databa
func (user *User) UserDelete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de usuario en la database
func (user *User)Update() (err error) {
	statement := "update users set cuenta = $2, password = $3, nivel = $4 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Cuenta, user.Password, user.Nivel)
	return
}

// -----------------------------------------------------
// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// -----------------------------------------------------
// Get all users in the database and returns the list
  func Users() (users []Individ, err error) {
       var stq  string
       stq =   "SELECT u.id, u.cuenta, u.nivel, p.id, p.fname, p.lname,p.email, p.address, p.tele, p.mobil, p.tipo, p.photo,a.id, a.codigo, a.descripcion  FROM users u join persons p on u.person_id = p.id join apartas a on p.aparta_id = a.id order by p.fname, p.lname"

	rows, err := Db.Query(stq)
	if err != nil {
           log.Println(err)
           fmt.Println(err)
           return
	}
	defer rows.Close()
	for rows.Next() {
		indi := Individ{}
		if err = rows.Scan(&indi.Usu.Id,  &indi.Usu.Cuenta,  &indi.Usu.Nivel,&indi.Pers.Id,  &indi.Pers.Fname, &indi.Pers.Lname,&indi.Pers.Email, &indi.Pers.Address, &indi.Pers.Tele, &indi.Pers.Mobil, &indi.Pers.Tipo, &indi.Pers.Photo, &indi.Apto.Id, &indi.Apto.Codigo, &indi.Apto.Descripcion); err != nil {
              fmt.Println(indi.Usu, indi.Pers)
              log.Println(err)
              fmt.Println(err)
           }
	    users = append(users, indi)
	}
	return
 }
// -----------------------------------------------------
// -----------------------------------------------------
// Get selected users in the database and returns the list
  func SUsers(rsearch string) (users []Individ, err error) {
       var stqi, stqf, stq, stq1, stq2   string
        rsearch    =   strings.Trim(rsearch, " ")
        nCount     := strings.Count(rsearch, "@")
        arSt := strings.Split(rsearch, "@")
        if arSt[0] != ""{
            stq1  = " and  p.lname SIMILAR TO '" +arSt[0]+ "' "
        }
        if nCount >= 1 && arSt[1] != "" {
            stq2  = " and  p.fname SIMILAR TO '"+arSt[1]+ "' "
        }
        stqi =   "SELECT u.id, u.cuenta, u.nivel, p.fname, p.lname, p.address, p.tele, p.mobil, p.tipo, p.photo  FROM users u, persons p where u.person_id = p.id  "
               stqi =   stqi +   " and p.tipo = 'P' "
               stqf =    " order by p.lname, p.fname"
        stq =  stqi  + stq1 + stq2  + stqf

fmt.Println("SUsus ", stq)
	rows, err := Db.Query(stq)
	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		indi := Individ{}
                 var arPhoto []byte
		if err = rows.Scan(&indi.Usu.Id,  &indi.Usu.Cuenta,  &indi.Usu.Nivel, &indi.Pers.Fname, &indi.Pers.Lname, &indi.Pers.Address, &indi.Pers.Tele, &indi.Pers.Mobil, &indi.Pers.Tipo, &arPhoto); err != nil {

fmt.Println(indi.Usu, indi.Pers)
			return
		}
                indi.Pers.Photo = string(arPhoto)
		users = append(users, indi)
	}
	return
 }
// -----------------------------------------------------

// ===========================================================//

 func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("No se genera UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122 
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
 }

// -----------------------------------------------------



