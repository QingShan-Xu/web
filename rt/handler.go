package rt

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Handle func(w http.ResponseWriter, r *http.Request) error

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// handle returned error here.
		w.WriteHeader(503)
		w.Write([]byte("bad"))
	}
}

func (curRT *Router) Handler() http.Handler {
	// bindStruct := ds.ExtendStruct(curRT.Bind)

	// bindStruct.GetField()

	// newBind := bindStruct.Build().New()

	db := DB.Session(&gorm.Session{})
	if curRT.MODEL == nil && curRT.NoAutoMigrate {
		log.Fatalf("%s: MODEL required when NoAutoMigrate is true", curRT.Path)
	}
	// 迁移
	if curRT.NoAutoMigrate {
		if err := DB.AutoMigrate(curRT.MODEL); err != nil {
			log.Fatalf("%s: gorm AutoMigrate err: %v", curRT.Path, err)
		}
	}
	// WHERE语句
	if curRT.WHERE != nil && curRT.MODEL == nil {
		log.Fatalf("%s: MODEL required when WHERE not nil", curRT.Path)
	}
	for query, data := range curRT.WHERE {
		db = db.Where(query, data)
	}

	return nil
}
