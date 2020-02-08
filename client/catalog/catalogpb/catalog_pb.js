// source: catalog/catalogpb/catalog.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

goog.exportSymbol('proto.catalog.AppendItemDefinitionCategoryRequest', null, global);
goog.exportSymbol('proto.catalog.AppendItemDefinitionCategoryResponse', null, global);
goog.exportSymbol('proto.catalog.Category', null, global);
goog.exportSymbol('proto.catalog.Connection', null, global);
goog.exportSymbol('proto.catalog.CreateConnectionRqst', null, global);
goog.exportSymbol('proto.catalog.CreateConnectionRsp', null, global);
goog.exportSymbol('proto.catalog.Currency', null, global);
goog.exportSymbol('proto.catalog.DeleteCategoryRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteCategoryResponse', null, global);
goog.exportSymbol('proto.catalog.DeleteConnectionRqst', null, global);
goog.exportSymbol('proto.catalog.DeleteConnectionRsp', null, global);
goog.exportSymbol('proto.catalog.DeleteItemInstanceRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteItemInstanceResponse', null, global);
goog.exportSymbol('proto.catalog.DeleteItemManufacturerRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteItemManufacturerResponse', null, global);
goog.exportSymbol('proto.catalog.DeleteLocalisationRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteLocalisationResponse', null, global);
goog.exportSymbol('proto.catalog.DeleteManufacturerRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteManufacturerResponse', null, global);
goog.exportSymbol('proto.catalog.DeletePackageRequest', null, global);
goog.exportSymbol('proto.catalog.DeletePackageResponse', null, global);
goog.exportSymbol('proto.catalog.DeletePackageSupplierRequest', null, global);
goog.exportSymbol('proto.catalog.DeletePackageSupplierResponse', null, global);
goog.exportSymbol('proto.catalog.DeletePropertyDefinitionRequest', null, global);
goog.exportSymbol('proto.catalog.DeletePropertyDefinitionResponse', null, global);
goog.exportSymbol('proto.catalog.DeleteSupplierRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteSupplierResponse', null, global);
goog.exportSymbol('proto.catalog.DeleteUnitOfMeasureRequest', null, global);
goog.exportSymbol('proto.catalog.DeleteUnitOfMeasureResponse', null, global);
goog.exportSymbol('proto.catalog.Dimension', null, global);
goog.exportSymbol('proto.catalog.GetPackageRequest', null, global);
goog.exportSymbol('proto.catalog.GetPackageResponse', null, global);
goog.exportSymbol('proto.catalog.GetSupplierPackagesRequest', null, global);
goog.exportSymbol('proto.catalog.GetSupplierPackagesResponse', null, global);
goog.exportSymbol('proto.catalog.GetSupplierRequest', null, global);
goog.exportSymbol('proto.catalog.GetSupplierResponse', null, global);
goog.exportSymbol('proto.catalog.Inventory', null, global);
goog.exportSymbol('proto.catalog.ItemDefinition', null, global);
goog.exportSymbol('proto.catalog.ItemInstance', null, global);
goog.exportSymbol('proto.catalog.ItemManufacturer', null, global);
goog.exportSymbol('proto.catalog.Localisation', null, global);
goog.exportSymbol('proto.catalog.Manufacturer', null, global);
goog.exportSymbol('proto.catalog.Package', null, global);
goog.exportSymbol('proto.catalog.PackageSupplier', null, global);
goog.exportSymbol('proto.catalog.Price', null, global);
goog.exportSymbol('proto.catalog.PropertyDefinition', null, global);
goog.exportSymbol('proto.catalog.PropertyDefinition.Type', null, global);
goog.exportSymbol('proto.catalog.PropertyDefinitions', null, global);
goog.exportSymbol('proto.catalog.PropertyValue', null, global);
goog.exportSymbol('proto.catalog.PropertyValue.Booleans', null, global);
goog.exportSymbol('proto.catalog.PropertyValue.Dimensions', null, global);
goog.exportSymbol('proto.catalog.PropertyValue.Numerics', null, global);
goog.exportSymbol('proto.catalog.PropertyValue.Strings', null, global);
goog.exportSymbol('proto.catalog.PropertyValue.ValueCase', null, global);
goog.exportSymbol('proto.catalog.Reference', null, global);
goog.exportSymbol('proto.catalog.References', null, global);
goog.exportSymbol('proto.catalog.RemoveItemDefinitionCategoryRequest', null, global);
goog.exportSymbol('proto.catalog.RemoveItemDefinitionCategoryResponse', null, global);
goog.exportSymbol('proto.catalog.SaveCategoryRequest', null, global);
goog.exportSymbol('proto.catalog.SaveCategoryResponse', null, global);
goog.exportSymbol('proto.catalog.SaveInventoryRequest', null, global);
goog.exportSymbol('proto.catalog.SaveInventoryResponse', null, global);
goog.exportSymbol('proto.catalog.SaveItemDefinitionRequest', null, global);
goog.exportSymbol('proto.catalog.SaveItemDefinitionResponse', null, global);
goog.exportSymbol('proto.catalog.SaveItemInstanceRequest', null, global);
goog.exportSymbol('proto.catalog.SaveItemInstanceResponse', null, global);
goog.exportSymbol('proto.catalog.SaveItemManufacturerRequest', null, global);
goog.exportSymbol('proto.catalog.SaveItemManufacturerResponse', null, global);
goog.exportSymbol('proto.catalog.SaveLocalisationRequest', null, global);
goog.exportSymbol('proto.catalog.SaveLocalisationResponse', null, global);
goog.exportSymbol('proto.catalog.SaveManufacturerRequest', null, global);
goog.exportSymbol('proto.catalog.SaveManufacturerResponse', null, global);
goog.exportSymbol('proto.catalog.SavePackageRequest', null, global);
goog.exportSymbol('proto.catalog.SavePackageResponse', null, global);
goog.exportSymbol('proto.catalog.SavePackageSupplierRequest', null, global);
goog.exportSymbol('proto.catalog.SavePackageSupplierResponse', null, global);
goog.exportSymbol('proto.catalog.SavePropertyDefinitionRequest', null, global);
goog.exportSymbol('proto.catalog.SavePropertyDefinitionResponse', null, global);
goog.exportSymbol('proto.catalog.SaveSupplierRequest', null, global);
goog.exportSymbol('proto.catalog.SaveSupplierResponse', null, global);
goog.exportSymbol('proto.catalog.SaveUnitOfMeasureRequest', null, global);
goog.exportSymbol('proto.catalog.SaveUnitOfMeasureResponse', null, global);
goog.exportSymbol('proto.catalog.StoreType', null, global);
goog.exportSymbol('proto.catalog.Supplier', null, global);
goog.exportSymbol('proto.catalog.UnitOfMeasure', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Reference = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Reference, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Reference.displayName = 'proto.catalog.Reference';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.References = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.References.repeatedFields_, null);
};
goog.inherits(proto.catalog.References, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.References.displayName = 'proto.catalog.References';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Connection = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Connection, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Connection.displayName = 'proto.catalog.Connection';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.CreateConnectionRqst = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.CreateConnectionRqst, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.CreateConnectionRqst.displayName = 'proto.catalog.CreateConnectionRqst';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.CreateConnectionRsp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.CreateConnectionRsp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.CreateConnectionRsp.displayName = 'proto.catalog.CreateConnectionRsp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteConnectionRqst = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteConnectionRqst, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteConnectionRqst.displayName = 'proto.catalog.DeleteConnectionRqst';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteConnectionRsp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteConnectionRsp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteConnectionRsp.displayName = 'proto.catalog.DeleteConnectionRsp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyDefinition = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.PropertyDefinition.repeatedFields_, null);
};
goog.inherits(proto.catalog.PropertyDefinition, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyDefinition.displayName = 'proto.catalog.PropertyDefinition';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyDefinitions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.PropertyDefinitions.repeatedFields_, null);
};
goog.inherits(proto.catalog.PropertyDefinitions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyDefinitions.displayName = 'proto.catalog.PropertyDefinitions';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.ItemDefinition = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.ItemDefinition.repeatedFields_, null);
};
goog.inherits(proto.catalog.ItemDefinition, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.ItemDefinition.displayName = 'proto.catalog.ItemDefinition';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.AppendItemDefinitionCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.AppendItemDefinitionCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.AppendItemDefinitionCategoryRequest.displayName = 'proto.catalog.AppendItemDefinitionCategoryRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.AppendItemDefinitionCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.AppendItemDefinitionCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.AppendItemDefinitionCategoryResponse.displayName = 'proto.catalog.AppendItemDefinitionCategoryResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.RemoveItemDefinitionCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.RemoveItemDefinitionCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.RemoveItemDefinitionCategoryRequest.displayName = 'proto.catalog.RemoveItemDefinitionCategoryRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.RemoveItemDefinitionCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.RemoveItemDefinitionCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.RemoveItemDefinitionCategoryResponse.displayName = 'proto.catalog.RemoveItemDefinitionCategoryResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.UnitOfMeasure = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.UnitOfMeasure, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.UnitOfMeasure.displayName = 'proto.catalog.UnitOfMeasure';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Category = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Category, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Category.displayName = 'proto.catalog.Category';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Localisation = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.Localisation.repeatedFields_, null);
};
goog.inherits(proto.catalog.Localisation, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Localisation.displayName = 'proto.catalog.Localisation';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Inventory = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Inventory, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Inventory.displayName = 'proto.catalog.Inventory';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Price = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Price, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Price.displayName = 'proto.catalog.Price';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Package = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Package, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Package.displayName = 'proto.catalog.Package';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Supplier = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Supplier, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Supplier.displayName = 'proto.catalog.Supplier';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PackageSupplier = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.PackageSupplier, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PackageSupplier.displayName = 'proto.catalog.PackageSupplier';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Manufacturer = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Manufacturer, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Manufacturer.displayName = 'proto.catalog.Manufacturer';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.ItemManufacturer = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.ItemManufacturer, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.ItemManufacturer.displayName = 'proto.catalog.ItemManufacturer';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.Dimension = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.Dimension, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.Dimension.displayName = 'proto.catalog.Dimension';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyValue = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.catalog.PropertyValue.oneofGroups_);
};
goog.inherits(proto.catalog.PropertyValue, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyValue.displayName = 'proto.catalog.PropertyValue';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyValue.Booleans = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.PropertyValue.Booleans.repeatedFields_, null);
};
goog.inherits(proto.catalog.PropertyValue.Booleans, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyValue.Booleans.displayName = 'proto.catalog.PropertyValue.Booleans';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyValue.Numerics = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.PropertyValue.Numerics.repeatedFields_, null);
};
goog.inherits(proto.catalog.PropertyValue.Numerics, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyValue.Numerics.displayName = 'proto.catalog.PropertyValue.Numerics';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyValue.Strings = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.PropertyValue.Strings.repeatedFields_, null);
};
goog.inherits(proto.catalog.PropertyValue.Strings, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyValue.Strings.displayName = 'proto.catalog.PropertyValue.Strings';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.PropertyValue.Dimensions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.PropertyValue.Dimensions.repeatedFields_, null);
};
goog.inherits(proto.catalog.PropertyValue.Dimensions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.PropertyValue.Dimensions.displayName = 'proto.catalog.PropertyValue.Dimensions';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.ItemInstance = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.ItemInstance.repeatedFields_, null);
};
goog.inherits(proto.catalog.ItemInstance, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.ItemInstance.displayName = 'proto.catalog.ItemInstance';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveUnitOfMeasureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveUnitOfMeasureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveUnitOfMeasureRequest.displayName = 'proto.catalog.SaveUnitOfMeasureRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveUnitOfMeasureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveUnitOfMeasureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveUnitOfMeasureResponse.displayName = 'proto.catalog.SaveUnitOfMeasureResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SavePropertyDefinitionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SavePropertyDefinitionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SavePropertyDefinitionRequest.displayName = 'proto.catalog.SavePropertyDefinitionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SavePropertyDefinitionResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SavePropertyDefinitionResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SavePropertyDefinitionResponse.displayName = 'proto.catalog.SavePropertyDefinitionResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveItemDefinitionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveItemDefinitionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveItemDefinitionRequest.displayName = 'proto.catalog.SaveItemDefinitionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveItemDefinitionResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveItemDefinitionResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveItemDefinitionResponse.displayName = 'proto.catalog.SaveItemDefinitionResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveItemInstanceRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveItemInstanceRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveItemInstanceRequest.displayName = 'proto.catalog.SaveItemInstanceRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveItemInstanceResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveItemInstanceResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveItemInstanceResponse.displayName = 'proto.catalog.SaveItemInstanceResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveManufacturerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveManufacturerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveManufacturerRequest.displayName = 'proto.catalog.SaveManufacturerRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveManufacturerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveManufacturerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveManufacturerResponse.displayName = 'proto.catalog.SaveManufacturerResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveSupplierRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveSupplierRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveSupplierRequest.displayName = 'proto.catalog.SaveSupplierRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveSupplierResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveSupplierResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveSupplierResponse.displayName = 'proto.catalog.SaveSupplierResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveLocalisationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveLocalisationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveLocalisationRequest.displayName = 'proto.catalog.SaveLocalisationRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveLocalisationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveLocalisationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveLocalisationResponse.displayName = 'proto.catalog.SaveLocalisationResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveCategoryRequest.displayName = 'proto.catalog.SaveCategoryRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveCategoryResponse.displayName = 'proto.catalog.SaveCategoryResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveInventoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveInventoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveInventoryRequest.displayName = 'proto.catalog.SaveInventoryRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveInventoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveInventoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveInventoryResponse.displayName = 'proto.catalog.SaveInventoryResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SavePackageRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SavePackageRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SavePackageRequest.displayName = 'proto.catalog.SavePackageRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SavePackageResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SavePackageResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SavePackageResponse.displayName = 'proto.catalog.SavePackageResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SavePackageSupplierRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SavePackageSupplierRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SavePackageSupplierRequest.displayName = 'proto.catalog.SavePackageSupplierRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SavePackageSupplierResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SavePackageSupplierResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SavePackageSupplierResponse.displayName = 'proto.catalog.SavePackageSupplierResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveItemManufacturerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveItemManufacturerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveItemManufacturerRequest.displayName = 'proto.catalog.SaveItemManufacturerRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.SaveItemManufacturerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.SaveItemManufacturerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.SaveItemManufacturerResponse.displayName = 'proto.catalog.SaveItemManufacturerResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.GetSupplierRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.GetSupplierRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.GetSupplierRequest.displayName = 'proto.catalog.GetSupplierRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.GetSupplierResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.GetSupplierResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.GetSupplierResponse.displayName = 'proto.catalog.GetSupplierResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.GetSupplierPackagesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.GetSupplierPackagesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.GetSupplierPackagesRequest.displayName = 'proto.catalog.GetSupplierPackagesRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.GetSupplierPackagesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.catalog.GetSupplierPackagesResponse.repeatedFields_, null);
};
goog.inherits(proto.catalog.GetSupplierPackagesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.GetSupplierPackagesResponse.displayName = 'proto.catalog.GetSupplierPackagesResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.GetPackageRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.GetPackageRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.GetPackageRequest.displayName = 'proto.catalog.GetPackageRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.GetPackageResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.GetPackageResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.GetPackageResponse.displayName = 'proto.catalog.GetPackageResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeletePackageSupplierRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeletePackageSupplierRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeletePackageSupplierRequest.displayName = 'proto.catalog.DeletePackageSupplierRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeletePackageSupplierResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeletePackageSupplierResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeletePackageSupplierResponse.displayName = 'proto.catalog.DeletePackageSupplierResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeletePackageRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeletePackageRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeletePackageRequest.displayName = 'proto.catalog.DeletePackageRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeletePackageResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeletePackageResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeletePackageResponse.displayName = 'proto.catalog.DeletePackageResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteSupplierRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteSupplierRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteSupplierRequest.displayName = 'proto.catalog.DeleteSupplierRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteSupplierResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteSupplierResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteSupplierResponse.displayName = 'proto.catalog.DeleteSupplierResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeletePropertyDefinitionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeletePropertyDefinitionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeletePropertyDefinitionRequest.displayName = 'proto.catalog.DeletePropertyDefinitionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeletePropertyDefinitionResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeletePropertyDefinitionResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeletePropertyDefinitionResponse.displayName = 'proto.catalog.DeletePropertyDefinitionResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteUnitOfMeasureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteUnitOfMeasureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteUnitOfMeasureRequest.displayName = 'proto.catalog.DeleteUnitOfMeasureRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteUnitOfMeasureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteUnitOfMeasureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteUnitOfMeasureResponse.displayName = 'proto.catalog.DeleteUnitOfMeasureResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteItemInstanceRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteItemInstanceRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteItemInstanceRequest.displayName = 'proto.catalog.DeleteItemInstanceRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteItemInstanceResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteItemInstanceResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteItemInstanceResponse.displayName = 'proto.catalog.DeleteItemInstanceResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteManufacturerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteManufacturerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteManufacturerRequest.displayName = 'proto.catalog.DeleteManufacturerRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteManufacturerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteManufacturerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteManufacturerResponse.displayName = 'proto.catalog.DeleteManufacturerResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteItemManufacturerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteItemManufacturerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteItemManufacturerRequest.displayName = 'proto.catalog.DeleteItemManufacturerRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteItemManufacturerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteItemManufacturerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteItemManufacturerResponse.displayName = 'proto.catalog.DeleteItemManufacturerResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteCategoryRequest.displayName = 'proto.catalog.DeleteCategoryRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteCategoryResponse.displayName = 'proto.catalog.DeleteCategoryResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteLocalisationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteLocalisationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteLocalisationRequest.displayName = 'proto.catalog.DeleteLocalisationRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.catalog.DeleteLocalisationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.catalog.DeleteLocalisationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.catalog.DeleteLocalisationResponse.displayName = 'proto.catalog.DeleteLocalisationResponse';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Reference.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Reference.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Reference} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Reference.toObject = function(includeInstance, msg) {
  var f, obj = {
    refcolid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    refobjid: jspb.Message.getFieldWithDefault(msg, 2, ""),
    refdbname: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Reference}
 */
proto.catalog.Reference.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Reference;
  return proto.catalog.Reference.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Reference} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Reference}
 */
