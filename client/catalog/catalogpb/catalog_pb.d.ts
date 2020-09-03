import * as jspb from "google-protobuf"

export class Reference extends jspb.Message {
  getRefcolid(): string;
  setRefcolid(value: string): void;

  getRefobjid(): string;
  setRefobjid(value: string): void;

  getRefdbname(): string;
  setRefdbname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Reference.AsObject;
  static toObject(includeInstance: boolean, msg: Reference): Reference.AsObject;
  static serializeBinaryToWriter(message: Reference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Reference;
  static deserializeBinaryFromReader(message: Reference, reader: jspb.BinaryReader): Reference;
}

export namespace Reference {
  export type AsObject = {
    refcolid: string,
    refobjid: string,
    refdbname: string,
  }
}

export class References extends jspb.Message {
  getValuesList(): Array<Reference>;
  setValuesList(value: Array<Reference>): void;
  clearValuesList(): void;
  addValues(value?: Reference, index?: number): Reference;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): References.AsObject;
  static toObject(includeInstance: boolean, msg: References): References.AsObject;
  static serializeBinaryToWriter(message: References, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): References;
  static deserializeBinaryFromReader(message: References, reader: jspb.BinaryReader): References;
}

export namespace References {
  export type AsObject = {
    valuesList: Array<Reference.AsObject>,
  }
}

export class Connection extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getHost(): string;
  setHost(value: string): void;

  getStore(): StoreType;
  setStore(value: StoreType): void;

  getUser(): string;
  setUser(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  getTimeout(): number;
  setTimeout(value: number): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Connection.AsObject;
  static toObject(includeInstance: boolean, msg: Connection): Connection.AsObject;
  static serializeBinaryToWriter(message: Connection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Connection;
  static deserializeBinaryFromReader(message: Connection, reader: jspb.BinaryReader): Connection;
}

export namespace Connection {
  export type AsObject = {
    id: string,
    name: string,
    host: string,
    store: StoreType,
    user: string,
    password: string,
    port: number,
    timeout: number,
    options: string,
  }
}

export class CreateConnectionRqst extends jspb.Message {
  getConnection(): Connection | undefined;
  setConnection(value?: Connection): void;
  hasConnection(): boolean;
  clearConnection(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CreateConnectionRqst): CreateConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: CreateConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateConnectionRqst;
  static deserializeBinaryFromReader(message: CreateConnectionRqst, reader: jspb.BinaryReader): CreateConnectionRqst;
}

export namespace CreateConnectionRqst {
  export type AsObject = {
    connection?: Connection.AsObject,
  }
}

export class CreateConnectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CreateConnectionRsp): CreateConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: CreateConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateConnectionRsp;
  static deserializeBinaryFromReader(message: CreateConnectionRsp, reader: jspb.BinaryReader): CreateConnectionRsp;
}

export namespace CreateConnectionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteConnectionRqst extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteConnectionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteConnectionRqst): DeleteConnectionRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteConnectionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteConnectionRqst;
  static deserializeBinaryFromReader(message: DeleteConnectionRqst, reader: jspb.BinaryReader): DeleteConnectionRqst;
}

export namespace DeleteConnectionRqst {
  export type AsObject = {
    id: string,
  }
}

export class DeleteConnectionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteConnectionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteConnectionRsp): DeleteConnectionRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteConnectionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteConnectionRsp;
  static deserializeBinaryFromReader(message: DeleteConnectionRsp, reader: jspb.BinaryReader): DeleteConnectionRsp;
}

