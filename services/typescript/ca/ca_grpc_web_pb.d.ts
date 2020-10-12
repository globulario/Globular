import * as grpcWeb from 'grpc-web';

import * as ca_pb from './ca_pb';


export class CertificateAuthorityClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  signCertificate(
    request: ca_pb.SignCertificateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ca_pb.SignCertificateResponse) => void
  ): grpcWeb.ClientReadableStream<ca_pb.SignCertificateResponse>;

  getCaCertificate(
    request: ca_pb.GetCaCertificateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ca_pb.GetCaCertificateResponse) => void
  ): grpcWeb.ClientReadableStream<ca_pb.GetCaCertificateResponse>;

}

export class CertificateAuthorityPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  signCertificate(
    request: ca_pb.SignCertificateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ca_pb.SignCertificateResponse>;

  getCaCertificate(
    request: ca_pb.GetCaCertificateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ca_pb.GetCaCertificateResponse>;

}

