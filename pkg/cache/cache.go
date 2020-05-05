package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/STreeChin/contactapi/pkg/config"
	"github.com/STreeChin/contactapi/pkg/entities"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type cache struct {
	log  *logrus.Logger
	pool *redis.Pool
}

//NewCache instance
func NewCache(log *logrus.Logger, cfg config.Config) *cache {
	c := new(cache)
	c.log = log

	c.pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.GetCacheConfig().URL)
			if err != nil {
				return nil, errors.Wrap(err, "redis.Dial")
			}
			/*				if cfg.GetCacheConfig().Password != "" {
							_, err := c.Do("AUTH", cfg.GetCacheConfig().Password)
							if err != nil {
								err = c.Close()
								if err != nil {
									return nil, errors.Wrap(err, "redis.Close")
								}

								return nil, errors.Wrap(err, "redis.Do auth")
							}
						}*/

			return c, nil
		},
	}

	return c
}

func (c *cache) GetEmailByContactID(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	email, err := redis.String(conn.Do("get", key))
	return email, errors.Wrap(err, "redis get email")
}

func (c *cache) SetEmailByContactID(key, value string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("set", key, value)
	return errors.Wrap(err, "redis set email")
}

func (c *cache) GetOneContact(key string) (*entities.Contact, error) {
	conn := c.pool.Get()
	defer conn.Close()

	contact := new(entities.Contact)
	result, err := redis.Bytes(conn.Do("get", key))
	if err != nil {
		return contact, errors.Wrap(err, "redis get contact")
	}
	reader := bytes.NewReader(result)
	err = gob.NewDecoder(reader).Decode(contact)
	return contact, errors.Wrap(err, "redis get contact decode")
}

func (c *cache) SetOneContact(key string, contact *entities.Contact) error {
	conn := c.pool.Get()
	defer conn.Close()

	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(contact)
	if err != nil {
		return errors.Wrap(err, "redis set contact encode")
	}
	_, err = conn.Do("set", key, buffer.Bytes())
	return errors.Wrap(err, "redis set contact")
}

func (c *cache) DelOneContact(key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return errors.Wrap(err, "redis del contact")
}
