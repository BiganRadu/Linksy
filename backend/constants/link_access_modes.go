package constants

var LinkAccessModes = struct {
	Default     int
	IpWhiteList int
	IpBlackList int
}{
	Default:     0,
	IpWhiteList: 1,
	IpBlackList: 2,
}
