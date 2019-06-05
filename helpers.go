package main

import (
	"github.com/adonese/noebs/dashboard"
	"github.com/adonese/noebs/ebs_fields"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

func database(dialect string, fname string) *gorm.DB {
	db, err := gorm.Open(dialect, fname)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"details": "there's an error in connecting to DB",
		}).Info("there is an error in connecting to DB")
	}

	db.AutoMigrate(&dashboard.Transaction{})

	return db
}

type redisPurchaseFields map[string]interface{}


func structFields(s interface{}) (fields []map[string]interface{}){
	val := reflect.Indirect(reflect.ValueOf(s))
	// val is a reflect.Type now

	t := val.Type()

	for i := 0; i <= t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("json")
		tag := f.Tag.Get("binding")
		sType := f.Type

		if tag == ""{
			tag = "optional"
		}
		field := map[string]interface{}{
			"field": name,
			"is_required": tag,
			"type": sType,
		}
		fields = append(fields, field)

	}
	return fields
}

//endpointToFields the corresponding struct field for the endpoint string.
// we use simple string matching to capture the results
func endpointToFields(e string) ( interface{}) {
	if strings.Contains(e, "cashIn") {
		return ebs_fields.CashInFields{}
	}
	if strings.Contains(e, "cashOut") {
		return  ebs_fields.CashOutFields{}
	}
	if strings.Contains(e, "balance") {
		return  ebs_fields.BalanceFields{}
	}
	if strings.Contains(e, "billPayment") {
		return ebs_fields.BillPaymentFields{}
	}
	if strings.Contains(e, "cardTransfer") {
		return ebs_fields.CardTransferFields{}
	}
	if strings.Contains(e, "changePin") {
		return ebs_fields.ChangePINFields{}
	}
	if strings.Contains(e, "purchase") {
		return ebs_fields.PurchaseFields{}
	}
	return ebs_fields.GenericEBSResponseFields{}
}

//generateDoc generates API doc for this particular field
func generateDoc(e string) []map[string]interface{} {

	fields := endpointToFields(e)
	scheme := structFields(&fields)

	return scheme
}

//func redisOrNew(key string, data []map[string]interface{}) (string, error){
//	routes := getAllRoutes()
//
//	client := getRedis()
//
//	v, err := client.HMGet("doc")
//	if err == redis.Nil {
//		for _, r := range routes {
//			// get the particular fields for this route
//			doc := generateDoc(r["path"])
//			b, _ := json.Marshal(&r)
//			client.HSet(routes["path"], b)
//		}
//
//	}
//	client.Close()
//
//}

//getAllRoutes gets all routes for this particular engine
// perhaps, it could better be rewritten to explicitly show that
func getAllRoutes() []map[string]string{
	e := GetMainEngine()
	endpoints := e.Routes()
	var allRoutes []map[string]string
	for _, r := range endpoints{
		name := strings.TrimPrefix(r.Path, "/")
		mapping := map[string]string{
			"http_method": r.Method,
			"path": name,
		}
		allRoutes = append(allRoutes, mapping)
	}
	return allRoutes
}

func getRedis() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:               "127.0.0.1:6379",
		DB:                 0,

	})
	return client
}