export namespace DeleteConnectionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class PropertyDefinition extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getAbreviation(): string;
  setAbreviation(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getType(): PropertyDefinition.Type;
  setType(value: PropertyDefinition.Type): void;

  getProperties(): References | undefined;
  setProperties(value?: References): void;
  hasProperties(): boolean;
  clearProperties(): void;

  getChoicesList(): Array<string>;
  setChoicesList(value: Array<string>): void;
  clearChoicesList(): void;
  addChoices(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PropertyDefinition.AsObject;
  static toObject(includeInstance: boolean, msg: PropertyDefinition): PropertyDefinition.AsObject;
  static serializeBinaryToWriter(message: PropertyDefinition, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PropertyDefinition;
  static deserializeBinaryFromReader(message: PropertyDefinition, reader: jspb.BinaryReader): PropertyDefinition;
}

export namespace PropertyDefinition {
  export type AsObject = {
    id: string,
    name: string,
    languagecode: string,
    abreviation: string,
    description: string,
    type: PropertyDefinition.Type,
    properties?: References.AsObject,
    choicesList: Array<string>,
  }

  export enum Type { 
    numerical = 0,
    textual = 1,
    boolean = 2,
    dimentional = 3,
    enumeration = 4,
    aggregate = 5,
  }
}

export class PropertyDefinitions extends jspb.Message {
  getValuesList(): Array<PropertyDefinition>;
  setValuesList(value: Array<PropertyDefinition>): void;
  clearValuesList(): void;
  addValues(value?: PropertyDefinition, index?: number): PropertyDefinition;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PropertyDefinitions.AsObject;
  static toObject(includeInstance: boolean, msg: PropertyDefinitions): PropertyDefinitions.AsObject;
  static serializeBinaryToWriter(message: PropertyDefinitions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PropertyDefinitions;
  static deserializeBinaryFromReader(message: PropertyDefinitions, reader: jspb.BinaryReader): PropertyDefinitions;
}

export namespace PropertyDefinitions {
  export type AsObject = {
    valuesList: Array<PropertyDefinition.AsObject>,
  }
}

export class ItemDefinition extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getAbreviation(): string;
  setAbreviation(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getAliasList(): Array<string>;
  setAliasList(value: Array<string>): void;
  clearAliasList(): void;
  addAlias(value: string, index?: number): void;

  getKeywordsList(): Array<string>;
  setKeywordsList(value: Array<string>): void;
  clearKeywordsList(): void;
  addKeywords(value: string, index?: number): void;

  getProperties(): References | undefined;
  setProperties(value?: References): void;
  hasProperties(): boolean;
  clearProperties(): void;

  getReleadeditemdefintions(): References | undefined;
  setReleadeditemdefintions(value?: References): void;
  hasReleadeditemdefintions(): boolean;
  clearReleadeditemdefintions(): void;

  getEquivalentsitemdefintions(): References | undefined;
  setEquivalentsitemdefintions(value?: References): void;
  hasEquivalentsitemdefintions(): boolean;
  clearEquivalentsitemdefintions(): void;

  getCategories(): References | undefined;
  setCategories(value?: References): void;
  hasCategories(): boolean;
  clearCategories(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ItemDefinition.AsObject;
  static toObject(includeInstance: boolean, msg: ItemDefinition): ItemDefinition.AsObject;
  static serializeBinaryToWriter(message: ItemDefinition, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ItemDefinition;
  static deserializeBinaryFromReader(message: ItemDefinition, reader: jspb.BinaryReader): ItemDefinition;
}

export namespace ItemDefinition {
  export type AsObject = {
    id: string,
    name: string,
    languagecode: string,
    abreviation: string,
    description: string,
    aliasList: Array<string>,
    keywordsList: Array<string>,
    properties?: References.AsObject,
    releadeditemdefintions?: References.AsObject,
    equivalentsitemdefintions?: References.AsObject,
    categories?: References.AsObject,
  }
}

export class AppendItemDefinitionCategoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getCategory(): Reference | undefined;
  setCategory(value?: Reference): void;
  hasCategory(): boolean;
  clearCategory(): void;

  getItemdefinition(): Reference | undefined;
  setItemdefinition(value?: Reference): void;
  hasItemdefinition(): boolean;
  clearItemdefinition(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppendItemDefinitionCategoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AppendItemDefinitionCategoryRequest): AppendItemDefinitionCategoryRequest.AsObject;
  static serializeBinaryToWriter(message: AppendItemDefinitionCategoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppendItemDefinitionCategoryRequest;
  static deserializeBinaryFromReader(message: AppendItemDefinitionCategoryRequest, reader: jspb.BinaryReader): AppendItemDefinitionCategoryRequest;
}

export namespace AppendItemDefinitionCategoryRequest {
  export type AsObject = {
    connectionid: string,
    category?: Reference.AsObject,
    itemdefinition?: Reference.AsObject,
  }
}

export class AppendItemDefinitionCategoryResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppendItemDefinitionCategoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AppendItemDefinitionCategoryResponse): AppendItemDefinitionCategoryResponse.AsObject;
  static serializeBinaryToWriter(message: AppendItemDefinitionCategoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppendItemDefinitionCategoryResponse;
  static deserializeBinaryFromReader(message: AppendItemDefinitionCategoryResponse, reader: jspb.BinaryReader): AppendItemDefinitionCategoryResponse;
}

export namespace AppendItemDefinitionCategoryResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class RemoveItemDefinitionCategoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getCategory(): Reference | undefined;
  setCategory(value?: Reference): void;
  hasCategory(): boolean;
  clearCategory(): void;

  getItemdefinition(): Reference | undefined;
  setItemdefinition(value?: Reference): void;
  hasItemdefinition(): boolean;
  clearItemdefinition(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveItemDefinitionCategoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveItemDefinitionCategoryRequest): RemoveItemDefinitionCategoryRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveItemDefinitionCategoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveItemDefinitionCategoryRequest;
  static deserializeBinaryFromReader(message: RemoveItemDefinitionCategoryRequest, reader: jspb.BinaryReader): RemoveItemDefinitionCategoryRequest;
}

export namespace RemoveItemDefinitionCategoryRequest {
  export type AsObject = {
    connectionid: string,
    category?: Reference.AsObject,
    itemdefinition?: Reference.AsObject,
  }
}

export class RemoveItemDefinitionCategoryResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveItemDefinitionCategoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveItemDefinitionCategoryResponse): RemoveItemDefinitionCategoryResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveItemDefinitionCategoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveItemDefinitionCategoryResponse;
  static deserializeBinaryFromReader(message: RemoveItemDefinitionCategoryResponse, reader: jspb.BinaryReader): RemoveItemDefinitionCategoryResponse;
}

export namespace RemoveItemDefinitionCategoryResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class UnitOfMeasure extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getAbreviation(): string;
  setAbreviation(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnitOfMeasure.AsObject;
  static toObject(includeInstance: boolean, msg: UnitOfMeasure): UnitOfMeasure.AsObject;
  static serializeBinaryToWriter(message: UnitOfMeasure, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnitOfMeasure;
  static deserializeBinaryFromReader(message: UnitOfMeasure, reader: jspb.BinaryReader): UnitOfMeasure;
}

export namespace UnitOfMeasure {
  export type AsObject = {
    id: string,
    name: string,
    languagecode: string,
    abreviation: string,
    description: string,
  }
}

export class Category extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getCategories(): References | undefined;
  setCategories(value?: References): void;
  hasCategories(): boolean;
  clearCategories(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Category.AsObject;
  static toObject(includeInstance: boolean, msg: Category): Category.AsObject;
  static serializeBinaryToWriter(message: Category, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Category;
  static deserializeBinaryFromReader(message: Category, reader: jspb.BinaryReader): Category;
}

export namespace Category {
  export type AsObject = {
    id: string,
    name: string,
    languagecode: string,
    categories?: References.AsObject,
  }
}

export class Localisation extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getSublocalisations(): References | undefined;
  setSublocalisations(value?: References): void;
  hasSublocalisations(): boolean;
  clearSublocalisations(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Localisation.AsObject;
  static toObject(includeInstance: boolean, msg: Localisation): Localisation.AsObject;
  static serializeBinaryToWriter(message: Localisation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Localisation;
  static deserializeBinaryFromReader(message: Localisation, reader: jspb.BinaryReader): Localisation;
}

export namespace Localisation {
  export type AsObject = {
    id: string,
    name: string,
    languagecode: string,
    sublocalisations?: References.AsObject,
  }
}

export class Inventory extends jspb.Message {
  getSafetystock(): number;
  setSafetystock(value: number): void;

  getReorderquantity(): number;
  setReorderquantity(value: number): void;

  getQuantity(): number;
  setQuantity(value: number): void;

  getFactor(): number;
  setFactor(value: number): void;

  getLocalisationid(): string;
  setLocalisationid(value: string): void;

  getPacakgeid(): string;
  setPacakgeid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Inventory.AsObject;
  static toObject(includeInstance: boolean, msg: Inventory): Inventory.AsObject;
  static serializeBinaryToWriter(message: Inventory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Inventory;
  static deserializeBinaryFromReader(message: Inventory, reader: jspb.BinaryReader): Inventory;
}

export namespace Inventory {
  export type AsObject = {
    safetystock: number,
    reorderquantity: number,
    quantity: number,
    factor: number,
    localisationid: string,
    pacakgeid: string,
  }
}

export class Price extends jspb.Message {
  getValue(): number;
  setValue(value: number): void;

  getCurrency(): Currency;
  setCurrency(value: Currency): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Price.AsObject;
  static toObject(includeInstance: boolean, msg: Price): Price.AsObject;
  static serializeBinaryToWriter(message: Price, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Price;
  static deserializeBinaryFromReader(message: Price, reader: jspb.BinaryReader): Price;
}

export namespace Price {
  export type AsObject = {
    value: number,
    currency: Currency,
  }
}

export class SubPackage extends jspb.Message {
  getUnitofmeasure(): Reference | undefined;
  setUnitofmeasure(value?: Reference): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  getPackage(): Reference | undefined;
  setPackage(value?: Reference): void;
  hasPackage(): boolean;
  clearPackage(): void;

  getQuantity(): number;
  setQuantity(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubPackage.AsObject;
  static toObject(includeInstance: boolean, msg: SubPackage): SubPackage.AsObject;
  static serializeBinaryToWriter(message: SubPackage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubPackage;
  static deserializeBinaryFromReader(message: SubPackage, reader: jspb.BinaryReader): SubPackage;
}

export namespace SubPackage {
  export type AsObject = {
    unitofmeasure?: Reference.AsObject,
    pb_package?: Reference.AsObject,
    quantity: number,
  }
}

export class ItemInstancePackage extends jspb.Message {
  getUnitofmeasure(): Reference | undefined;
  setUnitofmeasure(value?: Reference): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  getIteminstance(): Reference | undefined;
  setIteminstance(value?: Reference): void;
  hasIteminstance(): boolean;
  clearIteminstance(): void;

  getQuantity(): number;
  setQuantity(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ItemInstancePackage.AsObject;
  static toObject(includeInstance: boolean, msg: ItemInstancePackage): ItemInstancePackage.AsObject;
  static serializeBinaryToWriter(message: ItemInstancePackage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ItemInstancePackage;
  static deserializeBinaryFromReader(message: ItemInstancePackage, reader: jspb.BinaryReader): ItemInstancePackage;
}

export namespace ItemInstancePackage {
  export type AsObject = {
    unitofmeasure?: Reference.AsObject,
    iteminstance?: Reference.AsObject,
    quantity: number,
  }
}

export class Package extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getSubpackagesList(): Array<SubPackage>;
  setSubpackagesList(value: Array<SubPackage>): void;
  clearSubpackagesList(): void;
  addSubpackages(value?: SubPackage, index?: number): SubPackage;

  getIteminstancesList(): Array<ItemInstancePackage>;
  setIteminstancesList(value: Array<ItemInstancePackage>): void;
  clearIteminstancesList(): void;
  addIteminstances(value?: ItemInstancePackage, index?: number): ItemInstancePackage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Package.AsObject;
  static toObject(includeInstance: boolean, msg: Package): Package.AsObject;
  static serializeBinaryToWriter(message: Package, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Package;
  static deserializeBinaryFromReader(message: Package, reader: jspb.BinaryReader): Package;
}

export namespace Package {
  export type AsObject = {
    id: string,
    name: string,
    languagecode: string,
    description: string,
    subpackagesList: Array<SubPackage.AsObject>,
    iteminstancesList: Array<ItemInstancePackage.AsObject>,
  }
}

export class Supplier extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Supplier.AsObject;
  static toObject(includeInstance: boolean, msg: Supplier): Supplier.AsObject;
  static serializeBinaryToWriter(message: Supplier, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Supplier;
  static deserializeBinaryFromReader(message: Supplier, reader: jspb.BinaryReader): Supplier;
}

export namespace Supplier {
  export type AsObject = {
    id: string,
    name: string,
  }
}

export class PackageSupplier extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getSupplier(): Reference | undefined;
  setSupplier(value?: Reference): void;
  hasSupplier(): boolean;
  clearSupplier(): void;

  getPackage(): Reference | undefined;
  setPackage(value?: Reference): void;
  hasPackage(): boolean;
  clearPackage(): void;

  getPrice(): Price | undefined;
  setPrice(value?: Price): void;
  hasPrice(): boolean;
  clearPrice(): void;

  getDate(): number;
  setDate(value: number): void;

  getQuantity(): number;
  setQuantity(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PackageSupplier.AsObject;
  static toObject(includeInstance: boolean, msg: PackageSupplier): PackageSupplier.AsObject;
  static serializeBinaryToWriter(message: PackageSupplier, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PackageSupplier;
  static deserializeBinaryFromReader(message: PackageSupplier, reader: jspb.BinaryReader): PackageSupplier;
}

export namespace PackageSupplier {
  export type AsObject = {
    id: string,
    supplier?: Reference.AsObject,
    pb_package?: Reference.AsObject,
    price?: Price.AsObject,
    date: number,
    quantity: number,
  }
}

export class Manufacturer extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Manufacturer.AsObject;
  static toObject(includeInstance: boolean, msg: Manufacturer): Manufacturer.AsObject;
  static serializeBinaryToWriter(message: Manufacturer, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Manufacturer;
  static deserializeBinaryFromReader(message: Manufacturer, reader: jspb.BinaryReader): Manufacturer;
}

export namespace Manufacturer {
  export type AsObject = {
    id: string,
    name: string,
  }
}

export class ItemManufacturer extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getManufacturer(): Reference | undefined;
  setManufacturer(value?: Reference): void;
  hasManufacturer(): boolean;
  clearManufacturer(): void;

  getItem(): Reference | undefined;
  setItem(value?: Reference): void;
  hasItem(): boolean;
  clearItem(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ItemManufacturer.AsObject;
  static toObject(includeInstance: boolean, msg: ItemManufacturer): ItemManufacturer.AsObject;
  static serializeBinaryToWriter(message: ItemManufacturer, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ItemManufacturer;
  static deserializeBinaryFromReader(message: ItemManufacturer, reader: jspb.BinaryReader): ItemManufacturer;
}

export namespace ItemManufacturer {
  export type AsObject = {
    id: string,
    manufacturer?: Reference.AsObject,
    item?: Reference.AsObject,
  }
}

export class Dimension extends jspb.Message {
  getUnitid(): string;
  setUnitid(value: string): void;

  getValue(): number;
  setValue(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Dimension.AsObject;
  static toObject(includeInstance: boolean, msg: Dimension): Dimension.AsObject;
  static serializeBinaryToWriter(message: Dimension, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Dimension;
  static deserializeBinaryFromReader(message: Dimension, reader: jspb.BinaryReader): Dimension;
}

export namespace Dimension {
  export type AsObject = {
    unitid: string,
    value: number,
  }
}

export class PropertyValue extends jspb.Message {
  getPropertydefinitionid(): string;
  setPropertydefinitionid(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getDimensionVal(): Dimension | undefined;
  setDimensionVal(value?: Dimension): void;
  hasDimensionVal(): boolean;
  clearDimensionVal(): void;
  hasDimensionVal(): boolean;

  getTextVal(): string;
  setTextVal(value: string): void;
  hasTextVal(): boolean;

  getNumberVal(): number;
  setNumberVal(value: number): void;
  hasNumberVal(): boolean;

  getBooleanVal(): boolean;
  setBooleanVal(value: boolean): void;
  hasBooleanVal(): boolean;

  getDimensionArr(): PropertyValue.Dimensions | undefined;
  setDimensionArr(value?: PropertyValue.Dimensions): void;
  hasDimensionArr(): boolean;
  clearDimensionArr(): void;
  hasDimensionArr(): boolean;

  getTextArr(): PropertyValue.Strings | undefined;
  setTextArr(value?: PropertyValue.Strings): void;
  hasTextArr(): boolean;
  clearTextArr(): void;
  hasTextArr(): boolean;

  getNumberArr(): PropertyValue.Numerics | undefined;
  setNumberArr(value?: PropertyValue.Numerics): void;
  hasNumberArr(): boolean;
  clearNumberArr(): void;
  hasNumberArr(): boolean;

  getBooleanArr(): PropertyValue.Booleans | undefined;
  setBooleanArr(value?: PropertyValue.Booleans): void;
  hasBooleanArr(): boolean;
  clearBooleanArr(): void;
  hasBooleanArr(): boolean;

  getValueCase(): PropertyValue.ValueCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PropertyValue.AsObject;
  static toObject(includeInstance: boolean, msg: PropertyValue): PropertyValue.AsObject;
  static serializeBinaryToWriter(message: PropertyValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PropertyValue;
  static deserializeBinaryFromReader(message: PropertyValue, reader: jspb.BinaryReader): PropertyValue;
}

export namespace PropertyValue {
  export type AsObject = {
    propertydefinitionid: string,
    languagecode: string,
    dimensionVal?: Dimension.AsObject,
    textVal: string,
    numberVal: number,
    booleanVal: boolean,
    dimensionArr?: PropertyValue.Dimensions.AsObject,
    textArr?: PropertyValue.Strings.AsObject,
    numberArr?: PropertyValue.Numerics.AsObject,
    booleanArr?: PropertyValue.Booleans.AsObject,
  }

  export class Booleans extends jspb.Message {
    getValuesList(): Array<boolean>;
    setValuesList(value: Array<boolean>): void;
    clearValuesList(): void;
    addValues(value: boolean, index?: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Booleans.AsObject;
    static toObject(includeInstance: boolean, msg: Booleans): Booleans.AsObject;
    static serializeBinaryToWriter(message: Booleans, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Booleans;
    static deserializeBinaryFromReader(message: Booleans, reader: jspb.BinaryReader): Booleans;
  }

  export namespace Booleans {
    export type AsObject = {
      valuesList: Array<boolean>,
    }
  }


  export class Numerics extends jspb.Message {
    getValuesList(): Array<number>;
    setValuesList(value: Array<number>): void;
    clearValuesList(): void;
    addValues(value: number, index?: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Numerics.AsObject;
    static toObject(includeInstance: boolean, msg: Numerics): Numerics.AsObject;
    static serializeBinaryToWriter(message: Numerics, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Numerics;
    static deserializeBinaryFromReader(message: Numerics, reader: jspb.BinaryReader): Numerics;
  }

  export namespace Numerics {
    export type AsObject = {
      valuesList: Array<number>,
    }
  }


  export class Strings extends jspb.Message {
    getValuesList(): Array<string>;
    setValuesList(value: Array<string>): void;
    clearValuesList(): void;
    addValues(value: string, index?: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Strings.AsObject;
    static toObject(includeInstance: boolean, msg: Strings): Strings.AsObject;
    static serializeBinaryToWriter(message: Strings, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Strings;
    static deserializeBinaryFromReader(message: Strings, reader: jspb.BinaryReader): Strings;
  }

  export namespace Strings {
    export type AsObject = {
      valuesList: Array<string>,
    }
  }


  export class Dimensions extends jspb.Message {
    getValuesList(): Array<PropertyValue.Dimensions>;
    setValuesList(value: Array<PropertyValue.Dimensions>): void;
    clearValuesList(): void;
    addValues(value?: PropertyValue.Dimensions, index?: number): PropertyValue.Dimensions;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Dimensions.AsObject;
    static toObject(includeInstance: boolean, msg: Dimensions): Dimensions.AsObject;
    static serializeBinaryToWriter(message: Dimensions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Dimensions;
    static deserializeBinaryFromReader(message: Dimensions, reader: jspb.BinaryReader): Dimensions;
  }

  export namespace Dimensions {
    export type AsObject = {
      valuesList: Array<PropertyValue.Dimensions.AsObject>,
    }
  }


  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    DIMENSION_VAL = 3,
    TEXT_VAL = 4,
    NUMBER_VAL = 5,
    BOOLEAN_VAL = 6,
    DIMENSION_ARR = 7,
    TEXT_ARR = 8,
    NUMBER_ARR = 9,
    BOOLEAN_ARR = 10,
  }
}

export class ItemInstance extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getItemdefinitionid(): string;
  setItemdefinitionid(value: string): void;

  getValuesList(): Array<PropertyValue>;
  setValuesList(value: Array<PropertyValue>): void;
  clearValuesList(): void;
  addValues(value?: PropertyValue, index?: number): PropertyValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ItemInstance.AsObject;
  static toObject(includeInstance: boolean, msg: ItemInstance): ItemInstance.AsObject;
  static serializeBinaryToWriter(message: ItemInstance, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ItemInstance;
  static deserializeBinaryFromReader(message: ItemInstance, reader: jspb.BinaryReader): ItemInstance;
}

export namespace ItemInstance {
  export type AsObject = {
    id: string,
    itemdefinitionid: string,
    valuesList: Array<PropertyValue.AsObject>,
  }
}

export class SaveUnitOfMeasureRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getUnitofmeasure(): UnitOfMeasure | undefined;
  setUnitofmeasure(value?: UnitOfMeasure): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveUnitOfMeasureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveUnitOfMeasureRequest): SaveUnitOfMeasureRequest.AsObject;
  static serializeBinaryToWriter(message: SaveUnitOfMeasureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveUnitOfMeasureRequest;
  static deserializeBinaryFromReader(message: SaveUnitOfMeasureRequest, reader: jspb.BinaryReader): SaveUnitOfMeasureRequest;
}

export namespace SaveUnitOfMeasureRequest {
  export type AsObject = {
    connectionid: string,
    unitofmeasure?: UnitOfMeasure.AsObject,
  }
}

export class SaveUnitOfMeasureResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveUnitOfMeasureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveUnitOfMeasureResponse): SaveUnitOfMeasureResponse.AsObject;
  static serializeBinaryToWriter(message: SaveUnitOfMeasureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveUnitOfMeasureResponse;
  static deserializeBinaryFromReader(message: SaveUnitOfMeasureResponse, reader: jspb.BinaryReader): SaveUnitOfMeasureResponse;
}

export namespace SaveUnitOfMeasureResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveInventoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getInventory(): Inventory | undefined;
  setInventory(value?: Inventory): void;
  hasInventory(): boolean;
  clearInventory(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveInventoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveInventoryRequest): SaveInventoryRequest.AsObject;
  static serializeBinaryToWriter(message: SaveInventoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveInventoryRequest;
  static deserializeBinaryFromReader(message: SaveInventoryRequest, reader: jspb.BinaryReader): SaveInventoryRequest;
}

export namespace SaveInventoryRequest {
  export type AsObject = {
    connectionid: string,
    inventory?: Inventory.AsObject,
  }
}

export class SaveInventoryResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveInventoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveInventoryResponse): SaveInventoryResponse.AsObject;
  static serializeBinaryToWriter(message: SaveInventoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveInventoryResponse;
  static deserializeBinaryFromReader(message: SaveInventoryResponse, reader: jspb.BinaryReader): SaveInventoryResponse;
}

export namespace SaveInventoryResponse {
  export type AsObject = {
    id: string,
  }
}

export class SavePropertyDefinitionRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPropertydefinition(): PropertyDefinition | undefined;
  setPropertydefinition(value?: PropertyDefinition): void;
  hasPropertydefinition(): boolean;
  clearPropertydefinition(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SavePropertyDefinitionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SavePropertyDefinitionRequest): SavePropertyDefinitionRequest.AsObject;
  static serializeBinaryToWriter(message: SavePropertyDefinitionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SavePropertyDefinitionRequest;
  static deserializeBinaryFromReader(message: SavePropertyDefinitionRequest, reader: jspb.BinaryReader): SavePropertyDefinitionRequest;
}

export namespace SavePropertyDefinitionRequest {
  export type AsObject = {
    connectionid: string,
    propertydefinition?: PropertyDefinition.AsObject,
  }
}

export class SavePropertyDefinitionResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SavePropertyDefinitionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SavePropertyDefinitionResponse): SavePropertyDefinitionResponse.AsObject;
  static serializeBinaryToWriter(message: SavePropertyDefinitionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SavePropertyDefinitionResponse;
  static deserializeBinaryFromReader(message: SavePropertyDefinitionResponse, reader: jspb.BinaryReader): SavePropertyDefinitionResponse;
}

export namespace SavePropertyDefinitionResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveItemDefinitionRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getItemdefinition(): ItemDefinition | undefined;
  setItemdefinition(value?: ItemDefinition): void;
  hasItemdefinition(): boolean;
  clearItemdefinition(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveItemDefinitionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveItemDefinitionRequest): SaveItemDefinitionRequest.AsObject;
  static serializeBinaryToWriter(message: SaveItemDefinitionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveItemDefinitionRequest;
  static deserializeBinaryFromReader(message: SaveItemDefinitionRequest, reader: jspb.BinaryReader): SaveItemDefinitionRequest;
}

export namespace SaveItemDefinitionRequest {
  export type AsObject = {
    connectionid: string,
    itemdefinition?: ItemDefinition.AsObject,
  }
}

export class SaveItemDefinitionResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveItemDefinitionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveItemDefinitionResponse): SaveItemDefinitionResponse.AsObject;
  static serializeBinaryToWriter(message: SaveItemDefinitionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveItemDefinitionResponse;
  static deserializeBinaryFromReader(message: SaveItemDefinitionResponse, reader: jspb.BinaryReader): SaveItemDefinitionResponse;
}

export namespace SaveItemDefinitionResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveItemInstanceRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getIteminstance(): ItemInstance | undefined;
  setIteminstance(value?: ItemInstance): void;
  hasIteminstance(): boolean;
  clearIteminstance(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveItemInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveItemInstanceRequest): SaveItemInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: SaveItemInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveItemInstanceRequest;
  static deserializeBinaryFromReader(message: SaveItemInstanceRequest, reader: jspb.BinaryReader): SaveItemInstanceRequest;
}

export namespace SaveItemInstanceRequest {
  export type AsObject = {
    connectionid: string,
    iteminstance?: ItemInstance.AsObject,
  }
}

export class SaveItemInstanceResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveItemInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveItemInstanceResponse): SaveItemInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: SaveItemInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveItemInstanceResponse;
  static deserializeBinaryFromReader(message: SaveItemInstanceResponse, reader: jspb.BinaryReader): SaveItemInstanceResponse;
}

export namespace SaveItemInstanceResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveManufacturerRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getManufacturer(): Manufacturer | undefined;
  setManufacturer(value?: Manufacturer): void;
  hasManufacturer(): boolean;
  clearManufacturer(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveManufacturerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveManufacturerRequest): SaveManufacturerRequest.AsObject;
  static serializeBinaryToWriter(message: SaveManufacturerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveManufacturerRequest;
  static deserializeBinaryFromReader(message: SaveManufacturerRequest, reader: jspb.BinaryReader): SaveManufacturerRequest;
}

export namespace SaveManufacturerRequest {
  export type AsObject = {
    connectionid: string,
    manufacturer?: Manufacturer.AsObject,
  }
}

export class SaveManufacturerResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveManufacturerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveManufacturerResponse): SaveManufacturerResponse.AsObject;
  static serializeBinaryToWriter(message: SaveManufacturerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveManufacturerResponse;
  static deserializeBinaryFromReader(message: SaveManufacturerResponse, reader: jspb.BinaryReader): SaveManufacturerResponse;
}

export namespace SaveManufacturerResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveSupplierRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getSupplier(): Supplier | undefined;
  setSupplier(value?: Supplier): void;
  hasSupplier(): boolean;
  clearSupplier(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveSupplierRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveSupplierRequest): SaveSupplierRequest.AsObject;
  static serializeBinaryToWriter(message: SaveSupplierRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveSupplierRequest;
  static deserializeBinaryFromReader(message: SaveSupplierRequest, reader: jspb.BinaryReader): SaveSupplierRequest;
}

export namespace SaveSupplierRequest {
  export type AsObject = {
    connectionid: string,
    supplier?: Supplier.AsObject,
  }
}

export class SaveSupplierResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveSupplierResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveSupplierResponse): SaveSupplierResponse.AsObject;
  static serializeBinaryToWriter(message: SaveSupplierResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveSupplierResponse;
  static deserializeBinaryFromReader(message: SaveSupplierResponse, reader: jspb.BinaryReader): SaveSupplierResponse;
}

export namespace SaveSupplierResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveLocalisationRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getLocalisation(): Localisation | undefined;
  setLocalisation(value?: Localisation): void;
  hasLocalisation(): boolean;
  clearLocalisation(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveLocalisationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveLocalisationRequest): SaveLocalisationRequest.AsObject;
  static serializeBinaryToWriter(message: SaveLocalisationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveLocalisationRequest;
  static deserializeBinaryFromReader(message: SaveLocalisationRequest, reader: jspb.BinaryReader): SaveLocalisationRequest;
}

export namespace SaveLocalisationRequest {
  export type AsObject = {
    connectionid: string,
    localisation?: Localisation.AsObject,
  }
}

export class SaveLocalisationResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveLocalisationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveLocalisationResponse): SaveLocalisationResponse.AsObject;
  static serializeBinaryToWriter(message: SaveLocalisationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveLocalisationResponse;
  static deserializeBinaryFromReader(message: SaveLocalisationResponse, reader: jspb.BinaryReader): SaveLocalisationResponse;
}

export namespace SaveLocalisationResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveCategoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getCategory(): Category | undefined;
  setCategory(value?: Category): void;
  hasCategory(): boolean;
  clearCategory(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveCategoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveCategoryRequest): SaveCategoryRequest.AsObject;
  static serializeBinaryToWriter(message: SaveCategoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveCategoryRequest;
  static deserializeBinaryFromReader(message: SaveCategoryRequest, reader: jspb.BinaryReader): SaveCategoryRequest;
}

export namespace SaveCategoryRequest {
  export type AsObject = {
    connectionid: string,
    category?: Category.AsObject,
  }
}

export class SaveCategoryResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveCategoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveCategoryResponse): SaveCategoryResponse.AsObject;
  static serializeBinaryToWriter(message: SaveCategoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveCategoryResponse;
  static deserializeBinaryFromReader(message: SaveCategoryResponse, reader: jspb.BinaryReader): SaveCategoryResponse;
}

export namespace SaveCategoryResponse {
  export type AsObject = {
    id: string,
  }
}

export class SavePackageRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPackage(): Package | undefined;
  setPackage(value?: Package): void;
  hasPackage(): boolean;
  clearPackage(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SavePackageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SavePackageRequest): SavePackageRequest.AsObject;
  static serializeBinaryToWriter(message: SavePackageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SavePackageRequest;
  static deserializeBinaryFromReader(message: SavePackageRequest, reader: jspb.BinaryReader): SavePackageRequest;
}

export namespace SavePackageRequest {
  export type AsObject = {
    connectionid: string,
    pb_package?: Package.AsObject,
  }
}

export class SavePackageResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SavePackageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SavePackageResponse): SavePackageResponse.AsObject;
  static serializeBinaryToWriter(message: SavePackageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SavePackageResponse;
  static deserializeBinaryFromReader(message: SavePackageResponse, reader: jspb.BinaryReader): SavePackageResponse;
}

export namespace SavePackageResponse {
  export type AsObject = {
    id: string,
  }
}

export class SavePackageSupplierRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPackagesupplier(): PackageSupplier | undefined;
  setPackagesupplier(value?: PackageSupplier): void;
  hasPackagesupplier(): boolean;
  clearPackagesupplier(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SavePackageSupplierRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SavePackageSupplierRequest): SavePackageSupplierRequest.AsObject;
  static serializeBinaryToWriter(message: SavePackageSupplierRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SavePackageSupplierRequest;
  static deserializeBinaryFromReader(message: SavePackageSupplierRequest, reader: jspb.BinaryReader): SavePackageSupplierRequest;
}

export namespace SavePackageSupplierRequest {
  export type AsObject = {
    connectionid: string,
    packagesupplier?: PackageSupplier.AsObject,
  }
}

export class SavePackageSupplierResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SavePackageSupplierResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SavePackageSupplierResponse): SavePackageSupplierResponse.AsObject;
  static serializeBinaryToWriter(message: SavePackageSupplierResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SavePackageSupplierResponse;
  static deserializeBinaryFromReader(message: SavePackageSupplierResponse, reader: jspb.BinaryReader): SavePackageSupplierResponse;
}

export namespace SavePackageSupplierResponse {
  export type AsObject = {
    id: string,
  }
}

export class SaveItemManufacturerRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getItemmanafacturer(): ItemManufacturer | undefined;
  setItemmanafacturer(value?: ItemManufacturer): void;
  hasItemmanafacturer(): boolean;
  clearItemmanafacturer(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveItemManufacturerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveItemManufacturerRequest): SaveItemManufacturerRequest.AsObject;
  static serializeBinaryToWriter(message: SaveItemManufacturerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveItemManufacturerRequest;
  static deserializeBinaryFromReader(message: SaveItemManufacturerRequest, reader: jspb.BinaryReader): SaveItemManufacturerRequest;
}

export namespace SaveItemManufacturerRequest {
  export type AsObject = {
    connectionid: string,
    itemmanafacturer?: ItemManufacturer.AsObject,
  }
}

export class SaveItemManufacturerResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveItemManufacturerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveItemManufacturerResponse): SaveItemManufacturerResponse.AsObject;
  static serializeBinaryToWriter(message: SaveItemManufacturerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveItemManufacturerResponse;
  static deserializeBinaryFromReader(message: SaveItemManufacturerResponse, reader: jspb.BinaryReader): SaveItemManufacturerResponse;
}

export namespace SaveItemManufacturerResponse {
  export type AsObject = {
    id: string,
  }
}

export class GetSupplierRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getSupplierid(): string;
  setSupplierid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupplierRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupplierRequest): GetSupplierRequest.AsObject;
  static serializeBinaryToWriter(message: GetSupplierRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupplierRequest;
  static deserializeBinaryFromReader(message: GetSupplierRequest, reader: jspb.BinaryReader): GetSupplierRequest;
}

export namespace GetSupplierRequest {
  export type AsObject = {
    connectionid: string,
    supplierid: string,
  }
}

export class GetSupplierResponse extends jspb.Message {
  getSupplier(): Supplier | undefined;
  setSupplier(value?: Supplier): void;
  hasSupplier(): boolean;
  clearSupplier(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupplierResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupplierResponse): GetSupplierResponse.AsObject;
  static serializeBinaryToWriter(message: GetSupplierResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupplierResponse;
  static deserializeBinaryFromReader(message: GetSupplierResponse, reader: jspb.BinaryReader): GetSupplierResponse;
}

export namespace GetSupplierResponse {
  export type AsObject = {
    supplier?: Supplier.AsObject,
  }
}

export class Suppliers extends jspb.Message {
  getSuppliersList(): Array<Supplier>;
  setSuppliersList(value: Array<Supplier>): void;
  clearSuppliersList(): void;
  addSuppliers(value?: Supplier, index?: number): Supplier;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Suppliers.AsObject;
  static toObject(includeInstance: boolean, msg: Suppliers): Suppliers.AsObject;
  static serializeBinaryToWriter(message: Suppliers, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Suppliers;
  static deserializeBinaryFromReader(message: Suppliers, reader: jspb.BinaryReader): Suppliers;
}

export namespace Suppliers {
  export type AsObject = {
    suppliersList: Array<Supplier.AsObject>,
  }
}

export class GetSupplierPackagesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getSupplierid(): string;
  setSupplierid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupplierPackagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupplierPackagesRequest): GetSupplierPackagesRequest.AsObject;
  static serializeBinaryToWriter(message: GetSupplierPackagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupplierPackagesRequest;
  static deserializeBinaryFromReader(message: GetSupplierPackagesRequest, reader: jspb.BinaryReader): GetSupplierPackagesRequest;
}

export namespace GetSupplierPackagesRequest {
  export type AsObject = {
    connectionid: string,
    supplierid: string,
  }
}

export class GetSupplierPackagesResponse extends jspb.Message {
  getPackagessupplierList(): Array<PackageSupplier>;
  setPackagessupplierList(value: Array<PackageSupplier>): void;
  clearPackagessupplierList(): void;
  addPackagessupplier(value?: PackageSupplier, index?: number): PackageSupplier;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupplierPackagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupplierPackagesResponse): GetSupplierPackagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetSupplierPackagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupplierPackagesResponse;
  static deserializeBinaryFromReader(message: GetSupplierPackagesResponse, reader: jspb.BinaryReader): GetSupplierPackagesResponse;
}

export namespace GetSupplierPackagesResponse {
  export type AsObject = {
    packagessupplierList: Array<PackageSupplier.AsObject>,
  }
}

export class GetSuppliersRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSuppliersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSuppliersRequest): GetSuppliersRequest.AsObject;
  static serializeBinaryToWriter(message: GetSuppliersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSuppliersRequest;
  static deserializeBinaryFromReader(message: GetSuppliersRequest, reader: jspb.BinaryReader): GetSuppliersRequest;
}

export namespace GetSuppliersRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetSuppliersResponse extends jspb.Message {
  getSuppliersList(): Array<Supplier>;
  setSuppliersList(value: Array<Supplier>): void;
  clearSuppliersList(): void;
  addSuppliers(value?: Supplier, index?: number): Supplier;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSuppliersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSuppliersResponse): GetSuppliersResponse.AsObject;
  static serializeBinaryToWriter(message: GetSuppliersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSuppliersResponse;
  static deserializeBinaryFromReader(message: GetSuppliersResponse, reader: jspb.BinaryReader): GetSuppliersResponse;
}

export namespace GetSuppliersResponse {
  export type AsObject = {
    suppliersList: Array<Supplier.AsObject>,
  }
}

export class Manufacturers extends jspb.Message {
  getManufacturersList(): Array<Manufacturer>;
  setManufacturersList(value: Array<Manufacturer>): void;
  clearManufacturersList(): void;
  addManufacturers(value?: Manufacturer, index?: number): Manufacturer;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Manufacturers.AsObject;
  static toObject(includeInstance: boolean, msg: Manufacturers): Manufacturers.AsObject;
  static serializeBinaryToWriter(message: Manufacturers, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Manufacturers;
  static deserializeBinaryFromReader(message: Manufacturers, reader: jspb.BinaryReader): Manufacturers;
}

export namespace Manufacturers {
  export type AsObject = {
    manufacturersList: Array<Manufacturer.AsObject>,
  }
}

export class GetManufacturerRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getManufacturerid(): string;
  setManufacturerid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetManufacturerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetManufacturerRequest): GetManufacturerRequest.AsObject;
  static serializeBinaryToWriter(message: GetManufacturerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetManufacturerRequest;
  static deserializeBinaryFromReader(message: GetManufacturerRequest, reader: jspb.BinaryReader): GetManufacturerRequest;
}

export namespace GetManufacturerRequest {
  export type AsObject = {
    connectionid: string,
    manufacturerid: string,
  }
}

export class GetManufacturerResponse extends jspb.Message {
  getManufacturer(): Manufacturer | undefined;
  setManufacturer(value?: Manufacturer): void;
  hasManufacturer(): boolean;
  clearManufacturer(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetManufacturerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetManufacturerResponse): GetManufacturerResponse.AsObject;
  static serializeBinaryToWriter(message: GetManufacturerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetManufacturerResponse;
  static deserializeBinaryFromReader(message: GetManufacturerResponse, reader: jspb.BinaryReader): GetManufacturerResponse;
}

export namespace GetManufacturerResponse {
  export type AsObject = {
    manufacturer?: Manufacturer.AsObject,
  }
}

export class GetManufacturersRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetManufacturersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetManufacturersRequest): GetManufacturersRequest.AsObject;
  static serializeBinaryToWriter(message: GetManufacturersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetManufacturersRequest;
  static deserializeBinaryFromReader(message: GetManufacturersRequest, reader: jspb.BinaryReader): GetManufacturersRequest;
}

export namespace GetManufacturersRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetManufacturersResponse extends jspb.Message {
  getManufacturersList(): Array<Manufacturer>;
  setManufacturersList(value: Array<Manufacturer>): void;
  clearManufacturersList(): void;
  addManufacturers(value?: Manufacturer, index?: number): Manufacturer;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetManufacturersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetManufacturersResponse): GetManufacturersResponse.AsObject;
  static serializeBinaryToWriter(message: GetManufacturersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetManufacturersResponse;
  static deserializeBinaryFromReader(message: GetManufacturersResponse, reader: jspb.BinaryReader): GetManufacturersResponse;
}

export namespace GetManufacturersResponse {
  export type AsObject = {
    manufacturersList: Array<Manufacturer.AsObject>,
  }
}

export class Packages extends jspb.Message {
  getPackagesList(): Array<Package>;
  setPackagesList(value: Array<Package>): void;
  clearPackagesList(): void;
  addPackages(value?: Package, index?: number): Package;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Packages.AsObject;
  static toObject(includeInstance: boolean, msg: Packages): Packages.AsObject;
  static serializeBinaryToWriter(message: Packages, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Packages;
  static deserializeBinaryFromReader(message: Packages, reader: jspb.BinaryReader): Packages;
}

export namespace Packages {
  export type AsObject = {
    packagesList: Array<Package.AsObject>,
  }
}

export class GetPackageRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPackageid(): string;
  setPackageid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPackageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPackageRequest): GetPackageRequest.AsObject;
  static serializeBinaryToWriter(message: GetPackageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPackageRequest;
  static deserializeBinaryFromReader(message: GetPackageRequest, reader: jspb.BinaryReader): GetPackageRequest;
}

export namespace GetPackageRequest {
  export type AsObject = {
    connectionid: string,
    packageid: string,
  }
}

export class GetPackageResponse extends jspb.Message {
  getPacakge(): Package | undefined;
  setPacakge(value?: Package): void;
  hasPacakge(): boolean;
  clearPacakge(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPackageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPackageResponse): GetPackageResponse.AsObject;
  static serializeBinaryToWriter(message: GetPackageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPackageResponse;
  static deserializeBinaryFromReader(message: GetPackageResponse, reader: jspb.BinaryReader): GetPackageResponse;
}

export namespace GetPackageResponse {
  export type AsObject = {
    pacakge?: Package.AsObject,
  }
}

export class GetPackagesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPackagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPackagesRequest): GetPackagesRequest.AsObject;
  static serializeBinaryToWriter(message: GetPackagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPackagesRequest;
  static deserializeBinaryFromReader(message: GetPackagesRequest, reader: jspb.BinaryReader): GetPackagesRequest;
}

export namespace GetPackagesRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetPackagesResponse extends jspb.Message {
  getPackagesList(): Array<Package>;
  setPackagesList(value: Array<Package>): void;
  clearPackagesList(): void;
  addPackages(value?: Package, index?: number): Package;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPackagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPackagesResponse): GetPackagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetPackagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPackagesResponse;
  static deserializeBinaryFromReader(message: GetPackagesResponse, reader: jspb.BinaryReader): GetPackagesResponse;
}

export namespace GetPackagesResponse {
  export type AsObject = {
    packagesList: Array<Package.AsObject>,
  }
}

export class Localisations extends jspb.Message {
  getLocalisationsList(): Array<Localisation>;
  setLocalisationsList(value: Array<Localisation>): void;
  clearLocalisationsList(): void;
  addLocalisations(value?: Localisation, index?: number): Localisation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Localisations.AsObject;
  static toObject(includeInstance: boolean, msg: Localisations): Localisations.AsObject;
  static serializeBinaryToWriter(message: Localisations, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Localisations;
  static deserializeBinaryFromReader(message: Localisations, reader: jspb.BinaryReader): Localisations;
}

export namespace Localisations {
  export type AsObject = {
    localisationsList: Array<Localisation.AsObject>,
  }
}

export class GetLocalisationRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getLocalisationid(): string;
  setLocalisationid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLocalisationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLocalisationRequest): GetLocalisationRequest.AsObject;
  static serializeBinaryToWriter(message: GetLocalisationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLocalisationRequest;
  static deserializeBinaryFromReader(message: GetLocalisationRequest, reader: jspb.BinaryReader): GetLocalisationRequest;
}

export namespace GetLocalisationRequest {
  export type AsObject = {
    connectionid: string,
    localisationid: string,
  }
}

export class GetLocalisationResponse extends jspb.Message {
  getLocalisation(): Localisation | undefined;
  setLocalisation(value?: Localisation): void;
  hasLocalisation(): boolean;
  clearLocalisation(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLocalisationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLocalisationResponse): GetLocalisationResponse.AsObject;
  static serializeBinaryToWriter(message: GetLocalisationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLocalisationResponse;
  static deserializeBinaryFromReader(message: GetLocalisationResponse, reader: jspb.BinaryReader): GetLocalisationResponse;
}

export namespace GetLocalisationResponse {
  export type AsObject = {
    localisation?: Localisation.AsObject,
  }
}

export class GetLocalisationsRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLocalisationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLocalisationsRequest): GetLocalisationsRequest.AsObject;
  static serializeBinaryToWriter(message: GetLocalisationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLocalisationsRequest;
  static deserializeBinaryFromReader(message: GetLocalisationsRequest, reader: jspb.BinaryReader): GetLocalisationsRequest;
}

export namespace GetLocalisationsRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetLocalisationsResponse extends jspb.Message {
  getLocalisationsList(): Array<Localisation>;
  setLocalisationsList(value: Array<Localisation>): void;
  clearLocalisationsList(): void;
  addLocalisations(value?: Localisation, index?: number): Localisation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLocalisationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLocalisationsResponse): GetLocalisationsResponse.AsObject;
  static serializeBinaryToWriter(message: GetLocalisationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLocalisationsResponse;
  static deserializeBinaryFromReader(message: GetLocalisationsResponse, reader: jspb.BinaryReader): GetLocalisationsResponse;
}

export namespace GetLocalisationsResponse {
  export type AsObject = {
    localisationsList: Array<Localisation.AsObject>,
  }
}

export class UnitOfMeasures extends jspb.Message {
  getUnitofmeasuresList(): Array<UnitOfMeasure>;
  setUnitofmeasuresList(value: Array<UnitOfMeasure>): void;
  clearUnitofmeasuresList(): void;
  addUnitofmeasures(value?: UnitOfMeasure, index?: number): UnitOfMeasure;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnitOfMeasures.AsObject;
  static toObject(includeInstance: boolean, msg: UnitOfMeasures): UnitOfMeasures.AsObject;
  static serializeBinaryToWriter(message: UnitOfMeasures, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnitOfMeasures;
  static deserializeBinaryFromReader(message: UnitOfMeasures, reader: jspb.BinaryReader): UnitOfMeasures;
}

export namespace UnitOfMeasures {
  export type AsObject = {
    unitofmeasuresList: Array<UnitOfMeasure.AsObject>,
  }
}

export class GetUnitOfMeasureRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getUnitofmeasureid(): string;
  setUnitofmeasureid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUnitOfMeasureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUnitOfMeasureRequest): GetUnitOfMeasureRequest.AsObject;
  static serializeBinaryToWriter(message: GetUnitOfMeasureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUnitOfMeasureRequest;
  static deserializeBinaryFromReader(message: GetUnitOfMeasureRequest, reader: jspb.BinaryReader): GetUnitOfMeasureRequest;
}

export namespace GetUnitOfMeasureRequest {
  export type AsObject = {
    connectionid: string,
    unitofmeasureid: string,
  }
}

export class GetUnitOfMeasureResponse extends jspb.Message {
  getUnitofmeasure(): UnitOfMeasure | undefined;
  setUnitofmeasure(value?: UnitOfMeasure): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUnitOfMeasureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUnitOfMeasureResponse): GetUnitOfMeasureResponse.AsObject;
  static serializeBinaryToWriter(message: GetUnitOfMeasureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUnitOfMeasureResponse;
  static deserializeBinaryFromReader(message: GetUnitOfMeasureResponse, reader: jspb.BinaryReader): GetUnitOfMeasureResponse;
}

export namespace GetUnitOfMeasureResponse {
  export type AsObject = {
    unitofmeasure?: UnitOfMeasure.AsObject,
  }
}

export class GetUnitOfMeasuresRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUnitOfMeasuresRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUnitOfMeasuresRequest): GetUnitOfMeasuresRequest.AsObject;
  static serializeBinaryToWriter(message: GetUnitOfMeasuresRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUnitOfMeasuresRequest;
  static deserializeBinaryFromReader(message: GetUnitOfMeasuresRequest, reader: jspb.BinaryReader): GetUnitOfMeasuresRequest;
}

export namespace GetUnitOfMeasuresRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetUnitOfMeasuresResponse extends jspb.Message {
  getUnitofmeasuresList(): Array<UnitOfMeasure>;
  setUnitofmeasuresList(value: Array<UnitOfMeasure>): void;
  clearUnitofmeasuresList(): void;
  addUnitofmeasures(value?: UnitOfMeasure, index?: number): UnitOfMeasure;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUnitOfMeasuresResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUnitOfMeasuresResponse): GetUnitOfMeasuresResponse.AsObject;
  static serializeBinaryToWriter(message: GetUnitOfMeasuresResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUnitOfMeasuresResponse;
  static deserializeBinaryFromReader(message: GetUnitOfMeasuresResponse, reader: jspb.BinaryReader): GetUnitOfMeasuresResponse;
}

export namespace GetUnitOfMeasuresResponse {
  export type AsObject = {
    unitofmeasuresList: Array<UnitOfMeasure.AsObject>,
  }
}

export class Inventories extends jspb.Message {
  getInventoriesList(): Array<Inventory>;
  setInventoriesList(value: Array<Inventory>): void;
  clearInventoriesList(): void;
  addInventories(value?: Inventory, index?: number): Inventory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Inventories.AsObject;
  static toObject(includeInstance: boolean, msg: Inventories): Inventories.AsObject;
  static serializeBinaryToWriter(message: Inventories, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Inventories;
  static deserializeBinaryFromReader(message: Inventories, reader: jspb.BinaryReader): Inventories;
}

export namespace Inventories {
  export type AsObject = {
    inventoriesList: Array<Inventory.AsObject>,
  }
}

export class GetInventoriesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInventoriesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetInventoriesRequest): GetInventoriesRequest.AsObject;
  static serializeBinaryToWriter(message: GetInventoriesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInventoriesRequest;
  static deserializeBinaryFromReader(message: GetInventoriesRequest, reader: jspb.BinaryReader): GetInventoriesRequest;
}

export namespace GetInventoriesRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetInventoriesResponse extends jspb.Message {
  getInventoriesList(): Array<Inventory>;
  setInventoriesList(value: Array<Inventory>): void;
  clearInventoriesList(): void;
  addInventories(value?: Inventory, index?: number): Inventory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInventoriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetInventoriesResponse): GetInventoriesResponse.AsObject;
  static serializeBinaryToWriter(message: GetInventoriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInventoriesResponse;
  static deserializeBinaryFromReader(message: GetInventoriesResponse, reader: jspb.BinaryReader): GetInventoriesResponse;
}

export namespace GetInventoriesResponse {
  export type AsObject = {
    inventoriesList: Array<Inventory.AsObject>,
  }
}

export class Categories extends jspb.Message {
  getCategoriesList(): Array<Category>;
  setCategoriesList(value: Array<Category>): void;
  clearCategoriesList(): void;
  addCategories(value?: Category, index?: number): Category;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Categories.AsObject;
  static toObject(includeInstance: boolean, msg: Categories): Categories.AsObject;
  static serializeBinaryToWriter(message: Categories, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Categories;
  static deserializeBinaryFromReader(message: Categories, reader: jspb.BinaryReader): Categories;
}

export namespace Categories {
  export type AsObject = {
    categoriesList: Array<Category.AsObject>,
  }
}

export class GetCategoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getCategoryid(): string;
  setCategoryid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCategoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCategoryRequest): GetCategoryRequest.AsObject;
  static serializeBinaryToWriter(message: GetCategoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCategoryRequest;
  static deserializeBinaryFromReader(message: GetCategoryRequest, reader: jspb.BinaryReader): GetCategoryRequest;
}

export namespace GetCategoryRequest {
  export type AsObject = {
    connectionid: string,
    categoryid: string,
  }
}

export class GetCategoryResponse extends jspb.Message {
  getCategory(): Category | undefined;
  setCategory(value?: Category): void;
  hasCategory(): boolean;
  clearCategory(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCategoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCategoryResponse): GetCategoryResponse.AsObject;
  static serializeBinaryToWriter(message: GetCategoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCategoryResponse;
  static deserializeBinaryFromReader(message: GetCategoryResponse, reader: jspb.BinaryReader): GetCategoryResponse;
}

export namespace GetCategoryResponse {
  export type AsObject = {
    category?: Category.AsObject,
  }
}

export class GetCategoriesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCategoriesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCategoriesRequest): GetCategoriesRequest.AsObject;
  static serializeBinaryToWriter(message: GetCategoriesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCategoriesRequest;
  static deserializeBinaryFromReader(message: GetCategoriesRequest, reader: jspb.BinaryReader): GetCategoriesRequest;
}

export namespace GetCategoriesRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetCategoriesResponse extends jspb.Message {
  getCategoriesList(): Array<Category>;
  setCategoriesList(value: Array<Category>): void;
  clearCategoriesList(): void;
  addCategories(value?: Category, index?: number): Category;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCategoriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCategoriesResponse): GetCategoriesResponse.AsObject;
  static serializeBinaryToWriter(message: GetCategoriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCategoriesResponse;
  static deserializeBinaryFromReader(message: GetCategoriesResponse, reader: jspb.BinaryReader): GetCategoriesResponse;
}

export namespace GetCategoriesResponse {
  export type AsObject = {
    categoriesList: Array<Category.AsObject>,
  }
}

export class ItemInstances extends jspb.Message {
  getIteminstancesList(): Array<ItemInstance>;
  setIteminstancesList(value: Array<ItemInstance>): void;
  clearIteminstancesList(): void;
  addIteminstances(value?: ItemInstance, index?: number): ItemInstance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ItemInstances.AsObject;
  static toObject(includeInstance: boolean, msg: ItemInstances): ItemInstances.AsObject;
  static serializeBinaryToWriter(message: ItemInstances, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ItemInstances;
  static deserializeBinaryFromReader(message: ItemInstances, reader: jspb.BinaryReader): ItemInstances;
}

export namespace ItemInstances {
  export type AsObject = {
    iteminstancesList: Array<ItemInstance.AsObject>,
  }
}

export class GetItemInstanceRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getIteminstanceid(): string;
  setIteminstanceid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemInstanceRequest): GetItemInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: GetItemInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemInstanceRequest;
  static deserializeBinaryFromReader(message: GetItemInstanceRequest, reader: jspb.BinaryReader): GetItemInstanceRequest;
}

export namespace GetItemInstanceRequest {
  export type AsObject = {
    connectionid: string,
    iteminstanceid: string,
  }
}

export class GetItemInstanceResponse extends jspb.Message {
  getIteminstance(): ItemInstance | undefined;
  setIteminstance(value?: ItemInstance): void;
  hasIteminstance(): boolean;
  clearIteminstance(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemInstanceResponse): GetItemInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: GetItemInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemInstanceResponse;
  static deserializeBinaryFromReader(message: GetItemInstanceResponse, reader: jspb.BinaryReader): GetItemInstanceResponse;
}

export namespace GetItemInstanceResponse {
  export type AsObject = {
    iteminstance?: ItemInstance.AsObject,
  }
}

export class GetItemInstancesRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemInstancesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemInstancesRequest): GetItemInstancesRequest.AsObject;
  static serializeBinaryToWriter(message: GetItemInstancesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemInstancesRequest;
  static deserializeBinaryFromReader(message: GetItemInstancesRequest, reader: jspb.BinaryReader): GetItemInstancesRequest;
}

export namespace GetItemInstancesRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetItemInstancesResponse extends jspb.Message {
  getIteminstancesList(): Array<ItemInstance>;
  setIteminstancesList(value: Array<ItemInstance>): void;
  clearIteminstancesList(): void;
  addIteminstances(value?: ItemInstance, index?: number): ItemInstance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemInstancesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemInstancesResponse): GetItemInstancesResponse.AsObject;
  static serializeBinaryToWriter(message: GetItemInstancesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemInstancesResponse;
  static deserializeBinaryFromReader(message: GetItemInstancesResponse, reader: jspb.BinaryReader): GetItemInstancesResponse;
}

export namespace GetItemInstancesResponse {
  export type AsObject = {
    iteminstancesList: Array<ItemInstance.AsObject>,
  }
}

export class ItemDefinitions extends jspb.Message {
  getItemdefinitionsList(): Array<ItemDefinition>;
  setItemdefinitionsList(value: Array<ItemDefinition>): void;
  clearItemdefinitionsList(): void;
  addItemdefinitions(value?: ItemDefinition, index?: number): ItemDefinition;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ItemDefinitions.AsObject;
  static toObject(includeInstance: boolean, msg: ItemDefinitions): ItemDefinitions.AsObject;
  static serializeBinaryToWriter(message: ItemDefinitions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ItemDefinitions;
  static deserializeBinaryFromReader(message: ItemDefinitions, reader: jspb.BinaryReader): ItemDefinitions;
}

export namespace ItemDefinitions {
  export type AsObject = {
    itemdefinitionsList: Array<ItemDefinition.AsObject>,
  }
}

export class GetItemDefinitionRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getItemdefinitionid(): string;
  setItemdefinitionid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemDefinitionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemDefinitionRequest): GetItemDefinitionRequest.AsObject;
  static serializeBinaryToWriter(message: GetItemDefinitionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemDefinitionRequest;
  static deserializeBinaryFromReader(message: GetItemDefinitionRequest, reader: jspb.BinaryReader): GetItemDefinitionRequest;
}

export namespace GetItemDefinitionRequest {
  export type AsObject = {
    connectionid: string,
    itemdefinitionid: string,
  }
}

export class GetItemDefinitionResponse extends jspb.Message {
  getItemdefinition(): ItemDefinition | undefined;
  setItemdefinition(value?: ItemDefinition): void;
  hasItemdefinition(): boolean;
  clearItemdefinition(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemDefinitionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemDefinitionResponse): GetItemDefinitionResponse.AsObject;
  static serializeBinaryToWriter(message: GetItemDefinitionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemDefinitionResponse;
  static deserializeBinaryFromReader(message: GetItemDefinitionResponse, reader: jspb.BinaryReader): GetItemDefinitionResponse;
}

export namespace GetItemDefinitionResponse {
  export type AsObject = {
    itemdefinition?: ItemDefinition.AsObject,
  }
}

export class GetItemDefinitionsRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getOptions(): string;
  setOptions(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemDefinitionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemDefinitionsRequest): GetItemDefinitionsRequest.AsObject;
  static serializeBinaryToWriter(message: GetItemDefinitionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemDefinitionsRequest;
  static deserializeBinaryFromReader(message: GetItemDefinitionsRequest, reader: jspb.BinaryReader): GetItemDefinitionsRequest;
}

export namespace GetItemDefinitionsRequest {
  export type AsObject = {
    connectionid: string,
    query: string,
    options: string,
  }
}

export class GetItemDefinitionsResponse extends jspb.Message {
  getItemdefinitionsList(): Array<ItemDefinition>;
  setItemdefinitionsList(value: Array<ItemDefinition>): void;
  clearItemdefinitionsList(): void;
  addItemdefinitions(value?: ItemDefinition, index?: number): ItemDefinition;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetItemDefinitionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetItemDefinitionsResponse): GetItemDefinitionsResponse.AsObject;
  static serializeBinaryToWriter(message: GetItemDefinitionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetItemDefinitionsResponse;
  static deserializeBinaryFromReader(message: GetItemDefinitionsResponse, reader: jspb.BinaryReader): GetItemDefinitionsResponse;
}

export namespace GetItemDefinitionsResponse {
  export type AsObject = {
    itemdefinitionsList: Array<ItemDefinition.AsObject>,
  }
}

export class DeletePackageSupplierRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPackagesupplier(): PackageSupplier | undefined;
  setPackagesupplier(value?: PackageSupplier): void;
  hasPackagesupplier(): boolean;
  clearPackagesupplier(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePackageSupplierRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePackageSupplierRequest): DeletePackageSupplierRequest.AsObject;
  static serializeBinaryToWriter(message: DeletePackageSupplierRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePackageSupplierRequest;
  static deserializeBinaryFromReader(message: DeletePackageSupplierRequest, reader: jspb.BinaryReader): DeletePackageSupplierRequest;
}

export namespace DeletePackageSupplierRequest {
  export type AsObject = {
    connectionid: string,
    packagesupplier?: PackageSupplier.AsObject,
  }
}

export class DeletePackageSupplierResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePackageSupplierResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePackageSupplierResponse): DeletePackageSupplierResponse.AsObject;
  static serializeBinaryToWriter(message: DeletePackageSupplierResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePackageSupplierResponse;
  static deserializeBinaryFromReader(message: DeletePackageSupplierResponse, reader: jspb.BinaryReader): DeletePackageSupplierResponse;
}

export namespace DeletePackageSupplierResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeletePackageRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPackage(): Package | undefined;
  setPackage(value?: Package): void;
  hasPackage(): boolean;
  clearPackage(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePackageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePackageRequest): DeletePackageRequest.AsObject;
  static serializeBinaryToWriter(message: DeletePackageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePackageRequest;
  static deserializeBinaryFromReader(message: DeletePackageRequest, reader: jspb.BinaryReader): DeletePackageRequest;
}

export namespace DeletePackageRequest {
  export type AsObject = {
    connectionid: string,
    pb_package?: Package.AsObject,
  }
}

export class DeletePackageResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePackageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePackageResponse): DeletePackageResponse.AsObject;
  static serializeBinaryToWriter(message: DeletePackageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePackageResponse;
  static deserializeBinaryFromReader(message: DeletePackageResponse, reader: jspb.BinaryReader): DeletePackageResponse;
}

