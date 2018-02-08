package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db}
}

func (s *UserService) Can(userID uint, route string) (ok bool) {
	ps := s.Permissions(userID)
	// 这里直接从permissions中找是否有合适，也可以直接从数据库中搜索
	for _, p := range ps {
		if p.Route == route {
			return true
		}
	}
	return
}

func (s *UserService) Permissions(userID uint) (aps []Permission) {
	u := User{}
	s.db.First(&u, userID)
	roles := []Role{}
	s.db.Model(&u).Related(&roles, "Roles")
	aps = []Permission{}
	for _, r := range roles {
		ps := []Permission{}
		s.db.Model(&r).Related(&ps, "Permissions")
		aps = append(aps, ps...)
	}
	return
}

type User struct {
	gorm.Model
	Name  string `gorm:"name"`
	Roles []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
	Users       []User       `gorm:"many2many:user_roles"`
}

type Permission struct {
	gorm.Model
	Name    string `gorm:"name"`
	Roles   []Role `gorm:"many2many:role_permissions"`
	Subject string `gorm:"subject"`
	Action  string `gorm:"action"`
	Route   string
}

type Article struct {
	gorm.Model
	Title   string `gorm:"title"`
	Content string `gorm:"content"`
}

func NewAuthorizeMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			route := ctx.Request().Method + ctx.Path()
			//log.Printf("%s",ctx.Path())
			if NewUserService(db).Can(1, route) {
				return next(ctx)
			}
			return ctx.JSON(http.StatusForbidden, "no permission")
			return next(ctx)
		}
	}
}

func main() {
	db, err := gorm.Open("mysql", "root:@/rbac?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		panic("can not connect to db")
	}

	// Add table suffix when create tables
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Role{}, &Permission{}, &Article{})

	e := echo.New()
	am := NewAuthorizeMiddleware(db)

	e.GET("/test/:name", index)
	e.POST("/test", index)
	e.Use(am)
	ClearAndPersistRoutes(e, db)
	e.Logger.Fatal(e.Start(":8800"))

}

func ClearAndPersistRoutes(e *echo.Echo, db *gorm.DB) {
	db.DropTable(&Permission{})
	db.CreateTable(&Permission{})
	for _, r := range e.Routes() {
		p := Permission{}
		p.Route = r.Method + r.Path
		db.Create(&p)
	}
}

func index(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, ctx.Echo().Routes())
}