proto.catalog.Reference.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRefcolid(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setRefobjid(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setRefdbname(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Reference.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Reference.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Reference} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Reference.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRefcolid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRefobjid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getRefdbname();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string refColId = 1;
 * @return {string}
 */
proto.catalog.Reference.prototype.getRefcolid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Reference.prototype.setRefcolid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string refObjId = 2;
 * @return {string}
 */
proto.catalog.Reference.prototype.getRefobjid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Reference.prototype.setRefobjid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string refDbName = 3;
 * @return {string}
 */
proto.catalog.Reference.prototype.getRefdbname = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.Reference.prototype.setRefdbname = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.References.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.References.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.References.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.References} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.References.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: jspb.Message.toObjectList(msg.getValuesList(),
    proto.catalog.Reference.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.References}
 */
proto.catalog.References.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.References;
  return proto.catalog.References.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.References} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.References}
 */
proto.catalog.References.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.addValues(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.References.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.References.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.References} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.References.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * repeated Reference values = 1;
 * @return {!Array<!proto.catalog.Reference>}
 */
proto.catalog.References.prototype.getValuesList = function() {
  return /** @type{!Array<!proto.catalog.Reference>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.catalog.Reference, 1));
};


/** @param {!Array<!proto.catalog.Reference>} value */
proto.catalog.References.prototype.setValuesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.catalog.Reference=} opt_value
 * @param {number=} opt_index
 * @return {!proto.catalog.Reference}
 */
proto.catalog.References.prototype.addValues = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.catalog.Reference, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.References.prototype.clearValuesList = function() {
  this.setValuesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Connection.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Connection.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Connection} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Connection.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    host: jspb.Message.getFieldWithDefault(msg, 3, ""),
    store: jspb.Message.getFieldWithDefault(msg, 5, 0),
    user: jspb.Message.getFieldWithDefault(msg, 6, ""),
    password: jspb.Message.getFieldWithDefault(msg, 7, ""),
    port: jspb.Message.getFieldWithDefault(msg, 8, 0),
    timeout: jspb.Message.getFieldWithDefault(msg, 9, 0),
    options: jspb.Message.getFieldWithDefault(msg, 10, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Connection}
 */
proto.catalog.Connection.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Connection;
  return proto.catalog.Connection.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Connection} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Connection}
 */
proto.catalog.Connection.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setHost(value);
      break;
    case 5:
      var value = /** @type {!proto.catalog.StoreType} */ (reader.readEnum());
      msg.setStore(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setUser(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
      break;
    case 8:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPort(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setTimeout(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setOptions(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Connection.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Connection.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Connection} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Connection.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getHost();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getStore();
  if (f !== 0.0) {
    writer.writeEnum(
      5,
      f
    );
  }
  f = message.getUser();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getPort();
  if (f !== 0) {
    writer.writeInt32(
      8,
      f
    );
  }
  f = message.getTimeout();
  if (f !== 0) {
    writer.writeInt32(
      9,
      f
    );
  }
  f = message.getOptions();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Connection.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Connection.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.Connection.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Connection.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string host = 3;
 * @return {string}
 */
proto.catalog.Connection.prototype.getHost = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.Connection.prototype.setHost = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional StoreType store = 5;
 * @return {!proto.catalog.StoreType}
 */
proto.catalog.Connection.prototype.getStore = function() {
  return /** @type {!proto.catalog.StoreType} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/** @param {!proto.catalog.StoreType} value */
proto.catalog.Connection.prototype.setStore = function(value) {
  jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional string user = 6;
 * @return {string}
 */
proto.catalog.Connection.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.catalog.Connection.prototype.setUser = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string password = 7;
 * @return {string}
 */
proto.catalog.Connection.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.catalog.Connection.prototype.setPassword = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional int32 port = 8;
 * @return {number}
 */
proto.catalog.Connection.prototype.getPort = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/** @param {number} value */
proto.catalog.Connection.prototype.setPort = function(value) {
  jspb.Message.setProto3IntField(this, 8, value);
};


/**
 * optional int32 timeout = 9;
 * @return {number}
 */
proto.catalog.Connection.prototype.getTimeout = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/** @param {number} value */
proto.catalog.Connection.prototype.setTimeout = function(value) {
  jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * optional string options = 10;
 * @return {string}
 */
proto.catalog.Connection.prototype.getOptions = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.catalog.Connection.prototype.setOptions = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.CreateConnectionRqst.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.CreateConnectionRqst.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.CreateConnectionRqst} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.CreateConnectionRqst.toObject = function(includeInstance, msg) {
  var f, obj = {
    connection: (f = msg.getConnection()) && proto.catalog.Connection.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.CreateConnectionRqst}
 */
proto.catalog.CreateConnectionRqst.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.CreateConnectionRqst;
  return proto.catalog.CreateConnectionRqst.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.CreateConnectionRqst} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.CreateConnectionRqst}
 */
proto.catalog.CreateConnectionRqst.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.Connection;
      reader.readMessage(value,proto.catalog.Connection.deserializeBinaryFromReader);
      msg.setConnection(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.CreateConnectionRqst.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.CreateConnectionRqst.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.CreateConnectionRqst} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.CreateConnectionRqst.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnection();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.catalog.Connection.serializeBinaryToWriter
    );
  }
};


/**
 * optional Connection connection = 1;
 * @return {?proto.catalog.Connection}
 */
proto.catalog.CreateConnectionRqst.prototype.getConnection = function() {
  return /** @type{?proto.catalog.Connection} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Connection, 1));
};


/** @param {?proto.catalog.Connection|undefined} value */
proto.catalog.CreateConnectionRqst.prototype.setConnection = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.CreateConnectionRqst.prototype.clearConnection = function() {
  this.setConnection(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.CreateConnectionRqst.prototype.hasConnection = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.CreateConnectionRsp.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.CreateConnectionRsp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.CreateConnectionRsp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.CreateConnectionRsp.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.CreateConnectionRsp}
 */
proto.catalog.CreateConnectionRsp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.CreateConnectionRsp;
  return proto.catalog.CreateConnectionRsp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.CreateConnectionRsp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.CreateConnectionRsp}
 */
proto.catalog.CreateConnectionRsp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.CreateConnectionRsp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.CreateConnectionRsp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.CreateConnectionRsp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.CreateConnectionRsp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.CreateConnectionRsp.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.CreateConnectionRsp.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteConnectionRqst.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteConnectionRqst.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteConnectionRqst} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteConnectionRqst.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteConnectionRqst}
 */
proto.catalog.DeleteConnectionRqst.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteConnectionRqst;
  return proto.catalog.DeleteConnectionRqst.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteConnectionRqst} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteConnectionRqst}
 */
proto.catalog.DeleteConnectionRqst.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteConnectionRqst.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteConnectionRqst.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteConnectionRqst} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteConnectionRqst.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.DeleteConnectionRqst.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteConnectionRqst.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteConnectionRsp.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteConnectionRsp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteConnectionRsp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteConnectionRsp.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteConnectionRsp}
 */
proto.catalog.DeleteConnectionRsp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteConnectionRsp;
  return proto.catalog.DeleteConnectionRsp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteConnectionRsp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteConnectionRsp}
 */
proto.catalog.DeleteConnectionRsp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteConnectionRsp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteConnectionRsp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteConnectionRsp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteConnectionRsp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteConnectionRsp.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteConnectionRsp.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.PropertyDefinition.repeatedFields_ = [9];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyDefinition.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyDefinition.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyDefinition} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyDefinition.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 3, ""),
    abreviation: jspb.Message.getFieldWithDefault(msg, 4, ""),
    description: jspb.Message.getFieldWithDefault(msg, 5, ""),
    type: jspb.Message.getFieldWithDefault(msg, 6, 0),
    properties: (f = msg.getProperties()) && proto.catalog.PropertyDefinitions.toObject(includeInstance, f),
    propertiesids: (f = msg.getPropertiesids()) && proto.catalog.References.toObject(includeInstance, f),
    choicesList: (f = jspb.Message.getRepeatedField(msg, 9)) == null ? undefined : f
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyDefinition}
 */
proto.catalog.PropertyDefinition.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyDefinition;
  return proto.catalog.PropertyDefinition.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyDefinition} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyDefinition}
 */
proto.catalog.PropertyDefinition.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setAbreviation(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 6:
      var value = /** @type {!proto.catalog.PropertyDefinition.Type} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 7:
      var value = new proto.catalog.PropertyDefinitions;
      reader.readMessage(value,proto.catalog.PropertyDefinitions.deserializeBinaryFromReader);
      msg.setProperties(value);
      break;
    case 8:
      var value = new proto.catalog.References;
      reader.readMessage(value,proto.catalog.References.deserializeBinaryFromReader);
      msg.setPropertiesids(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addChoices(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyDefinition.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyDefinition.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyDefinition} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyDefinition.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAbreviation();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      6,
      f
    );
  }
  f = message.getProperties();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.catalog.PropertyDefinitions.serializeBinaryToWriter
    );
  }
  f = message.getPropertiesids();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.catalog.References.serializeBinaryToWriter
    );
  }
  f = message.getChoicesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.catalog.PropertyDefinition.Type = {
  NUMERICAL: 0,
  TEXTUAL: 1,
  BOOLEAN: 2,
  DIMENTIONAL: 3,
  ENUMERATION: 4,
  AGGREGATE: 5
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.PropertyDefinition.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.PropertyDefinition.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.PropertyDefinition.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.PropertyDefinition.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string languageCode = 3;
 * @return {string}
 */
proto.catalog.PropertyDefinition.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.PropertyDefinition.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string abreviation = 4;
 * @return {string}
 */
proto.catalog.PropertyDefinition.prototype.getAbreviation = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.catalog.PropertyDefinition.prototype.setAbreviation = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string description = 5;
 * @return {string}
 */
proto.catalog.PropertyDefinition.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.catalog.PropertyDefinition.prototype.setDescription = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional Type type = 6;
 * @return {!proto.catalog.PropertyDefinition.Type}
 */
proto.catalog.PropertyDefinition.prototype.getType = function() {
  return /** @type {!proto.catalog.PropertyDefinition.Type} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/** @param {!proto.catalog.PropertyDefinition.Type} value */
proto.catalog.PropertyDefinition.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 6, value);
};


/**
 * optional PropertyDefinitions properties = 7;
 * @return {?proto.catalog.PropertyDefinitions}
 */
