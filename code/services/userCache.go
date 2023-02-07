package services

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type UserService struct {
	cache *cache.Cache
}

var userServices *UserService

func (u UserService) Get(userId string) string {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(userId)
	if !ok {
		return ""
	}
	//list to string
	list := sessionContext.([]string)
	var result string
	for _, v := range list {
		result += v + "------------------------\n"
	}
	return result
}

func (u UserService) Set(userId string, question, reply string) {
	// 列表，最多保存4个
	//如果满了，删除最早的一个
	//如果没有满，直接添加
	listOut := make([]string, 4)
	value := question + "\n" + reply

	raw, ok := u.cache.Get(userId)
	if ok {
		listOut = raw.([]string)
		if len(listOut) == 4 {
			listOut = listOut[1:]
		}
		listOut = append(listOut, value)
	} else {
		listOut = append(listOut, value)
	}
	u.cache.Set(userId, listOut, time.Minute*5)
}

func (u UserService) Clear(userId string) bool {
	u.cache.Delete(userId)
	return true
}

type UserCacheInterface interface {
	Get(userId string) string
	Set(userId string, question, reply string)
	Clear(userId string) bool
}

func GetUserCache() UserCacheInterface {
	if userServices == nil {
		userServices = &UserService{cache: cache.New(10*time.Minute, 10*time.Minute)}
	}
	return userServices
}