export namespace DeletePackageResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteSupplierRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getSupplier(): Supplier | undefined;
  setSupplier(value?: Supplier): void;
  hasSupplier(): boolean;
  clearSupplier(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSupplierRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSupplierRequest): DeleteSupplierRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteSupplierRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSupplierRequest;
  static deserializeBinaryFromReader(message: DeleteSupplierRequest, reader: jspb.BinaryReader): DeleteSupplierRequest;
}

export namespace DeleteSupplierRequest {
  export type AsObject = {
    connectionid: string,
    supplier?: Supplier.AsObject,
  }
}

export class DeleteSupplierResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSupplierResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSupplierResponse): DeleteSupplierResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteSupplierResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSupplierResponse;
  static deserializeBinaryFromReader(message: DeleteSupplierResponse, reader: jspb.BinaryReader): DeleteSupplierResponse;
}

export namespace DeleteSupplierResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeletePropertyDefinitionRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getPropertydefinition(): PropertyDefinition | undefined;
  setPropertydefinition(value?: PropertyDefinition): void;
  hasPropertydefinition(): boolean;
  clearPropertydefinition(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePropertyDefinitionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePropertyDefinitionRequest): DeletePropertyDefinitionRequest.AsObject;
  static serializeBinaryToWriter(message: DeletePropertyDefinitionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePropertyDefinitionRequest;
  static deserializeBinaryFromReader(message: DeletePropertyDefinitionRequest, reader: jspb.BinaryReader): DeletePropertyDefinitionRequest;
}

