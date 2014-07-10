package post

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/missdeer/KellyBackend/cache"
	"github.com/missdeer/KellyBackend/modules/models"
)

func MemcachedGetInt64(key string) (ret int64, err error) {
	var val *memcache.Item
	if val, err = cache.Mc.Get(key); err != nil {
		return 0, err
	}

	ret, err = strconv.ParseInt(string(val.Value), 10, 64)
	return ret, err
}

func MemcachedSetInt64(key string, val int64) (err error) {
	buf := []byte(strconv.FormatInt(val, 10))
	err = cache.Mc.Set(&memcache.Item{Key: key, Value: buf})
	return err
}

func MemcachedGetPosts(key string, posts *[]models.Post) (err error) {
	var p *memcache.Item
	if p, err = cache.Mc.Get(key); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(p.Value)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(posts)
	return err
}

func MemcachedSetPosts(key string, posts *[]models.Post) (err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(posts); err != nil {
		return err
	}
	PostsCache := &memcache.Item{Key: key, Value: buf.Bytes()}
	err = cache.Mc.Set(PostsCache)
	return err
}

func MemcachedGetTopics(key string, topics *[]models.Topic) (err error) {
	var t *memcache.Item
	if t, err = cache.Mc.Get(key); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(t.Value)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&topics)
	return err
}

func MemcachedSetTopics(key string, topics *[]models.Topic) (err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(&topics); err != nil {
		return err
	}

	TopicsCache := &memcache.Item{Key: key, Value: buf.Bytes()}
	err = cache.Mc.Set(TopicsCache)
	return err
}

func MemcachedGetCategories(key string, categories *[]models.Category) (err error) {
	var c *memcache.Item
	if c, err = cache.Mc.Get(key); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(c.Value)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&categories)
	return err
}

func MemcachedSetCategories(key string, categories *[]models.Category) (err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(&categories); err != nil {
		return err
	}
	err = cache.Mc.Set(&memcache.Item{Key: key, Value: buf.Bytes()})
	return err
}