proto.catalog.PropertyDefinition.prototype.getProperties = function() {
  return /** @type{?proto.catalog.PropertyDefinitions} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyDefinitions, 7));
};


/** @param {?proto.catalog.PropertyDefinitions|undefined} value */
proto.catalog.PropertyDefinition.prototype.setProperties = function(value) {
  jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyDefinition.prototype.clearProperties = function() {
  this.setProperties(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyDefinition.prototype.hasProperties = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional References propertiesIds = 8;
 * @return {?proto.catalog.References}
 */
proto.catalog.PropertyDefinition.prototype.getPropertiesids = function() {
  return /** @type{?proto.catalog.References} */ (
    jspb.Message.getWrapperField(this, proto.catalog.References, 8));
};


/** @param {?proto.catalog.References|undefined} value */
proto.catalog.PropertyDefinition.prototype.setPropertiesids = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyDefinition.prototype.clearPropertiesids = function() {
  this.setPropertiesids(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyDefinition.prototype.hasPropertiesids = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * repeated string choices = 9;
 * @return {!Array<string>}
 */
proto.catalog.PropertyDefinition.prototype.getChoicesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/** @param {!Array<string>} value */
proto.catalog.PropertyDefinition.prototype.setChoicesList = function(value) {
  jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.catalog.PropertyDefinition.prototype.addChoices = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.PropertyDefinition.prototype.clearChoicesList = function() {
  this.setChoicesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.PropertyDefinitions.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyDefinitions.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyDefinitions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyDefinitions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyDefinitions.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: jspb.Message.toObjectList(msg.getValuesList(),
    proto.catalog.PropertyDefinition.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyDefinitions}
 */
proto.catalog.PropertyDefinitions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyDefinitions;
  return proto.catalog.PropertyDefinitions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyDefinitions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyDefinitions}
 */
proto.catalog.PropertyDefinitions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.PropertyDefinition;
      reader.readMessage(value,proto.catalog.PropertyDefinition.deserializeBinaryFromReader);
      msg.addValues(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyDefinitions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyDefinitions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyDefinitions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyDefinitions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.catalog.PropertyDefinition.serializeBinaryToWriter
    );
  }
};


/**
 * repeated PropertyDefinition values = 1;
 * @return {!Array<!proto.catalog.PropertyDefinition>}
 */
proto.catalog.PropertyDefinitions.prototype.getValuesList = function() {
  return /** @type{!Array<!proto.catalog.PropertyDefinition>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.catalog.PropertyDefinition, 1));
};


/** @param {!Array<!proto.catalog.PropertyDefinition>} value */
proto.catalog.PropertyDefinitions.prototype.setValuesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.catalog.PropertyDefinition=} opt_value
 * @param {number=} opt_index
 * @return {!proto.catalog.PropertyDefinition}
 */
proto.catalog.PropertyDefinitions.prototype.addValues = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.catalog.PropertyDefinition, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.PropertyDefinitions.prototype.clearValuesList = function() {
  this.setValuesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.ItemDefinition.repeatedFields_ = [6,7];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.ItemDefinition.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.ItemDefinition.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.ItemDefinition} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.ItemDefinition.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 3, ""),
    abreviation: jspb.Message.getFieldWithDefault(msg, 4, ""),
    description: jspb.Message.getFieldWithDefault(msg, 5, ""),
    aliasList: (f = jspb.Message.getRepeatedField(msg, 6)) == null ? undefined : f,
    keywordsList: (f = jspb.Message.getRepeatedField(msg, 7)) == null ? undefined : f,
    properties: (f = msg.getProperties()) && proto.catalog.PropertyDefinitions.toObject(includeInstance, f),
    propertiesids: (f = msg.getPropertiesids()) && proto.catalog.References.toObject(includeInstance, f),
    releadeditemdefintionsrefs: (f = msg.getReleadeditemdefintionsrefs()) && proto.catalog.References.toObject(includeInstance, f),
    equivalentsitemdefintionsrefs: (f = msg.getEquivalentsitemdefintionsrefs()) && proto.catalog.References.toObject(includeInstance, f),
    categories: (f = msg.getCategories()) && proto.catalog.Reference.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.ItemDefinition}
 */
proto.catalog.ItemDefinition.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.ItemDefinition;
  return proto.catalog.ItemDefinition.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.ItemDefinition} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.ItemDefinition}
 */
proto.catalog.ItemDefinition.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setAbreviation(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.addAlias(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.addKeywords(value);
      break;
    case 8:
      var value = new proto.catalog.PropertyDefinitions;
      reader.readMessage(value,proto.catalog.PropertyDefinitions.deserializeBinaryFromReader);
      msg.setProperties(value);
      break;
    case 9:
      var value = new proto.catalog.References;
      reader.readMessage(value,proto.catalog.References.deserializeBinaryFromReader);
      msg.setPropertiesids(value);
      break;
    case 10:
      var value = new proto.catalog.References;
      reader.readMessage(value,proto.catalog.References.deserializeBinaryFromReader);
      msg.setReleadeditemdefintionsrefs(value);
      break;
    case 11:
      var value = new proto.catalog.References;
      reader.readMessage(value,proto.catalog.References.deserializeBinaryFromReader);
      msg.setEquivalentsitemdefintionsrefs(value);
      break;
    case 12:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setCategories(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.ItemDefinition.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.ItemDefinition.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.ItemDefinition} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.ItemDefinition.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAbreviation();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getAliasList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      6,
      f
    );
  }
  f = message.getKeywordsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      7,
      f
    );
  }
  f = message.getProperties();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.catalog.PropertyDefinitions.serializeBinaryToWriter
    );
  }
  f = message.getPropertiesids();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.catalog.References.serializeBinaryToWriter
    );
  }
  f = message.getReleadeditemdefintionsrefs();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      proto.catalog.References.serializeBinaryToWriter
    );
  }
  f = message.getEquivalentsitemdefintionsrefs();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      proto.catalog.References.serializeBinaryToWriter
    );
  }
  f = message.getCategories();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.ItemDefinition.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.ItemDefinition.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.ItemDefinition.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.ItemDefinition.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string languageCode = 3;
 * @return {string}
 */
proto.catalog.ItemDefinition.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.ItemDefinition.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string abreviation = 4;
 * @return {string}
 */
proto.catalog.ItemDefinition.prototype.getAbreviation = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.catalog.ItemDefinition.prototype.setAbreviation = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string description = 5;
 * @return {string}
 */
proto.catalog.ItemDefinition.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.catalog.ItemDefinition.prototype.setDescription = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * repeated string alias = 6;
 * @return {!Array<string>}
 */
proto.catalog.ItemDefinition.prototype.getAliasList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 6));
};


/** @param {!Array<string>} value */
proto.catalog.ItemDefinition.prototype.setAliasList = function(value) {
  jspb.Message.setField(this, 6, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.catalog.ItemDefinition.prototype.addAlias = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 6, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.ItemDefinition.prototype.clearAliasList = function() {
  this.setAliasList([]);
};


/**
 * repeated string keyWords = 7;
 * @return {!Array<string>}
 */
proto.catalog.ItemDefinition.prototype.getKeywordsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 7));
};


/** @param {!Array<string>} value */
proto.catalog.ItemDefinition.prototype.setKeywordsList = function(value) {
  jspb.Message.setField(this, 7, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.catalog.ItemDefinition.prototype.addKeywords = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 7, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.ItemDefinition.prototype.clearKeywordsList = function() {
  this.setKeywordsList([]);
};


/**
 * optional PropertyDefinitions properties = 8;
 * @return {?proto.catalog.PropertyDefinitions}
 */
proto.catalog.ItemDefinition.prototype.getProperties = function() {
  return /** @type{?proto.catalog.PropertyDefinitions} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyDefinitions, 8));
};


/** @param {?proto.catalog.PropertyDefinitions|undefined} value */
proto.catalog.ItemDefinition.prototype.setProperties = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemDefinition.prototype.clearProperties = function() {
  this.setProperties(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemDefinition.prototype.hasProperties = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional References propertiesIds = 9;
 * @return {?proto.catalog.References}
 */
proto.catalog.ItemDefinition.prototype.getPropertiesids = function() {
  return /** @type{?proto.catalog.References} */ (
    jspb.Message.getWrapperField(this, proto.catalog.References, 9));
};


/** @param {?proto.catalog.References|undefined} value */
proto.catalog.ItemDefinition.prototype.setPropertiesids = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemDefinition.prototype.clearPropertiesids = function() {
  this.setPropertiesids(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemDefinition.prototype.hasPropertiesids = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional References releadedItemDefintionsRefs = 10;
 * @return {?proto.catalog.References}
 */
proto.catalog.ItemDefinition.prototype.getReleadeditemdefintionsrefs = function() {
  return /** @type{?proto.catalog.References} */ (
    jspb.Message.getWrapperField(this, proto.catalog.References, 10));
};


/** @param {?proto.catalog.References|undefined} value */
proto.catalog.ItemDefinition.prototype.setReleadeditemdefintionsrefs = function(value) {
  jspb.Message.setWrapperField(this, 10, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemDefinition.prototype.clearReleadeditemdefintionsrefs = function() {
  this.setReleadeditemdefintionsrefs(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemDefinition.prototype.hasReleadeditemdefintionsrefs = function() {
  return jspb.Message.getField(this, 10) != null;
};


/**
 * optional References equivalentsItemDefintionsRefs = 11;
 * @return {?proto.catalog.References}
 */
proto.catalog.ItemDefinition.prototype.getEquivalentsitemdefintionsrefs = function() {
  return /** @type{?proto.catalog.References} */ (
    jspb.Message.getWrapperField(this, proto.catalog.References, 11));
};


/** @param {?proto.catalog.References|undefined} value */
proto.catalog.ItemDefinition.prototype.setEquivalentsitemdefintionsrefs = function(value) {
  jspb.Message.setWrapperField(this, 11, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemDefinition.prototype.clearEquivalentsitemdefintionsrefs = function() {
  this.setEquivalentsitemdefintionsrefs(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemDefinition.prototype.hasEquivalentsitemdefintionsrefs = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional Reference categories = 12;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.ItemDefinition.prototype.getCategories = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 12));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.ItemDefinition.prototype.setCategories = function(value) {
  jspb.Message.setWrapperField(this, 12, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemDefinition.prototype.clearCategories = function() {
  this.setCategories(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemDefinition.prototype.hasCategories = function() {
  return jspb.Message.getField(this, 12) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.AppendItemDefinitionCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.AppendItemDefinitionCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.AppendItemDefinitionCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    category: (f = msg.getCategory()) && proto.catalog.Reference.toObject(includeInstance, f),
    itemdefinition: (f = msg.getItemdefinition()) && proto.catalog.Reference.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.AppendItemDefinitionCategoryRequest}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.AppendItemDefinitionCategoryRequest;
  return proto.catalog.AppendItemDefinitionCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.AppendItemDefinitionCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.AppendItemDefinitionCategoryRequest}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setCategory(value);
      break;
    case 3:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setItemdefinition(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.AppendItemDefinitionCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.AppendItemDefinitionCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.AppendItemDefinitionCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCategory();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
  f = message.getItemdefinition();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Reference category = 2;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.getCategory = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 2));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.setCategory = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.clearCategory = function() {
  this.setCategory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.hasCategory = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Reference itemDefinition = 3;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.getItemdefinition = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 3));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.setItemdefinition = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.clearItemdefinition = function() {
  this.setItemdefinition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.AppendItemDefinitionCategoryRequest.prototype.hasItemdefinition = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.AppendItemDefinitionCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.AppendItemDefinitionCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.AppendItemDefinitionCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.AppendItemDefinitionCategoryResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.AppendItemDefinitionCategoryResponse}
 */
proto.catalog.AppendItemDefinitionCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.AppendItemDefinitionCategoryResponse;
  return proto.catalog.AppendItemDefinitionCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.AppendItemDefinitionCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.AppendItemDefinitionCategoryResponse}
 */
proto.catalog.AppendItemDefinitionCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.AppendItemDefinitionCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.AppendItemDefinitionCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.AppendItemDefinitionCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.AppendItemDefinitionCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.AppendItemDefinitionCategoryResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.AppendItemDefinitionCategoryResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.RemoveItemDefinitionCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.RemoveItemDefinitionCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    category: (f = msg.getCategory()) && proto.catalog.Reference.toObject(includeInstance, f),
    itemdefinition: (f = msg.getItemdefinition()) && proto.catalog.Reference.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.RemoveItemDefinitionCategoryRequest}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.RemoveItemDefinitionCategoryRequest;
  return proto.catalog.RemoveItemDefinitionCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.RemoveItemDefinitionCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.RemoveItemDefinitionCategoryRequest}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setCategory(value);
      break;
    case 3:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setItemdefinition(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.RemoveItemDefinitionCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.RemoveItemDefinitionCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCategory();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
  f = message.getItemdefinition();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Reference category = 2;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.getCategory = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 2));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.setCategory = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.clearCategory = function() {
  this.setCategory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.hasCategory = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Reference itemDefinition = 3;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.getItemdefinition = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 3));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.setItemdefinition = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.clearItemdefinition = function() {
  this.setItemdefinition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.RemoveItemDefinitionCategoryRequest.prototype.hasItemdefinition = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.RemoveItemDefinitionCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.RemoveItemDefinitionCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.RemoveItemDefinitionCategoryResponse}
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.RemoveItemDefinitionCategoryResponse;
  return proto.catalog.RemoveItemDefinitionCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.RemoveItemDefinitionCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.RemoveItemDefinitionCategoryResponse}
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.RemoveItemDefinitionCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.RemoveItemDefinitionCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.RemoveItemDefinitionCategoryResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.RemoveItemDefinitionCategoryResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.UnitOfMeasure.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.UnitOfMeasure.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.UnitOfMeasure} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.UnitOfMeasure.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 3, ""),
    abreviation: jspb.Message.getFieldWithDefault(msg, 4, ""),
    description: jspb.Message.getFieldWithDefault(msg, 5, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.UnitOfMeasure}
 */
proto.catalog.UnitOfMeasure.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.UnitOfMeasure;
  return proto.catalog.UnitOfMeasure.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.UnitOfMeasure} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.UnitOfMeasure}
 */
proto.catalog.UnitOfMeasure.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setAbreviation(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.UnitOfMeasure.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.UnitOfMeasure.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.UnitOfMeasure} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.UnitOfMeasure.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAbreviation();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.UnitOfMeasure.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.UnitOfMeasure.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.UnitOfMeasure.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.UnitOfMeasure.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string languageCode = 3;
 * @return {string}
 */
proto.catalog.UnitOfMeasure.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.UnitOfMeasure.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string abreviation = 4;
 * @return {string}
 */
proto.catalog.UnitOfMeasure.prototype.getAbreviation = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.catalog.UnitOfMeasure.prototype.setAbreviation = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string description = 5;
 * @return {string}
 */
proto.catalog.UnitOfMeasure.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.catalog.UnitOfMeasure.prototype.setDescription = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Category.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Category.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Category} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Category.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 3, ""),
    categories: (f = msg.getCategories()) && proto.catalog.References.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Category}
 */
proto.catalog.Category.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Category;
  return proto.catalog.Category.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Category} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Category}
 */
proto.catalog.Category.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 4:
      var value = new proto.catalog.References;
      reader.readMessage(value,proto.catalog.References.deserializeBinaryFromReader);
      msg.setCategories(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Category.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Category.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Category} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Category.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getCategories();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.catalog.References.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Category.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Category.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.Category.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Category.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string languageCode = 3;
 * @return {string}
 */
proto.catalog.Category.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.Category.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional References categories = 4;
 * @return {?proto.catalog.References}
 */
proto.catalog.Category.prototype.getCategories = function() {
  return /** @type{?proto.catalog.References} */ (
    jspb.Message.getWrapperField(this, proto.catalog.References, 4));
};


/** @param {?proto.catalog.References|undefined} value */
proto.catalog.Category.prototype.setCategories = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.Category.prototype.clearCategories = function() {
  this.setCategories(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.Category.prototype.hasCategories = function() {
  return jspb.Message.getField(this, 4) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.Localisation.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Localisation.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Localisation.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Localisation} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Localisation.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 3, ""),
    sublocalisationsList: jspb.Message.toObjectList(msg.getSublocalisationsList(),
    proto.catalog.Localisation.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Localisation}
 */
proto.catalog.Localisation.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Localisation;
  return proto.catalog.Localisation.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Localisation} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Localisation}
 */
proto.catalog.Localisation.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 4:
      var value = new proto.catalog.Localisation;
      reader.readMessage(value,proto.catalog.Localisation.deserializeBinaryFromReader);
      msg.addSublocalisations(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Localisation.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Localisation.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Localisation} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Localisation.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getSublocalisationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.catalog.Localisation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Localisation.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Localisation.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.Localisation.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Localisation.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string languageCode = 3;
 * @return {string}
 */
proto.catalog.Localisation.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.Localisation.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated Localisation subLocalisations = 4;
 * @return {!Array<!proto.catalog.Localisation>}
 */