export namespace DeletePropertyDefinitionRequest {
  export type AsObject = {
    connectionid: string,
    propertydefinition?: PropertyDefinition.AsObject,
  }
}

export class DeletePropertyDefinitionResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePropertyDefinitionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePropertyDefinitionResponse): DeletePropertyDefinitionResponse.AsObject;
  static serializeBinaryToWriter(message: DeletePropertyDefinitionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePropertyDefinitionResponse;
  static deserializeBinaryFromReader(message: DeletePropertyDefinitionResponse, reader: jspb.BinaryReader): DeletePropertyDefinitionResponse;
}

export namespace DeletePropertyDefinitionResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteUnitOfMeasureRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getUnitofmeasure(): UnitOfMeasure | undefined;
  setUnitofmeasure(value?: UnitOfMeasure): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUnitOfMeasureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUnitOfMeasureRequest): DeleteUnitOfMeasureRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteUnitOfMeasureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUnitOfMeasureRequest;
  static deserializeBinaryFromReader(message: DeleteUnitOfMeasureRequest, reader: jspb.BinaryReader): DeleteUnitOfMeasureRequest;
}

export namespace DeleteUnitOfMeasureRequest {
  export type AsObject = {
    connectionid: string,
    unitofmeasure?: UnitOfMeasure.AsObject,
  }
}

export class DeleteUnitOfMeasureResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUnitOfMeasureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUnitOfMeasureResponse): DeleteUnitOfMeasureResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteUnitOfMeasureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUnitOfMeasureResponse;
  static deserializeBinaryFromReader(message: DeleteUnitOfMeasureResponse, reader: jspb.BinaryReader): DeleteUnitOfMeasureResponse;
}

export namespace DeleteUnitOfMeasureResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteItemInstanceRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getInstance(): ItemInstance | undefined;
  setInstance(value?: ItemInstance): void;
  hasInstance(): boolean;
  clearInstance(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteItemInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteItemInstanceRequest): DeleteItemInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteItemInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteItemInstanceRequest;
  static deserializeBinaryFromReader(message: DeleteItemInstanceRequest, reader: jspb.BinaryReader): DeleteItemInstanceRequest;
}

export namespace DeleteItemInstanceRequest {
  export type AsObject = {
    connectionid: string,
    instance?: ItemInstance.AsObject,
  }
}

export class DeleteItemInstanceResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteItemInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteItemInstanceResponse): DeleteItemInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteItemInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteItemInstanceResponse;
  static deserializeBinaryFromReader(message: DeleteItemInstanceResponse, reader: jspb.BinaryReader): DeleteItemInstanceResponse;
}

