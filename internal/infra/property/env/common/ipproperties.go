package common

type IpProperties struct {
	publicIpDetectorAddresses       []string
	cachedPrivateIpExpireAfterMilis int
	cachedPublicIpExpireAfterMillis int
}

func NewIpProperties() *IpProperties {
	list := make([]string, 4)
	list[0] = "https://checkip.amazonaws.com"
	list[1] = "https://whatismyip.akamai.com"
	list[2] = "https://ifconfig.me/ip"
	list[3] = "https://myip.dnsomatic.com"

	return &IpProperties{
		publicIpDetectorAddresses:       list,
		cachedPrivateIpExpireAfterMilis: 60 * 1000,
		cachedPublicIpExpireAfterMillis: 60 * 1000,
	}
}
