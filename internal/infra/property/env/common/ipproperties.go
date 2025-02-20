package common

type IpProperties struct {
	PublicIpDetectorAddresses       []string
	CachedPrivateIpExpireAfterMilis int
	CachedPublicIpExpireAfterMillis int
}

func NewIpProperties() *IpProperties {
	list := make([]string, 4)
	list[0] = "https://checkip.amazonaws.com"
	list[1] = "https://whatismyip.akamai.com"
	list[2] = "https://ifconfig.me/ip"
	list[3] = "https://myip.dnsomatic.com"

	return &IpProperties{
		PublicIpDetectorAddresses:       list,
		CachedPrivateIpExpireAfterMilis: 60 * 1000,
		CachedPublicIpExpireAfterMillis: 60 * 1000,
	}
}
