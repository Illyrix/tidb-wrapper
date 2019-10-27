package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestBlockScanner(t *testing.T) {
	code := `package test

import "fmt"

type student struct {
        Name string
        Age  int
}

func printStu(){
        stus := []student{
                {Name: "zhou", Age: 24},
                {Name: "li", Age: 23},
                {Name: "wang", Age: 22},
        }


        for _, stu := range stus {
                // stu是对象的副本而不是对象本身
                stu.Age = stu.Age+10
        }

        for _, stu := range stus {
                fmt.Println(stu.Age)
        }
}

func hello(){
        print("hello")
}
`

	src := []byte(code)
	fileSet, afile, err := parse(src)
	assert.Equal(t, nil, err)

	expected := [][]int{
		{10, 18},
		{23, 23},
		{18, 21},
		{23, 25},
		{28, 30},
	}

	scanner := NewBlockScanner(fileSet, src, func(startLine, endLine, blockNum int) (insert string) {
		assert.Equal(t, expected[blockNum][0], startLine)
		assert.Equal(t, expected[blockNum][1], endLine)
		return "//aaa"
	})

	ast.Walk(scanner, afile)

	res := string(scanner.edit.Bytes())

	expectedRes := `package test

import "fmt"

type student struct {
        Name string
        Age  int
}

func printStu(){//aaa
        stus := []student{
                {Name: "zhou", Age: 24},
                {Name: "li", Age: 23},
                {Name: "wang", Age: 22},
        }


        for _, stu := range stus {//aaa
                // stu是对象的副本而不是对象本身
                stu.Age = stu.Age+10
        }

        //aaafor _, stu := range stus {//aaa
                fmt.Println(stu.Age)
        }
}

func hello(){//aaa
        print("hello")
}
`
	assert.Equal(t, expectedRes, res)
}

func TestSimplePrint(t *testing.T) {
	code := `package cmd

import (
	"strconv"
	"wrapper_test/trace_util_0"
)

func hello(a int, digest_0 string) string {//[0]
	b := 10

	if a > 10 {//[2]
		a++
	} else { //[3]
		a--
	}

	b-- // [1]

	return strconv.Itoa(a + b)
}
`

	src := []byte(code)
	fileSet, afile, err := parse(src)
	assert.Equal(t, nil, err)
	scanner := NewBlockScanner(fileSet, src, func(startLine, endLine, blockNum int) (insert string) {
		fmt.Printf("block: %d, start: %d, end: %d\n", blockNum, startLine, endLine)
		return "---aaa---"
	})

	ast.Walk(scanner, afile)

	fmt.Println(string(scanner.edit.Bytes()))
}