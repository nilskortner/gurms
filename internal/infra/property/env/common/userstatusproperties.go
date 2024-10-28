package common

type UserStatusProperties struct {
	cacheUserSessionsStatus        bool
	userSessionsStatusCacheMaxSize int
	userSessionsStatusExpireAfter  int
}

func NewUserStatusProperties() *UserStatusProperties {
	return &UserStatusProperties{
		cacheUserSessionsStatus:        true,
		userSessionsStatusCacheMaxSize: -1,
		userSessionsStatusExpireAfter:  60,
	}
}
