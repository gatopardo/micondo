package controller

import (
        "fmt"
        "strings" 
          "strconv"
//	   "github.com/jung-kurt/gofpdf"
  )

      const(
              limit       = 15
              margenlat   = 3
            )
      var (
            TotalCount  int
            offset      int
            posact      int
           )

//-------------------------------------------------------------
func roundU(val float64) int {
    if val > 0 { return int(val+1.0) }
    return int(val)
}
// ---------------------------------------------------
  func atoi32( str string) (nr uint32,err error){
        i, errn := strconv.Atoi(str)
        nr  = uint32(i)
        err =  errn
        return
    }
// ---------------------------------------------------
func rpad(s string,pad string, plength int)string{
    for i:=len(s);i<plength;i++{
        s= s + pad
    }
    return s
}
func  money2uint64(st string)(val uint64, err error){
       parts :=   strings.Split(st, ".")
       fmt.Println(st, parts)
       if len(parts) == 1 {
            parts = append(parts, "00")
        }else  {
            parts[1] = rpad(parts[1], "0",2)
        }
       fmt.Println(st, parts)
       str   :=  strings.Join( parts, "" )
       str    =  strings.Join( strings.Split(str, ","), "" )
       val, err  = strconv.ParseUint(str, 10, 64)
       return 
    }
// ---------------------------------------------------
func getNumberOfButtonsForPagination(TotalCount int, limit int) int {
    num := (int)(TotalCount / limit)
    if (TotalCount%limit > 0) {
        num++
    }
    return num
}
// ---------------------------------------------------
func createSliceForBtns(number int, posact int) []int {
    var sliceOfBtn []int
    lffin := margenlat
    rtini := number   -  margenlat  + 1
    inilf := posact   -  margenlat
    finrt := posact   +  margenlat
    if inilf < 1 {
       inilf = 1
      }
    if finrt > number  {
       finrt =  number
      }
    if lffin  > inilf  {
       lffin  = inilf - 1
    }
    if rtini  < finrt  {
        rtini = finrt  + 1
    }
    for i := 1; i <= lffin; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    for i := inilf; i <= finrt; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    for i := rtini; i <= number; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    return sliceOfBtn
}
// ---------------------------------------------------
