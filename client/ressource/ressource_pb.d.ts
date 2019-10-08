import * as jspb from "google-protobuf"

export class Account extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Account.AsObject;
  static toObject(includeInstance: boolean, msg: Account): Account.AsObject;
  static serializeBinaryToWriter(message: Account, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Account;
  static deserializeBinaryFromReader(message: Account, reader: jspb.BinaryReader): Account;
}

export namespace Account {
  export type AsObject = {
    name: string,
    email: string,
    password: string,
  }
}

export class Role extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getActionsList(): Array<string>;
  setActionsList(value: Array<string>): void;
  clearActionsList(): void;
  addActions(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Role.AsObject;
  static toObject(includeInstance: boolean, msg: Role): Role.AsObject;
  static serializeBinaryToWriter(message: Role, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Role;
  static deserializeBinaryFromReader(message: Role, reader: jspb.BinaryReader): Role;
}

export namespace Role {
  export type AsObject = {
    name: string,
    actionsList: Array<string>,
  }
}

export class RegisterAccountRqst extends jspb.Message {
  getAccount(): Account | undefined;
  setAccount(value?: Account): void;
  hasAccount(): boolean;
  clearAccount(): void;

  getPassword(): string;
  setPassword(value: string): void;

  getConfirmPassword(): string;
  setConfirmPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterAccountRqst.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterAccountRqst): RegisterAccountRqst.AsObject;
  static serializeBinaryToWriter(message: RegisterAccountRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterAccountRqst;
  static deserializeBinaryFromReader(message: RegisterAccountRqst, reader: jspb.BinaryReader): RegisterAccountRqst;
}

export namespace RegisterAccountRqst {
  export type AsObject = {
    account?: Account.AsObject,
    password: string,
    confirmPassword: string,
  }
}

export class RegisterAccountRsp extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterAccountRsp.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterAccountRsp): RegisterAccountRsp.AsObject;
  static serializeBinaryToWriter(message: RegisterAccountRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterAccountRsp;
  static deserializeBinaryFromReader(message: RegisterAccountRsp, reader: jspb.BinaryReader): RegisterAccountRsp;
}

export namespace RegisterAccountRsp {
  export type AsObject = {
    result: string,
  }
}

export class DeleteAccountRqst extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountRqst): DeleteAccountRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteAccountRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountRqst;
  static deserializeBinaryFromReader(message: DeleteAccountRqst, reader: jspb.BinaryReader): DeleteAccountRqst;
}

export namespace DeleteAccountRqst {
  export type AsObject = {
    name: string,
  }
}

export class DeleteAccountRsp extends jspb.Message {
  getResult(): string;
  setResult(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAccountRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAccountRsp): DeleteAccountRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteAccountRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAccountRsp;
  static deserializeBinaryFromReader(message: DeleteAccountRsp, reader: jspb.BinaryReader): DeleteAccountRsp;
}

export namespace DeleteAccountRsp {
  export type AsObject = {
    result: string,
  }
}

export class AuthenticateRqst extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticateRqst.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticateRqst): AuthenticateRqst.AsObject;
  static serializeBinaryToWriter(message: AuthenticateRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticateRqst;
  static deserializeBinaryFromReader(message: AuthenticateRqst, reader: jspb.BinaryReader): AuthenticateRqst;
}

export namespace AuthenticateRqst {
  export type AsObject = {
    name: string,
    password: string,
  }
}

export class AuthenticateRsp extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticateRsp.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticateRsp): AuthenticateRsp.AsObject;
  static serializeBinaryToWriter(message: AuthenticateRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticateRsp;
  static deserializeBinaryFromReader(message: AuthenticateRsp, reader: jspb.BinaryReader): AuthenticateRsp;
}

export namespace AuthenticateRsp {
  export type AsObject = {
    token: string,
  }
}