export namespace DeleteItemInstanceResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteManufacturerRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getManufacturer(): Manufacturer | undefined;
  setManufacturer(value?: Manufacturer): void;
  hasManufacturer(): boolean;
  clearManufacturer(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteManufacturerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteManufacturerRequest): DeleteManufacturerRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteManufacturerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteManufacturerRequest;
  static deserializeBinaryFromReader(message: DeleteManufacturerRequest, reader: jspb.BinaryReader): DeleteManufacturerRequest;
}

export namespace DeleteManufacturerRequest {
  export type AsObject = {
    connectionid: string,
    manufacturer?: Manufacturer.AsObject,
  }
}

export class DeleteManufacturerResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteManufacturerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteManufacturerResponse): DeleteManufacturerResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteManufacturerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteManufacturerResponse;
  static deserializeBinaryFromReader(message: DeleteManufacturerResponse, reader: jspb.BinaryReader): DeleteManufacturerResponse;
}

export namespace DeleteManufacturerResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteItemManufacturerRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getItemmanufacturer(): ItemManufacturer | undefined;
  setItemmanufacturer(value?: ItemManufacturer): void;
  hasItemmanufacturer(): boolean;
  clearItemmanufacturer(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteItemManufacturerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteItemManufacturerRequest): DeleteItemManufacturerRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteItemManufacturerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteItemManufacturerRequest;
  static deserializeBinaryFromReader(message: DeleteItemManufacturerRequest, reader: jspb.BinaryReader): DeleteItemManufacturerRequest;
}

