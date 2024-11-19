package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	log "k8s.io/klog/v2"
)

type flagSlice []string

func (t *flagSlice) String() string {
	return fmt.Sprintf("%v", *t)
}

func (t *flagSlice) Set(value string) error {
	*t = append(*t, value)
	return nil
}

type CmdLineOpts struct {
	etcdEndpoints             string
	etcdPrefix                string
	etcdKeyfile               string
	etcdCertfile              string
	etcdCAFile                string
	etcdUsername              string
	etcdPassword              string
	version                   bool
	kubeSubnetMgr             bool
	kubeApiUrl                string
	kubeAnnotationPrefix      string
	kubeConfigFile            string
	iface                     flagSlice
	ifaceRegex                flagSlice
	ipMasq                    bool
	ifaceCanReach             string
	subnetFile                string
	publicIP                  string
	publicIPv6                string
	subnetLeaseRenewMargin    int
	healthzIP                 string
	healthzPort               int
	iptablesResyncSeconds     int
	iptablesForwardRules      bool
	netConfPath               string
	setNodeNetworkUnavailable bool
}

var (
	opts           CmdLineOpts
	errInterrupted = errors.New("interrupted")
	errCanceled    = errors.New("canceled")
	flatFlags      = flag.NewFlagSet("flat", flag.ExitOnError)
)

func init() {
	flatFlags.StringVar(&opts.etcdEndpoints, "etcd-endpoints", "http://127.0.0.1:4001,http://127.0.0.1:2379", "etcd 엔드포인트 목록")
	flatFlags.StringVar(&opts.etcdPrefix, "etcd-prefix", "/kube-centos/network", "etcd 키의 접두사")
	flatFlags.StringVar(&opts.etcdKeyfile, "etcd-keyfile", "", "etcd 클라이언트 키 파일")
	flatFlags.StringVar(&opts.etcdCertfile, "etcd-certfile", "", "etcd 클라이언트 인증서 파일")
	flatFlags.StringVar(&opts.etcdCAFile, "etcd-cafile", "", "etcd CA 파일")
	flatFlags.StringVar(&opts.etcdUsername, "etcd-username", "", "etcd 사용자 이름")
	flatFlags.StringVar(&opts.etcdPassword, "etcd-password", "", "etcd 비밀번호")
	flatFlags.BoolVar(&opts.version, "version", false, "버전 정보 출력")
	flatFlags.BoolVar(&opts.kubeSubnetMgr, "kube-subnet-mgr", false, "Kubernetes 서브넷 관리자")
	flatFlags.StringVar(&opts.kubeApiUrl, "kube-api-url", "", "Kubernetes API 서버 URL")
	flatFlags.StringVar(&opts.kubeAnnotationPrefix, "kube-annotation-prefix", "kube-centos.network", "Kubernetes 주석 접두사")
	flatFlags.StringVar(&opts.kubeConfigFile, "kube-config-file", "", "Kubernetes 구성 파일")
	flatFlags.Var(&opts.iface, "iface", "사용할 인터페이스")
	flatFlags.Var(&opts.ifaceRegex, "iface-regex", "인터페이스를 선택하는 정규식")
	flatFlags.BoolVar(&opts.ipMasq, "ip-masq", false, "IP 마스커레이드")
	flatFlags.StringVar(&opts.ifaceCanReach, "iface-can-reach", "", "인터페이스가 도달할 수 있는 IP 주소")
	flatFlags.StringVar(&opts.subnetFile, "subnet-file", "", "서브넷 파일")
	flatFlags.StringVar(&opts.publicIP, "public-ip", "", "공용 IP 주소")
	flatFlags.StringVar(&opts.publicIPv6, "public-ipv6", "", "공용 IPv6 주소")
	flatFlags.IntVar(&opts.subnetLeaseRenewMargin, "subnet-lease-renew-margin", 60, "서브넷 임대 갱신 여유 시간")
	flatFlags.StringVar(&opts.healthzIP, "healthz-ip", "0.0.0.0", "healthz 서버가 수신할 IP 주소")
	flatFlags.IntVar(&opts.healthzPort, "healthz-port", 0, "healthz 서버가 수신할 포트")
	flatFlags.IntVar(&opts.iptablesResyncSeconds, "iptables-resync-period", 0, "iptables 재동기화 주기")
	flatFlags.BoolVar(&opts.iptablesForwardRules, "iptables-forward-rules", false, "iptables 전달 규칙")
	flatFlags.StringVar(&opts.netConfPath, "net-conf-path", "/etc/kube-flat/net-conf.json", "네트워크 구성 파일 경로")
	flatFlags.BoolVar(&opts.setNodeNetworkUnavailable, "set-node-network-unavailable", true, "노드 네트워크를 사용할 수 없음으로 설정")

	log.InitFlags(nil)

	err := flag.Set("logtostderr", "true")
	if err != nil {
		log.Error("Can't set the logtostderr flag", err)
		os.Exit(0)
	}

	copyFlag("v")
	copyFlag("vmodule")
	copyFlag("log_backtrace_at")

	flatFlags.Usage = usage

	err = flatFlags.Parse(os.Args[1:])
	if err != nil {
		log.Error("플래그를 파싱할 수 ", err)
		os.Exit(1)
	}

	log.Info("초기화가 완료되었습니다.")
}

func copyFlag(name string) {
	flatFlags.Var(flag.Lookup(name).Value, flag.Lookup(name).Name, flag.Lookup(name).Usage)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flatFlags.PrintDefaults()
	os.Exit(1)
}

func main() {

}
