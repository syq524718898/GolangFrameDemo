package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

var addresses = []string{"http://127.0.0.1:9200", "http://127.0.0.1:9201","http://127.0.0.1:9202"}

func main()  {
	Search()
}

func Index() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err!=nil{
		panic(err)
	}
	// Index creates or updates a document in an index
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"title": "你看到外面的世界是什么样的？",
		"content": "外面的世界真的很精彩",
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		panic(err)
	}

	res, err := es.Index("index_test", &buf, es.Index.WithDocumentType("doc"))

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Search() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err!=nil{
		panic(err)
	}

	// info
	res, err := es.Info()
	if err!=nil{
		panic(err)
	}
	fmt.Println(res.String())
	// search - highlight
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "世界",
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags" : []string{"<font color='red'>"},
			"post_tags" : []string{"</font>"},
			"fields" : map[string]interface{}{
				"title" : map[string]interface{}{},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}
	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("index_test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func DeleteByQuery() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	// DeleteByQuery deletes documents matching the provided query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "外面",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}
	index := []string{"index_test"}
	res, err := es.DeleteByQuery(index, &buf)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Delete() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	// Delete removes a document from the index
	res, err := es.Delete("index_test", "1")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

// 创建文档
func Create() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	// Create creates a new document in the index.
	// Returns a 409 response when a document with a same ID already exists in the index.
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"title": "你看到外面的世界是什么样的？",
		"content": "外面的世界真的很精彩",
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		panic(err)
	}
	res, err := es.Create("index_test", "1", &buf, es.Create.WithDocumentType("doc"))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Get() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	res, err := es.Get("index_test", "1")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Update() {

	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	// Update updates a document with a script or partial document.
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"doc": map[string]interface{}{
			"title": "更新你看到外面的世界是什么样的？",
			"content": "更新外面的世界真的很精彩",
		},
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		panic(err)
	}
	res, err := es.Update("index_test", "1", &buf, es.Update.WithDocumentType("doc"))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}


func UpdateByQuery() {
	config := elasticsearch.Config{
		Addresses:             addresses,
		Username:              "",
		Password:              "",
		CloudID:               "",
		APIKey:                "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	// UpdateByQuery performs an update on every document in the index without changing the source,
	// for example to pick up a mapping change.
	index := []string{"demo"}
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "外面",
			},
		},
		// 根据搜索条件更新title
		/*
		   "script": map[string]interface{}{
		       "source": "ctx._source['title']='更新你看到外面的世界是什么样的？'",
		   },
		*/
		// 根据搜索条件更新title、content
		/*
		   "script": map[string]interface{}{
		       "source": "ctx._source=params",
		       "params": map[string]interface{}{
		           "title": "外面的世界真的很精彩",
		           "content": "你看到外面的世界是什么样的？",
		       },
		       "lang": "painless",
		   },
		*/
		// 根据搜索条件更新title、content
		"script": map[string]interface{}{
			"source": "ctx._source.title=params.title;ctx._source.content=params.content;",
			"params": map[string]interface{}{
				"title": "看看外面的世界真的很精彩",
				"content": "他们和你看到外面的世界是什么样的？",
			},
			"lang": "painless",
		},
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		panic(err)
	}
	res, err := es.UpdateByQuery(
		index,
		es.UpdateByQuery.WithDocumentType("index_test"),
		es.UpdateByQuery.WithBody(&buf),
		es.UpdateByQuery.WithContext(context.Background()),
		es.UpdateByQuery.WithPretty(),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}