export namespace DeleteItemManufacturerRequest {
  export type AsObject = {
    connectionid: string,
    itemmanufacturer?: ItemManufacturer.AsObject,
  }
}

export class DeleteItemManufacturerResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteItemManufacturerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteItemManufacturerResponse): DeleteItemManufacturerResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteItemManufacturerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteItemManufacturerResponse;
  static deserializeBinaryFromReader(message: DeleteItemManufacturerResponse, reader: jspb.BinaryReader): DeleteItemManufacturerResponse;
}

export namespace DeleteItemManufacturerResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteCategoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getCategory(): Category | undefined;
  setCategory(value?: Category): void;
  hasCategory(): boolean;
  clearCategory(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCategoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteCategoryRequest): DeleteCategoryRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteCategoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCategoryRequest;
  static deserializeBinaryFromReader(message: DeleteCategoryRequest, reader: jspb.BinaryReader): DeleteCategoryRequest;
}

export namespace DeleteCategoryRequest {
  export type AsObject = {
    connectionid: string,
    category?: Category.AsObject,
  }
}

export class DeleteCategoryResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCategoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteCategoryResponse): DeleteCategoryResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteCategoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCategoryResponse;
  static deserializeBinaryFromReader(message: DeleteCategoryResponse, reader: jspb.BinaryReader): DeleteCategoryResponse;
}

