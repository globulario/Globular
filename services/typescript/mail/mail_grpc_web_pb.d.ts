import * as grpcWeb from 'grpc-web';

import * as mail_pb from './mail_pb';


export class MailServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: mail_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: mail_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<mail_pb.StopResponse>;

  createConnection(
    request: mail_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: mail_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<mail_pb.CreateConnectionRsp>;

  deleteConnection(
    request: mail_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: mail_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<mail_pb.DeleteConnectionRsp>;

  sendEmail(
    request: mail_pb.SendEmailRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: mail_pb.SendEmailRsp) => void
  ): grpcWeb.ClientReadableStream<mail_pb.SendEmailRsp>;

}

export class MailServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: mail_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<mail_pb.StopResponse>;

  createConnection(
    request: mail_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<mail_pb.CreateConnectionRsp>;

  deleteConnection(
    request: mail_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<mail_pb.DeleteConnectionRsp>;

  sendEmail(
    request: mail_pb.SendEmailRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<mail_pb.SendEmailRsp>;

}

