package DnsCompare

import (
	"encoding/json"
	"errors"
	"os"

	//"crypto/x509"
	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
	//"github.com/ameshkov/dnscrypt/v2"
)

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
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	servers := make([]DNSServer, 0)
	err = decoder.Decode(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}
func WriteServers(configFile string, servers []DNSServer) error {
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

func resolve(address string, dnsserver string, dnsType uint16) ([]string, error) {
	options := upstream.Options{
		Timeout:            0,
		InsecureSkipVerify: false,
	}
	up, err := upstream.AddressToUpstream(dnsserver, options)
	if err != nil {
		return nil, err
	}

	req := new(dns.Msg)
	req.SetQuestion(address, dnsType)
	req.RecursionDesired = true

	res, err := up.Exchange(req)
	if err != nil {
		return nil, err
	}
	if res.Truncated {
		return nil, errors.New("DNS Response is Truncated")
	}

	results := make([]string, 0)
	for _, record := range res.Answer {
		if a, ok := record.(*dns.A); ok {
			results = append(results, a.A.String())
		}
	}
	return results, nil

}

func (srv *DNSServer) ResolveUDP(address string) ([]string, error) {
	if len(srv.DNS) > 0 {
		return resolve(address, srv.DNS[0], dns.TypeA)
	}
	return nil, nil
}

func (srv *DNSServer) ResolveTCP(address string) ([]string, error) {
	if len(srv.DNS) > 0 {
		return resolve(address, "tcp://"+srv.DNS[0], dns.TypeA)
	}
	return nil, nil
}
func (srv *DNSServer) ResolveDOT(address string) ([]string, error) {
	return nil, errors.New("not implemented")
}
func (srv *DNSServer) ResolveDOH(address string) ([]string, error) {
	return nil, errors.New("not implemented")
}
