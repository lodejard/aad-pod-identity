package main

import (
	"os"

	"github.com/Azure/aad-pod-identity/pkg/k8s"
	server "github.com/Azure/aad-pod-identity/pkg/nmi/server"
	"github.com/Azure/aad-pod-identity/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	defaultMetadataIP                         = "169.254.169.254"
	defaultMetadataPort                       = "80"
	defaultNmiPort                            = "2579"
	defaultIPTableUpdateTimeIntervalInSeconds = 60
)

var (
	debug                              = pflag.Bool("debug", true, "sets log to debug level")
	versionInfo                        = pflag.Bool("version", false, "prints the version information")
	nmiPort                            = pflag.String("nmi-port", defaultNmiPort, "NMI application port")
	metadataIP                         = pflag.String("metadata-ip", defaultMetadataIP, "instance metadata host ip")
	metadataPort                       = pflag.String("metadata-port", defaultMetadataPort, "instance metadata host ip")
	hostIP                             = pflag.String("host-ip", "", "host IP address")
	nodename                           = pflag.String("node", "", "node name")
	ipTableUpdateTimeIntervalInSeconds = pflag.Int("ipt-update-interval-sec", defaultIPTableUpdateTimeIntervalInSeconds, "update interval of iptables")
	forceNamespaced                    = pflag.Bool("forceNamespaced", false, "Forces mic to namespace identities, binding, and assignment")
	micNamespace                       = pflag.String("MICNamespace", "default", "MIC namespace to short circuit MIC token requests")
)

func main() {
	pflag.Parse()
	if *versionInfo {
		version.PrintVersionAndExit()
	}
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Infof("Starting nmi process. Version: %v. Build date: %v", version.NMIVersion, version.BuildDate)
	client, err := k8s.NewKubeClient()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	*forceNamespaced = *forceNamespaced || "true" == os.Getenv("FORCENAMESPACED")
	s := server.NewServer(*forceNamespaced, *micNamespace)
	s.KubeClient = client
	s.MetadataIP = *metadataIP
	s.MetadataPort = *metadataPort
	s.NMIPort = *nmiPort
	s.HostIP = *hostIP
	s.NodeName = *nodename
	s.IPTableUpdateTimeIntervalInSeconds = *ipTableUpdateTimeIntervalInSeconds

	if err := s.Run(); err != nil {
		log.Fatalf("%s", err)
	}
}
