package main

import (
	"context"
	"github.com/wonderivan/logger"
	"gopkg.in/olivere/elastic.v6"
	"reflect"
	"testgin/api/config"
	"testgin/api/models"
)

var (

	//can match field
	fields = []string{"author_name", "title", "class", "sub_class"}
)

func saveToElastic(info *models.VInfo) error {
	indexServe := config.Set.ElsClient.Index().
		Index(config.Set.Index).Type(config.Set.Type).Id(info.Vid).
		BodyJson(info)
	_, err = indexServe.Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func deleteElasticData(id string) error {
	if _, err := config.Set.ElsClient.Delete().Index(config.Set.Index).Type(config.Set.Type).
		Id(id).Do(context.Background()); err != nil {
		return err
	}
	return nil
}

func searchElastic(matchData string) (interface{}, error) {
	m := elastic.NewMultiMatchQuery(matchData, fields...)
	result, err := config.Set.ElsClient.Search().
		Index(config.Set.Index).Type(config.Set.Type).
		Query(m).
		From(0).Size(15).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		logger.Debug(err)
		return nil, err
	}

	var slice []interface{}
	for _, his := range result.Each(reflect.TypeOf(&models.VideoInfo{})) {
		slice = append(slice, his)
	}
	return slice, nil
}
