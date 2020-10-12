import * as grpcWeb from 'grpc-web';

import * as file_pb from './file_pb';


export class FileServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: file_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.StopResponse>;

  readDir(
    request: file_pb.ReadDirRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<file_pb.ReadDirResponse>;

  createDir(
    request: file_pb.CreateDirRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.CreateDirResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.CreateDirResponse>;

  deleteDir(
    request: file_pb.DeleteDirRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.DeleteDirResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.DeleteDirResponse>;

  rename(
    request: file_pb.RenameRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.RenameResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.RenameResponse>;

  createAchive(
    request: file_pb.CreateArchiveRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.CreateArchiveResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.CreateArchiveResponse>;

  getFileInfo(
    request: file_pb.GetFileInfoRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.GetFileInfoResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.GetFileInfoResponse>;

  readFile(
    request: file_pb.ReadFileRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<file_pb.ReadFileResponse>;

  deleteFile(
    request: file_pb.DeleteFileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.DeleteFileResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.DeleteFileResponse>;

  getThumbnails(
    request: file_pb.GetThumbnailsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<file_pb.GetThumbnailsResponse>;

  writeExcelFile(
    request: file_pb.WriteExcelFileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: file_pb.WriteExcelFileResponse) => void
  ): grpcWeb.ClientReadableStream<file_pb.WriteExcelFileResponse>;

}

export class FileServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: file_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.StopResponse>;

  readDir(
    request: file_pb.ReadDirRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<file_pb.ReadDirResponse>;

  createDir(
    request: file_pb.CreateDirRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.CreateDirResponse>;

  deleteDir(
    request: file_pb.DeleteDirRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.DeleteDirResponse>;

  rename(
    request: file_pb.RenameRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.RenameResponse>;

  createAchive(
    request: file_pb.CreateArchiveRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.CreateArchiveResponse>;

  getFileInfo(
    request: file_pb.GetFileInfoRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.GetFileInfoResponse>;

  readFile(
    request: file_pb.ReadFileRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<file_pb.ReadFileResponse>;

  deleteFile(
    request: file_pb.DeleteFileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.DeleteFileResponse>;

  getThumbnails(
    request: file_pb.GetThumbnailsRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<file_pb.GetThumbnailsResponse>;

  writeExcelFile(
    request: file_pb.WriteExcelFileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<file_pb.WriteExcelFileResponse>;

}