export namespace DeleteCategoryResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteLocalisationRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getLocalisation(): Localisation | undefined;
  setLocalisation(value?: Localisation): void;
  hasLocalisation(): boolean;
  clearLocalisation(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteLocalisationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteLocalisationRequest): DeleteLocalisationRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteLocalisationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteLocalisationRequest;
  static deserializeBinaryFromReader(message: DeleteLocalisationRequest, reader: jspb.BinaryReader): DeleteLocalisationRequest;
}

export namespace DeleteLocalisationRequest {
  export type AsObject = {
    connectionid: string,
    localisation?: Localisation.AsObject,
  }
}

export class DeleteLocalisationResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteLocalisationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteLocalisationResponse): DeleteLocalisationResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteLocalisationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteLocalisationResponse;
  static deserializeBinaryFromReader(message: DeleteLocalisationResponse, reader: jspb.BinaryReader): DeleteLocalisationResponse;
}

export namespace DeleteLocalisationResponse {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteInventoryRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getInventory(): Inventory | undefined;
  setInventory(value?: Inventory): void;
  hasInventory(): boolean;
  clearInventory(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteInventoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteInventoryRequest): DeleteInventoryRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteInventoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteInventoryRequest;
  static deserializeBinaryFromReader(message: DeleteInventoryRequest, reader: jspb.BinaryReader): DeleteInventoryRequest;
}

export namespace DeleteInventoryRequest {
  export type AsObject = {
    connectionid: string,
    inventory?: Inventory.AsObject,
  }
}

export class DeleteInventoryResponse extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteInventoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteInventoryResponse): DeleteInventoryResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteInventoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteInventoryResponse;
  static deserializeBinaryFromReader(message: DeleteInventoryResponse, reader: jspb.BinaryReader): DeleteInventoryResponse;
}

export namespace DeleteInventoryResponse {
  export type AsObject = {
    result: boolean,
  }
}

export enum StoreType { 
  MONGO = 0,
}
export enum Currency { 
  us = 0,
  can = 1,
  euro = 2,
}
