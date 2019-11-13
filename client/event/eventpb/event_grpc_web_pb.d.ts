import * as grpcWeb from 'grpc-web';

import {
  PublishRequest,
  PublishResponse,
  SubscribeRequest,
  SubscribeResponse,
  UnSubscribeRequest,
  UnSubscribeResponse} from './event_pb';

export class EventServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  subscribe(
    request: SubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<SubscribeResponse>;

  unSubscribe(
    request: UnSubscribeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UnSubscribeResponse) => void
  ): grpcWeb.ClientReadableStream<UnSubscribeResponse>;

  publish(
    request: PublishRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PublishResponse) => void
  ): grpcWeb.ClientReadableStream<PublishResponse>;

}

export class EventServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  subscribe(
    request: SubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<SubscribeResponse>;

  unSubscribe(
    request: UnSubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UnSubscribeResponse>;

  publish(
    request: PublishRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<PublishResponse>;

}

