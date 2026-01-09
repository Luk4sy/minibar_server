package redis_article

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"minibar_server/global"
	"minibar_server/utils/date"
	"strconv"
)

type articleCacheType string

const (
	articleLookCache    articleCacheType = "article_look_key"
	articleDiggCache    articleCacheType = "article_digg_key"
	articleCollectCache articleCacheType = "article_collect_key"
)

func set(t articleCacheType, articleID uint, increment bool) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	if !increment {
		num--
	} else {
		num++
	}

	global.Redis.HSet(string(t), strconv.Itoa(int(articleID)), num)
}

func SetLookCache(articleID uint, increment bool) {
	set(articleLookCache, articleID, increment)
}
func SetDiggCache(articleID uint, increment bool) {
	set(articleDiggCache, articleID, increment)
}
func SetCollectCache(articleID uint, increment bool) {
	set(articleCollectCache, articleID, increment)
}

func get(t articleCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	return num
}

func GetLookCache(articleID uint) int {
	return get(articleLookCache, articleID)
}
func GetDiggCache(articleID uint) int {
	return get(articleDiggCache, articleID)
}
func GetCollectCache(articleID uint) int {
	return get(articleCollectCache, articleID)
}

func GetAll(t articleCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(string(t)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, numS := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		iN, err := strconv.Atoi(numS)
		if err != nil {
			continue
		}
		mps[uint(iK)] = iN
	}
	return mps
}

func GetAllLookCache() (mps map[uint]int) {
	return GetAll(articleLookCache)
}
func GetAllDiggCache() (mps map[uint]int) {
	return GetAll(articleDiggCache)
}
func GetAllCollectCache() (mps map[uint]int) {
	return GetAll(articleCollectCache)
}

func SetUserArticleHistoryCache(articleID, userId uint) {
	key := fmt.Sprintf("history_%d", userId)
	field := fmt.Sprintf("%d", articleID)

	endTime := date.GetNowAfter()
	err := global.Redis.HSet(key, field, "").Err()
	if err != nil {
		logrus.Error(err)
		return
	}
	err = global.Redis.ExpireAt(key, endTime).Err()
	if err != nil {
		logrus.Error(err)
		return
	}
}

func GetUserArticleHistoryCache(articleID, userId uint) (ok bool) {
	key := fmt.Sprintf("history_%d", userId)
	field := fmt.Sprintf("%d", articleID)
	err := global.Redis.HGet(key, field).Err()
	if err != nil {
		return false
	}
	return true
}
