import * as grpcWeb from 'grpc-web';

import {
  AppendItemDefinitionCategoryRequest,
  AppendItemDefinitionCategoryResponse,
  CreateConnectionRqst,
  CreateConnectionRsp,
  DeleteCategoryRequest,
  DeleteCategoryResponse,
  DeleteConnectionRqst,
  DeleteConnectionRsp,
  DeleteItemInstanceRequest,
  DeleteItemInstanceResponse,
  DeleteItemManufacturerRequest,
  DeleteItemManufacturerResponse,
  DeleteLocalisationRequest,
  DeleteLocalisationResponse,
  DeleteManufacturerRequest,
  DeleteManufacturerResponse,
  DeletePackageRequest,
  DeletePackageResponse,
  DeletePackageSupplierRequest,
  DeletePackageSupplierResponse,
  DeletePropertyDefinitionRequest,
  DeletePropertyDefinitionResponse,
  DeleteSupplierRequest,
  DeleteSupplierResponse,
  DeleteUnitOfMeasureRequest,
  DeleteUnitOfMeasureResponse,
  GetPackageRequest,
  GetPackageResponse,
  GetSupplierPackagesRequest,
  GetSupplierPackagesResponse,
  GetSupplierRequest,
  GetSupplierResponse,
  RemoveItemDefinitionCategoryRequest,
  RemoveItemDefinitionCategoryResponse,
  SaveCategoryRequest,
  SaveCategoryResponse,
  SaveInventoryRequest,
  SaveInventoryResponse,
  SaveItemDefinitionRequest,
  SaveItemDefinitionResponse,
  SaveItemInstanceRequest,
  SaveItemInstanceResponse,
  SaveItemManufacturerRequest,
  SaveItemManufacturerResponse,
  SaveLocalisationRequest,
  SaveLocalisationResponse,
  SaveManufacturerRequest,
  SaveManufacturerResponse,
  SavePackageRequest,
  SavePackageResponse,
  SavePackageSupplierRequest,
  SavePackageSupplierResponse,
  SavePropertyDefinitionRequest,
  SavePropertyDefinitionResponse,
  SaveSupplierRequest,
  SaveSupplierResponse,
  SaveUnitOfMeasureRequest,
  SaveUnitOfMeasureResponse} from './catalog_pb';