proto.catalog.Localisation.prototype.getSublocalisationsList = function() {
  return /** @type{!Array<!proto.catalog.Localisation>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.catalog.Localisation, 4));
};


/** @param {!Array<!proto.catalog.Localisation>} value */
proto.catalog.Localisation.prototype.setSublocalisationsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.catalog.Localisation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.catalog.Localisation}
 */
proto.catalog.Localisation.prototype.addSublocalisations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.catalog.Localisation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.Localisation.prototype.clearSublocalisationsList = function() {
  this.setSublocalisationsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Inventory.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Inventory.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Inventory} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Inventory.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    safetystock: jspb.Message.getFieldWithDefault(msg, 2, 0),
    reorderqte: jspb.Message.getFieldWithDefault(msg, 3, 0),
    qte: jspb.Message.getFieldWithDefault(msg, 4, 0),
    factor: jspb.Message.getFloatingPointFieldWithDefault(msg, 5, 0.0),
    unitofmeasureid: jspb.Message.getFieldWithDefault(msg, 6, ""),
    localisationid: jspb.Message.getFieldWithDefault(msg, 7, ""),
    iteminstance: (f = msg.getIteminstance()) && proto.catalog.Reference.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Inventory}
 */
proto.catalog.Inventory.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Inventory;
  return proto.catalog.Inventory.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Inventory} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Inventory}
 */
proto.catalog.Inventory.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setSafetystock(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setReorderqte(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setQte(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readDouble());
      msg.setFactor(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setUnitofmeasureid(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setLocalisationid(value);
      break;
    case 8:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setIteminstance(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Inventory.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Inventory.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Inventory} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Inventory.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSafetystock();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getReorderqte();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
  f = message.getQte();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
  f = message.getFactor();
  if (f !== 0.0) {
    writer.writeDouble(
      5,
      f
    );
  }
  f = message.getUnitofmeasureid();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getLocalisationid();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getIteminstance();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Inventory.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Inventory.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 safetyStock = 2;
 * @return {number}
 */
proto.catalog.Inventory.prototype.getSafetystock = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.catalog.Inventory.prototype.setSafetystock = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int64 reorderQte = 3;
 * @return {number}
 */
proto.catalog.Inventory.prototype.getReorderqte = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {number} value */
proto.catalog.Inventory.prototype.setReorderqte = function(value) {
  jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional int64 qte = 4;
 * @return {number}
 */
proto.catalog.Inventory.prototype.getQte = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/** @param {number} value */
proto.catalog.Inventory.prototype.setQte = function(value) {
  jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional double factor = 5;
 * @return {number}
 */
proto.catalog.Inventory.prototype.getFactor = function() {
  return /** @type {number} */ (jspb.Message.getFloatingPointFieldWithDefault(this, 5, 0.0));
};


/** @param {number} value */
proto.catalog.Inventory.prototype.setFactor = function(value) {
  jspb.Message.setProto3FloatField(this, 5, value);
};


/**
 * optional string unitOfMeasureId = 6;
 * @return {string}
 */
proto.catalog.Inventory.prototype.getUnitofmeasureid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.catalog.Inventory.prototype.setUnitofmeasureid = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string localisationId = 7;
 * @return {string}
 */
proto.catalog.Inventory.prototype.getLocalisationid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.catalog.Inventory.prototype.setLocalisationid = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional Reference itemInstance = 8;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.Inventory.prototype.getIteminstance = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 8));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.Inventory.prototype.setIteminstance = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.Inventory.prototype.clearIteminstance = function() {
  this.setIteminstance(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.Inventory.prototype.hasIteminstance = function() {
  return jspb.Message.getField(this, 8) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Price.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Price.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Price} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Price.toObject = function(includeInstance, msg) {
  var f, obj = {
    value: jspb.Message.getFloatingPointFieldWithDefault(msg, 1, 0.0),
    currency: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Price}
 */
proto.catalog.Price.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Price;
  return proto.catalog.Price.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Price} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Price}
 */
proto.catalog.Price.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readDouble());
      msg.setValue(value);
      break;
    case 2:
      var value = /** @type {!proto.catalog.Currency} */ (reader.readEnum());
      msg.setCurrency(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Price.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Price.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Price} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Price.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValue();
  if (f !== 0.0) {
    writer.writeDouble(
      1,
      f
    );
  }
  f = message.getCurrency();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional double value = 1;
 * @return {number}
 */
proto.catalog.Price.prototype.getValue = function() {
  return /** @type {number} */ (jspb.Message.getFloatingPointFieldWithDefault(this, 1, 0.0));
};


/** @param {number} value */
proto.catalog.Price.prototype.setValue = function(value) {
  jspb.Message.setProto3FloatField(this, 1, value);
};


/**
 * optional Currency currency = 2;
 * @return {!proto.catalog.Currency}
 */
proto.catalog.Price.prototype.getCurrency = function() {
  return /** @type {!proto.catalog.Currency} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.catalog.Currency} value */
proto.catalog.Price.prototype.setCurrency = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Package.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Package.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Package} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Package.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, ""),
    itemdefinitions: (f = msg.getItemdefinitions()) && proto.catalog.References.toObject(includeInstance, f),
    unitofmeasure: (f = msg.getUnitofmeasure()) && proto.catalog.Reference.toObject(includeInstance, f),
    qte: jspb.Message.getFieldWithDefault(msg, 6, 0),
    inventory: (f = msg.getInventory()) && proto.catalog.Reference.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Package}
 */
proto.catalog.Package.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Package;
  return proto.catalog.Package.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Package} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Package}
 */
proto.catalog.Package.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 4:
      var value = new proto.catalog.References;
      reader.readMessage(value,proto.catalog.References.deserializeBinaryFromReader);
      msg.setItemdefinitions(value);
      break;
    case 5:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setUnitofmeasure(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setQte(value);
      break;
    case 7:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setInventory(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Package.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Package.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Package} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Package.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getItemdefinitions();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.catalog.References.serializeBinaryToWriter
    );
  }
  f = message.getUnitofmeasure();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
  f = message.getQte();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getInventory();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Package.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Package.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string languageCode = 2;
 * @return {string}
 */
proto.catalog.Package.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Package.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.catalog.Package.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.catalog.Package.prototype.setDescription = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional References itemDefinitions = 4;
 * @return {?proto.catalog.References}
 */
proto.catalog.Package.prototype.getItemdefinitions = function() {
  return /** @type{?proto.catalog.References} */ (
    jspb.Message.getWrapperField(this, proto.catalog.References, 4));
};


/** @param {?proto.catalog.References|undefined} value */
proto.catalog.Package.prototype.setItemdefinitions = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.Package.prototype.clearItemdefinitions = function() {
  this.setItemdefinitions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.Package.prototype.hasItemdefinitions = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Reference unitOfMeasure = 5;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.Package.prototype.getUnitofmeasure = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 5));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.Package.prototype.setUnitofmeasure = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.Package.prototype.clearUnitofmeasure = function() {
  this.setUnitofmeasure(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.Package.prototype.hasUnitofmeasure = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional int64 qte = 6;
 * @return {number}
 */
proto.catalog.Package.prototype.getQte = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/** @param {number} value */
proto.catalog.Package.prototype.setQte = function(value) {
  jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional Reference inventory = 7;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.Package.prototype.getInventory = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 7));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.Package.prototype.setInventory = function(value) {
  jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.Package.prototype.clearInventory = function() {
  this.setInventory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.Package.prototype.hasInventory = function() {
  return jspb.Message.getField(this, 7) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Supplier.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Supplier.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Supplier} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Supplier.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Supplier}
 */
proto.catalog.Supplier.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Supplier;
  return proto.catalog.Supplier.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Supplier} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Supplier}
 */
proto.catalog.Supplier.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Supplier.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Supplier.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Supplier} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Supplier.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Supplier.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Supplier.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.Supplier.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Supplier.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PackageSupplier.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PackageSupplier.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PackageSupplier} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PackageSupplier.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    supplier: (f = msg.getSupplier()) && proto.catalog.Reference.toObject(includeInstance, f),
    pb_package: (f = msg.getPackage()) && proto.catalog.Reference.toObject(includeInstance, f),
    price: (f = msg.getPrice()) && proto.catalog.Price.toObject(includeInstance, f),
    date: jspb.Message.getFieldWithDefault(msg, 5, 0),
    qte: jspb.Message.getFieldWithDefault(msg, 6, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PackageSupplier}
 */
proto.catalog.PackageSupplier.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PackageSupplier;
  return proto.catalog.PackageSupplier.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PackageSupplier} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PackageSupplier}
 */
proto.catalog.PackageSupplier.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setSupplier(value);
      break;
    case 3:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setPackage(value);
      break;
    case 4:
      var value = new proto.catalog.Price;
      reader.readMessage(value,proto.catalog.Price.deserializeBinaryFromReader);
      msg.setPrice(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setDate(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setQte(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PackageSupplier.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PackageSupplier.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PackageSupplier} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PackageSupplier.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSupplier();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
  f = message.getPackage();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
  f = message.getPrice();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.catalog.Price.serializeBinaryToWriter
    );
  }
  f = message.getDate();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getQte();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.PackageSupplier.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.PackageSupplier.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Reference supplier = 2;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.PackageSupplier.prototype.getSupplier = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 2));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.PackageSupplier.prototype.setSupplier = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PackageSupplier.prototype.clearSupplier = function() {
  this.setSupplier(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PackageSupplier.prototype.hasSupplier = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Reference package = 3;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.PackageSupplier.prototype.getPackage = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 3));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.PackageSupplier.prototype.setPackage = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PackageSupplier.prototype.clearPackage = function() {
  this.setPackage(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PackageSupplier.prototype.hasPackage = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Price price = 4;
 * @return {?proto.catalog.Price}
 */
proto.catalog.PackageSupplier.prototype.getPrice = function() {
  return /** @type{?proto.catalog.Price} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Price, 4));
};


/** @param {?proto.catalog.Price|undefined} value */
proto.catalog.PackageSupplier.prototype.setPrice = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PackageSupplier.prototype.clearPrice = function() {
  this.setPrice(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PackageSupplier.prototype.hasPrice = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional int64 date = 5;
 * @return {number}
 */
proto.catalog.PackageSupplier.prototype.getDate = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/** @param {number} value */
proto.catalog.PackageSupplier.prototype.setDate = function(value) {
  jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional int64 qte = 6;
 * @return {number}
 */
proto.catalog.PackageSupplier.prototype.getQte = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/** @param {number} value */
proto.catalog.PackageSupplier.prototype.setQte = function(value) {
  jspb.Message.setProto3IntField(this, 6, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Manufacturer.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Manufacturer.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Manufacturer} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Manufacturer.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Manufacturer}
 */
proto.catalog.Manufacturer.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Manufacturer;
  return proto.catalog.Manufacturer.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Manufacturer} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Manufacturer}
 */
proto.catalog.Manufacturer.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Manufacturer.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Manufacturer.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Manufacturer} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Manufacturer.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.Manufacturer.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Manufacturer.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.catalog.Manufacturer.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.Manufacturer.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.ItemManufacturer.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.ItemManufacturer.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.ItemManufacturer} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.ItemManufacturer.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    manufacturer: (f = msg.getManufacturer()) && proto.catalog.Reference.toObject(includeInstance, f),
    item: (f = msg.getItem()) && proto.catalog.Reference.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.ItemManufacturer}
 */
proto.catalog.ItemManufacturer.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.ItemManufacturer;
  return proto.catalog.ItemManufacturer.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.ItemManufacturer} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.ItemManufacturer}
 */
proto.catalog.ItemManufacturer.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setManufacturer(value);
      break;
    case 3:
      var value = new proto.catalog.Reference;
      reader.readMessage(value,proto.catalog.Reference.deserializeBinaryFromReader);
      msg.setItem(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.ItemManufacturer.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.ItemManufacturer.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.ItemManufacturer} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.ItemManufacturer.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getManufacturer();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
  f = message.getItem();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.catalog.Reference.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.ItemManufacturer.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.ItemManufacturer.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Reference manufacturer = 2;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.ItemManufacturer.prototype.getManufacturer = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 2));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.ItemManufacturer.prototype.setManufacturer = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemManufacturer.prototype.clearManufacturer = function() {
  this.setManufacturer(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemManufacturer.prototype.hasManufacturer = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Reference item = 3;
 * @return {?proto.catalog.Reference}
 */
proto.catalog.ItemManufacturer.prototype.getItem = function() {
  return /** @type{?proto.catalog.Reference} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Reference, 3));
};


/** @param {?proto.catalog.Reference|undefined} value */
proto.catalog.ItemManufacturer.prototype.setItem = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.ItemManufacturer.prototype.clearItem = function() {
  this.setItem(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.ItemManufacturer.prototype.hasItem = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.Dimension.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.Dimension.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.Dimension} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Dimension.toObject = function(includeInstance, msg) {
  var f, obj = {
    unitid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    value: jspb.Message.getFloatingPointFieldWithDefault(msg, 2, 0.0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.Dimension}
 */
proto.catalog.Dimension.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.Dimension;
  return proto.catalog.Dimension.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.Dimension} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.Dimension}
 */
proto.catalog.Dimension.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUnitid(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readDouble());
      msg.setValue(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.Dimension.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.Dimension.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.Dimension} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.Dimension.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUnitid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getValue();
  if (f !== 0.0) {
    writer.writeDouble(
      2,
      f
    );
  }
};


/**
 * optional string unitId = 1;
 * @return {string}
 */
proto.catalog.Dimension.prototype.getUnitid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.Dimension.prototype.setUnitid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional double value = 2;
 * @return {number}
 */
proto.catalog.Dimension.prototype.getValue = function() {
  return /** @type {number} */ (jspb.Message.getFloatingPointFieldWithDefault(this, 2, 0.0));
};


/** @param {number} value */
proto.catalog.Dimension.prototype.setValue = function(value) {
  jspb.Message.setProto3FloatField(this, 2, value);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.catalog.PropertyValue.oneofGroups_ = [[3,4,5,6,7,8,9,10]];

/**
 * @enum {number}
 */
proto.catalog.PropertyValue.ValueCase = {
  VALUE_NOT_SET: 0,
  DIMENSION_VAL: 3,
  TEXT_VAL: 4,
  NUMBER_VAL: 5,
  BOOLEAN_VAL: 6,
  DIMENSION_ARR: 7,
  TEXT_ARR: 8,
  NUMBER_ARR: 9,
  BOOLEAN_ARR: 10
};

/**
 * @return {proto.catalog.PropertyValue.ValueCase}
 */
proto.catalog.PropertyValue.prototype.getValueCase = function() {
  return /** @type {proto.catalog.PropertyValue.ValueCase} */(jspb.Message.computeOneofCase(this, proto.catalog.PropertyValue.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyValue.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyValue.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyValue} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.toObject = function(includeInstance, msg) {
  var f, obj = {
    propertydefinitionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    languagecode: jspb.Message.getFieldWithDefault(msg, 2, ""),
    dimensionVal: (f = msg.getDimensionVal()) && proto.catalog.Dimension.toObject(includeInstance, f),
    textVal: jspb.Message.getFieldWithDefault(msg, 4, ""),
    numberVal: jspb.Message.getFloatingPointFieldWithDefault(msg, 5, 0.0),
    booleanVal: jspb.Message.getBooleanFieldWithDefault(msg, 6, false),
    dimensionArr: (f = msg.getDimensionArr()) && proto.catalog.PropertyValue.Dimensions.toObject(includeInstance, f),
    textArr: (f = msg.getTextArr()) && proto.catalog.PropertyValue.Strings.toObject(includeInstance, f),
    numberArr: (f = msg.getNumberArr()) && proto.catalog.PropertyValue.Numerics.toObject(includeInstance, f),
    booleanArr: (f = msg.getBooleanArr()) && proto.catalog.PropertyValue.Booleans.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyValue}
 */
proto.catalog.PropertyValue.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyValue;
  return proto.catalog.PropertyValue.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyValue} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyValue}
 */
proto.catalog.PropertyValue.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPropertydefinitionid(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setLanguagecode(value);
      break;
    case 3:
      var value = new proto.catalog.Dimension;
      reader.readMessage(value,proto.catalog.Dimension.deserializeBinaryFromReader);
      msg.setDimensionVal(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setTextVal(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readDouble());
      msg.setNumberVal(value);
      break;
    case 6:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setBooleanVal(value);
      break;
    case 7:
      var value = new proto.catalog.PropertyValue.Dimensions;
      reader.readMessage(value,proto.catalog.PropertyValue.Dimensions.deserializeBinaryFromReader);
      msg.setDimensionArr(value);
      break;
    case 8:
      var value = new proto.catalog.PropertyValue.Strings;
      reader.readMessage(value,proto.catalog.PropertyValue.Strings.deserializeBinaryFromReader);
      msg.setTextArr(value);
      break;
    case 9:
      var value = new proto.catalog.PropertyValue.Numerics;
      reader.readMessage(value,proto.catalog.PropertyValue.Numerics.deserializeBinaryFromReader);
      msg.setNumberArr(value);
      break;
    case 10:
      var value = new proto.catalog.PropertyValue.Booleans;
      reader.readMessage(value,proto.catalog.PropertyValue.Booleans.deserializeBinaryFromReader);
      msg.setBooleanArr(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyValue.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyValue.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyValue} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPropertydefinitionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLanguagecode();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getDimensionVal();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.catalog.Dimension.serializeBinaryToWriter
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 4));
  if (f != null) {
    writer.writeString(
      4,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeDouble(
      5,
      f
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 6));
  if (f != null) {
    writer.writeBool(
      6,
      f
    );
  }
  f = message.getDimensionArr();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.catalog.PropertyValue.Dimensions.serializeBinaryToWriter
    );
  }
  f = message.getTextArr();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.catalog.PropertyValue.Strings.serializeBinaryToWriter
    );
  }
  f = message.getNumberArr();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.catalog.PropertyValue.Numerics.serializeBinaryToWriter
    );
  }
  f = message.getBooleanArr();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      proto.catalog.PropertyValue.Booleans.serializeBinaryToWriter
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.PropertyValue.Booleans.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyValue.Booleans.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyValue.Booleans.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyValue.Booleans} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Booleans.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: (f = jspb.Message.getRepeatedBooleanField(msg, 1)) == null ? undefined : f
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyValue.Booleans}
 */
