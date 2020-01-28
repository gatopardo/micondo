package model

import (
        "database/sql"
	"time"
//        "log"
//	"fmt"
)

// *****************************************************************************
// Egreso
// *****************************************************************************.

// Egreso table contains the information for each egres
type Egreso struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
        Descripcion      string      `db:"dscripcion" bson:"dscripcion"`
	CreatedAt   time.Time       `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" bson:"updated_at"`
}

type EgresoN struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
        Period      time.Time       `db:"period" bson:"period"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
	Tipo             string     `db:"tcodigo" bson:"tcodigo,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           int64     `db:"amount" bson:"amount"`
        Descripcion      string      `db:"dscripcion" bson:"dscripcion"`
	CreatedAt time.Time        `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" bson:"updated_at"`
}

// -----------------------------------------
// EgresById tenemos el egreso dado id
func (egres * EgresoN)EgresById() (err error) {
        stq  :=   "SELECT e.id, e.period_id,p.inicio,e.tipo_id, t.codigo, e.fecha, e.amount, e.description, e.created_at, e.updated_at FROM egresos e, periods p  WHERE e.period_id = p.id  e.tipo_id = t.id  and e.id=$1"
		err = Db.QueryRow(stq, &egres.Id). Scan(&egres.Id, &egres.PeriodId,&egres.Period,  &egres.TipoId, &egres.Tipo, &egres.Fecha, &egres.Amount, &egres.Descripcion,  &egres.CreatedAt, &egres.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// EgresCreate crear egreso
func (e * EgresoN)EgresCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO egresos ( period_id, tipo_id, fecha, amount, description, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6, $7) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
         }
         defer stmt.Close()
         var id uint32
         err = stmt.QueryRow(&e.PeriodId, &e.TipoId,&e.Fecha, &e.Amount, &e.Descripcion,  now, now ).Scan(&id)
         if err == nil {
              e.Id = id
         }
	 return standardizeError(err)
  }

// -----------------------------------------------------
 func  (egres * Egreso)EgresDeleteById()( err error){
         stqd :=  "DELETE FROM egresos where id = $1"
           _, err = Db.Exec(stqd, egres.Id) 
         return
       }

// -----------------------------------------------------
// Delete egres from databa
func (egres *Egreso) EgresDelete() (err error) {
	statement := "delete from egresos where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(egres.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de egres en la database
func (egres *EgresoN)EgresUpdate(stq string) (err error) {
        _, err = Db.Exec(stq ) 
        return standardizeError(err)
}

// -----------------------------------------------------
// Delete all egresos from database
func EgresDeleteAll() (err error) {
	statement := "delete from egresos"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in egresos
  func EgresCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM egresos "
	rows, err := Db.Query(stq)
	if err != nil {
		return
	}
	defer rows.Close()
        for rows.Next() {
            err = rows.Scan(&count)
	    if err != nil {
	         return
	    }
        }
	return
 }
// -------------------------------------------------------------
// Get egresos from a period 
  func EgresLim(id uint32 ) (egresos []EgresoN, err error) {
        stq :=   "SELECT e.id, e.period_id, p.inicio,  e.tipo_id, t.codigo, e.fecha, e.amount, e.description, e.created_at, e.updated_at FROM egresos e, periods p,  tipos t where e.period_id = p.id  and e.tipo_id = t.id and p.id = $1 order by p.inicio, e.fecha "
	rows, err := Db.Query(stq, id)
	if err != nil {
            return
	}
        defer rows.Close()
        for rows.Next() {
           egres := EgresoN{}
           if err = rows.Scan(&egres.Id,&egres.PeriodId,&egres.Period, &egres.TipoId, &egres.Tipo, &egres.Fecha, &egres.Amount, &egres.Descripcion,  &egres.CreatedAt, &egres.UpdatedAt); err != nil {
                  return
            }
           egresos = append(egresos, egres)
         }
       return
 }
// -------------------------------------------------------------
// Get all egresos in the database and returns the list
  func Egress() (egresos []EgresoN, err error) {
        stq :=   "SELECT e.id, e.period_id, p.inicio,  a.codigo, e.tipo_id, t.codigo, e.fecha, e.amount, e.description, e.created_at, e.updated_at FROM egresos e, periods p,  tipos t where e.period_id = p.id  and e.Tipo_id = t.id order by p.inicio"
	rows, err := Db.Query(stq)
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            egres := EgresoN{}
           if err = rows.Scan(&egres.Id,&egres.PeriodId,&egres.Period, &egres.TipoId, &egres.Tipo,  &egres.Fecha,&egres.Amount, &egres.Descripcion,  &egres.CreatedAt, &egres.UpdatedAt); err != nil {
                  return
             }
             egresos = append(egresos, egres)
	}
        return
 }
// -------------------------------------------------------------
// Get all egresos per a period in the database and returns the list
  func (egre * EgresoN)EgresPer() (egresos []EgresoN, err error) {
        stq :=   "SELECT e.id, e.period_id, p.inicio,  e.tipo_id, t.codigo, e.fecha, e.amount, e.description, e.created_at, e.updated_at FROM egresos e, periods p,  tipos t where e.period_id = p.id and e.tipo_id = t.id and e.period_id = $1 order by p.inicio"
	rows, err := Db.Query(stq, egre.PeriodId )
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            e := EgresoN{}
           if err = rows.Scan(&e.Id, &e.PeriodId, &e.Period, &e.TipoId, &e.Tipo, &e.Fecha, &e.Amount, &e.Descripcion, &e.CreatedAt, &e.UpdatedAt); err != nil {
                  return
             }
             egresos = append(egresos, e)
	}
        return
 }
// -------------------------------------------------------------
