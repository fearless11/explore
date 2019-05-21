package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Stu struct {
	Name string
	ID   int
}

// simpleTest basic template usage
// output: vera ID is 11
func simpleTest() {
	stu := Stu{Name: "vera", ID: 11}
	// 创建模板new, 解析模板字符串parse
	tmpl, err := template.New("test").Parse("{{.Name}} ID is {{ .ID}}")
	if err != nil {
		panic(err)
	}
	// 渲染模板数据
	// 什么数据根据该模板渲染后输出到哪里
	err = tmpl.Execute(os.Stdout, stu)
	if err != nil {
		panic(err)
	}
}

type Fruit struct {
	Name   string
	Status string
}

// sampleTest template object usage
/* output:
第一种：指定模板对象名sample渲染数据
sample: apple is good
第二种: 自定义选择某个模板对象并传递参数
apple is good
*/
func sampleTest() {
	apple := Fruit{"apple", "good"}
	// 包括创建和解析模板文件两步
	t, err := template.ParseFiles("sample.tmpl")
	if err != nil {
		panic(err)
	}
	// 指定模板对象名进行渲染数据
	fmt.Println("\n第一种：指定模板对象名sample渲染数据")
	err = t.ExecuteTemplate(os.Stdout, "sample", apple)
	if err != nil {
		panic(err)
	}

	// 克隆模板对象等价于创建模板,同时指定模板对象,传递参数
	fmt.Println("\n第二种: 自定义选择某个模板对象并传递参数")
	tmpl, err := template.Must(t.Clone()).Parse(`{{ template "apple.message" .}}`)
	if err != nil {
		panic(err)
	}
	// 渲染模板数据
	tmpl.Execute(os.Stdout, apple)
}

// customFuncs template custom funcs usage
/* output:
模板函数
Name have Tom--Vera--Mary
*/
func customFuncs() {
	tmplStr := `Name have {{ block "list.test" .}}{{ joinT . "--"}}{{end}}`
	funcs := template.FuncMap{"joinT": strings.Join}
	data := []string{"Tom", "Vera", "Mary"}

	basicTmpl, err := template.New("basic").Funcs(funcs).Parse(tmplStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n模板函数")
	basicTmpl.Execute(os.Stdout, data)
}

// globParseFile  glob是通配符意思
/* output:
pattern: *.tmpl
通配glob解析模板
just do it
模板克隆
Everything is possible
*/
func globParseFile() {

	pattern := filepath.Join(".", "*.tmpl")
	fmt.Println("\npattern:", pattern)
	fmt.Println("通配glob解析模板")
	tmpl := template.Must(template.ParseGlob(pattern))
	err := tmpl.ExecuteTemplate(os.Stdout, "glob", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n模板克隆")
	t, _ := template.Must(tmpl.Clone()).Parse(`{{ template "good" .}}`)
	t.Execute(os.Stdout, nil)
}

func main() {
	//模板三步： 创建、解析、渲染
	simpleTest()
	sampleTest()
	customFuncs()
	globParseFile()
}
