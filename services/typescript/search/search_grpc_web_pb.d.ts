import * as grpcWeb from 'grpc-web';

import * as search_pb from './search_pb';


export class SearchServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: search_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.StopResponse>;

  getEngineVersion(
    request: search_pb.GetEngineVersionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.GetEngineVersionResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.GetEngineVersionResponse>;

  indexJsonObject(
    request: search_pb.IndexJsonObjectRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.IndexJsonObjectResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.IndexJsonObjectResponse>;

  indexFile(
    request: search_pb.IndexFileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.IndexFileResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.IndexFileResponse>;

  indexDir(
    request: search_pb.IndexDirRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.IndexDirResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.IndexDirResponse>;

  count(
    request: search_pb.CountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.CountResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.CountResponse>;

  deleteDocument(
    request: search_pb.DeleteDocumentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.DeleteDocumentResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.DeleteDocumentResponse>;

  searchDocuments(
    request: search_pb.SearchDocumentsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: search_pb.SearchDocumentsResponse) => void
  ): grpcWeb.ClientReadableStream<search_pb.SearchDocumentsResponse>;

}

export class SearchServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: search_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.StopResponse>;

  getEngineVersion(
    request: search_pb.GetEngineVersionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.GetEngineVersionResponse>;

  indexJsonObject(
    request: search_pb.IndexJsonObjectRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.IndexJsonObjectResponse>;

  indexFile(
    request: search_pb.IndexFileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.IndexFileResponse>;

  indexDir(
    request: search_pb.IndexDirRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.IndexDirResponse>;

  count(
    request: search_pb.CountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.CountResponse>;

  deleteDocument(
    request: search_pb.DeleteDocumentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.DeleteDocumentResponse>;

  searchDocuments(
    request: search_pb.SearchDocumentsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<search_pb.SearchDocumentsResponse>;

}

