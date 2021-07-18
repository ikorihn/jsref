package jsref_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	jsref "github.com/r57ty7/jsref"
	"github.com/r57ty7/jsref/provider"
)

func Example() {
	var v interface{}
	f, err := os.Open("testfile/foo.yaml")
	if err != nil {
		log.Printf("ファイルエラー: %v", err)
		return
	}
	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("ファイルreadエラー: %v", err)
		return
	}
	if err := yaml.Unmarshal(src, &v); err != nil {
		log.Printf("%s", err)
		return
	}

	// External reference
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("url parse: %s", err)
		return
	}
	mp := provider.NewFS(wd + "/testfile/")

	res := jsref.New()
	res.AddProvider(mp) // Register the provider

	data := []struct {
		Ptr     string
		Options []jsref.Option
	}{
		{
			Ptr: "#/foo/0", // "bar"
		},
		{
			Ptr: "#/foo/1", // "baz"
		},
		{
			Ptr: "#/foo/2", // "quux" (resolves via `mp`)
		},
		{
			Ptr: "#/foo", // ["bar",{"$ref":"#/sub"},{"$ref":"obj2#/sub"}]
		},
		{
			Ptr: "#/foo", // ["bar","baz","quux"]
			// experimental option to resolve all resulting values
			Options: []jsref.Option{jsref.WithRecursiveResolution(true)},
		},
	}
	for _, set := range data {
		result, err := res.Resolve(v, set.Ptr, set.Options...)
		if err != nil { // failed to resolve
			fmt.Printf("err: %s\n", err)
			continue
		}
		b, _ := json.Marshal(result)
		fmt.Printf("%s -> %s\n", set.Ptr, string(b))
	}

	// OUTPUT:
	// #/foo/0 -> "bar"
	// #/foo/1 -> "baz"
	// #/foo/2 -> "piyo"
	// #/foo -> ["bar",{"$ref":"#/sub"},{"$ref":"file://./sub.yaml#/other"}]
	// #/foo -> ["bar","baz","piyo"]
}