export class CatalogServiceClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

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

  saveUnitOfMeasure(
    request: SaveUnitOfMeasureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveUnitOfMeasureResponse) => void
  ): grpcWeb.ClientReadableStream<SaveUnitOfMeasureResponse>;

  savePropertyDefinition(
    request: SavePropertyDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SavePropertyDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<SavePropertyDefinitionResponse>;

  saveItemDefinition(
    request: SaveItemDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveItemDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<SaveItemDefinitionResponse>;

  saveItemInstance(
    request: SaveItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<SaveItemInstanceResponse>;

  saveManufacturer(
    request: SaveManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<SaveManufacturerResponse>;

  saveSupplier(
    request: SaveSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<SaveSupplierResponse>;

  saveLocalisation(
    request: SaveLocalisationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveLocalisationResponse) => void
  ): grpcWeb.ClientReadableStream<SaveLocalisationResponse>;

  saveInventory(
    request: SaveInventoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveInventoryResponse) => void
  ): grpcWeb.ClientReadableStream<SaveInventoryResponse>;

  savePackage(
    request: SavePackageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SavePackageResponse) => void
  ): grpcWeb.ClientReadableStream<SavePackageResponse>;

  savePackageSupplier(
    request: SavePackageSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SavePackageSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<SavePackageSupplierResponse>;

  saveItemManufacturer(
    request: SaveItemManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveItemManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<SaveItemManufacturerResponse>;

  saveCategory(
    request: SaveCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<SaveCategoryResponse>;

  appendItemDefinitionCategory(
    request: AppendItemDefinitionCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AppendItemDefinitionCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<AppendItemDefinitionCategoryResponse>;

  removeItemDefinitionCategory(
    request: RemoveItemDefinitionCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: RemoveItemDefinitionCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<RemoveItemDefinitionCategoryResponse>;

  getSupplier(
    request: GetSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<GetSupplierResponse>;

  getSupplierPackages(
    request: GetSupplierPackagesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetSupplierPackagesResponse) => void
  ): grpcWeb.ClientReadableStream<GetSupplierPackagesResponse>;

  getPackage(
    request: GetPackageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetPackageResponse) => void
  ): grpcWeb.ClientReadableStream<GetPackageResponse>;

  deletePackage(
    request: DeletePackageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeletePackageResponse) => void
  ): grpcWeb.ClientReadableStream<DeletePackageResponse>;

  deletePackageSupplier(
    request: DeletePackageSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeletePackageSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<DeletePackageSupplierResponse>;

  deleteSupplier(
    request: DeleteSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteSupplierResponse>;

  deletePropertyDefinition(
    request: DeletePropertyDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeletePropertyDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<DeletePropertyDefinitionResponse>;

  deleteUnitOfMeasure(
    request: DeleteUnitOfMeasureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteUnitOfMeasureResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteUnitOfMeasureResponse>;

  deleteItemInstance(
    request: DeleteItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteItemInstanceResponse>;

  deleteManufacturer(
    request: DeleteManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteManufacturerResponse>;

  deleteItemManufacturer(
    request: DeleteItemManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteItemManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteItemManufacturerResponse>;

  deleteCategory(
    request: DeleteCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteCategoryResponse>;

  deleteLocalisation(
    request: DeleteLocalisationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteLocalisationResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteLocalisationResponse>;

}

export class CatalogServicePromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  createConnection(
    request: CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<CreateConnectionRsp>;

  deleteConnection(
    request: DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteConnectionRsp>;

  saveUnitOfMeasure(
    request: SaveUnitOfMeasureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveUnitOfMeasureResponse>;

  savePropertyDefinition(
    request: SavePropertyDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SavePropertyDefinitionResponse>;

  saveItemDefinition(
    request: SaveItemDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveItemDefinitionResponse>;

  saveItemInstance(
    request: SaveItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveItemInstanceResponse>;

  saveManufacturer(
    request: SaveManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveManufacturerResponse>;

  saveSupplier(
    request: SaveSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveSupplierResponse>;

  saveLocalisation(
    request: SaveLocalisationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveLocalisationResponse>;

  saveInventory(
    request: SaveInventoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveInventoryResponse>;

  savePackage(
    request: SavePackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SavePackageResponse>;

  savePackageSupplier(
    request: SavePackageSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SavePackageSupplierResponse>;

  saveItemManufacturer(
    request: SaveItemManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveItemManufacturerResponse>;

  saveCategory(
    request: SaveCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveCategoryResponse>;

  appendItemDefinitionCategory(
    request: AppendItemDefinitionCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AppendItemDefinitionCategoryResponse>;

  removeItemDefinitionCategory(
    request: RemoveItemDefinitionCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<RemoveItemDefinitionCategoryResponse>;

  getSupplier(
    request: GetSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetSupplierResponse>;

  getSupplierPackages(
    request: GetSupplierPackagesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetSupplierPackagesResponse>;

  getPackage(
    request: GetPackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetPackageResponse>;

  deletePackage(
    request: DeletePackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeletePackageResponse>;

  deletePackageSupplier(
    request: DeletePackageSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeletePackageSupplierResponse>;

  deleteSupplier(
    request: DeleteSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteSupplierResponse>;

  deletePropertyDefinition(
    request: DeletePropertyDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeletePropertyDefinitionResponse>;

  deleteUnitOfMeasure(
    request: DeleteUnitOfMeasureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteUnitOfMeasureResponse>;

  deleteItemInstance(
    request: DeleteItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteItemInstanceResponse>;

  deleteManufacturer(
    request: DeleteManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteManufacturerResponse>;

  deleteItemManufacturer(
    request: DeleteItemManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteItemManufacturerResponse>;

  deleteCategory(
    request: DeleteCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteCategoryResponse>;

  deleteLocalisation(
    request: DeleteLocalisationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteLocalisationResponse>;

}