proto.catalog.PropertyValue.Booleans.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyValue.Booleans;
  return proto.catalog.PropertyValue.Booleans.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyValue.Booleans} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyValue.Booleans}
 */
proto.catalog.PropertyValue.Booleans.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Array<boolean>} */ (reader.readPackedBool());
      msg.setValuesList(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyValue.Booleans.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyValue.Booleans.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyValue.Booleans} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Booleans.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writePackedBool(
      1,
      f
    );
  }
};


/**
 * repeated bool values = 1;
 * @return {!Array<boolean>}
 */
proto.catalog.PropertyValue.Booleans.prototype.getValuesList = function() {
  return /** @type {!Array<boolean>} */ (jspb.Message.getRepeatedBooleanField(this, 1));
};


/** @param {!Array<boolean>} value */
proto.catalog.PropertyValue.Booleans.prototype.setValuesList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {boolean} value
 * @param {number=} opt_index
 */
proto.catalog.PropertyValue.Booleans.prototype.addValues = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.PropertyValue.Booleans.prototype.clearValuesList = function() {
  this.setValuesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.PropertyValue.Numerics.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyValue.Numerics.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyValue.Numerics.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyValue.Numerics} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Numerics.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: (f = jspb.Message.getRepeatedFloatingPointField(msg, 1)) == null ? undefined : f
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyValue.Numerics}
 */
proto.catalog.PropertyValue.Numerics.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyValue.Numerics;
  return proto.catalog.PropertyValue.Numerics.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyValue.Numerics} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyValue.Numerics}
 */
proto.catalog.PropertyValue.Numerics.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Array<number>} */ (reader.readPackedDouble());
      msg.setValuesList(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyValue.Numerics.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyValue.Numerics.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyValue.Numerics} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Numerics.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writePackedDouble(
      1,
      f
    );
  }
};


/**
 * repeated double values = 1;
 * @return {!Array<number>}
 */
proto.catalog.PropertyValue.Numerics.prototype.getValuesList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedFloatingPointField(this, 1));
};


/** @param {!Array<number>} value */
proto.catalog.PropertyValue.Numerics.prototype.setValuesList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {number} value
 * @param {number=} opt_index
 */
proto.catalog.PropertyValue.Numerics.prototype.addValues = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.PropertyValue.Numerics.prototype.clearValuesList = function() {
  this.setValuesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.PropertyValue.Strings.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyValue.Strings.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyValue.Strings.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyValue.Strings} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Strings.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyValue.Strings}
 */
proto.catalog.PropertyValue.Strings.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyValue.Strings;
  return proto.catalog.PropertyValue.Strings.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyValue.Strings} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyValue.Strings}
 */
proto.catalog.PropertyValue.Strings.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addValues(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyValue.Strings.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyValue.Strings.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyValue.Strings} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Strings.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
};


/**
 * repeated string values = 1;
 * @return {!Array<string>}
 */
proto.catalog.PropertyValue.Strings.prototype.getValuesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.catalog.PropertyValue.Strings.prototype.setValuesList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.catalog.PropertyValue.Strings.prototype.addValues = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.PropertyValue.Strings.prototype.clearValuesList = function() {
  this.setValuesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.PropertyValue.Dimensions.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.PropertyValue.Dimensions.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.PropertyValue.Dimensions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.PropertyValue.Dimensions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Dimensions.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: jspb.Message.toObjectList(msg.getValuesList(),
    proto.catalog.PropertyValue.Dimensions.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.PropertyValue.Dimensions}
 */
proto.catalog.PropertyValue.Dimensions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.PropertyValue.Dimensions;
  return proto.catalog.PropertyValue.Dimensions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.PropertyValue.Dimensions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.PropertyValue.Dimensions}
 */
proto.catalog.PropertyValue.Dimensions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.PropertyValue.Dimensions;
      reader.readMessage(value,proto.catalog.PropertyValue.Dimensions.deserializeBinaryFromReader);
      msg.addValues(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.PropertyValue.Dimensions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.PropertyValue.Dimensions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.PropertyValue.Dimensions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.PropertyValue.Dimensions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.catalog.PropertyValue.Dimensions.serializeBinaryToWriter
    );
  }
};


/**
 * repeated Dimensions values = 1;
 * @return {!Array<!proto.catalog.PropertyValue.Dimensions>}
 */
proto.catalog.PropertyValue.Dimensions.prototype.getValuesList = function() {
  return /** @type{!Array<!proto.catalog.PropertyValue.Dimensions>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.catalog.PropertyValue.Dimensions, 1));
};


/** @param {!Array<!proto.catalog.PropertyValue.Dimensions>} value */
proto.catalog.PropertyValue.Dimensions.prototype.setValuesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.catalog.PropertyValue.Dimensions=} opt_value
 * @param {number=} opt_index
 * @return {!proto.catalog.PropertyValue.Dimensions}
 */
proto.catalog.PropertyValue.Dimensions.prototype.addValues = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.catalog.PropertyValue.Dimensions, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.PropertyValue.Dimensions.prototype.clearValuesList = function() {
  this.setValuesList([]);
};


/**
 * optional string propertyDefinitionId = 1;
 * @return {string}
 */
proto.catalog.PropertyValue.prototype.getPropertydefinitionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.PropertyValue.prototype.setPropertydefinitionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string languageCode = 2;
 * @return {string}
 */
proto.catalog.PropertyValue.prototype.getLanguagecode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.PropertyValue.prototype.setLanguagecode = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional Dimension dimension_val = 3;
 * @return {?proto.catalog.Dimension}
 */
