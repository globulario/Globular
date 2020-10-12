import * as grpcWeb from 'grpc-web';

import * as event_pb from './event_pb';


export class EventServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: event_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: event_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<event_pb.StopResponse>;

  onEvent(
    request: event_pb.OnEventRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<event_pb.OnEventResponse>;

  quit(
    request: event_pb.QuitRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: event_pb.QuitResponse) => void
  ): grpcWeb.ClientReadableStream<event_pb.QuitResponse>;

  subscribe(
    request: event_pb.SubscribeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: event_pb.SubscribeResponse) => void
  ): grpcWeb.ClientReadableStream<event_pb.SubscribeResponse>;

  unSubscribe(
    request: event_pb.UnSubscribeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: event_pb.UnSubscribeResponse) => void
  ): grpcWeb.ClientReadableStream<event_pb.UnSubscribeResponse>;

  publish(
    request: event_pb.PublishRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: event_pb.PublishResponse) => void
  ): grpcWeb.ClientReadableStream<event_pb.PublishResponse>;

}

export class EventServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: event_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<event_pb.StopResponse>;

  onEvent(
    request: event_pb.OnEventRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<event_pb.OnEventResponse>;

  quit(
    request: event_pb.QuitRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<event_pb.QuitResponse>;

  subscribe(
    request: event_pb.SubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<event_pb.SubscribeResponse>;

  unSubscribe(
    request: event_pb.UnSubscribeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<event_pb.UnSubscribeResponse>;

  publish(
    request: event_pb.PublishRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<event_pb.PublishResponse>;

}

