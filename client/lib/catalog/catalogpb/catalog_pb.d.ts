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

  getProperties(): PropertyDefinitions | undefined;
  setProperties(value?: PropertyDefinitions): void;
  hasProperties(): boolean;
  clearProperties(): void;

  getPropertiesids(): References | undefined;
  setPropertiesids(value?: References): void;
  hasPropertiesids(): boolean;
  clearPropertiesids(): void;

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
    properties?: PropertyDefinitions.AsObject,
    propertiesids?: References.AsObject,
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

  getProperties(): PropertyDefinitions | undefined;
  setProperties(value?: PropertyDefinitions): void;
  hasProperties(): boolean;
  clearProperties(): void;

  getPropertiesids(): References | undefined;
  setPropertiesids(value?: References): void;
  hasPropertiesids(): boolean;
  clearPropertiesids(): void;

  getReleadeditemdefintionsrefs(): References | undefined;
  setReleadeditemdefintionsrefs(value?: References): void;
  hasReleadeditemdefintionsrefs(): boolean;
  clearReleadeditemdefintionsrefs(): void;

  getEquivalentsitemdefintionsrefs(): References | undefined;
  setEquivalentsitemdefintionsrefs(value?: References): void;
  hasEquivalentsitemdefintionsrefs(): boolean;
  clearEquivalentsitemdefintionsrefs(): void;

  getCategories(): Reference | undefined;
  setCategories(value?: Reference): void;
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
    properties?: PropertyDefinitions.AsObject,
    propertiesids?: References.AsObject,
    releadeditemdefintionsrefs?: References.AsObject,
    equivalentsitemdefintionsrefs?: References.AsObject,
    categories?: Reference.AsObject,
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

  getSublocalisationsList(): Array<Localisation>;
  setSublocalisationsList(value: Array<Localisation>): void;
  clearSublocalisationsList(): void;
  addSublocalisations(value?: Localisation, index?: number): Localisation;

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
    sublocalisationsList: Array<Localisation.AsObject>,
  }
}

export class Inventory extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getSafetystock(): number;
  setSafetystock(value: number): void;

  getReorderqte(): number;
  setReorderqte(value: number): void;

  getQte(): number;
  setQte(value: number): void;

  getFactor(): number;
  setFactor(value: number): void;

  getUnitofmeasureid(): string;
  setUnitofmeasureid(value: string): void;

  getLocalisationid(): string;
  setLocalisationid(value: string): void;

  getIteminstance(): Reference | undefined;
  setIteminstance(value?: Reference): void;
  hasIteminstance(): boolean;
  clearIteminstance(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Inventory.AsObject;
  static toObject(includeInstance: boolean, msg: Inventory): Inventory.AsObject;
  static serializeBinaryToWriter(message: Inventory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Inventory;
  static deserializeBinaryFromReader(message: Inventory, reader: jspb.BinaryReader): Inventory;
}

export namespace Inventory {
  export type AsObject = {
    id: string,
    safetystock: number,
    reorderqte: number,
    qte: number,
    factor: number,
    unitofmeasureid: string,
    localisationid: string,
    iteminstance?: Reference.AsObject,
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

export class Package extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getItemdefinitions(): References | undefined;
  setItemdefinitions(value?: References): void;
  hasItemdefinitions(): boolean;
  clearItemdefinitions(): void;

  getUnitofmeasure(): Reference | undefined;
  setUnitofmeasure(value?: Reference): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  getQte(): number;
  setQte(value: number): void;

  getInventory(): Reference | undefined;
  setInventory(value?: Reference): void;
  hasInventory(): boolean;
  clearInventory(): void;

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
    languagecode: string,
    description: string,
    itemdefinitions?: References.AsObject,
    unitofmeasure?: Reference.AsObject,
    qte: number,
    inventory?: Reference.AsObject,
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

  getQte(): number;
  setQte(value: number): void;

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
    qte: number,
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
  getPacakgessupplierList(): Array<PackageSupplier>;
  setPacakgessupplierList(value: Array<PackageSupplier>): void;
  clearPacakgessupplierList(): void;
  addPacakgessupplier(value?: PackageSupplier, index?: number): PackageSupplier;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupplierPackagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupplierPackagesResponse): GetSupplierPackagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetSupplierPackagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupplierPackagesResponse;
  static deserializeBinaryFromReader(message: GetSupplierPackagesResponse, reader: jspb.BinaryReader): GetSupplierPackagesResponse;
}

export namespace GetSupplierPackagesResponse {
  export type AsObject = {
    pacakgessupplierList: Array<PackageSupplier.AsObject>,
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

export enum StoreType { 
  MONGO = 0,
}
export enum Currency { 
  us = 0,
  can = 1,
  euro = 2,
}
