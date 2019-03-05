package sessions

import (
	"time"
	"github.com/go-redis/redis"
	"encoding/json"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	rs := RedisStore{
		client,
		sessionDuration,
	}
	return &rs
}

//Store implementation

//saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//marshal the session state to JSON
	sessionStateJson, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	//grab redis key of the session id.
	key := sid.getRedisKey()
	//save the given state into redis store.
	err = rs.Client.Set(key,sessionStateJson,rs.SessionDuration).Err()
	if err != nil {
		return err
	}
	return nil
}

//Fetches session state data from redis store and
//populates the current session state object.
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {

	//create a redis pipeline to send multiple requests in just one network round trip.
	pipe := rs.Client.Pipeline()
	key := sid.getRedisKey()
	//push get and expire commands into the pipe.
	//get the previously-saved session state data from redis,
	cmd := pipe.Get(key)
	//Expire() resets the session expiry time.
	pipe.Expire(key, rs.SessionDuration)
	_, err := pipe.Exec()
	if err == redis.Nil {
		return ErrStateNotFound
	} else if err != nil {
		return err
	}
	jsonString, err := cmd.Result()
	if err != nil {
		return err
	}
	//unmarshal the session state fetched from redis store.
	err = json.Unmarshal([]byte(jsonString), &sessionState)
	if err != nil {
		return err
	}
	return err
}

//removes the entry associated with given session ID.
func (rs *RedisStore) Delete(sid SessionID) error {
	//delete the session state with given sessionID
	err := rs.Client.Del(sid.getRedisKey()).Err()
	if err != nil {
		return err
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
