package redis

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gmlexx/redigomock"
	"github.com/moira-alert/moira-alert"
	"time"
)

var event = moira_alert.EventData{
	Metric:    "generate.event.1",
	State:     "OK",
	OldState:  "WARN",
	TriggerID: trigger.ID,
}

var trigger = moira_alert.TriggerData{
	ID:         "triggerID-0000000000001",
	Name:       "test trigger 1",
	Targets:    []string{"test.target.1"},
	WarnValue:  10,
	ErrorValue: 20,
	Tags:       []string{"test-tag-1"},
}

var subscription = moira_alert.SubscriptionData{
	ID:                "subscriptionID-00000000000001",
	Enabled:           true,
	Tags:              []string{"test-tag-1"},
	Contacts:          []string{contact.ID},
	ThrottlingEnabled: true,
}

var contact = moira_alert.ContactData{
	ID:    "ContactID-000000000000001",
	Type:  "mega-sender",
	Value: "mail1@example.com",
}

func InitFake(logger moira_alert.Logger) *DbConnector {
	fakeRedis := redigomock.NewFakeRedis()
	expectEvent(fakeRedis)
	pool := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return fakeRedis, nil
		}}

	db := DbConnector{
		logger: logger,
		pool:   &pool,
	}
	c := db.pool.Get()
	defer c.Close()
	c.Do("FLUSHDB")

	testContactString, _ := json.Marshal(contact)
	c.Do("SET", fmt.Sprintf("moira-contact:%s", contact.ID), testContactString)

	testSubscriptionString, _ := json.Marshal(subscription)
	c.Do("SET", fmt.Sprintf("moira-subscription:%s", subscription.ID), testSubscriptionString)
	c.Do("SADD", fmt.Sprintf("moira-tag-subscriptions:%s", subscription.Tags[0]), subscription.ID)

	testTriggerString, _ := json.Marshal(trigger)
	c.Do("SET", fmt.Sprintf("moira-trigger:%s", trigger.ID), testTriggerString)

	for _, tag := range trigger.Tags {
		c.Do("SADD", fmt.Sprintf("moira-trigger-tags:%s", trigger.ID), tag)
	}

	return &db
}

//Duty hack. Need to realize BRPOP command in redigomock
func expectEvent(fakeRedis *redigomock.Conn) {
	eventString, _ := json.Marshal(event)
	alreadyGet := false
	fakeRedis.Command("BRPOP", "moira-trigger-events", 1).ExpectCallback(func(args []interface{}) (interface{}, error) {
		if !alreadyGet {
			result := make([]interface{}, 0)
			result = append(result, []byte("key"))
			result = append(result, eventString)
			alreadyGet = true
			return result, nil
		}
		return nil, nil
	})
}