package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"os"
	"strings"
)

var (
	//dsn = "host=192.168.11.241 user=postgres password=123456 dbname=%s port=5432 sslmode=disable"
	//dsn       = "host=47.100.89.235 user=postgres password=dodo2023 dbname=%s port=5432 sslmode=disable"
	//dsn       = "host=35.241.123.213 user=postgres password=defo2023 dbname=%s port=5432 sslmode=disable"
	tableName = ""
	dbName    = ""
	outpath   = ""
)

func getDsn(dbName string) string {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	psw := os.Getenv("PG_PASSWORD")
	port := os.Getenv("PG_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, psw, dbName, port)
	return dsn
}

func main() {
	//flag.StringVar(&dsn, "dsn", dsn, "数据库地址")
	flag.StringVar(&dbName, "db", dbName, "数据库名")
	flag.StringVar(&tableName, "tb", tableName, "表名")
	flag.StringVar(&outpath, "o", outpath, "输出路径")
	flag.Parse()

	if tableName == "" || outpath == "" {
		panic(fmt.Sprintf("\ntableName:%s,\noutpath:%s\n", tableName, outpath))
	}

	dsn := getDsn(dbName)
	conf := gen.Config{
		ModelPkgPath:   outpath,
		FieldNullable:  true,
		FieldSignable:  true,
		FieldCoverable: false,

		//FieldWithIndexTag: true,
		//FieldWithTypeTag:  true,
	}
	//conf.WithModelNameStrategy(func(tableName string) (modelName string) {
	//	return strings.ToUpper(tableName[0:1]) + tableName[1:]
	//})

	g := gen.NewGenerator(conf)

	fmt.Printf("db dsn=%s\n", dsn)
	db, _ := gorm.Open(postgres.Open(dsn))
	g.UseDB(db)

	rr := g.GenerateModel(tableName,
		gen.FieldGORMTagReg("^.*$", func(tag field.GormTag) field.GormTag {
			tag.Remove("default")
			tag.Remove("comment")
			return tag
		}),
	)
	for _, v := range rr.Fields {
		if v.Type == "int64" && strings.HasSuffix(strings.ToLower(v.ColumnName), "id") {
			v.Tag.Set("json", snakeToCamel(v.ColumnName)+",string")
		} else {
			v.Tag.Set("json", snakeToCamel(v.ColumnName))
		}
	}

	g.Execute()
}

// snakeToCamel 将蛇形字符串转换为驼峰命名
func snakeToCamel(str string) string {
	// 分割字符串
	words := strings.Split(str, "_")

	// 将单词转换为驼峰形式
	for i, word := range words {
		// 只有非第一个单词首字母大写
		if i > 0 {
			words[i] = strings.Title(word)
		} else {
			words[i] = strings.ToLower(word) // 确保第一个单词是小写的
		}
	}

	// 合并单词
	return strings.Join(words, "")
}
