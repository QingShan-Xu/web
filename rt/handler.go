package rt

import (
	"fmt"
	"net/http"

	"github.com/QingShan-Xu/web/ds"
	"gorm.io/gorm"
)

func (curRT *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bindData interface{}
	var err error

	if curRT.Bind != nil {
		dataBinder := NewDataBinder()
		bindData, err = dataBinder.BindData(curRT, r)
		if err != nil {
			// todo
			return
		}
	}

	bindReader := ds.NewReader(bindData)

	db := DB.Session(&gorm.Session{})

	scopes := []func(db *gorm.DB) *gorm.DB{}
	for _, scope := range curRT.SCOPES {
		scopes = append(scopes, scope(bindReader))
	}

	db.Scopes(scopes...)

	var aaaa = interface{}(nil)
	if err := db.Find(&aaaa).Error; err != nil {
		print(1)
	}

	fmt.Printf("%+v,%v", bindData, err)

	// if err := h(); err != nil {
	// 	// handle returned error here.
	// 	w.WriteHeader(503)
	// 	w.Write([]byte("bad"))
	// }
}

// func (curRT *Router) Handler(w http.ResponseWriter, r *http.Request) {

// // WHERE语句
// if curRT.WHERE != nil && curRT.MODEL == nil {
// 	log.Fatalf("%s: MODEL required when WHERE not nil", curRT.Path)
// }
// for query, data := range curRT.WHERE {
// 	db = db.Where(query, data)
// }

// }
