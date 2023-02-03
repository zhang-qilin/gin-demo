/*
* @Time ： 2023-02-03 13:32
* @Auth ： 张齐林
* @File ：ORM_GORM_1&2.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/gorm/schema"
)

var (
	gormDb       *gorm.DB
	gormResponse GormResponse
)

// 特别注意：结构体名称为Product，创建的表名称为products

// Product 定义数据库表对应的结构体
type Product struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Number         string    `gorm:"unique" json:"number"`                       // 商品编号(唯一)
	Category       string    `gorm:"type:varchar(256);not null" json:"category"` // 商品类别
	Name           string    `gorm:"type:varchar(20);not null" json:"name"`      // 商品名称
	MadeIn         string    `gorm:"type:varchar(128);not null" json:"made_in"`  // 生产地
	ProductionTime time.Time `grom:"type:datetime" json:"production_time"`
}

// TableName 获取表名 【如果是要自定义结构体在数据库中所对应的名称可以为相对应的方法中添加一个TableName()的“专属”方法来设置结构体在数据库中所对应的表名称】
// func (Product) TableName() string {
// 	return "zhanqilin"
// }

type GormResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func init() {
	var err error
	sqlStr := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=true&loc=Local"
	gormDb, err = gorm.Open(mysql.Open(sqlStr), &gorm.Config{

		/*
			在我的个人编程规范里面，表名应该是单数而不是复数。名字代表其数据的性质（名称），而其里面的内容注定是不只有一条（不然也不用mysql了对吧）。所以没必要整个加长，增加存储负担（多打几个无意义的字母，你多1我多1，整个世界的网络会被浪费多少流量，硬盘会被浪费多少存储空间啊）。

			而在gorm里面默认的复数形式不是简单的加s或es,是根据其词的意义来加的，对特殊的词进行了特殊处理。通过阅读代码，发现有几个词是不会变形的
			具体可以参考gorm项目下在gorm.Config{}里面有一个配置叫NamingStrategy,和这个配置有关系，通过对这个参数if的判断找到用它的地方，有一个TableName方法
			在inflection.Plural下找到了person变people的秘密
			具体代码如下：

				var irregularInflections = IrregularSlice{
					{"person", "people"},
					{"man", "men"},
					{"child", "children"},
					{"sex", "sexes"},
					{"move", "moves"},
					{"mombie", "mombies"},
				}

				var uncountableInflections = []string{"equipment", "information", "rice", "money", "species", "series", "fish", "sheep", "jeans", "police"}

		*/
		// 在v1的gorm中，关闭表名复数很简单
		// gormDb.SingularTable(true)
		// 在v2的gorm中，稍微复杂一点，需要用到 NamingStategy
		//
		// NamingStrategy: schema.NamingStrategy{
		// 	SingularTable: true,
		// },
		//
	})
	if err != nil {
		fmt.Println("数据库连接出现了问题：", err)
		return
	}

	// 通过GORM的 Migrator 接口的 HasTable() 方法和 Migrator()方法来设置程序启动时判断系统中是否结构体中是否有想相对应的数据表，
	// 没有的话 gormDb.Migrator().HasTable(&结构体名称{}) 和 gormDb.Migrator().HasTable("结构体在数据库中所对应数据表名称")，返回 true 反之返回 false
	/* ========== 友情提示 ==========
		建议创建数据表还是程序员自己设计比较好，
		因为随着时间的推移，
		业务数据量增加到一定的时间时，
		可能会出现数据表操作慢滴问题...
	 ========== 友情提示 ==========*/

	fmt.Println("结构体Product在数据库所对应的数据表: ", gormDb.Migrator().HasTable(&Product{}))
	fmt.Println("结构体Product在数据库所对应的数据表products: ", gormDb.Migrator().HasTable("products"))

	if !gormDb.Migrator().HasTable(&Product{}) || !gormDb.Migrator().HasTable("products") {
		fmt.Println("Product结构体对应的数据表products不存在于数据库中现在开始自动创建先对于的数据表...")
		err = gormDb.Migrator().CreateTable(&Product{})
		if err != nil {
			return
		}
		if gormDb.Migrator().HasTable(&Product{}) || gormDb.Migrator().HasTable("products") {
			fmt.Println("Product结构体对应的数据表products创建成功...")
		} else {
			fmt.Println("表创建失败, err: ", err)
		}
	} else {
		fmt.Println("Product结构体对应的数据表products已存在...")
	}

}

func main() {
	r := gin.Default()
	// 数据库的CRUD 对应Gin的Post、Get、Put、Delete方法
	r.POST("gorm/insert", gormInsertData) // 创建数据
	r.GET("gorm/get", gormGetData)        // 查询数据(单条记录)
	r.Run(":9090")
}

// Catch_Exception 用于捕获异常
func Catch_Exception(c *gin.Context) {
	// ============================= 捕获异常 =============================
	// defer func() {
	err := recover()
	if err != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "错误"
		gormResponse.Data = err
		c.JSON(http.StatusOK, gormResponse)
	}
	// }()
	// ===================================================================
}

func gormGetData(c *gin.Context) {
	defer Catch_Exception(c)
	number := c.Query("number")
	var product Product
	tx := gormDb.Where("number = ?", number).First(&product)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询失败"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}

	gormResponse.Code = http.StatusOK
	gormResponse.Message = "查询成功"
	gormResponse.Data = product
	c.JSON(http.StatusOK, gormResponse)
}

func gormInsertData(c *gin.Context) {
	// // ============================= 捕获异常 =============================
	// defer func() {
	// 	err := recover()
	// 	if err != nil {
	// 		gormResponse.Code = http.StatusBadRequest
	// 		gormResponse.Message = "错误"
	// 		gormResponse.Data = err
	// 		c.JSON(http.StatusOK, gormResponse)
	// 	}
	// }()
	// // ===================================================================
	defer Catch_Exception(c)
	var p Product
	err := c.Bind(&p)
	if err != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "参数错误"
		gormResponse.Data = err
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	tx := gormDb.Create(&p)
	fmt.Println("Insert Failed, err: ", tx.Error)
	fmt.Println(tx)
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "写入成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	fmt.Printf("Insert Failed, err: %s\n\n", tx.Error)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "写入失败"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
	fmt.Println(tx)
}
