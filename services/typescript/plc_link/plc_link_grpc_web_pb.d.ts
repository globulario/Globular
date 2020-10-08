import * as grpcWeb from 'grpc-web';

import {
  LinkRqst,
  LinkRsp,
  ResumeRqst,
  ResumeRsp,
  StopRequest,
  StopResponse,
  SuspendRqst,
  SuspendRsp,
  UnLinkRqst,
  UnLinkRsp} from './plc_link_pb';

export class PlcLinkServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

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

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

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

