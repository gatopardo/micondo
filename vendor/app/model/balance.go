package model

import (
        "database/sql"
	"time"
//        "log"
//	"fmt"
)

// *****************************************************************************
// Balance
// *****************************************************************************.

// Balance table contains the information for each balan
type Balance struct {
	Id               uint32     `db:"id" bson:"id,omitempty"`
	PeriodId         uint32     `db:"balanid" bson:"balanid,omitempty"`
        Amount           uint64     `db:"amount" bson:"amount"` 
        Cuota            uint64     `db:"cuota" bson:"cuota"` 
	CreatedAt   time.Time       `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" bson:"updated_at"`
}

type BalanceN struct {
	Id             uint32      `db:"id" bson:"id,omitempty"`
	PeriodId       uint32      `db:"periodid" bson:"id,omitempty"`
	Period    time.Time        `db:"period" bson:"period"`
        Amount         uint64      `db:"amount" bson:"amount"` 
        Cuota          uint64      `db:"cuota" bson:"cuota"` 
	CreatedAt time.Time        `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" bson:"updated_at"`
}

// -----------------------------------------
// BalanById tenemos el balance dado id
func (balan * BalanceN)BalanById() (err error) {
        stq  :=   "SELECT b.id, b.period_id,p.inicio, b.amount, b.cuota, b.created_at, b.updated_at FROM balances b, periods p  WHERE b.period_id = p.id and b.id=$1"
		err = Db.QueryRow(stq, &balan.Id). Scan(&balan.Id, &balan.PeriodId,&balan.Period, &balan.Amount, &balan.Cuota, &balan.CreatedAt, &balan.UpdatedAt)

	return  standardizeError(err)
}

// --------------------------------------------------------

// BalanByPeriod gets balan information from period
func (balan *Balance)BalanByPeriod() ( error) {
	var err error
        stq  :=   "SELECT id,period_id, amount, cuota, created_at, updated_at FROM balances WHERE period_id=$1"
		err = Db.QueryRow(stq, &balan.PeriodId).Scan(&balan.Id,&balan.PeriodId, &balan.Amount, &balan.Cuota, &balan.CreatedAt, &balan.UpdatedAt)

	return   standardizeError(err)
}

// -----------------------------------------------------
// BalanCreate crear balance
func (b *Balance)BalanCreate() error {
         var err error
         var stmt  *sql.Stmt
         stq := "INSERT INTO balances ( period_id, amount, cuota, created_at, updated_at ) VALUES ($1,$2,$3,$4, $5) returning id"
	 now  := time.Now()
         if stmt, err = Db.Prepare(stq ); err != nil  {
	          return standardizeError(err)
         }
         defer stmt.Close()
         var id uint32
         err = stmt.QueryRow(&b.PeriodId,&b.Amount, &b.Cuota, now, now ).Scan(&id)
         if err == nil {
               b.Id = id
         }
	 return standardizeError(err)
  }

// -----------------------------------------------------
 func  (balan * Balance)BalanDeleteById()( err error){
         stqd :=  "DELETE FROM balances where id = $1"
           _, err = Db.Exec(stqd, balan.Id) 
         return
       }

// -----------------------------------------------------
// Delete balan from databa
func (balan *Balance) BalanDelete() (err error) {
	statement := "delete from balances where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(balan.Id)
	return
}

// -----------------------------------------------------
// Actualizar informacion de balance en la database
func (balan *BalanceN)BalanUpdate(stq string) (err error) {
        _, err = Db.Exec(stq ) 
        return standardizeError(err)
}

// -----------------------------------------------------
// Delete all balances from database
func BalanDeleteAll() (err error) {
	statement := "delete from balances"
	_, err = Db.Exec(statement)
	return
}

// -------------------------------------------------------------
// Get number of records in balances
  func BalansCount( ) ( count int) {
        stq :=  "SELECT COUNT(*) as count FROM balances "
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
// Get limit records from offset
  func BalansLim(lim int , offs int) (balances []Balance, err error) {
        stq :=   "SELECT id, period_id, amount, cuota, created_at, updated_at FROM balances order by level, cuota LIMIT $1 OFFSET $2"
	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
            return
	}
        defer rows.Close()
        for rows.Next() {
           balan := Balance{}
           if err = rows.Scan(&balan.Id,&balan.PeriodId, &balan.Amount, &balan.Cuota, &balan.CreatedAt, &balan.UpdatedAt); err != nil {
                  return
            }
           balances = append(balances, balan)
         }
       return
 }
// -------------------------------------------------------------
// Get limit records from offset
  func BalanPerLim(lim int , offs int) (balances []BalanceN, err error) {
        stq :=   "SELECT b.id, b.period_id, p.inicio, b.amount, b.cuota, b.created_at, b.updated_at FROM balances b, periods p where p.id = b.period_id order by p.inicio LIMIT $1 OFFSET $2"
	rows, err := Db.Query(stq, lim, offs)
	if err != nil {
            return
	}
        defer rows.Close()
        for rows.Next() {
           balan := BalanceN{}
           if err = rows.Scan(&balan.Id,&balan.PeriodId, &balan.Period,&balan.Amount, &balan.Cuota, &balan.CreatedAt, &balan.UpdatedAt); err != nil {
            }else{
               balances = append(balances, balan)
           }
         }
       return
 }
// -------------------------------------------------------------
// Get all balances in the database and returns the list
  func Balans() (balances []BalanceN, err error) {
        stq :=   "SELECT b.id, b.period_id, p.inicio, b.amount, b.cuota, b.created_at, b.updated_at FROM balances b, periods p where b.period_id = p.id order by p.inicio"
	rows, err := Db.Query(stq)
	if err != nil {
            return
	}
	defer rows.Close()
	for rows.Next() {
            balan := BalanceN{}
            if err = rows.Scan(&balan.Id,&balan.PeriodId, &balan.Period, &balan.Amount, &balan.Cuota, &balan.CreatedAt, &balan.UpdatedAt); err != nil {
                  return
             }
             balances = append(balances, balan)
	}
        return
 }
// -------------------------------------------------------------
