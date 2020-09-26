package main

import (
	"context"
	"os"

	"os/exec"

	"io/ioutil"

	"log"
	"net"
	"os/signal"

	"strconv"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/ca/capb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (self *Globule) startCertificateAuthorityService() error {
	// The Certificate Authority
	certificate_authority_server, err := self.startInternalService(string(capb.File_services_proto_ca_proto.Services().Get(0).FullName()),
		capb.File_services_proto_ca_proto.Path(), self.CertificateAuthorityPort, self.CertificateAuthorityProxy, false, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.CertificateAuthorityPort))
		if err != nil {
			log.Fatalf("could not certificate authority signing  service %s: %s", self.Name, err)
		}

		capb.RegisterCertificateAuthorityServer(certificate_authority_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {

				// no web-rpc server.
				if err := certificate_authority_server.Serve(lis); err != nil {
					log.Println(err)

				}

			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
		}()
	}
	return err

}

func (self *Globule) signCertificate(client_csr string) (string, error) {

	// first of all I will save the incomming file into a temporary file...
	client_csr_path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()
	err := ioutil.WriteFile(client_csr_path, []byte(client_csr), 0644)
	if err != nil {
		return "", err

	}

	client_crt_path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+self.CertPassword)
	args = append(args, "-days")
	args = append(args, Utility.ToString(self.CertExpirationDelay))
	args = append(args, "-in")
	args = append(args, client_csr_path)
	args = append(args, "-CA")
	args = append(args, self.creds+string(os.PathSeparator)+"ca.crt") // use certificate
	args = append(args, "-CAkey")
	args = append(args, self.creds+string(os.PathSeparator)+"ca.key") // and private key to sign the incommin csr
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, client_crt_path)
	args = append(args, "-extfile")
	args = append(args, self.creds+string(os.PathSeparator)+"san.conf")
	args = append(args, "-extensions")
	args = append(args, "v3_req")
	err = exec.Command(cmd, args...).Run()
	if err != nil {
		return "", err
	}

	// I will read back the crt file.
	client_crt, err := ioutil.ReadFile(client_crt_path)

	// remove the tow temporary files.
	defer os.Remove(client_crt_path)
	defer os.Remove(client_csr_path)

	if err != nil {
		return "", err
	}

	return string(client_crt), nil

}

// Signed certificate request (CSR)
func (self *Globule) SignCertificate(ctx context.Context, rqst *capb.SignCertificateRequest) (*capb.SignCertificateResponse, error) {

	client_crt, err := self.signCertificate(rqst.Csr)

	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

	}

	return &capb.SignCertificateResponse{
		Crt: client_crt,
	}, nil
}

// Return the Authority Trust Certificate. (ca.crt)
func (self *Globule) GetCaCertificate(ctx context.Context, rqst *capb.GetCaCertificateRequest) (*capb.GetCaCertificateResponse, error) {

	ca_crt, err := ioutil.ReadFile(self.creds + string(os.PathSeparator) + "ca.crt")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &capb.GetCaCertificateResponse{
		Ca: string(ca_crt),
	}, nil
}
