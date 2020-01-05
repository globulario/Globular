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

export class RefreshTokenRqst extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshTokenRqst.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshTokenRqst): RefreshTokenRqst.AsObject;
  static serializeBinaryToWriter(message: RefreshTokenRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshTokenRqst;
  static deserializeBinaryFromReader(message: RefreshTokenRqst, reader: jspb.BinaryReader): RefreshTokenRqst;
}

export namespace RefreshTokenRqst {
  export type AsObject = {
    token: string,
  }
}

export class RefreshTokenRsp extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshTokenRsp.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshTokenRsp): RefreshTokenRsp.AsObject;
  static serializeBinaryToWriter(message: RefreshTokenRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshTokenRsp;
  static deserializeBinaryFromReader(message: RefreshTokenRsp, reader: jspb.BinaryReader): RefreshTokenRsp;
}

export namespace RefreshTokenRsp {
  export type AsObject = {
    token: string,
  }
}

export class AddAccountRoleRqst extends jspb.Message {
  getAccountid(): string;
  setAccountid(value: string): void;

  getRoleid(): string;
  setRoleid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAccountRoleRqst.AsObject;
  static toObject(includeInstance: boolean, msg: AddAccountRoleRqst): AddAccountRoleRqst.AsObject;
  static serializeBinaryToWriter(message: AddAccountRoleRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAccountRoleRqst;
  static deserializeBinaryFromReader(message: AddAccountRoleRqst, reader: jspb.BinaryReader): AddAccountRoleRqst;
}

export namespace AddAccountRoleRqst {
  export type AsObject = {
    accountid: string,
    roleid: string,
  }
}

export class AddAccountRoleRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAccountRoleRsp.AsObject;
  static toObject(includeInstance: boolean, msg: AddAccountRoleRsp): AddAccountRoleRsp.AsObject;
  static serializeBinaryToWriter(message: AddAccountRoleRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAccountRoleRsp;
  static deserializeBinaryFromReader(message: AddAccountRoleRsp, reader: jspb.BinaryReader): AddAccountRoleRsp;
}

export namespace AddAccountRoleRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class RemoveAccountRoleRqst extends jspb.Message {
  getAccountid(): string;
  setAccountid(value: string): void;

  getRoleid(): string;
  setRoleid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveAccountRoleRqst.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveAccountRoleRqst): RemoveAccountRoleRqst.AsObject;
  static serializeBinaryToWriter(message: RemoveAccountRoleRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveAccountRoleRqst;
  static deserializeBinaryFromReader(message: RemoveAccountRoleRqst, reader: jspb.BinaryReader): RemoveAccountRoleRqst;
}

export namespace RemoveAccountRoleRqst {
  export type AsObject = {
    accountid: string,
    roleid: string,
  }
}

export class RemoveAccountRoleRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveAccountRoleRsp.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveAccountRoleRsp): RemoveAccountRoleRsp.AsObject;
  static serializeBinaryToWriter(message: RemoveAccountRoleRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveAccountRoleRsp;
  static deserializeBinaryFromReader(message: RemoveAccountRoleRsp, reader: jspb.BinaryReader): RemoveAccountRoleRsp;
}

export namespace RemoveAccountRoleRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class CreateRoleRqst extends jspb.Message {
  getRole(): Role | undefined;
  setRole(value?: Role): void;
  hasRole(): boolean;
  clearRole(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoleRqst.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoleRqst): CreateRoleRqst.AsObject;
  static serializeBinaryToWriter(message: CreateRoleRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoleRqst;
  static deserializeBinaryFromReader(message: CreateRoleRqst, reader: jspb.BinaryReader): CreateRoleRqst;
}

export namespace CreateRoleRqst {
  export type AsObject = {
    role?: Role.AsObject,
  }
}

export class CreateRoleRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoleRsp.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoleRsp): CreateRoleRsp.AsObject;
  static serializeBinaryToWriter(message: CreateRoleRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoleRsp;
  static deserializeBinaryFromReader(message: CreateRoleRsp, reader: jspb.BinaryReader): CreateRoleRsp;
}

export namespace CreateRoleRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class DeleteRoleRqst extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleRqst.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleRqst): DeleteRoleRqst.AsObject;
  static serializeBinaryToWriter(message: DeleteRoleRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleRqst;
  static deserializeBinaryFromReader(message: DeleteRoleRqst, reader: jspb.BinaryReader): DeleteRoleRqst;
}

export namespace DeleteRoleRqst {
  export type AsObject = {
    roleid: string,
  }
}

export class DeleteRoleRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleRsp.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleRsp): DeleteRoleRsp.AsObject;
  static serializeBinaryToWriter(message: DeleteRoleRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleRsp;
  static deserializeBinaryFromReader(message: DeleteRoleRsp, reader: jspb.BinaryReader): DeleteRoleRsp;
}

export namespace DeleteRoleRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class AddRoleActionRqst extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  getAction(): string;
  setAction(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddRoleActionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: AddRoleActionRqst): AddRoleActionRqst.AsObject;
  static serializeBinaryToWriter(message: AddRoleActionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddRoleActionRqst;
  static deserializeBinaryFromReader(message: AddRoleActionRqst, reader: jspb.BinaryReader): AddRoleActionRqst;
}

export namespace AddRoleActionRqst {
  export type AsObject = {
    roleid: string,
    action: string,
  }
}

export class AddRoleActionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddRoleActionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: AddRoleActionRsp): AddRoleActionRsp.AsObject;
  static serializeBinaryToWriter(message: AddRoleActionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddRoleActionRsp;
  static deserializeBinaryFromReader(message: AddRoleActionRsp, reader: jspb.BinaryReader): AddRoleActionRsp;
}

export namespace AddRoleActionRsp {
  export type AsObject = {
    result: boolean,
  }
}

export class RemoveRoleActionRqst extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  getAction(): string;
  setAction(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveRoleActionRqst.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveRoleActionRqst): RemoveRoleActionRqst.AsObject;
  static serializeBinaryToWriter(message: RemoveRoleActionRqst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveRoleActionRqst;
  static deserializeBinaryFromReader(message: RemoveRoleActionRqst, reader: jspb.BinaryReader): RemoveRoleActionRqst;
}

export namespace RemoveRoleActionRqst {
  export type AsObject = {
    roleid: string,
    action: string,
  }
}

export class RemoveRoleActionRsp extends jspb.Message {
  getResult(): boolean;
  setResult(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveRoleActionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveRoleActionRsp): RemoveRoleActionRsp.AsObject;
  static serializeBinaryToWriter(message: RemoveRoleActionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveRoleActionRsp;
  static deserializeBinaryFromReader(message: RemoveRoleActionRsp, reader: jspb.BinaryReader): RemoveRoleActionRsp;
}

export namespace RemoveRoleActionRsp {
  export type AsObject = {
    result: boolean,
  }
}

