package dbseed

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type seed struct {
	fn   func(conn *gorm.DB) error
	conn *gorm.DB
}

func Run(conn *gorm.DB, fn ...func(conn *gorm.DB) error) error {
	conn = conn.Session(&gorm.Session{Logger: logger.Discard})

	seeds := make([]seed, len(fn)+2)
	seeds[0] = seed{
		fn:   seedRole,
		conn: conn,
	}
	seeds[1] = seed{
		fn:   seedUser,
		conn: conn,
	}

	for i, f := range fn {
		seeds[i] = seed{
			fn:   f,
			conn: conn,
		}
	}

	for _, s := range seeds {
		now := time.Now()

		err := s.fn(s.conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Seeded [%s] %s \n", getFnName(s.fn), time.Since(now))
	}

	return nil
}

func getFnName(fn func(*gorm.DB) error) string {
	fullname := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	name := strings.Split(fullname, ".")
	return name[len(name)-1]
}
