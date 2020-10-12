import * as grpcWeb from 'grpc-web';

import * as plc_link_pb from './plc_link_pb';


export class PlcLinkServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: plc_link_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_link_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<plc_link_pb.StopResponse>;

  link(
    request: plc_link_pb.LinkRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_link_pb.LinkRsp) => void
  ): grpcWeb.ClientReadableStream<plc_link_pb.LinkRsp>;

  unLink(
    request: plc_link_pb.UnLinkRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_link_pb.UnLinkRsp) => void
  ): grpcWeb.ClientReadableStream<plc_link_pb.UnLinkRsp>;

  suspend(
    request: plc_link_pb.SuspendRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_link_pb.SuspendRsp) => void
  ): grpcWeb.ClientReadableStream<plc_link_pb.SuspendRsp>;

  resume(
    request: plc_link_pb.ResumeRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: plc_link_pb.ResumeRsp) => void
  ): grpcWeb.ClientReadableStream<plc_link_pb.ResumeRsp>;

}

export class PlcLinkServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: plc_link_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_link_pb.StopResponse>;

  link(
    request: plc_link_pb.LinkRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_link_pb.LinkRsp>;

  unLink(
    request: plc_link_pb.UnLinkRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_link_pb.UnLinkRsp>;

  suspend(
    request: plc_link_pb.SuspendRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_link_pb.SuspendRsp>;

  resume(
    request: plc_link_pb.ResumeRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<plc_link_pb.ResumeRsp>;

}

