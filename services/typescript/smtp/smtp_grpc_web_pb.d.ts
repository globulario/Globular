import * as grpcWeb from 'grpc-web';

import {
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  SendEmailRqst,
  SendEmailRsp,
  SendEmailWithAttachementsRqst,
  SendEmailWithAttachementsRsp,
  StopRequest,
  StopResponse} from './smtp_pb';

export class SmtpServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

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

  sendEmail(
    request: SendEmailRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SendEmailRsp) => void
  ): grpcWeb.ClientReadableStream<SendEmailRsp>;

}

export class SmtpServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  sendEmail(
    request: SendEmailRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<SendEmailRsp>;

}

