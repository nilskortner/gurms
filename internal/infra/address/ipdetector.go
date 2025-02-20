package address

import (
	"fmt"
	"gurms/internal/infra/property"
	"net"
	"time"
)

var ErrNoAvailableAddressFound error = fmt.Errorf(
	"failed to detect the public IP of the local node because there is no available IP")

type IpDetector struct {
	propertiesManager       *property.GurmsPropertiesManager
	cachedPrivateIp         string
	privateIpLastUpdateDate int64
	cachedPublicIp          string
	publicIpLastUpdateDate  int64
}

func NewIpDetector(propertiesManager *property.GurmsPropertiesManager) *IpDetector {
	return &IpDetector{
		propertiesManager: propertiesManager,
	}
}

func (d *IpDetector) queryPrivateIp() (string, error) {
	ipProperties := d.propertiesManager.LocalGurmsProperties.Ip
	cachedPrivateIpExpireAfterMillis := ipProperties.CachedPrivateIpExpireAfterMilis
	localCachedPrivatedIp := d.cachedPrivateIp
	if cachedPrivateIpExpireAfterMillis > 0 &&
		localCachedPrivatedIp != "" &&
		time.Now().UnixMilli()-d.privateIpLastUpdateDate < int64(cachedPrivateIpExpireAfterMillis) {
		return localCachedPrivatedIp, nil
	}
	conn, err := net.Dial("udp", "8.8.8.8:10002")
	if err != nil {
		return "", fmt.Errorf("failed to detect local IP: %v", err)
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddress.IP.String()
	if !localAddress.IP.IsPrivate() {
		return ip, fmt.Errorf("the new IP address (%s) is not a site local IP address")
	}
	d.privateIpLastUpdateDate = time.Now().UnixMilli()
	d.cachedPrivateIp = ip
	return ip, nil
}

func (d *IpDetector) queryPublicIp() (string, error) {
	ipProperties := d.propertiesManager.LocalGurmsProperties.Ip
	cachedPublicIpExpireAfterMillis := ipProperties.CachedPublicIpExpireAfterMillis
	localCachedPublicIp := d.cachedPublicIp
	if cachedPublicIpExpireAfterMillis > 0 &&
		localCachedPublicIp != "" &&
		time.Now().UnixMilli()-d.privateIpLastUpdateDate < int64(cachedPublicIpExpireAfterMillis) {
		return localCachedPublicIp, nil
	}
	ipDetectorAddresses := ipProperties.PublicIpDetectorAddresses
	if len(ipDetectorAddresses) == 0 {
		return "", fmt.Errorf(
			"failed to detect the public IP of the local node because no IP detector address is specified")
	}

}
