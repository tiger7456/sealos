package install

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/fanux/sealos/net"
	"github.com/wonderivan/logger"
)

const defaultConfigPath = "/.sealos"
const defaultConfigFile = "/config.yaml"

// SealConfig for ~/.sealos/config.yaml
type SealConfig struct {
	Masters []string
	Nodes   []string
	//config from kubeadm.cfg. ex. cluster.local
	DnsDomain         string
	ApiServerCertSANs []string

	//SSHConfig
	User       string
	Passwd     string
	PrivateKey string
	PkPassword string
	//ApiServer ex. apiserver.cluster.local
	ApiServerDomain string
	Network         string
	VIP             string
	PkgURL          string
	Version         string
	Repo            string
	PodCIDR         string
	PodIPV6CIDR         string
	SvcCIDR         string
	SvcIPV6CIDR         string
	//certs location
	CertPath     string
	CertEtcdPath string
	//lvscare images
	LvscareName string
	LvscareTag  string
	AliOss
}

//Dump is
func (c *SealConfig) Dump(path string) {
	home, _ := os.UserHomeDir()
	if path == "" {
		path = home + defaultConfigPath + defaultConfigFile
	}
	MasterIPs = ParseIPs(MasterIPs)
	c.Masters = MasterIPs
	NodeIPs = ParseIPs(NodeIPs)
	c.Nodes = ParseIPs(NodeIPs)
	c.User = SSHConfig.User
	c.Passwd = SSHConfig.Password
	c.PrivateKey = SSHConfig.PkFile
	c.PkPassword = SSHConfig.PkPassword
	c.ApiServerDomain = ApiServer
	c.Network = Network
	c.VIP = VIP
	c.PkgURL = PkgUrl
	c.Version = Version
	c.Repo = Repo
	c.SvcCIDR = SvcCIDR
	c.PodCIDR = PodCIDR
	c.SvcIPV6CIDR = SvcIPV6CIDR
	c.PodIPV6CIDR = PodIPV6CIDR

	c.DnsDomain = DnsDomain
	c.ApiServerCertSANs = ApiServerCertSANs
	c.CertPath = CertPath
	c.CertEtcdPath = CertEtcdPath
	//lvscare
	c.LvscareName = LvscareImage.Image
	c.LvscareTag = LvscareImage.Tag
	// oss
	c.AliOss.AccessKeyId = AccessKeyId
	c.AliOss.AccessKeySecrets = AccessKeySecrets
	c.AliOss.OssEndpoint = OssEndpoint
	c.AliOss.BucketName = BucketName
	c.AliOss.ObjectPath = ObjectPath
	y, err := yaml.Marshal(c)
	if err != nil {
		logger.Error("dump config file failed: %s", err)
	}

	err = os.MkdirAll(home+defaultConfigPath, os.ModePerm)
	if err != nil {
		logger.Warn("create default sealos config dir failed, please create it by your self mkdir -p /root/.sealos && touch /root/.sealos/config.yaml")
	}

	if err = ioutil.WriteFile(path, y, 0644); err != nil {
		logger.Warn("write to file %s failed: %s", path, err)
	}
}

func Dump(path string, content interface{}) error {
	y, err := yaml.Marshal(content)
	if err != nil {
		logger.Error("dump config file failed: %s", err)
		return err
	}
	home, _ := os.UserHomeDir()
	err = os.MkdirAll(home+defaultConfigPath, os.ModePerm)
	if err != nil {
		logger.Error("create dump dir failed %s", err)
		return err
	}

	ioutil.WriteFile(path, y, 0644)
	return nil
}

//Load is
func (c *SealConfig) Load(path string) (err error) {
	if path == "" {
		home, _ := os.UserHomeDir()
		path = home + defaultConfigPath + defaultConfigFile
	}

	y, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %s failed %w", path, err)
	}

	err = yaml.Unmarshal(y, c)
	if err != nil {
		return fmt.Errorf("unmarshal config file failed: %w", err)
	}

	MasterIPs = c.Masters
	NodeIPs = c.Nodes
	SSHConfig.User = c.User
	SSHConfig.Password = c.Passwd
	SSHConfig.PkFile = c.PrivateKey
	SSHConfig.PkPassword = c.PkPassword
	ApiServer = c.ApiServerDomain
	Network = c.Network
	VIP = c.VIP
	PkgUrl = c.PkgURL
	Version = c.Version
	Repo = c.Repo
	PodCIDR = c.PodCIDR
	SvcCIDR = c.SvcCIDR
	DnsDomain = c.DnsDomain
	ApiServerCertSANs = c.ApiServerCertSANs
	CertPath = c.CertPath
	CertEtcdPath = c.CertEtcdPath
	//lvscare
	LvscareImage.Image = c.LvscareName
	LvscareImage.Tag = c.LvscareTag

	// 优先使用使用命令行， 再使用配置文件
	if AccessKeyId == "" || AccessKeySecrets == "" ||
		OssEndpoint == "" || BucketName == "" || ObjectPath == "" {
		AccessKeyId = c.AliOss.AccessKeyId
		AccessKeySecrets = c.AliOss.AccessKeySecrets
		OssEndpoint = c.AliOss.OssEndpoint
		BucketName = c.AliOss.BucketName
		ObjectPath = c.AliOss.ObjectPath
	}
	return
}

func Load(path string, content interface{}) error {
	y, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("read config file %s failed %s", path, err)
		os.Exit(0)
	}

	err = yaml.Unmarshal(y, content)
	if err != nil {
		logger.Error("unmarshal config file failed: %s", err)
	}
	return nil
}

func (c *SealConfig) ShowDefaultConfig() {
	home, _ := os.UserHomeDir()
	c.Masters = []string{"192.168.0.2", "192.168.0.2", "192.168.0.2"}
	c.Nodes = []string{"192.168.0.3", "192.168.0.4"}
	c.User = "root"
	c.Passwd = "123456"
	c.PrivateKey = home + "/.ssh/id_rsa"
	c.ApiServerDomain = "apiserver.cluster.local"
	c.Network = net.CALICO
	c.VIP = "10.103.97.2"
	c.PkgURL = home + "/kube1.17.13.tar.gz"
	c.Version = "v1.17.13"
	c.Repo = "k8s.gcr.io"
	c.PodCIDR = "100.64.0.0/10"
	c.PodIPV6CIDR = "2001:1010:43::/108"
	c.SvcCIDR = "10.96.0.0/12"
	c.SvcIPV6CIDR = "2002:1010:43::/108"
	c.ApiServerCertSANs = []string{"apiserver.cluster.local", "127.0.0.1"}
	c.CertPath = home + "/.sealos/pki"
	c.CertEtcdPath = home + "/.sealos/pki/etcd"
	c.LvscareName = "fanux/lvscare"
	c.LvscareTag = "latest"

	y, err := yaml.Marshal(c)
	if err != nil {
		logger.Error("marshal config file failed: %s", err)
	}

	logger.Info("\n\n%s\n\n", string(y))
	logger.Info("Please save above config in ~/.sealos/config.yaml and edit values on your own")
}
