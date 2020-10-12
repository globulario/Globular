import * as grpcWeb from 'grpc-web';

import * as catalog_pb from './catalog_pb';


export class CatalogServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: catalog_pb.StopRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.StopResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.StopResponse>;

  createConnection(
    request: catalog_pb.CreateConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.CreateConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.CreateConnectionRsp>;

  deleteConnection(
    request: catalog_pb.DeleteConnectionRqst,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteConnectionRsp) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteConnectionRsp>;

  saveUnitOfMeasure(
    request: catalog_pb.SaveUnitOfMeasureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveUnitOfMeasureResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveUnitOfMeasureResponse>;

  savePropertyDefinition(
    request: catalog_pb.SavePropertyDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SavePropertyDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SavePropertyDefinitionResponse>;

  saveItemDefinition(
    request: catalog_pb.SaveItemDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveItemDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveItemDefinitionResponse>;

  saveItemInstance(
    request: catalog_pb.SaveItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveItemInstanceResponse>;

  saveInventory(
    request: catalog_pb.SaveInventoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveInventoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveInventoryResponse>;

  saveManufacturer(
    request: catalog_pb.SaveManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveManufacturerResponse>;

  saveSupplier(
    request: catalog_pb.SaveSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveSupplierResponse>;

  saveLocalisation(
    request: catalog_pb.SaveLocalisationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveLocalisationResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveLocalisationResponse>;

  savePackage(
    request: catalog_pb.SavePackageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SavePackageResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SavePackageResponse>;

  savePackageSupplier(
    request: catalog_pb.SavePackageSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SavePackageSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SavePackageSupplierResponse>;

  saveItemManufacturer(
    request: catalog_pb.SaveItemManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveItemManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveItemManufacturerResponse>;

  saveCategory(
    request: catalog_pb.SaveCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.SaveCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.SaveCategoryResponse>;

  appendItemDefinitionCategory(
    request: catalog_pb.AppendItemDefinitionCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.AppendItemDefinitionCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.AppendItemDefinitionCategoryResponse>;

  removeItemDefinitionCategory(
    request: catalog_pb.RemoveItemDefinitionCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.RemoveItemDefinitionCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.RemoveItemDefinitionCategoryResponse>;

  getSupplier(
    request: catalog_pb.GetSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetSupplierResponse>;

  getSuppliers(
    request: catalog_pb.GetSuppliersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetSuppliersResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetSuppliersResponse>;

  getManufacturer(
    request: catalog_pb.GetManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetManufacturerResponse>;

  getManufacturers(
    request: catalog_pb.GetManufacturersRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetManufacturersResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetManufacturersResponse>;

  getSupplierPackages(
    request: catalog_pb.GetSupplierPackagesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetSupplierPackagesResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetSupplierPackagesResponse>;

  getPackage(
    request: catalog_pb.GetPackageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetPackageResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetPackageResponse>;

  getPackages(
    request: catalog_pb.GetPackagesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetPackagesResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetPackagesResponse>;

  getUnitOfMeasure(
    request: catalog_pb.GetUnitOfMeasureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetUnitOfMeasureResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetUnitOfMeasureResponse>;

  getUnitOfMeasures(
    request: catalog_pb.GetUnitOfMeasuresRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetUnitOfMeasuresResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetUnitOfMeasuresResponse>;

  getItemDefinition(
    request: catalog_pb.GetItemDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetItemDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetItemDefinitionResponse>;

  getItemDefinitions(
    request: catalog_pb.GetItemDefinitionsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetItemDefinitionsResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetItemDefinitionsResponse>;

  getItemInstance(
    request: catalog_pb.GetItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetItemInstanceResponse>;

  getItemInstances(
    request: catalog_pb.GetItemInstancesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetItemInstancesResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetItemInstancesResponse>;

  getLocalisation(
    request: catalog_pb.GetLocalisationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetLocalisationResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetLocalisationResponse>;

  getLocalisations(
    request: catalog_pb.GetLocalisationsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetLocalisationsResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetLocalisationsResponse>;

  getCategory(
    request: catalog_pb.GetCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetCategoryResponse>;

  getCategories(
    request: catalog_pb.GetCategoriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetCategoriesResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetCategoriesResponse>;

  getInventories(
    request: catalog_pb.GetInventoriesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.GetInventoriesResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.GetInventoriesResponse>;

  deleteInventory(
    request: catalog_pb.DeleteInventoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteInventoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteInventoryResponse>;

  deletePackage(
    request: catalog_pb.DeletePackageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeletePackageResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeletePackageResponse>;

  deletePackageSupplier(
    request: catalog_pb.DeletePackageSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeletePackageSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeletePackageSupplierResponse>;

  deleteSupplier(
    request: catalog_pb.DeleteSupplierRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteSupplierResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteSupplierResponse>;

  deletePropertyDefinition(
    request: catalog_pb.DeletePropertyDefinitionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeletePropertyDefinitionResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeletePropertyDefinitionResponse>;

  deleteUnitOfMeasure(
    request: catalog_pb.DeleteUnitOfMeasureRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteUnitOfMeasureResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteUnitOfMeasureResponse>;

  deleteItemInstance(
    request: catalog_pb.DeleteItemInstanceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteItemInstanceResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteItemInstanceResponse>;

  deleteManufacturer(
    request: catalog_pb.DeleteManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteManufacturerResponse>;

  deleteItemManufacturer(
    request: catalog_pb.DeleteItemManufacturerRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteItemManufacturerResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteItemManufacturerResponse>;

  deleteCategory(
    request: catalog_pb.DeleteCategoryRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteCategoryResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteCategoryResponse>;

  deleteLocalisation(
    request: catalog_pb.DeleteLocalisationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: catalog_pb.DeleteLocalisationResponse) => void
  ): grpcWeb.ClientReadableStream<catalog_pb.DeleteLocalisationResponse>;

}

export class CatalogServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  stop(
    request: catalog_pb.StopRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.StopResponse>;

  createConnection(
    request: catalog_pb.CreateConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.CreateConnectionRsp>;

  deleteConnection(
    request: catalog_pb.DeleteConnectionRqst,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteConnectionRsp>;

  saveUnitOfMeasure(
    request: catalog_pb.SaveUnitOfMeasureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveUnitOfMeasureResponse>;

  savePropertyDefinition(
    request: catalog_pb.SavePropertyDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SavePropertyDefinitionResponse>;

  saveItemDefinition(
    request: catalog_pb.SaveItemDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveItemDefinitionResponse>;

  saveItemInstance(
    request: catalog_pb.SaveItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveItemInstanceResponse>;

  saveInventory(
    request: catalog_pb.SaveInventoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveInventoryResponse>;

  saveManufacturer(
    request: catalog_pb.SaveManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveManufacturerResponse>;

  saveSupplier(
    request: catalog_pb.SaveSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveSupplierResponse>;

  saveLocalisation(
    request: catalog_pb.SaveLocalisationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveLocalisationResponse>;

  savePackage(
    request: catalog_pb.SavePackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SavePackageResponse>;

  savePackageSupplier(
    request: catalog_pb.SavePackageSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SavePackageSupplierResponse>;

  saveItemManufacturer(
    request: catalog_pb.SaveItemManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveItemManufacturerResponse>;

  saveCategory(
    request: catalog_pb.SaveCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.SaveCategoryResponse>;

  appendItemDefinitionCategory(
    request: catalog_pb.AppendItemDefinitionCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.AppendItemDefinitionCategoryResponse>;

  removeItemDefinitionCategory(
    request: catalog_pb.RemoveItemDefinitionCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.RemoveItemDefinitionCategoryResponse>;

  getSupplier(
    request: catalog_pb.GetSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetSupplierResponse>;

  getSuppliers(
    request: catalog_pb.GetSuppliersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetSuppliersResponse>;

  getManufacturer(
    request: catalog_pb.GetManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetManufacturerResponse>;

  getManufacturers(
    request: catalog_pb.GetManufacturersRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetManufacturersResponse>;

  getSupplierPackages(
    request: catalog_pb.GetSupplierPackagesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetSupplierPackagesResponse>;

  getPackage(
    request: catalog_pb.GetPackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetPackageResponse>;

  getPackages(
    request: catalog_pb.GetPackagesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetPackagesResponse>;

  getUnitOfMeasure(
    request: catalog_pb.GetUnitOfMeasureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetUnitOfMeasureResponse>;

  getUnitOfMeasures(
    request: catalog_pb.GetUnitOfMeasuresRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetUnitOfMeasuresResponse>;

  getItemDefinition(
    request: catalog_pb.GetItemDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetItemDefinitionResponse>;

  getItemDefinitions(
    request: catalog_pb.GetItemDefinitionsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetItemDefinitionsResponse>;

  getItemInstance(
    request: catalog_pb.GetItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetItemInstanceResponse>;

  getItemInstances(
    request: catalog_pb.GetItemInstancesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetItemInstancesResponse>;

  getLocalisation(
    request: catalog_pb.GetLocalisationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetLocalisationResponse>;

  getLocalisations(
    request: catalog_pb.GetLocalisationsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetLocalisationsResponse>;

  getCategory(
    request: catalog_pb.GetCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetCategoryResponse>;

  getCategories(
    request: catalog_pb.GetCategoriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetCategoriesResponse>;

  getInventories(
    request: catalog_pb.GetInventoriesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.GetInventoriesResponse>;

  deleteInventory(
    request: catalog_pb.DeleteInventoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteInventoryResponse>;

  deletePackage(
    request: catalog_pb.DeletePackageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeletePackageResponse>;

  deletePackageSupplier(
    request: catalog_pb.DeletePackageSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeletePackageSupplierResponse>;

  deleteSupplier(
    request: catalog_pb.DeleteSupplierRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteSupplierResponse>;

  deletePropertyDefinition(
    request: catalog_pb.DeletePropertyDefinitionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeletePropertyDefinitionResponse>;

  deleteUnitOfMeasure(
    request: catalog_pb.DeleteUnitOfMeasureRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteUnitOfMeasureResponse>;

  deleteItemInstance(
    request: catalog_pb.DeleteItemInstanceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteItemInstanceResponse>;

  deleteManufacturer(
    request: catalog_pb.DeleteManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteManufacturerResponse>;

  deleteItemManufacturer(
    request: catalog_pb.DeleteItemManufacturerRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteItemManufacturerResponse>;

  deleteCategory(
    request: catalog_pb.DeleteCategoryRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteCategoryResponse>;

  deleteLocalisation(
    request: catalog_pb.DeleteLocalisationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<catalog_pb.DeleteLocalisationResponse>;

}

