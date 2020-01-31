import * as jspb from "google-protobuf"

export class Reference extends jspb.Message {
  getRefdbname(): string;
  setRefdbname(value: string): void;

  getRefobjid(): string;
  setRefobjid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Reference.AsObject;
  static toObject(includeInstance: boolean, msg: Reference): Reference.AsObject;
  static serializeBinaryToWriter(message: Reference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Reference;
  static deserializeBinaryFromReader(message: Reference, reader: jspb.BinaryReader): Reference;
}

export namespace Reference {
  export type AsObject = {
    refdbname: string,
    refobjid: string,
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

export class Language extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Language.AsObject;
  static toObject(includeInstance: boolean, msg: Language): Language.AsObject;
  static serializeBinaryToWriter(message: Language, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Language;
  static deserializeBinaryFromReader(message: Language, reader: jspb.BinaryReader): Language;
}

export namespace Language {
  export type AsObject = {
    code: string,
    name: string,
  }
}

export class PropertyDefinition extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getName(): string;
  setName(value: string): void;

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
    languagecode: string,
    name: string,
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

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getName(): string;
  setName(value: string): void;

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
    languagecode: string,
    name: string,
    abreviation: string,
    description: string,
    aliasList: Array<string>,
    keywordsList: Array<string>,
    properties?: PropertyDefinitions.AsObject,
    propertiesids?: References.AsObject,
    releadeditemdefintionsrefs?: References.AsObject,
    equivalentsitemdefintionsrefs?: References.AsObject,
  }
}

export class UnitOfMeasure extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getLanguagecode(): string;
  setLanguagecode(value: string): void;

  getName(): string;
  setName(value: string): void;

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
    languagecode: string,
    name: string,
    abreviation: string,
    description: string,
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

export class SaveUnitOfMesureRequest extends jspb.Message {
  getConnectionid(): string;
  setConnectionid(value: string): void;

  getUnitofmeasure(): UnitOfMeasure | undefined;
  setUnitofmeasure(value?: UnitOfMeasure): void;
  hasUnitofmeasure(): boolean;
  clearUnitofmeasure(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveUnitOfMesureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveUnitOfMesureRequest): SaveUnitOfMesureRequest.AsObject;
  static serializeBinaryToWriter(message: SaveUnitOfMesureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveUnitOfMesureRequest;
  static deserializeBinaryFromReader(message: SaveUnitOfMesureRequest, reader: jspb.BinaryReader): SaveUnitOfMesureRequest;
}

export namespace SaveUnitOfMesureRequest {
  export type AsObject = {
    connectionid: string,
    unitofmeasure?: UnitOfMeasure.AsObject,
  }
}

export class SaveUnitOfMesureResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveUnitOfMesureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveUnitOfMesureResponse): SaveUnitOfMesureResponse.AsObject;
  static serializeBinaryToWriter(message: SaveUnitOfMesureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveUnitOfMesureResponse;
  static deserializeBinaryFromReader(message: SaveUnitOfMesureResponse, reader: jspb.BinaryReader): SaveUnitOfMesureResponse;
}

export namespace SaveUnitOfMesureResponse {
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

export enum StoreType { 
  MONGO = 0,
}
