package app

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"gorm.io/gorm"
	"grest.dev/grest"
)

const CtxKey = "ctx"

type Ctx struct {
	Lang   string // bahasa yang digunakan oleh user ybs
	Action Action // informasi umum terkait request

	IsAsync bool     // for async use, autocommit
	mainTx  *gorm.DB // for normal use, commit & rollback from middleware
}

type Action struct {
	Method   string
	EndPoint string
	DataID   string
}

// Begin db transaction, dipanggil dari middleware sebelum masuk ke handler.
func (c *Ctx) TxBegin() error {
	mainTx, err := DB().Conn("main")
	if err != nil {
		return err
	}
	c.mainTx = mainTx.Begin()
	return nil
}

// Commit db transaction, dipanggil dari middleware setelah dari handler ketika response nya berhasil (2xx).
func (c *Ctx) TxCommit() {
	if c.mainTx != nil {
		c.mainTx.Commit()
	}

	// reset to nil to use gorm autocommit if use goroutine, etc
	c.mainTx = nil
}

// Rollback db transaction, dipanggil dari middleware setelah dari handler ketika response nya error (selain 2xx).
func (c *Ctx) TxRollback() {
	if c.mainTx != nil {
		c.mainTx.Rollback()
	}
	// reset to nil to use gorm autocommit if use goroutine, etc
	c.mainTx = nil
}

// Trans memberikan translasi atas key dan params sesuai dengan bahasa yang sedang digunakan user.
func (c Ctx) Trans(key string, params ...map[string]string) string {
	return Translator().Trans(c.Lang, key, params...)
}

// ValidateAuth melakukan validasi apakah permintaan dilakukan oleh user yang berwenang atau tidak.
func (c Ctx) ValidatePermission(aclKey string) error {
	// todo
	return nil
}

// ValidateParam melakukan validasi atas payload yang dikirim.
func (c Ctx) ValidateParam(v any) error {
	return Validator().ValidateStruct(v, c.Lang)
}

func (c Ctx) DB(connName ...string) (*gorm.DB, error) {
	if IS_USE_MOCK_DB {
		return Mock().DB()
	}
	// Control the transaction manually (set begin transaction, commit and rollback on middleware)
	if !c.IsAsync && c.mainTx != nil {
		return c.mainTx, nil
	}
	// Autocommit if use goroutine, etc
	return DB().Conn("main")
}

func (c Ctx) NotFoundError(err error, entity, key, value string) error {
	if err != nil && err == gorm.ErrRecordNotFound {
		return NewError(http.StatusNotFound, c.Trans("not_found",
			map[string]string{
				"entity": c.Trans(entity),
				"key":    c.Trans(key),
				"value":  value,
			},
		))
	}
	return nil
}

func (c Ctx) Hook(method, reason, id string, old any) {

	// kasih jeda 2 detik untuk memastikan db transaction nya sudah di commit
	time.Sleep(2 * time.Second)

	isFlat := false
	flat, ok := old.(interface{ IsFlat() bool })
	if ok {
		isFlat = flat.IsFlat()
	}

	oldData := old
	if !isFlat {
		oldData = grest.NewJSON(old).ToStructured().Data
	}
	oldJSON, _ := json.MarshalIndent(oldData, "", "  ")
	newJSON := []byte{}

	model := reflect.ValueOf(old)
	if m := model.MethodByName("Async"); m.IsValid() {
		useCase := m.Call([]reflect.Value{reflect.ValueOf(c)})
		if len(useCase) > 0 {
			if u := useCase[0].MethodByName("GetByID"); u.IsValid() {
				val := u.Call([]reflect.Value{reflect.ValueOf(id)})
				if len(val) > 0 {
					new := val[0].Interface()
					if !isFlat {
						new = grest.NewJSON(new).ToStructured().Data
					}
					newJSON, _ = json.MarshalIndent(new, "", "  ")
				}
			}
		}
	}
	_ = oldJSON
	_ = newJSON
}
