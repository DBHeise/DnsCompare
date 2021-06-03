package DnsCompare

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"

	//"crypto/x509"
	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
	"github.com/spf13/viper"

	//"github.com/ameshkov/dnscrypt/v2"

	"github.com/hashicorp/go-multierror"
)

var dnsRequestType uint16

type DNSServer struct {
	Name    string
	Website string
	Notes   string
	Filters bool
	Logs    bool
	DNS     []string
	DOT     []DOTSettings
	DOH     []string
}

type DOTSettings struct {
	Name string
	Port uint16
	IP   []string
	Keys struct {
		SPKI           string
		Thumbprint     string
		PublicKeyRSA   string
		PublicKeyECDSA string
	}
}

func ReadServers(configFile string) ([]DNSServer, error) {
	log.Trace().Str("configFile", configFile).Msg("Reading Server Config File")
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	servers := make([]DNSServer, 0)
	err = decoder.Decode(&servers)

	file.Close()

	if err != nil {
		return nil, err
	}
	return servers, nil
}
func WriteServers(configFile string, servers []DNSServer) error {
	log.Trace().Str("configFile", configFile).Int("ServerCount", len(servers)).Msg("Writing Server Config File")
	file, err := os.OpenFile(configFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(servers)
	if err != nil {
		return err
	}
	return nil
}

func resolve(domain string, dnsserver string, dnsType uint16) ([]string, error) {
	log.Trace().Str("Domain", domain).Str("Server", dnsserver).Uint16("RecordType", dnsType).Msg("Attempting to Resolve Domain")
	options := upstream.Options{
		Timeout:            0,
		InsecureSkipVerify: false,
	}
	up, err := upstream.AddressToUpstream(dnsserver, options)
	if err != nil {
		return nil, err
	}

	req := new(dns.Msg)
	req.SetQuestion(domain, dnsType)
	req.RecursionDesired = true

	res, err := up.Exchange(req)
	if err != nil {
		return nil, err
	}
	log.Trace().Str("DNS Response", res.String()).Msg("Server Response")
	if res.Truncated {
		return nil, errors.New("DNS Response is Truncated")
	}

	results := make([]string, 0)
	for _, record := range res.Answer {
		if a, ok := record.(*dns.A); ok {
			results = append(results, a.A.String())
		} else if aaaa, ok := record.(*dns.AAAA); ok {
			results = append(results, string(aaaa.AAAA))
		} else if mx, ok := record.(*dns.MX); ok {
			results = append(results, mx.Mx)
		} else if ns, ok := record.(*dns.NS); ok {
			results = append(results, ns.Ns)
		} else if cname, ok := record.(*dns.CNAME); ok {
			results = append(results, cname.Target)
		} else if ptr, ok := record.(*dns.PTR); ok {
			results = append(results, ptr.Ptr)
		} else if txt, ok := record.(*dns.TXT); ok {
			results = append(results, txt.Txt...)
		} else if srv, ok := record.(*dns.SRV); ok {
			results = append(results, srv.Target)
		} else if soa, ok := record.(*dns.SOA); ok {
			results = append(results, soa.Ns)
		} else if sig, ok := record.(*dns.SIG); ok {
			results = append(results, sig.Signature)
		}
	}
	return results, nil
}

func resolveList(domain string, dnsserver []string, dnsType uint16) ([]string, error) {
	ans := make([]string, 0)
	var err error
	for _, dnssrv := range dnsserver {
		r, e := resolve(domain, dnssrv, dnsType)
		if e != nil {
			err = multierror.Append(err, e)
		} else {
			ans = append(ans, r...)
		}
	}
	return ans, err
}

func SetDNSType(str string) {
	log.Trace().Str("DNSType", str).Msg("Setting DNS Request Type")
	switch strings.ToLower(str) {
	case "mx":
		dnsRequestType = dns.TypeMX
	case "ns":
		dnsRequestType = dns.TypeNS
	case "cname":
		dnsRequestType = dns.TypeCNAME
	case "ptr":
		dnsRequestType = dns.TypePTR
	case "txt":
		dnsRequestType = dns.TypeTXT
	case "srv":
		dnsRequestType = dns.TypeSRV
	case "soa":
		dnsRequestType = dns.TypeSOA
	case "sig":
		dnsRequestType = dns.TypeSIG
	case "aaaa":
		dnsRequestType = dns.TypeAAAA
	case "a":
		fallthrough
	default:
		dnsRequestType = dns.TypeA
	}
	log.Trace().Uint16("DNSType", dnsRequestType).Msg("Set DNS Request Type")
}

func IsIPv4(str string) bool {
	return net.ParseIP(str).To4() != nil
}
func IsIPv6(str string) bool {
	return net.ParseIP(str).To16() != nil && strings.Contains(str, ":")
}

func chooseServerAddress(method string, list []string) []string {
	switch strings.ToLower(method) {
	case "first":
		return []string{list[0]}
	case "first4":
		for idx := 0; idx < len(list); idx++ {
			if IsIPv4(list[idx]) {
				return []string{list[idx]}
			}
		}
	case "first6":
		for idx := 0; idx < len(list); idx++ {
			if IsIPv6(list[idx]) {
				return []string{list[idx]}
			}
		}
	case "all":
		return list
	case "all4":
		ans := make([]string, 0, len(list))
		for idx := 0; idx < len(list); idx++ {
			if IsIPv4(list[idx]) {
				ans = append(ans, list[idx])
			}
		}
		return ans
	case "all6":
		ans := make([]string, 0, len(list))
		for idx := 0; idx < len(list); idx++ {
			if IsIPv6(list[idx]) {
				ans = append(ans, list[idx])
			}
		}
		return ans
	case "random4":
		all4 := chooseServerAddress("all4", list)
		return []string{RandomFromStringSlice(all4)}
	case "random6":
		all6 := chooseServerAddress("all6", list)
		return []string{RandomFromStringSlice(all6)}
	case "random":
		fallthrough
	default:
		return []string{RandomFromStringSlice(list)}
	}
	return nil
}

func (srv *DNSServer) ResolveUDP(domain string) ([]string, error) {
	log.Trace().Str("Domain", domain).Str("Name", srv.Name).Msg("Attempting UDP DNS Resolution")
	if len(srv.DNS) > 0 {
		return resolveList(domain, chooseServerAddress(viper.GetString("PickMethod"), srv.DNS), dnsRequestType)
	}
	return nil, nil
}
func (srv *DNSServer) ResolveTCP(domain string) ([]string, error) {
	log.Trace().Str("Domain", domain).Str("Name", srv.Name).Msg("Attempting TCP DNS Resolution")
	if len(srv.DNS) > 0 {
		return resolve(domain, "tcp://"+RandomFromStringSlice(srv.DNS), dnsRequestType)
	}
	return nil, nil
}
func (srv *DNSServer) ResolveDOT(domain string) ([]string, error) {
	log.Trace().Str("Domain", domain).Str("Name", srv.Name).Msg("Attempting DoT DNS Resolution")
	if len(srv.DOT) > 0 {
		dot := RandomFromDOTSlice(srv.DOT)
		return resolve(domain, "tls://"+dot.Name+":"+strconv.Itoa(int(dot.Port)), dnsRequestType)
	}
	return nil, nil
}
func (srv *DNSServer) ResolveDOH(domain string) ([]string, error) {
	log.Trace().Str("Domain", domain).Str("Name", srv.Name).Msg("Attempting DoH DNS Resolution")
	if len(srv.DOH) > 0 {
		return resolveList(domain, chooseServerAddress(viper.GetString("PickMethod"), srv.DOH), dnsRequestType)
	}
	return nil, nil
}
