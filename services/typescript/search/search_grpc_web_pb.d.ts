import * as grpcWeb from 'grpc-web';

import {
  CountRequest,
  CountResponse,
  DeleteDocumentRequest,
  DeleteDocumentResponse,
  GetVersionRequest,
  GetVersionResponse,
  IndexDirRequest,
  IndexDirResponse,
  IndexFileRequest,
  IndexFileResponse,
  IndexJsonObjectRequest,
  IndexJsonObjectResponse,
  SearchDocumentsRequest,
  SearchDocumentsResponse,
  StopRequest,
  StopResponse} from './search_pb';

export class SearchServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

  getVersion(
    request: GetVersionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetVersionResponse) => void
  ): grpcWeb.ClientReadableStream<GetVersionResponse>;

  indexJsonObject(
    request: IndexJsonObjectRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: IndexJsonObjectResponse) => void
  ): grpcWeb.ClientReadableStream<IndexJsonObjectResponse>;

  indexFile(
    request: IndexFileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: IndexFileResponse) => void
  ): grpcWeb.ClientReadableStream<IndexFileResponse>;

  indexDir(
    request: IndexDirRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: IndexDirResponse) => void
  ): grpcWeb.ClientReadableStream<IndexDirResponse>;

  count(
    request: CountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CountResponse) => void
  ): grpcWeb.ClientReadableStream<CountResponse>;

  deleteDocument(
    request: DeleteDocumentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteDocumentResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteDocumentResponse>;

  searchDocuments(
    request: SearchDocumentsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SearchDocumentsResponse) => void
  ): grpcWeb.ClientReadableStream<SearchDocumentsResponse>;

}

export class SearchServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  getVersion(
    request: GetVersionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetVersionResponse>;

  indexJsonObject(
    request: IndexJsonObjectRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<IndexJsonObjectResponse>;

  indexFile(
    request: IndexFileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<IndexFileResponse>;

  indexDir(
    request: IndexDirRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<IndexDirResponse>;

  count(
    request: CountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<CountResponse>;

  deleteDocument(
    request: DeleteDocumentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteDocumentResponse>;

  searchDocuments(
    request: SearchDocumentsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SearchDocumentsResponse>;

}

