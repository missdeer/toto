package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/garyburd/redigo/redis"

	"github.com/missdeer/KellyBackend/modules/models"
)

func RedisGetInt64(key string) (ret int64, err error) {
	ret, err = redis.Int64(Rd.Do("GET", key))
	return ret, err
}

func RedisSetInt64(key string, val int64) (err error) {
	_, err = Rd.Do("SET", key, val)
	return err
}

func RedisGetString(key string) (ret string, err error) {
	ret, err = redis.String(Rd.Do("GET", key))
	return ret, err
}

func RedisSetString(key string, val *string) (err error) {
	_, err = Rd.Do("SET", key, *val)
	return err
}

func RedisGetPosts(key string, posts *[]models.Post) (err error) {
	p, err := redis.Bytes(Rd.Do("GET", key))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(p)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(posts)
	return err
}

func RedisSetPosts(key string, posts *[]models.Post) (err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(posts); err != nil {
		return err
	}

	_, err = Rd.Do("SET", key, buf.Bytes())
	return err
}

func RedisGetTopics(key string, topics *[]models.Topic) (err error) {
	t, err := redis.Bytes(Rd.Do("GET", key))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(t)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&topics)
	return err
}

func RedisSetTopics(key string, topics *[]models.Topic) (err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(&topics); err != nil {
		return err
	}

	_, err = Rd.Do("SET", key, buf.Bytes())
	return err
}

func RedisGetCategories(key string, categories *[]models.Category) (err error) {
	c, err := redis.Bytes(Rd.Do("GET", key))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(c)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&categories)
	return err
}

func RedisSetCategories(key string, categories *[]models.Category) (err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(&categories); err != nil {
		return err
	}
	_, err = Rd.Do("SET", key, buf.Bytes())
	return err
}
