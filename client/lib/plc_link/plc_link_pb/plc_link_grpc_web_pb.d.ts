import * as grpcWeb from 'grpc-web';

import {
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  LinkRqst,
  LinkRsp,
  ResumeRqst,
  ResumeRsp,
  SuspendRqst,
  SuspendRsp,
  UnLinkRqst,
  UnLinkRsp} from './plc_link_pb';

export class PlcLinkServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<DeleteConnectionRsp>;

  link(
    request: LinkRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: LinkRsp) => void
  ): grpcWeb.ClientReadableStream<LinkRsp>;

  unLink(
    request: UnLinkRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UnLinkRsp) => void
  ): grpcWeb.ClientReadableStream<UnLinkRsp>;

  suspend(
    request: SuspendRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SuspendRsp) => void
  ): grpcWeb.ClientReadableStream<SuspendRsp>;

  resume(
    request: ResumeRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ResumeRsp) => void
  ): grpcWeb.ClientReadableStream<ResumeRsp>;

}

export class PlcLinkServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  link(
    request: LinkRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<LinkRsp>;

  unLink(
    request: UnLinkRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<UnLinkRsp>;

  suspend(
    request: SuspendRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SuspendRsp>;

  resume(
    request: ResumeRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<ResumeRsp>;

}

