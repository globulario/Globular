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
  DeleteInventoryRequest,
  DeleteInventoryResponse,
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
  GetCategoriesRequest,
  GetCategoriesResponse,
  GetCategoryRequest,
  GetCategoryResponse,
  GetInventoriesRequest,
  GetInventoriesResponse,
  GetItemDefinitionRequest,
  GetItemDefinitionResponse,
  GetItemDefinitionsRequest,
  GetItemDefinitionsResponse,
  GetItemInstanceRequest,
  GetItemInstanceResponse,
  GetItemInstancesRequest,
  GetItemInstancesResponse,
  GetLocalisationRequest,
  GetLocalisationResponse,
  GetLocalisationsRequest,
  GetLocalisationsResponse,
  GetManufacturerRequest,
  GetManufacturerResponse,
  GetManufacturersRequest,
  GetManufacturersResponse,
  GetPackageRequest,
  GetPackageResponse,
  GetPackagesRequest,
  GetPackagesResponse,
  GetSupplierPackagesRequest,
  GetSupplierPackagesResponse,
  GetSupplierRequest,
  GetSupplierResponse,
  GetSuppliersRequest,
  GetSuppliersResponse,
  GetUnitOfMeasureRequest,
  GetUnitOfMeasureResponse,
  GetUnitOfMeasuresRequest,
  GetUnitOfMeasuresResponse,
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

  saveInventory(
    request: SaveInventoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SaveInventoryResponse) => void
  ): grpcWeb.ClientReadableStream<SaveInventoryResponse>;

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

  getSuppliers(
    request: GetSuppliersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetSuppliersResponse) => void
  ): grpcWeb.ClientReadableStream<GetSuppliersResponse>;

  getManufacturer(
    request: GetManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<GetManufacturerResponse>;

  getManufacturers(
    request: GetManufacturersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetManufacturersResponse) => void
  ): grpcWeb.ClientReadableStream<GetManufacturersResponse>;

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

  getPackages(
    request: GetPackagesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetPackagesResponse) => void
  ): grpcWeb.ClientReadableStream<GetPackagesResponse>;

  getUnitOfMeasure(
    request: GetUnitOfMeasureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetUnitOfMeasureResponse) => void
  ): grpcWeb.ClientReadableStream<GetUnitOfMeasureResponse>;

  getUnitOfMeasures(
    request: GetUnitOfMeasuresRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetUnitOfMeasuresResponse) => void
  ): grpcWeb.ClientReadableStream<GetUnitOfMeasuresResponse>;

  getItemDefinition(
    request: GetItemDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetItemDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<GetItemDefinitionResponse>;

  getItemDefinitions(
    request: GetItemDefinitionsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetItemDefinitionsResponse) => void
  ): grpcWeb.ClientReadableStream<GetItemDefinitionsResponse>;

  getItemInstance(
    request: GetItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<GetItemInstanceResponse>;

  getItemInstances(
    request: GetItemInstancesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetItemInstancesResponse) => void
  ): grpcWeb.ClientReadableStream<GetItemInstancesResponse>;

  getLocalisation(
    request: GetLocalisationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetLocalisationResponse) => void
  ): grpcWeb.ClientReadableStream<GetLocalisationResponse>;

  getLocalisations(
    request: GetLocalisationsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetLocalisationsResponse) => void
  ): grpcWeb.ClientReadableStream<GetLocalisationsResponse>;

  getCategory(
    request: GetCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<GetCategoryResponse>;

  getCategories(
    request: GetCategoriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetCategoriesResponse) => void
  ): grpcWeb.ClientReadableStream<GetCategoriesResponse>;

  getInventories(
    request: GetInventoriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetInventoriesResponse) => void
  ): grpcWeb.ClientReadableStream<GetInventoriesResponse>;

  deleteInventory(
    request: DeleteInventoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: DeleteInventoryResponse) => void
  ): grpcWeb.ClientReadableStream<DeleteInventoryResponse>;

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

  saveInventory(
    request: SaveInventoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SaveInventoryResponse>;

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

  getSuppliers(
    request: GetSuppliersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetSuppliersResponse>;

  getManufacturer(
    request: GetManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetManufacturerResponse>;

  getManufacturers(
    request: GetManufacturersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetManufacturersResponse>;

  getSupplierPackages(
    request: GetSupplierPackagesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetSupplierPackagesResponse>;

  getPackage(
    request: GetPackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetPackageResponse>;

  getPackages(
    request: GetPackagesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetPackagesResponse>;

  getUnitOfMeasure(
    request: GetUnitOfMeasureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetUnitOfMeasureResponse>;

  getUnitOfMeasures(
    request: GetUnitOfMeasuresRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetUnitOfMeasuresResponse>;

  getItemDefinition(
    request: GetItemDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetItemDefinitionResponse>;

  getItemDefinitions(
    request: GetItemDefinitionsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetItemDefinitionsResponse>;

  getItemInstance(
    request: GetItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetItemInstanceResponse>;

  getItemInstances(
    request: GetItemInstancesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetItemInstancesResponse>;

  getLocalisation(
    request: GetLocalisationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetLocalisationResponse>;

  getLocalisations(
    request: GetLocalisationsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetLocalisationsResponse>;

  getCategory(
    request: GetCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetCategoryResponse>;

  getCategories(
    request: GetCategoriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetCategoriesResponse>;

  getInventories(
    request: GetInventoriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetInventoriesResponse>;

  deleteInventory(
    request: DeleteInventoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<DeleteInventoryResponse>;

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

