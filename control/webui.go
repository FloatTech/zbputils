package control

import (
	"os"
	"time"

	"github.com/FloatTech/floatbox/process"
	sql "github.com/FloatTech/sqlite"
)

// User webui用户数据
type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

var (
	udb       = &sql.Sqlite{}
	udbFolder = "data/webui/"
)

func init() {
	_ = os.MkdirAll(udbFolder, 0755)
	udb.DBPath = udbFolder + "user.db"
	err := udb.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	go func() {
		process.GlobalInitMutex.Lock()
		process.SleepAbout1sTo2s()
		err := udb.Create("user", &User{})
		if err != nil {
			panic(err)
		}
	}()
}

// CreateOrUpdateUser 创建或修改用户密码
func CreateOrUpdateUser(u User) error {
	var fu User
	err := udb.Find("user", &fu, "WHERE username = '"+u.Username+"' AND password = '"+u.Password+"'")
	canFind := err == nil && fu.Username == u.Username
	if canFind {
		err = udb.Del("user", "WHERE username = '"+u.Username+"'")
		if err != nil {
			return err
		}
	}
	err = udb.Insert("user", &u)
	return err
}

// FindUser 查找webui账号
func FindUser(username, password string) (u User, err error) {
	err = udb.Find("user", &u, "WHERE username = '"+username+"' AND password = '"+password+"'")
	return
}
