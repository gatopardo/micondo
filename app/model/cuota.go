package model

import (
        "database/sql"
	"time"
        "log"
//	"fmt"
)

// *****************************************************************************
// Cuota
// *****************************************************************************.

// Cuota table contains the information for each cuot
type Cuota struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
	ApartaId         uint32     `db:"apartaid" bson:"apartaid,omitempty"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
        Fecha       time.Time       `db:"fecha" bson:"fecha"`
        Amount           uint64     `db:"amount" bson:"amount"`
	CreatedAt   time.Time       `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" bson:"updated_at"`
}

type CuotaN struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"periodid" bson:"periodid,omitempty"`
        Period        time.Time     `db:"period" bson:"period"`
	ApartaId         uint32     `db:"apartaid" bson:"apartaid,omitempty"`
	Apto             string     `db:"acodigo" bson:"acodigo,omitempty"`
	TipoId           uint32     `db:"tipoid" bson:"tipoid,omitempty"`
	Tipo             string     `db:"tdescripcion" bson:"tdescripcion,omitempty"`
        Fecha         time.Time     `db:"fecha" bson:"fecha"`
        Amount           uint64     `db:"amount" bson:"amount"`
	CreatedAt     time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" bson:"updated_at"`
}

// -----------------------------------------
// CuotById tenemos la cuota dado id
func (cuot * CuotaN)CuotById() (err error) {
        stq  :=   "SELECT c.id, c.period_id,p.inicio,c.aparta_id,a.codigo,c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p  WHERE c.period_id = p.id and c.aparta_id = a.id, c.tipo_id = t.id  and c.id=$1"
		err = Db.QueryRow(stq, &cuot.Id). Scan(&cuot.Id, &cuot.PeriodId,&cuot.Period,&cuot.ApartaId, &cuot.Apto, &cuot.TipoId, &cuot.Tipo, &cuot.Fecha, &cuot.Amount,  &cuot.CreatedAt, &cuot.UpdatedAt)

	return  standardizeError(err)
}

// -----------------------------------------------------
// CuotCreate crear cuota
func (c * CuotaN)CuotCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO cuotas ( period_id,aparta_id, tipo_id, fecha, amount,  created_at, updated_at ) VALUES ($1,$2,$3,$4, $5, $6, $7) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
         }
         defer stmt.Close()
         var id uint32
         err = stmt.QueryRow(&c.PeriodId,&c.ApartaId, &c.TipoId,&c.Fecha, &c.Amount,  now, now ).Scan(&id)
         if err == nil {
              c.Id = id
         }
	 return standardizeError(err)
  }

// -----------------------------------------------------
 func  (cuot * Cuota)CuotDeleteById()( err error){
         stqd :=  "DELETE FROM cuotas where id = $1"
           _, err = Db.Exec(stqd, cuot.Id) 
         return
       }

// -----------------------------------------------------
// Delete cuot from databa
func (cuot *Cuota) CuotDelete() (err error) {
	statement := "delete from cuotas where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(cuot.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de cuota en la database
func (cuot *CuotaN)CuotUpdate(stq string) (err error) {
        _, err = Db.Exec(stq ) 
        return standardizeError(err)
}

// -----------------------------------------------------
// Delete all cuotas from database
func CuotDeleteAll() (err error) {
	statement := "delete from cuotas"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in cuotas
  func CuotsCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM cuotas "
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
// Get cuotas from a period 
  func CuotLim(id uint32 ) (cuotas []CuotaN, err error) {
        stq :=   "SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, apartas a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.tipo_id = t.id and p.id = $1 order by p.inicio, c.fecha "
//	fmt.Println("CuotLim id ", id)
	rows, err := Db.Query(stq, id)
	if err != nil {
		log.Println(err)
//		fmt.Println(err)
            return
	}
//	fmt.Println("CuotLim middle  ")
        defer rows.Close()
        cuot := CuotaN{}
        for rows.Next() {
//           fmt.Println("CuotLim Next  ")
           if err = rows.Scan(&cuot.Id,&cuot.PeriodId,&cuot.Period, &cuot.ApartaId, &cuot.Apto, &cuot.TipoId, &cuot.Tipo, &cuot.Fecha, &cuot.Amount,  &cuot.CreatedAt, &cuot.UpdatedAt); err != nil {
		log.Println(err)
//		fmt.Println(err)
                 return
            }
//  fmt.Printf("%4d %4d %s %s %14d\n",cuot.Id,cuot.PeriodId, cuot.Period, cuot.Apto, cuot.Amount)
           cuotas = append(cuotas, cuot)
         }
//	fmt.Println("CuotLim fin  ", len(cuotas))
       return
 }



// -------------------------------------------------------------
// Get all cuotas in the database and returns the list
  func Cuots() (cuotas []CuotaN, err error) {
        stq :=   "SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, aparta a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.Tipo_id = t.id order by p.inicio"
	rows, err := Db.Query(stq)
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            cuot := CuotaN{}
           if err = rows.Scan(&cuot.Id,&cuot.PeriodId,&cuot.Period, &cuot.ApartaId, &cuot.Apto, &cuot.TipoId, &cuot.Tipo,  &cuot.Fecha,&cuot.Amount,  &cuot.CreatedAt, &cuot.UpdatedAt); err != nil {
                  return
             }
             cuotas = append(cuotas, cuot)
	}
        return
 }
// -------------------------------------------------------------
// Get all cuotas per a period in the database and returns the list
  func (cuot * CuotaN)CuotsPer() (cuotas []CuotaN, err error) {
        stq :=   "SELECT c.id, c.period_id, p.inicio, c.aparta_id, a.codigo, c.tipo_id, t.descripcion, c.fecha, c.amount, c.created_at, c.updated_at FROM cuotas c, periods p, apartas a, tipos t where c.period_id = p.id and c.aparta_id = a.id and c.tipo_id = t.id and c.period_id = $1 order by p.inicio"
	rows, err := Db.Query(stq, cuot.PeriodId )
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            c := CuotaN{}
           if err = rows.Scan(&c.Id,&c.PeriodId,&c.Period, &c.ApartaId, &c.Apto, &c.TipoId, &c.Tipo,  &c.Fecha,&c.Amount,  &c.CreatedAt, &c.UpdatedAt); err != nil {
                  return
             }
             cuotas = append(cuotas, c)
	}
        return
 }
// -------------------------------------------------------------
