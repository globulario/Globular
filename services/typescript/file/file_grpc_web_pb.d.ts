import * as grpcWeb from 'grpc-web';

import {
  CreateArchiveRequest,
  CreateArchiveResponse,
  CreateDirRequest,
  CreateDirResponse,
  DeleteDirRequest,
  DeleteDirResponse,
  DeleteFileRequest,
  DeleteFileResponse,
  GetFileInfoRequest,
  GetFileInfoResponse,
  GetThumbnailsRequest,
  GetThumbnailsResponse,
  ReadDirRequest,
  ReadDirResponse,
  ReadFileRequest,
  ReadFileResponse,
  RenameRequest,
  RenameResponse,
  SaveFileRequest,
  SaveFileResponse,
  StopRequest,
  StopResponse,
  WriteExcelFileRequest,
  WriteExcelFileResponse} from './file_pb';

export class FileServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: StopResponse) => void
  ): grpcWeb.ClientReadableStream<StopResponse>;

  readDir(
    request: ReadDirRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ReadDirResponse>;

  createDir(
    request: CreateDirRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateDirResponse) => void
  ): grpcWeb.ClientReadableStream<CreateDirResponse>;

  deleteDir(
    request: DeleteDirRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteDirResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteDirResponse>;

  rename(
    request: RenameRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RenameResponse) => void
  ): grpcWeb.ClientReadableStream<RenameResponse>;

  createAchive(
    request: CreateArchiveRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: CreateArchiveResponse) => void
  ): grpcWeb.ClientReadableStream<CreateArchiveResponse>;

  getFileInfo(
    request: GetFileInfoRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetFileInfoResponse) => void
  ): grpcWeb.ClientReadableStream<GetFileInfoResponse>;

  readFile(
    request: ReadFileRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ReadFileResponse>;

  deleteFile(
    request: DeleteFileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteFileResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteFileResponse>;

  getThumbnails(
    request: GetThumbnailsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GetThumbnailsResponse>;

  writeExcelFile(
    request: WriteExcelFileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: WriteExcelFileResponse) => void
  ): grpcWeb.ClientReadableStream<WriteExcelFileResponse>;

}

export class FileServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  stop(
    request: StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<StopResponse>;

  readDir(
    request: ReadDirRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ReadDirResponse>;

  createDir(
    request: CreateDirRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateDirResponse>;

  deleteDir(
    request: DeleteDirRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteDirResponse>;

  rename(
    request: RenameRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RenameResponse>;

  createAchive(
    request: CreateArchiveRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateArchiveResponse>;

  getFileInfo(
    request: GetFileInfoRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetFileInfoResponse>;

  readFile(
    request: ReadFileRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<ReadFileResponse>;

  deleteFile(
    request: DeleteFileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteFileResponse>;

  getThumbnails(
    request: GetThumbnailsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GetThumbnailsResponse>;

  writeExcelFile(
    request: WriteExcelFileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<WriteExcelFileResponse>;

}

