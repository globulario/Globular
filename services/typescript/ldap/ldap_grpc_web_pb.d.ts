import * as grpcWeb from 'grpc-web';

import * as ldap_pb from './ldap_pb';


export class LdapServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: ldap_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ldap_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<ldap_pb.StopResponse>;

  createConnection(
    request: ldap_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ldap_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<ldap_pb.CreateConnectionRsp>;

  deleteConnection(
    request: ldap_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ldap_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<ldap_pb.DeleteConnectionRsp>;

  close(
    request: ldap_pb.CloseRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ldap_pb.CloseRsp) => void
  ): grpcWeb.ClientReadableStream<ldap_pb.CloseRsp>;

  search(
    request: ldap_pb.SearchRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ldap_pb.SearchResp) => void
  ): grpcWeb.ClientReadableStream<ldap_pb.SearchResp>;

  authenticate(
    request: ldap_pb.AuthenticateRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ldap_pb.AuthenticateRsp) => void
  ): grpcWeb.ClientReadableStream<ldap_pb.AuthenticateRsp>;

}

export class LdapServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: ldap_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ldap_pb.StopResponse>;

  createConnection(
    request: ldap_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ldap_pb.CreateConnectionRsp>;

  deleteConnection(
    request: ldap_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ldap_pb.DeleteConnectionRsp>;

  close(
    request: ldap_pb.CloseRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ldap_pb.CloseRsp>;

  search(
    request: ldap_pb.SearchRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ldap_pb.SearchResp>;

  authenticate(
    request: ldap_pb.AuthenticateRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ldap_pb.AuthenticateRsp>;

}

