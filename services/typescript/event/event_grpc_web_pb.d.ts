import * as grpcWeb from 'grpc-web';

import {
  OnEventRequest,
  OnEventResponse,
  PublishRequest,
  PublishResponse,
  QuitRequest,
  QuitResponse,
  StopRequest,
  StopResponse,
  SubscribeRequest,
  SubscribeResponse,
  UnSubscribeRequest,
  UnSubscribeResponse} from './event_pb';

export class EventServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

  onEvent(
    request: OnEventRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<OnEventResponse>;

  quit(
    request: QuitRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: QuitResponse) => void
  ): grpcWeb.ClientReadableStream<QuitResponse>;

  subscribe(
    request: SubscribeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SubscribeResponse) => void
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

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  onEvent(
    request: OnEventRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<OnEventResponse>;

  quit(
    request: QuitRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<QuitResponse>;

  subscribe(
    request: SubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SubscribeResponse>;

  unSubscribe(
    request: UnSubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UnSubscribeResponse>;

  publish(
    request: PublishRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<PublishResponse>;

}

