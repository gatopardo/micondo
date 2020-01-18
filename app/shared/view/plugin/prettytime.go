package plugin

import (
	"html/template"
	"time"
        "fmt"
        "strings"
)

func commas(s string) string {
    if len(s) <= 3 {
        return s
    } else {
        return commas(s[0:len(s)-3]) + "," + s[len(s)-3:]
    }
}

func Format64() template.FuncMap {
      f := make(template.FuncMap)

      f["FORMAT64"] = func(dat uint64)string {
          str  := fmt.Sprintf("%d", dat)
          lon  := len(str)  - 2
          sini := str[:lon]
          sfin := str[lon:]
          scom := commas(sini)
          sres :=  scom + "." + sfin 
          return sres
     }

     return f
   }

func toString(f float64) string {
    parts := strings.Split(fmt.Sprintf("%.2f", f), ".")
    if parts[0][0] == '-' {
        return "-" + commas(parts[0][1:]) + "." + parts[1]
    }
    return commas(parts[0]) + "." + parts[1]
}

// PrettyTime returns a template.FuncMap
// * PRETTYTIME outputs a nice time format
func PrettyTime() template.FuncMap {
	f := make(template.FuncMap)

	f["PRETTYTIME"] = func(t time.Time) string {
		return t.Format("3:04 PM 01/02/2006")
	}

	return f
}

// * DATEFORMAT outputs a nice time format
func DateFormat() template.FuncMap {
	f := make(template.FuncMap)
        layout := "2006/01/02"
	f["DATEFORMAT"] = func(t time.Time) string {
		return t.Format(layout)
	}

	return f
 }