proto.catalog.PropertyValue.prototype.getDimensionVal = function() {
  return /** @type{?proto.catalog.Dimension} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Dimension, 3));
};


/** @param {?proto.catalog.Dimension|undefined} value */
proto.catalog.PropertyValue.prototype.setDimensionVal = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearDimensionVal = function() {
  this.setDimensionVal(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasDimensionVal = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string text_val = 4;
 * @return {string}
 */
proto.catalog.PropertyValue.prototype.getTextVal = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.catalog.PropertyValue.prototype.setTextVal = function(value) {
  jspb.Message.setOneofField(this, 4, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearTextVal = function() {
  jspb.Message.setOneofField(this, 4, proto.catalog.PropertyValue.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasTextVal = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional double number_val = 5;
 * @return {number}
 */
proto.catalog.PropertyValue.prototype.getNumberVal = function() {
  return /** @type {number} */ (jspb.Message.getFloatingPointFieldWithDefault(this, 5, 0.0));
};


/** @param {number} value */
proto.catalog.PropertyValue.prototype.setNumberVal = function(value) {
  jspb.Message.setOneofField(this, 5, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearNumberVal = function() {
  jspb.Message.setOneofField(this, 5, proto.catalog.PropertyValue.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasNumberVal = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional bool boolean_val = 6;
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.getBooleanVal = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 6, false));
};


/** @param {boolean} value */
proto.catalog.PropertyValue.prototype.setBooleanVal = function(value) {
  jspb.Message.setOneofField(this, 6, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearBooleanVal = function() {
  jspb.Message.setOneofField(this, 6, proto.catalog.PropertyValue.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasBooleanVal = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional Dimensions dimension_arr = 7;
 * @return {?proto.catalog.PropertyValue.Dimensions}
 */
proto.catalog.PropertyValue.prototype.getDimensionArr = function() {
  return /** @type{?proto.catalog.PropertyValue.Dimensions} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyValue.Dimensions, 7));
};


/** @param {?proto.catalog.PropertyValue.Dimensions|undefined} value */
proto.catalog.PropertyValue.prototype.setDimensionArr = function(value) {
  jspb.Message.setOneofWrapperField(this, 7, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearDimensionArr = function() {
  this.setDimensionArr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasDimensionArr = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional Strings text_arr = 8;
 * @return {?proto.catalog.PropertyValue.Strings}
 */
proto.catalog.PropertyValue.prototype.getTextArr = function() {
  return /** @type{?proto.catalog.PropertyValue.Strings} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyValue.Strings, 8));
};


/** @param {?proto.catalog.PropertyValue.Strings|undefined} value */
proto.catalog.PropertyValue.prototype.setTextArr = function(value) {
  jspb.Message.setOneofWrapperField(this, 8, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearTextArr = function() {
  this.setTextArr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasTextArr = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional Numerics number_arr = 9;
 * @return {?proto.catalog.PropertyValue.Numerics}
 */
proto.catalog.PropertyValue.prototype.getNumberArr = function() {
  return /** @type{?proto.catalog.PropertyValue.Numerics} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyValue.Numerics, 9));
};


/** @param {?proto.catalog.PropertyValue.Numerics|undefined} value */
proto.catalog.PropertyValue.prototype.setNumberArr = function(value) {
  jspb.Message.setOneofWrapperField(this, 9, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearNumberArr = function() {
  this.setNumberArr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasNumberArr = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional Booleans boolean_arr = 10;
 * @return {?proto.catalog.PropertyValue.Booleans}
 */
proto.catalog.PropertyValue.prototype.getBooleanArr = function() {
  return /** @type{?proto.catalog.PropertyValue.Booleans} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyValue.Booleans, 10));
};


/** @param {?proto.catalog.PropertyValue.Booleans|undefined} value */
proto.catalog.PropertyValue.prototype.setBooleanArr = function(value) {
  jspb.Message.setOneofWrapperField(this, 10, proto.catalog.PropertyValue.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.PropertyValue.prototype.clearBooleanArr = function() {
  this.setBooleanArr(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.PropertyValue.prototype.hasBooleanArr = function() {
  return jspb.Message.getField(this, 10) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.ItemInstance.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.ItemInstance.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.ItemInstance.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.ItemInstance} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.ItemInstance.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    itemdefinitionid: jspb.Message.getFieldWithDefault(msg, 2, ""),
    valuesList: jspb.Message.toObjectList(msg.getValuesList(),
    proto.catalog.PropertyValue.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.ItemInstance}
 */
proto.catalog.ItemInstance.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.ItemInstance;
  return proto.catalog.ItemInstance.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.ItemInstance} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.ItemInstance}
 */
proto.catalog.ItemInstance.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setItemdefinitionid(value);
      break;
    case 3:
      var value = new proto.catalog.PropertyValue;
      reader.readMessage(value,proto.catalog.PropertyValue.deserializeBinaryFromReader);
      msg.addValues(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.ItemInstance.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.ItemInstance.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.ItemInstance} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.ItemInstance.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getItemdefinitionid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.catalog.PropertyValue.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.ItemInstance.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.ItemInstance.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string itemDefinitionId = 2;
 * @return {string}
 */
proto.catalog.ItemInstance.prototype.getItemdefinitionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.ItemInstance.prototype.setItemdefinitionid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated PropertyValue values = 3;
 * @return {!Array<!proto.catalog.PropertyValue>}
 */
proto.catalog.ItemInstance.prototype.getValuesList = function() {
  return /** @type{!Array<!proto.catalog.PropertyValue>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.catalog.PropertyValue, 3));
};


/** @param {!Array<!proto.catalog.PropertyValue>} value */
proto.catalog.ItemInstance.prototype.setValuesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.catalog.PropertyValue=} opt_value
 * @param {number=} opt_index
 * @return {!proto.catalog.PropertyValue}
 */
proto.catalog.ItemInstance.prototype.addValues = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.catalog.PropertyValue, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.ItemInstance.prototype.clearValuesList = function() {
  this.setValuesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveUnitOfMeasureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveUnitOfMeasureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveUnitOfMeasureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveUnitOfMeasureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    unitofmeasure: (f = msg.getUnitofmeasure()) && proto.catalog.UnitOfMeasure.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveUnitOfMeasureRequest}
 */
proto.catalog.SaveUnitOfMeasureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveUnitOfMeasureRequest;
  return proto.catalog.SaveUnitOfMeasureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveUnitOfMeasureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveUnitOfMeasureRequest}
 */
proto.catalog.SaveUnitOfMeasureRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.UnitOfMeasure;
      reader.readMessage(value,proto.catalog.UnitOfMeasure.deserializeBinaryFromReader);
      msg.setUnitofmeasure(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveUnitOfMeasureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveUnitOfMeasureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveUnitOfMeasureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveUnitOfMeasureRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUnitofmeasure();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.UnitOfMeasure.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveUnitOfMeasureRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveUnitOfMeasureRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional UnitOfMeasure unitOfMeasure = 2;
 * @return {?proto.catalog.UnitOfMeasure}
 */
proto.catalog.SaveUnitOfMeasureRequest.prototype.getUnitofmeasure = function() {
  return /** @type{?proto.catalog.UnitOfMeasure} */ (
    jspb.Message.getWrapperField(this, proto.catalog.UnitOfMeasure, 2));
};


/** @param {?proto.catalog.UnitOfMeasure|undefined} value */
proto.catalog.SaveUnitOfMeasureRequest.prototype.setUnitofmeasure = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveUnitOfMeasureRequest.prototype.clearUnitofmeasure = function() {
  this.setUnitofmeasure(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveUnitOfMeasureRequest.prototype.hasUnitofmeasure = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveUnitOfMeasureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveUnitOfMeasureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveUnitOfMeasureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveUnitOfMeasureResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveUnitOfMeasureResponse}
 */
proto.catalog.SaveUnitOfMeasureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveUnitOfMeasureResponse;
  return proto.catalog.SaveUnitOfMeasureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveUnitOfMeasureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveUnitOfMeasureResponse}
 */
proto.catalog.SaveUnitOfMeasureResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveUnitOfMeasureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveUnitOfMeasureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveUnitOfMeasureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveUnitOfMeasureResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveUnitOfMeasureResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveUnitOfMeasureResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SavePropertyDefinitionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SavePropertyDefinitionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SavePropertyDefinitionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePropertyDefinitionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    propertydefinition: (f = msg.getPropertydefinition()) && proto.catalog.PropertyDefinition.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SavePropertyDefinitionRequest}
 */
proto.catalog.SavePropertyDefinitionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SavePropertyDefinitionRequest;
  return proto.catalog.SavePropertyDefinitionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SavePropertyDefinitionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SavePropertyDefinitionRequest}
 */
proto.catalog.SavePropertyDefinitionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.PropertyDefinition;
      reader.readMessage(value,proto.catalog.PropertyDefinition.deserializeBinaryFromReader);
      msg.setPropertydefinition(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SavePropertyDefinitionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SavePropertyDefinitionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SavePropertyDefinitionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePropertyDefinitionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPropertydefinition();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.PropertyDefinition.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SavePropertyDefinitionRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SavePropertyDefinitionRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional PropertyDefinition propertyDefinition = 2;
 * @return {?proto.catalog.PropertyDefinition}
 */
proto.catalog.SavePropertyDefinitionRequest.prototype.getPropertydefinition = function() {
  return /** @type{?proto.catalog.PropertyDefinition} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyDefinition, 2));
};


/** @param {?proto.catalog.PropertyDefinition|undefined} value */
proto.catalog.SavePropertyDefinitionRequest.prototype.setPropertydefinition = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SavePropertyDefinitionRequest.prototype.clearPropertydefinition = function() {
  this.setPropertydefinition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SavePropertyDefinitionRequest.prototype.hasPropertydefinition = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SavePropertyDefinitionResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SavePropertyDefinitionResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SavePropertyDefinitionResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePropertyDefinitionResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SavePropertyDefinitionResponse}
 */
proto.catalog.SavePropertyDefinitionResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SavePropertyDefinitionResponse;
  return proto.catalog.SavePropertyDefinitionResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SavePropertyDefinitionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SavePropertyDefinitionResponse}
 */
proto.catalog.SavePropertyDefinitionResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SavePropertyDefinitionResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SavePropertyDefinitionResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SavePropertyDefinitionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePropertyDefinitionResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SavePropertyDefinitionResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SavePropertyDefinitionResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveItemDefinitionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveItemDefinitionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveItemDefinitionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemDefinitionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    itemdefinition: (f = msg.getItemdefinition()) && proto.catalog.ItemDefinition.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveItemDefinitionRequest}
 */
proto.catalog.SaveItemDefinitionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveItemDefinitionRequest;
  return proto.catalog.SaveItemDefinitionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveItemDefinitionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveItemDefinitionRequest}
 */
proto.catalog.SaveItemDefinitionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.ItemDefinition;
      reader.readMessage(value,proto.catalog.ItemDefinition.deserializeBinaryFromReader);
      msg.setItemdefinition(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveItemDefinitionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveItemDefinitionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveItemDefinitionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemDefinitionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getItemdefinition();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.ItemDefinition.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveItemDefinitionRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveItemDefinitionRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ItemDefinition itemDefinition = 2;
 * @return {?proto.catalog.ItemDefinition}
 */
proto.catalog.SaveItemDefinitionRequest.prototype.getItemdefinition = function() {
  return /** @type{?proto.catalog.ItemDefinition} */ (
    jspb.Message.getWrapperField(this, proto.catalog.ItemDefinition, 2));
};


/** @param {?proto.catalog.ItemDefinition|undefined} value */
proto.catalog.SaveItemDefinitionRequest.prototype.setItemdefinition = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveItemDefinitionRequest.prototype.clearItemdefinition = function() {
  this.setItemdefinition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveItemDefinitionRequest.prototype.hasItemdefinition = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveItemDefinitionResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveItemDefinitionResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveItemDefinitionResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemDefinitionResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveItemDefinitionResponse}
 */
proto.catalog.SaveItemDefinitionResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveItemDefinitionResponse;
  return proto.catalog.SaveItemDefinitionResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveItemDefinitionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveItemDefinitionResponse}
 */
proto.catalog.SaveItemDefinitionResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveItemDefinitionResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveItemDefinitionResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveItemDefinitionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemDefinitionResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveItemDefinitionResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveItemDefinitionResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveItemInstanceRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveItemInstanceRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveItemInstanceRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemInstanceRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    iteminstance: (f = msg.getIteminstance()) && proto.catalog.ItemInstance.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveItemInstanceRequest}
 */
proto.catalog.SaveItemInstanceRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveItemInstanceRequest;
  return proto.catalog.SaveItemInstanceRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveItemInstanceRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveItemInstanceRequest}
 */
proto.catalog.SaveItemInstanceRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.ItemInstance;
      reader.readMessage(value,proto.catalog.ItemInstance.deserializeBinaryFromReader);
      msg.setIteminstance(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveItemInstanceRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveItemInstanceRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveItemInstanceRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemInstanceRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getIteminstance();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.ItemInstance.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveItemInstanceRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveItemInstanceRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ItemInstance itemInstance = 2;
 * @return {?proto.catalog.ItemInstance}
 */
proto.catalog.SaveItemInstanceRequest.prototype.getIteminstance = function() {
  return /** @type{?proto.catalog.ItemInstance} */ (
    jspb.Message.getWrapperField(this, proto.catalog.ItemInstance, 2));
};


/** @param {?proto.catalog.ItemInstance|undefined} value */
proto.catalog.SaveItemInstanceRequest.prototype.setIteminstance = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveItemInstanceRequest.prototype.clearIteminstance = function() {
  this.setIteminstance(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveItemInstanceRequest.prototype.hasIteminstance = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveItemInstanceResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveItemInstanceResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveItemInstanceResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemInstanceResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveItemInstanceResponse}
 */
proto.catalog.SaveItemInstanceResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveItemInstanceResponse;
  return proto.catalog.SaveItemInstanceResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveItemInstanceResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveItemInstanceResponse}
 */
proto.catalog.SaveItemInstanceResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveItemInstanceResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveItemInstanceResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveItemInstanceResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemInstanceResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveItemInstanceResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveItemInstanceResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveManufacturerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveManufacturerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveManufacturerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveManufacturerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    manufacturer: (f = msg.getManufacturer()) && proto.catalog.Manufacturer.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveManufacturerRequest}
 */
proto.catalog.SaveManufacturerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveManufacturerRequest;
  return proto.catalog.SaveManufacturerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveManufacturerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveManufacturerRequest}
 */
proto.catalog.SaveManufacturerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Manufacturer;
      reader.readMessage(value,proto.catalog.Manufacturer.deserializeBinaryFromReader);
      msg.setManufacturer(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveManufacturerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveManufacturerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveManufacturerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveManufacturerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getManufacturer();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Manufacturer.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveManufacturerRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveManufacturerRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Manufacturer manufacturer = 2;
 * @return {?proto.catalog.Manufacturer}
 */
proto.catalog.SaveManufacturerRequest.prototype.getManufacturer = function() {
  return /** @type{?proto.catalog.Manufacturer} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Manufacturer, 2));
};


/** @param {?proto.catalog.Manufacturer|undefined} value */
proto.catalog.SaveManufacturerRequest.prototype.setManufacturer = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveManufacturerRequest.prototype.clearManufacturer = function() {
  this.setManufacturer(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveManufacturerRequest.prototype.hasManufacturer = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveManufacturerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveManufacturerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveManufacturerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveManufacturerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveManufacturerResponse}
 */
proto.catalog.SaveManufacturerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveManufacturerResponse;
  return proto.catalog.SaveManufacturerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveManufacturerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveManufacturerResponse}
 */
proto.catalog.SaveManufacturerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveManufacturerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveManufacturerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveManufacturerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveManufacturerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveManufacturerResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveManufacturerResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveSupplierRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveSupplierRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveSupplierRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveSupplierRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    supplier: (f = msg.getSupplier()) && proto.catalog.Supplier.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveSupplierRequest}
 */
proto.catalog.SaveSupplierRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveSupplierRequest;
  return proto.catalog.SaveSupplierRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveSupplierRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveSupplierRequest}
 */
proto.catalog.SaveSupplierRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Supplier;
      reader.readMessage(value,proto.catalog.Supplier.deserializeBinaryFromReader);
      msg.setSupplier(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveSupplierRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveSupplierRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveSupplierRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveSupplierRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSupplier();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Supplier.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveSupplierRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveSupplierRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Supplier supplier = 2;
 * @return {?proto.catalog.Supplier}
 */
proto.catalog.SaveSupplierRequest.prototype.getSupplier = function() {
  return /** @type{?proto.catalog.Supplier} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Supplier, 2));
};


/** @param {?proto.catalog.Supplier|undefined} value */
proto.catalog.SaveSupplierRequest.prototype.setSupplier = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveSupplierRequest.prototype.clearSupplier = function() {
  this.setSupplier(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveSupplierRequest.prototype.hasSupplier = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveSupplierResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveSupplierResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveSupplierResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveSupplierResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveSupplierResponse}
 */
proto.catalog.SaveSupplierResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveSupplierResponse;
  return proto.catalog.SaveSupplierResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveSupplierResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveSupplierResponse}
 */
proto.catalog.SaveSupplierResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveSupplierResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveSupplierResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveSupplierResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveSupplierResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveSupplierResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveSupplierResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveLocalisationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveLocalisationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveLocalisationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveLocalisationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    localisation: (f = msg.getLocalisation()) && proto.catalog.Localisation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveLocalisationRequest}
 */
proto.catalog.SaveLocalisationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveLocalisationRequest;
  return proto.catalog.SaveLocalisationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveLocalisationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveLocalisationRequest}
 */
proto.catalog.SaveLocalisationRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Localisation;
      reader.readMessage(value,proto.catalog.Localisation.deserializeBinaryFromReader);
      msg.setLocalisation(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveLocalisationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveLocalisationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveLocalisationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveLocalisationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLocalisation();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Localisation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveLocalisationRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveLocalisationRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Localisation localisation = 2;
 * @return {?proto.catalog.Localisation}
 */
proto.catalog.SaveLocalisationRequest.prototype.getLocalisation = function() {
  return /** @type{?proto.catalog.Localisation} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Localisation, 2));
};


/** @param {?proto.catalog.Localisation|undefined} value */
proto.catalog.SaveLocalisationRequest.prototype.setLocalisation = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveLocalisationRequest.prototype.clearLocalisation = function() {
  this.setLocalisation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveLocalisationRequest.prototype.hasLocalisation = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveLocalisationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveLocalisationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveLocalisationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveLocalisationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveLocalisationResponse}
 */
proto.catalog.SaveLocalisationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveLocalisationResponse;
  return proto.catalog.SaveLocalisationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveLocalisationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveLocalisationResponse}
 */
proto.catalog.SaveLocalisationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveLocalisationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveLocalisationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveLocalisationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveLocalisationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveLocalisationResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveLocalisationResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    category: (f = msg.getCategory()) && proto.catalog.Category.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveCategoryRequest}
 */
proto.catalog.SaveCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveCategoryRequest;
  return proto.catalog.SaveCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveCategoryRequest}
 */
proto.catalog.SaveCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Category;
      reader.readMessage(value,proto.catalog.Category.deserializeBinaryFromReader);
      msg.setCategory(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCategory();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Category.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveCategoryRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveCategoryRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Category category = 2;
 * @return {?proto.catalog.Category}
 */
proto.catalog.SaveCategoryRequest.prototype.getCategory = function() {
  return /** @type{?proto.catalog.Category} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Category, 2));
};


/** @param {?proto.catalog.Category|undefined} value */
proto.catalog.SaveCategoryRequest.prototype.setCategory = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveCategoryRequest.prototype.clearCategory = function() {
  this.setCategory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveCategoryRequest.prototype.hasCategory = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveCategoryResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveCategoryResponse}
 */
proto.catalog.SaveCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveCategoryResponse;
  return proto.catalog.SaveCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveCategoryResponse}
 */
proto.catalog.SaveCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveCategoryResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveCategoryResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveInventoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveInventoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveInventoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveInventoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    inventory: (f = msg.getInventory()) && proto.catalog.Inventory.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveInventoryRequest}
 */
proto.catalog.SaveInventoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveInventoryRequest;
  return proto.catalog.SaveInventoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveInventoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveInventoryRequest}
 */
proto.catalog.SaveInventoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Inventory;
      reader.readMessage(value,proto.catalog.Inventory.deserializeBinaryFromReader);
      msg.setInventory(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveInventoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveInventoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveInventoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveInventoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getInventory();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Inventory.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveInventoryRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveInventoryRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Inventory inventory = 2;
 * @return {?proto.catalog.Inventory}
 */
proto.catalog.SaveInventoryRequest.prototype.getInventory = function() {
  return /** @type{?proto.catalog.Inventory} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Inventory, 2));
};


/** @param {?proto.catalog.Inventory|undefined} value */
proto.catalog.SaveInventoryRequest.prototype.setInventory = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveInventoryRequest.prototype.clearInventory = function() {
  this.setInventory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveInventoryRequest.prototype.hasInventory = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveInventoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveInventoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveInventoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveInventoryResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveInventoryResponse}
 */
proto.catalog.SaveInventoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveInventoryResponse;
  return proto.catalog.SaveInventoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveInventoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveInventoryResponse}
 */
proto.catalog.SaveInventoryResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveInventoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveInventoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveInventoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveInventoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveInventoryResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveInventoryResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SavePackageRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SavePackageRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SavePackageRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    pb_package: (f = msg.getPackage()) && proto.catalog.Package.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SavePackageRequest}
 */
proto.catalog.SavePackageRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SavePackageRequest;
  return proto.catalog.SavePackageRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SavePackageRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SavePackageRequest}
 */
proto.catalog.SavePackageRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Package;
      reader.readMessage(value,proto.catalog.Package.deserializeBinaryFromReader);
      msg.setPackage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SavePackageRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SavePackageRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SavePackageRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPackage();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Package.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SavePackageRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SavePackageRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Package package = 2;
 * @return {?proto.catalog.Package}
 */
proto.catalog.SavePackageRequest.prototype.getPackage = function() {
  return /** @type{?proto.catalog.Package} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Package, 2));
};


/** @param {?proto.catalog.Package|undefined} value */
proto.catalog.SavePackageRequest.prototype.setPackage = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SavePackageRequest.prototype.clearPackage = function() {
  this.setPackage(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SavePackageRequest.prototype.hasPackage = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SavePackageResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SavePackageResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SavePackageResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SavePackageResponse}
 */
proto.catalog.SavePackageResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SavePackageResponse;
  return proto.catalog.SavePackageResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SavePackageResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SavePackageResponse}
 */
proto.catalog.SavePackageResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SavePackageResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SavePackageResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SavePackageResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SavePackageResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SavePackageResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SavePackageSupplierRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SavePackageSupplierRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SavePackageSupplierRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageSupplierRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    packagesupplier: (f = msg.getPackagesupplier()) && proto.catalog.PackageSupplier.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SavePackageSupplierRequest}
 */
proto.catalog.SavePackageSupplierRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SavePackageSupplierRequest;
  return proto.catalog.SavePackageSupplierRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SavePackageSupplierRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SavePackageSupplierRequest}
 */
proto.catalog.SavePackageSupplierRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.PackageSupplier;
      reader.readMessage(value,proto.catalog.PackageSupplier.deserializeBinaryFromReader);
      msg.setPackagesupplier(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SavePackageSupplierRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SavePackageSupplierRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SavePackageSupplierRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageSupplierRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPackagesupplier();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.PackageSupplier.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SavePackageSupplierRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SavePackageSupplierRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional PackageSupplier packageSupplier = 2;
 * @return {?proto.catalog.PackageSupplier}
 */
proto.catalog.SavePackageSupplierRequest.prototype.getPackagesupplier = function() {
  return /** @type{?proto.catalog.PackageSupplier} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PackageSupplier, 2));
};


/** @param {?proto.catalog.PackageSupplier|undefined} value */
proto.catalog.SavePackageSupplierRequest.prototype.setPackagesupplier = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SavePackageSupplierRequest.prototype.clearPackagesupplier = function() {
  this.setPackagesupplier(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SavePackageSupplierRequest.prototype.hasPackagesupplier = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SavePackageSupplierResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SavePackageSupplierResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SavePackageSupplierResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageSupplierResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SavePackageSupplierResponse}
 */
proto.catalog.SavePackageSupplierResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SavePackageSupplierResponse;
  return proto.catalog.SavePackageSupplierResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SavePackageSupplierResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SavePackageSupplierResponse}
 */
proto.catalog.SavePackageSupplierResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SavePackageSupplierResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SavePackageSupplierResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SavePackageSupplierResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SavePackageSupplierResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SavePackageSupplierResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SavePackageSupplierResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveItemManufacturerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveItemManufacturerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveItemManufacturerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemManufacturerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    itemmanafacturer: (f = msg.getItemmanafacturer()) && proto.catalog.ItemManufacturer.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveItemManufacturerRequest}
 */
proto.catalog.SaveItemManufacturerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveItemManufacturerRequest;
  return proto.catalog.SaveItemManufacturerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveItemManufacturerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveItemManufacturerRequest}
 */
proto.catalog.SaveItemManufacturerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.ItemManufacturer;
      reader.readMessage(value,proto.catalog.ItemManufacturer.deserializeBinaryFromReader);
      msg.setItemmanafacturer(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveItemManufacturerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveItemManufacturerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveItemManufacturerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemManufacturerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getItemmanafacturer();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.ItemManufacturer.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.SaveItemManufacturerRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveItemManufacturerRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ItemManufacturer itemManafacturer = 2;
 * @return {?proto.catalog.ItemManufacturer}
 */
proto.catalog.SaveItemManufacturerRequest.prototype.getItemmanafacturer = function() {
  return /** @type{?proto.catalog.ItemManufacturer} */ (
    jspb.Message.getWrapperField(this, proto.catalog.ItemManufacturer, 2));
};


/** @param {?proto.catalog.ItemManufacturer|undefined} value */
proto.catalog.SaveItemManufacturerRequest.prototype.setItemmanafacturer = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.SaveItemManufacturerRequest.prototype.clearItemmanafacturer = function() {
  this.setItemmanafacturer(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.SaveItemManufacturerRequest.prototype.hasItemmanafacturer = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.SaveItemManufacturerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.SaveItemManufacturerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.SaveItemManufacturerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemManufacturerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.SaveItemManufacturerResponse}
 */
proto.catalog.SaveItemManufacturerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.SaveItemManufacturerResponse;
  return proto.catalog.SaveItemManufacturerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.SaveItemManufacturerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.SaveItemManufacturerResponse}
 */
proto.catalog.SaveItemManufacturerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.SaveItemManufacturerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.SaveItemManufacturerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.SaveItemManufacturerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.SaveItemManufacturerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.catalog.SaveItemManufacturerResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.SaveItemManufacturerResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.GetSupplierRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.GetSupplierRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.GetSupplierRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    supplierid: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.GetSupplierRequest}
 */
proto.catalog.GetSupplierRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.GetSupplierRequest;
  return proto.catalog.GetSupplierRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.GetSupplierRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.GetSupplierRequest}
 */
proto.catalog.GetSupplierRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setSupplierid(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.GetSupplierRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.GetSupplierRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.GetSupplierRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSupplierid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.GetSupplierRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.GetSupplierRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string supplierId = 2;
 * @return {string}
 */
proto.catalog.GetSupplierRequest.prototype.getSupplierid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.GetSupplierRequest.prototype.setSupplierid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.GetSupplierResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.GetSupplierResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.GetSupplierResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    supplier: (f = msg.getSupplier()) && proto.catalog.Supplier.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.GetSupplierResponse}
 */
proto.catalog.GetSupplierResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.GetSupplierResponse;
  return proto.catalog.GetSupplierResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.GetSupplierResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.GetSupplierResponse}
 */
proto.catalog.GetSupplierResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.Supplier;
      reader.readMessage(value,proto.catalog.Supplier.deserializeBinaryFromReader);
      msg.setSupplier(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.GetSupplierResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.GetSupplierResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.GetSupplierResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSupplier();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.catalog.Supplier.serializeBinaryToWriter
    );
  }
};


/**
 * optional Supplier supplier = 1;
 * @return {?proto.catalog.Supplier}
 */
proto.catalog.GetSupplierResponse.prototype.getSupplier = function() {
  return /** @type{?proto.catalog.Supplier} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Supplier, 1));
};


/** @param {?proto.catalog.Supplier|undefined} value */
proto.catalog.GetSupplierResponse.prototype.setSupplier = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.GetSupplierResponse.prototype.clearSupplier = function() {
  this.setSupplier(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.GetSupplierResponse.prototype.hasSupplier = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.GetSupplierPackagesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.GetSupplierPackagesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.GetSupplierPackagesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierPackagesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    supplierid: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.GetSupplierPackagesRequest}
 */
proto.catalog.GetSupplierPackagesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.GetSupplierPackagesRequest;
  return proto.catalog.GetSupplierPackagesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.GetSupplierPackagesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.GetSupplierPackagesRequest}
 */
proto.catalog.GetSupplierPackagesRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setSupplierid(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.GetSupplierPackagesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.GetSupplierPackagesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.GetSupplierPackagesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierPackagesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSupplierid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.GetSupplierPackagesRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.GetSupplierPackagesRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string supplierId = 2;
 * @return {string}
 */
proto.catalog.GetSupplierPackagesRequest.prototype.getSupplierid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.GetSupplierPackagesRequest.prototype.setSupplierid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.catalog.GetSupplierPackagesResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.GetSupplierPackagesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.GetSupplierPackagesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.GetSupplierPackagesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierPackagesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    pacakgessupplierList: jspb.Message.toObjectList(msg.getPacakgessupplierList(),
    proto.catalog.PackageSupplier.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.GetSupplierPackagesResponse}
 */
proto.catalog.GetSupplierPackagesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.GetSupplierPackagesResponse;
  return proto.catalog.GetSupplierPackagesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.GetSupplierPackagesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.GetSupplierPackagesResponse}
 */
proto.catalog.GetSupplierPackagesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.PackageSupplier;
      reader.readMessage(value,proto.catalog.PackageSupplier.deserializeBinaryFromReader);
      msg.addPacakgessupplier(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.GetSupplierPackagesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.GetSupplierPackagesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.GetSupplierPackagesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetSupplierPackagesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPacakgessupplierList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.catalog.PackageSupplier.serializeBinaryToWriter
    );
  }
};


/**
 * repeated PackageSupplier pacakgesSupplier = 1;
 * @return {!Array<!proto.catalog.PackageSupplier>}
 */
proto.catalog.GetSupplierPackagesResponse.prototype.getPacakgessupplierList = function() {
  return /** @type{!Array<!proto.catalog.PackageSupplier>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.catalog.PackageSupplier, 1));
};


/** @param {!Array<!proto.catalog.PackageSupplier>} value */
proto.catalog.GetSupplierPackagesResponse.prototype.setPacakgessupplierList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.catalog.PackageSupplier=} opt_value
 * @param {number=} opt_index
 * @return {!proto.catalog.PackageSupplier}
 */
proto.catalog.GetSupplierPackagesResponse.prototype.addPacakgessupplier = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.catalog.PackageSupplier, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.catalog.GetSupplierPackagesResponse.prototype.clearPacakgessupplierList = function() {
  this.setPacakgessupplierList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.GetPackageRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.GetPackageRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.GetPackageRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetPackageRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    packageid: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.GetPackageRequest}
 */
proto.catalog.GetPackageRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.GetPackageRequest;
  return proto.catalog.GetPackageRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.GetPackageRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.GetPackageRequest}
 */
proto.catalog.GetPackageRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPackageid(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.GetPackageRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.GetPackageRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.GetPackageRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetPackageRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPackageid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.GetPackageRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.GetPackageRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string packageId = 2;
 * @return {string}
 */
proto.catalog.GetPackageRequest.prototype.getPackageid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.catalog.GetPackageRequest.prototype.setPackageid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.GetPackageResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.GetPackageResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.GetPackageResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetPackageResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    pacakge: (f = msg.getPacakge()) && proto.catalog.Package.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.GetPackageResponse}
 */
proto.catalog.GetPackageResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.GetPackageResponse;
  return proto.catalog.GetPackageResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.GetPackageResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.GetPackageResponse}
 */
proto.catalog.GetPackageResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.catalog.Package;
      reader.readMessage(value,proto.catalog.Package.deserializeBinaryFromReader);
      msg.setPacakge(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.GetPackageResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.GetPackageResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.GetPackageResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.GetPackageResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPacakge();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.catalog.Package.serializeBinaryToWriter
    );
  }
};


/**
 * optional Package pacakge = 1;
 * @return {?proto.catalog.Package}
 */
proto.catalog.GetPackageResponse.prototype.getPacakge = function() {
  return /** @type{?proto.catalog.Package} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Package, 1));
};


/** @param {?proto.catalog.Package|undefined} value */
proto.catalog.GetPackageResponse.prototype.setPacakge = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.GetPackageResponse.prototype.clearPacakge = function() {
  this.setPacakge(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.GetPackageResponse.prototype.hasPacakge = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeletePackageSupplierRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeletePackageSupplierRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeletePackageSupplierRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageSupplierRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    packagesupplier: (f = msg.getPackagesupplier()) && proto.catalog.PackageSupplier.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeletePackageSupplierRequest}
 */
proto.catalog.DeletePackageSupplierRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeletePackageSupplierRequest;
  return proto.catalog.DeletePackageSupplierRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeletePackageSupplierRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeletePackageSupplierRequest}
 */
proto.catalog.DeletePackageSupplierRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.PackageSupplier;
      reader.readMessage(value,proto.catalog.PackageSupplier.deserializeBinaryFromReader);
      msg.setPackagesupplier(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeletePackageSupplierRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeletePackageSupplierRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeletePackageSupplierRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageSupplierRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPackagesupplier();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.PackageSupplier.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeletePackageSupplierRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeletePackageSupplierRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional PackageSupplier packageSupplier = 2;
 * @return {?proto.catalog.PackageSupplier}
 */
proto.catalog.DeletePackageSupplierRequest.prototype.getPackagesupplier = function() {
  return /** @type{?proto.catalog.PackageSupplier} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PackageSupplier, 2));
};


/** @param {?proto.catalog.PackageSupplier|undefined} value */
proto.catalog.DeletePackageSupplierRequest.prototype.setPackagesupplier = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeletePackageSupplierRequest.prototype.clearPackagesupplier = function() {
  this.setPackagesupplier(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeletePackageSupplierRequest.prototype.hasPackagesupplier = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeletePackageSupplierResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeletePackageSupplierResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeletePackageSupplierResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageSupplierResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeletePackageSupplierResponse}
 */
proto.catalog.DeletePackageSupplierResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeletePackageSupplierResponse;
  return proto.catalog.DeletePackageSupplierResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeletePackageSupplierResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeletePackageSupplierResponse}
 */
proto.catalog.DeletePackageSupplierResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeletePackageSupplierResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeletePackageSupplierResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeletePackageSupplierResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageSupplierResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeletePackageSupplierResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeletePackageSupplierResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeletePackageRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeletePackageRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeletePackageRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    pb_package: (f = msg.getPackage()) && proto.catalog.Package.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeletePackageRequest}
 */
proto.catalog.DeletePackageRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeletePackageRequest;
  return proto.catalog.DeletePackageRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeletePackageRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeletePackageRequest}
 */
proto.catalog.DeletePackageRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Package;
      reader.readMessage(value,proto.catalog.Package.deserializeBinaryFromReader);
      msg.setPackage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeletePackageRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeletePackageRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeletePackageRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPackage();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Package.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeletePackageRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeletePackageRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Package package = 2;
 * @return {?proto.catalog.Package}
 */
proto.catalog.DeletePackageRequest.prototype.getPackage = function() {
  return /** @type{?proto.catalog.Package} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Package, 2));
};


/** @param {?proto.catalog.Package|undefined} value */
proto.catalog.DeletePackageRequest.prototype.setPackage = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeletePackageRequest.prototype.clearPackage = function() {
  this.setPackage(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeletePackageRequest.prototype.hasPackage = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeletePackageResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeletePackageResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeletePackageResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeletePackageResponse}
 */
proto.catalog.DeletePackageResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeletePackageResponse;
  return proto.catalog.DeletePackageResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeletePackageResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeletePackageResponse}
 */
proto.catalog.DeletePackageResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeletePackageResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeletePackageResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeletePackageResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePackageResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeletePackageResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeletePackageResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteSupplierRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteSupplierRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteSupplierRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteSupplierRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    supplier: (f = msg.getSupplier()) && proto.catalog.Supplier.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteSupplierRequest}
 */
proto.catalog.DeleteSupplierRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteSupplierRequest;
  return proto.catalog.DeleteSupplierRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteSupplierRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteSupplierRequest}
 */
proto.catalog.DeleteSupplierRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Supplier;
      reader.readMessage(value,proto.catalog.Supplier.deserializeBinaryFromReader);
      msg.setSupplier(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteSupplierRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteSupplierRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteSupplierRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteSupplierRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSupplier();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Supplier.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteSupplierRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteSupplierRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Supplier supplier = 2;
 * @return {?proto.catalog.Supplier}
 */
proto.catalog.DeleteSupplierRequest.prototype.getSupplier = function() {
  return /** @type{?proto.catalog.Supplier} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Supplier, 2));
};


/** @param {?proto.catalog.Supplier|undefined} value */
proto.catalog.DeleteSupplierRequest.prototype.setSupplier = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteSupplierRequest.prototype.clearSupplier = function() {
  this.setSupplier(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteSupplierRequest.prototype.hasSupplier = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteSupplierResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteSupplierResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteSupplierResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteSupplierResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteSupplierResponse}
 */
proto.catalog.DeleteSupplierResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteSupplierResponse;
  return proto.catalog.DeleteSupplierResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteSupplierResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteSupplierResponse}
 */
proto.catalog.DeleteSupplierResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteSupplierResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteSupplierResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteSupplierResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteSupplierResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteSupplierResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteSupplierResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeletePropertyDefinitionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeletePropertyDefinitionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeletePropertyDefinitionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePropertyDefinitionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    propertydefinition: (f = msg.getPropertydefinition()) && proto.catalog.PropertyDefinition.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeletePropertyDefinitionRequest}
 */
proto.catalog.DeletePropertyDefinitionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeletePropertyDefinitionRequest;
  return proto.catalog.DeletePropertyDefinitionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeletePropertyDefinitionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeletePropertyDefinitionRequest}
 */
proto.catalog.DeletePropertyDefinitionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.PropertyDefinition;
      reader.readMessage(value,proto.catalog.PropertyDefinition.deserializeBinaryFromReader);
      msg.setPropertydefinition(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeletePropertyDefinitionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeletePropertyDefinitionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeletePropertyDefinitionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePropertyDefinitionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPropertydefinition();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.PropertyDefinition.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeletePropertyDefinitionRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeletePropertyDefinitionRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional PropertyDefinition propertyDefinition = 2;
 * @return {?proto.catalog.PropertyDefinition}
 */
proto.catalog.DeletePropertyDefinitionRequest.prototype.getPropertydefinition = function() {
  return /** @type{?proto.catalog.PropertyDefinition} */ (
    jspb.Message.getWrapperField(this, proto.catalog.PropertyDefinition, 2));
};


/** @param {?proto.catalog.PropertyDefinition|undefined} value */
proto.catalog.DeletePropertyDefinitionRequest.prototype.setPropertydefinition = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeletePropertyDefinitionRequest.prototype.clearPropertydefinition = function() {
  this.setPropertydefinition(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeletePropertyDefinitionRequest.prototype.hasPropertydefinition = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeletePropertyDefinitionResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeletePropertyDefinitionResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeletePropertyDefinitionResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePropertyDefinitionResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeletePropertyDefinitionResponse}
 */
proto.catalog.DeletePropertyDefinitionResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeletePropertyDefinitionResponse;
  return proto.catalog.DeletePropertyDefinitionResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeletePropertyDefinitionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeletePropertyDefinitionResponse}
 */
proto.catalog.DeletePropertyDefinitionResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeletePropertyDefinitionResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeletePropertyDefinitionResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeletePropertyDefinitionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeletePropertyDefinitionResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeletePropertyDefinitionResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeletePropertyDefinitionResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteUnitOfMeasureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteUnitOfMeasureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteUnitOfMeasureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    unitofmeasure: (f = msg.getUnitofmeasure()) && proto.catalog.UnitOfMeasure.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteUnitOfMeasureRequest}
 */
proto.catalog.DeleteUnitOfMeasureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteUnitOfMeasureRequest;
  return proto.catalog.DeleteUnitOfMeasureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteUnitOfMeasureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteUnitOfMeasureRequest}
 */
proto.catalog.DeleteUnitOfMeasureRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.UnitOfMeasure;
      reader.readMessage(value,proto.catalog.UnitOfMeasure.deserializeBinaryFromReader);
      msg.setUnitofmeasure(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteUnitOfMeasureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteUnitOfMeasureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteUnitOfMeasureRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUnitofmeasure();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.UnitOfMeasure.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional UnitOfMeasure unitOfMeasure = 2;
 * @return {?proto.catalog.UnitOfMeasure}
 */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.getUnitofmeasure = function() {
  return /** @type{?proto.catalog.UnitOfMeasure} */ (
    jspb.Message.getWrapperField(this, proto.catalog.UnitOfMeasure, 2));
};


/** @param {?proto.catalog.UnitOfMeasure|undefined} value */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.setUnitofmeasure = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.clearUnitofmeasure = function() {
  this.setUnitofmeasure(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteUnitOfMeasureRequest.prototype.hasUnitofmeasure = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteUnitOfMeasureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteUnitOfMeasureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteUnitOfMeasureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteUnitOfMeasureResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteUnitOfMeasureResponse}
 */
proto.catalog.DeleteUnitOfMeasureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteUnitOfMeasureResponse;
  return proto.catalog.DeleteUnitOfMeasureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteUnitOfMeasureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteUnitOfMeasureResponse}
 */
proto.catalog.DeleteUnitOfMeasureResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteUnitOfMeasureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteUnitOfMeasureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteUnitOfMeasureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteUnitOfMeasureResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteUnitOfMeasureResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteUnitOfMeasureResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteItemInstanceRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteItemInstanceRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteItemInstanceRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemInstanceRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    instance: (f = msg.getInstance()) && proto.catalog.ItemInstance.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteItemInstanceRequest}
 */
proto.catalog.DeleteItemInstanceRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteItemInstanceRequest;
  return proto.catalog.DeleteItemInstanceRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteItemInstanceRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteItemInstanceRequest}
 */
proto.catalog.DeleteItemInstanceRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.ItemInstance;
      reader.readMessage(value,proto.catalog.ItemInstance.deserializeBinaryFromReader);
      msg.setInstance(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteItemInstanceRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteItemInstanceRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteItemInstanceRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemInstanceRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getInstance();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.ItemInstance.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteItemInstanceRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteItemInstanceRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ItemInstance instance = 2;
 * @return {?proto.catalog.ItemInstance}
 */
proto.catalog.DeleteItemInstanceRequest.prototype.getInstance = function() {
  return /** @type{?proto.catalog.ItemInstance} */ (
    jspb.Message.getWrapperField(this, proto.catalog.ItemInstance, 2));
};


/** @param {?proto.catalog.ItemInstance|undefined} value */
proto.catalog.DeleteItemInstanceRequest.prototype.setInstance = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteItemInstanceRequest.prototype.clearInstance = function() {
  this.setInstance(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteItemInstanceRequest.prototype.hasInstance = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteItemInstanceResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteItemInstanceResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteItemInstanceResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemInstanceResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteItemInstanceResponse}
 */
proto.catalog.DeleteItemInstanceResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteItemInstanceResponse;
  return proto.catalog.DeleteItemInstanceResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteItemInstanceResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteItemInstanceResponse}
 */
proto.catalog.DeleteItemInstanceResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteItemInstanceResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteItemInstanceResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteItemInstanceResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemInstanceResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteItemInstanceResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteItemInstanceResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteManufacturerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteManufacturerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteManufacturerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteManufacturerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    manufacturer: (f = msg.getManufacturer()) && proto.catalog.Manufacturer.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteManufacturerRequest}
 */
proto.catalog.DeleteManufacturerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteManufacturerRequest;
  return proto.catalog.DeleteManufacturerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteManufacturerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteManufacturerRequest}
 */
proto.catalog.DeleteManufacturerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Manufacturer;
      reader.readMessage(value,proto.catalog.Manufacturer.deserializeBinaryFromReader);
      msg.setManufacturer(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteManufacturerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteManufacturerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteManufacturerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteManufacturerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getManufacturer();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Manufacturer.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteManufacturerRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteManufacturerRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Manufacturer manufacturer = 2;
 * @return {?proto.catalog.Manufacturer}
 */
proto.catalog.DeleteManufacturerRequest.prototype.getManufacturer = function() {
  return /** @type{?proto.catalog.Manufacturer} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Manufacturer, 2));
};


/** @param {?proto.catalog.Manufacturer|undefined} value */
proto.catalog.DeleteManufacturerRequest.prototype.setManufacturer = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteManufacturerRequest.prototype.clearManufacturer = function() {
  this.setManufacturer(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteManufacturerRequest.prototype.hasManufacturer = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteManufacturerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteManufacturerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteManufacturerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteManufacturerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteManufacturerResponse}
 */
proto.catalog.DeleteManufacturerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteManufacturerResponse;
  return proto.catalog.DeleteManufacturerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteManufacturerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteManufacturerResponse}
 */
proto.catalog.DeleteManufacturerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteManufacturerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteManufacturerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteManufacturerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteManufacturerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteManufacturerResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteManufacturerResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteItemManufacturerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteItemManufacturerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteItemManufacturerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemManufacturerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    itemmanufacturer: (f = msg.getItemmanufacturer()) && proto.catalog.ItemManufacturer.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteItemManufacturerRequest}
 */
proto.catalog.DeleteItemManufacturerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteItemManufacturerRequest;
  return proto.catalog.DeleteItemManufacturerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteItemManufacturerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteItemManufacturerRequest}
 */
proto.catalog.DeleteItemManufacturerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.ItemManufacturer;
      reader.readMessage(value,proto.catalog.ItemManufacturer.deserializeBinaryFromReader);
      msg.setItemmanufacturer(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteItemManufacturerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteItemManufacturerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteItemManufacturerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemManufacturerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getItemmanufacturer();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.ItemManufacturer.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteItemManufacturerRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteItemManufacturerRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ItemManufacturer itemManufacturer = 2;
 * @return {?proto.catalog.ItemManufacturer}
 */
proto.catalog.DeleteItemManufacturerRequest.prototype.getItemmanufacturer = function() {
  return /** @type{?proto.catalog.ItemManufacturer} */ (
    jspb.Message.getWrapperField(this, proto.catalog.ItemManufacturer, 2));
};


/** @param {?proto.catalog.ItemManufacturer|undefined} value */
proto.catalog.DeleteItemManufacturerRequest.prototype.setItemmanufacturer = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteItemManufacturerRequest.prototype.clearItemmanufacturer = function() {
  this.setItemmanufacturer(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteItemManufacturerRequest.prototype.hasItemmanufacturer = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteItemManufacturerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteItemManufacturerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteItemManufacturerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemManufacturerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteItemManufacturerResponse}
 */
proto.catalog.DeleteItemManufacturerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteItemManufacturerResponse;
  return proto.catalog.DeleteItemManufacturerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteItemManufacturerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteItemManufacturerResponse}
 */
proto.catalog.DeleteItemManufacturerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteItemManufacturerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteItemManufacturerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteItemManufacturerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteItemManufacturerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteItemManufacturerResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteItemManufacturerResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    category: (f = msg.getCategory()) && proto.catalog.Category.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteCategoryRequest}
 */
proto.catalog.DeleteCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteCategoryRequest;
  return proto.catalog.DeleteCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteCategoryRequest}
 */
proto.catalog.DeleteCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Category;
      reader.readMessage(value,proto.catalog.Category.deserializeBinaryFromReader);
      msg.setCategory(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCategory();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Category.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteCategoryRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteCategoryRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Category category = 2;
 * @return {?proto.catalog.Category}
 */
proto.catalog.DeleteCategoryRequest.prototype.getCategory = function() {
  return /** @type{?proto.catalog.Category} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Category, 2));
};


/** @param {?proto.catalog.Category|undefined} value */
proto.catalog.DeleteCategoryRequest.prototype.setCategory = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteCategoryRequest.prototype.clearCategory = function() {
  this.setCategory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteCategoryRequest.prototype.hasCategory = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteCategoryResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteCategoryResponse}
 */
proto.catalog.DeleteCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteCategoryResponse;
  return proto.catalog.DeleteCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteCategoryResponse}
 */
proto.catalog.DeleteCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteCategoryResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteCategoryResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteLocalisationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteLocalisationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteLocalisationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteLocalisationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    connectionid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    localisation: (f = msg.getLocalisation()) && proto.catalog.Localisation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteLocalisationRequest}
 */
proto.catalog.DeleteLocalisationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteLocalisationRequest;
  return proto.catalog.DeleteLocalisationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteLocalisationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteLocalisationRequest}
 */
proto.catalog.DeleteLocalisationRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setConnectionid(value);
      break;
    case 2:
      var value = new proto.catalog.Localisation;
      reader.readMessage(value,proto.catalog.Localisation.deserializeBinaryFromReader);
      msg.setLocalisation(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteLocalisationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteLocalisationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteLocalisationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteLocalisationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConnectionid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLocalisation();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.catalog.Localisation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string connectionId = 1;
 * @return {string}
 */
proto.catalog.DeleteLocalisationRequest.prototype.getConnectionid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.catalog.DeleteLocalisationRequest.prototype.setConnectionid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Localisation localisation = 2;
 * @return {?proto.catalog.Localisation}
 */
proto.catalog.DeleteLocalisationRequest.prototype.getLocalisation = function() {
  return /** @type{?proto.catalog.Localisation} */ (
    jspb.Message.getWrapperField(this, proto.catalog.Localisation, 2));
};


/** @param {?proto.catalog.Localisation|undefined} value */
proto.catalog.DeleteLocalisationRequest.prototype.setLocalisation = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.catalog.DeleteLocalisationRequest.prototype.clearLocalisation = function() {
  this.setLocalisation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.catalog.DeleteLocalisationRequest.prototype.hasLocalisation = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.catalog.DeleteLocalisationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.catalog.DeleteLocalisationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.catalog.DeleteLocalisationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteLocalisationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.catalog.DeleteLocalisationResponse}
 */
proto.catalog.DeleteLocalisationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.catalog.DeleteLocalisationResponse;
  return proto.catalog.DeleteLocalisationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.catalog.DeleteLocalisationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.catalog.DeleteLocalisationResponse}
 */
proto.catalog.DeleteLocalisationResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.catalog.DeleteLocalisationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.catalog.DeleteLocalisationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.catalog.DeleteLocalisationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.catalog.DeleteLocalisationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool result = 1;
 * @return {boolean}
 */
proto.catalog.DeleteLocalisationResponse.prototype.getResult = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.catalog.DeleteLocalisationResponse.prototype.setResult = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * @enum {number}
 */
proto.catalog.StoreType = {
  MONGO: 0
};

/**
 * @enum {number}
 */
proto.catalog.Currency = {
  US: 0,
  CAN: 1,
  EURO: 2
};

goog.object.extend(exports, proto.catalog);
