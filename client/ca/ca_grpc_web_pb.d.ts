import * as grpcWeb from 'grpc-web';

import {
  GetCaCertificateRequest,
  GetCaCertificateResponse,
  SignCertificateRequest,
  SignCertificateResponse} from './ca_pb';

export class CertificateAuthorityClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  signCertificate(
    request: SignCertificateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SignCertificateResponse) => void
  ): grpcWeb.ClientReadableStream<SignCertificateResponse>;

  getCaCertificate(
    request: GetCaCertificateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetCaCertificateResponse) => void
  ): grpcWeb.ClientReadableStream<GetCaCertificateResponse>;

}

export class CertificateAuthorityPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  signCertificate(
    request: SignCertificateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SignCertificateResponse>;

  getCaCertificate(
    request: GetCaCertificateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetCaCertificateResponse>;

}

