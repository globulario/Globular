package catalog_client

import (
	"strconv"

	"context"

	"github.com/davecourtois/Globular/services/golang/catalog/catalogpb"
	globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"github.com/davecourtois/Utility"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// catalog Client Service
////////////////////////////////////////////////////////////////////////////////

type Catalog_Client struct {
	cc *grpc.ClientConn
	c  catalogpb.CatalogServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The port of the client.
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewCatalogService_Client(address string, id string) (*Catalog_Client, error) {
	client := new(Catalog_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}

	client.c = catalogpb.NewCatalogServiceClient(client.cc)

	return client, nil
}

func (self *Catalog_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *Catalog_Client) GetDomain() string {
	return self.domain
}

func (self *Catalog_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *Catalog_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Catalog_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Catalog_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Catalog_Client) SetPort(port int) {
	self.port = port
}

// Set the client instance id.
func (self *Catalog_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *Catalog_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Catalog_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Catalog_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Catalog_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Catalog_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Catalog_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Catalog_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Catalog_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Catalog_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Catalog_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////////////////// API ////////////////////////
// Stop the service.
func (self *Catalog_Client) StopService() {
	self.c.Stop(globular.GetClientContext(self), &catalogpb.StopRequest{})
}

// Create a new datastore connection.
func (self *Catalog_Client) CreateConnection(connectionId string, name string, host string, port float64, storeType float64, user string, pwd string, timeout float64, options string) error {
	rqst := &catalogpb.CreateConnectionRqst{
		Connection: &catalogpb.Connection{
			Id:       connectionId,
			Name:     name,
			Host:     host,
			Port:     int32(Utility.ToInt(port)),
			Store:    catalogpb.StoreType(storeType),
			User:     user,
			Password: pwd,
			Timeout:  int32(Utility.ToInt(timeout)),
			Options:  options,
		},
	}

	_, err := self.c.CreateConnection(globular.GetClientContext(self), rqst)
	return err
}

/**
 * Create a new unit of measure
 */
func (self *Catalog_Client) SaveUnitOfMesure(connectionId string, id string, languageCode string, name string, abreviation string, description string) error {
	rqst := &catalogpb.SaveUnitOfMeasureRequest{
		ConnectionId: connectionId,
		UnitOfMeasure: &catalogpb.UnitOfMeasure{
			Id:           id,
			LanguageCode: languageCode,
			Name:         name,
			Description:  description,
			Abreviation:  abreviation,
		},
	}

	_, err := self.c.SaveUnitOfMeasure(globular.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Save item property definition.
 */
func (self *Catalog_Client) SavePropertyDefinition(connectionId string, id string, languageCode string, name string, abreviation string, description string, valueType float64) error {
	rqst := &catalogpb.SavePropertyDefinitionRequest{
		ConnectionId: connectionId,
		PropertyDefinition: &catalogpb.PropertyDefinition{
			Id:           id,
			LanguageCode: languageCode,
			Name:         name,
			Description:  description,
			Abreviation:  abreviation,
			Type:         catalogpb.PropertyDefinition_Type(int32(Utility.ToInt(valueType))),
		},
	}

	_, err := self.c.SavePropertyDefinition(globular.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Save item property definition.
 */
func (self *Catalog_Client) SaveItemDefinition(connectionId string, id string, languageCode string, name string, abreviation string, description string, properties_str string, properties_ids_str string) error {

	properties := new(catalogpb.References)

	jsonpb.UnmarshalString(properties_str, properties)
	jsonpb.UnmarshalString(properties_ids_str, properties)

	rqst := &catalogpb.SaveItemDefinitionRequest{
		ConnectionId: connectionId,
		ItemDefinition: &catalogpb.ItemDefinition{
			Id:           id,
			LanguageCode: languageCode,
			Name:         name,
			Description:  description,
			Abreviation:  abreviation,
			Properties:   properties,
		},
	}

	_, err := self.c.SaveItemDefinition(globular.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Save item property definition.
 */
func (self *Catalog_Client) SaveItemInstance(connectionId string, jsonStr string) error {

	instance := new(catalogpb.ItemInstance)

	err := jsonpb.UnmarshalString(jsonStr, instance)
	if err != nil {
		return err
	}

	rqst := &catalogpb.SaveItemInstanceRequest{
		ItemInstance: instance,
		ConnectionId: connectionId,
	}

	_, err = self.c.SaveItemInstance(globular.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Save a Manufacturer whitout item.
 */
func (self *Catalog_Client) SaveManufacturer(connectionId string, id string, name string) error {
	rqst := &catalogpb.SaveManufacturerRequest{
		ConnectionId: connectionId,
		Manufacturer: &catalogpb.Manufacturer{
			Id:   id,
			Name: name,
		},
	}

	_, err := self.c.SaveManufacturer(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Save package, create it if it not already exist.
 * TODO append subPackage param and itemInstance...
 */
func (self *Catalog_Client) SavePackage(connectionId string, id string, name string, languageCode string, description string, inventories []*catalogpb.Inventory) error {

	// The request.
	rqst := &catalogpb.SavePackageRequest{
		Package: &catalogpb.Package{
			Id:            id,
			Name:          name,
			LanguageCode:  languageCode,
			Description:   description,
			Subpackages:   make([]*catalogpb.SubPackage, 0),
			ItemInstances: make([]*catalogpb.ItemInstancePackage, 0),
		},
		ConnectionId: connectionId,
	}

	_, err := self.c.SavePackage(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Save a supplier.
 */
func (self *Catalog_Client) SaveSupplier(connectionId string, id string, name string) error {
	rqst := &catalogpb.SaveSupplierRequest{
		ConnectionId: connectionId,
		Supplier: &catalogpb.Supplier{
			Id:   id,
			Name: name,
		},
	}

	_, err := self.c.SaveSupplier(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Save package supplier.
 */
func (self *Catalog_Client) SavePackageSupplier(connectionId string, id string, supplier_ref_str string, packege_ref_str string, price_str string, date int64, qty int64) error {

	// Supplier Ref.
	supplierRef := new(catalogpb.Reference)
	err := jsonpb.UnmarshalString(supplier_ref_str, supplierRef)
	if err != nil {
		return err
	}

	// Pacakge Ref.
	packageRef := new(catalogpb.Reference)
	err = jsonpb.UnmarshalString(packege_ref_str, packageRef)
	if err != nil {
		return err
	}

	price := new(catalogpb.Price)
	err = jsonpb.UnmarshalString(price_str, price)
	if err != nil {
		return err
	}

	rqst := new(catalogpb.SavePackageSupplierRequest)
	rqst.ConnectionId = connectionId
	rqst.PackageSupplier = &catalogpb.PackageSupplier{Id: id, Supplier: supplierRef, Package: packageRef, Price: price, Date: date, Quantity: qty}

	_, err = self.c.SavePackageSupplier(globular.GetClientContext(self), rqst)
	return err
}

/**
 * Save Item Manufacturer.
 */
func (self *Catalog_Client) SaveItemManufacturer(connectionId string, id string, manufacturer_ref_str string, item_ref_str string) error {

	// Supplier Ref.
	manufacturerRef := new(catalogpb.Reference)
	err := jsonpb.UnmarshalString(manufacturer_ref_str, manufacturerRef)
	if err != nil {
		return err
	}

	// Item Ref.
	itemRef := new(catalogpb.Reference)
	err = jsonpb.UnmarshalString(item_ref_str, itemRef)
	if err != nil {
		return err
	}

	rqst := new(catalogpb.SaveItemManufacturerRequest)
	rqst.ConnectionId = connectionId
	rqst.ItemManafacturer = &catalogpb.ItemManufacturer{Id: id, Manufacturer: manufacturerRef, Item: itemRef}

	_, err = self.c.SaveItemManufacturer(globular.GetClientContext(self), rqst)
	return err
}

/**
 * Save Item Manufacturer.
 */
func (self *Catalog_Client) SaveCategory(connectionId string, id string, name string, languageCode string, categories_str string) error {
	categories := new(catalogpb.References)
	jsonpb.UnmarshalString(categories_str, categories)

	rqst := &catalogpb.SaveCategoryRequest{
		ConnectionId: connectionId,
		Category: &catalogpb.Category{
			Id:           id,
			Name:         name,
			LanguageCode: languageCode,
			Categories:   categories,
		},
	}

	_, err := self.c.SaveCategory(globular.GetClientContext(self), rqst)
	return err
}

/**
 * Appen item defintion category.
 */
func (self *Catalog_Client) AppendItemDefinitionCategory(connectionId string, item_definition_ref_str string, category_ref_str string) error {
	// The item definition reference.
	itemDefinitionRef := new(catalogpb.Reference)
	err := jsonpb.UnmarshalString(item_definition_ref_str, itemDefinitionRef)
	if err != nil {
		return err
	}

	// The category reference.
	categoryRef := new(catalogpb.Reference)
	err = jsonpb.UnmarshalString(category_ref_str, categoryRef)
	if err != nil {
		return err
	}

	rqst := &catalogpb.AppendItemDefinitionCategoryRequest{
		ConnectionId:   connectionId,
		ItemDefinition: itemDefinitionRef,
		Category:       categoryRef,
	}

	_, err = self.c.AppendItemDefinitionCategory(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Remove item defintion category.
 */
func (self *Catalog_Client) RemoveItemDefinitionCategory(connectionId string, item_definition_ref_str string, category_ref_str string) error {
	// The item definition reference.
	itemDefinitionRef := new(catalogpb.Reference)
	err := jsonpb.UnmarshalString(item_definition_ref_str, itemDefinitionRef)
	if err != nil {
		return err
	}

	// The category reference.
	categoryRef := new(catalogpb.Reference)
	err = jsonpb.UnmarshalString(category_ref_str, categoryRef)
	if err != nil {
		return err
	}

	rqst := &catalogpb.RemoveItemDefinitionCategoryRequest{
		ConnectionId:   connectionId,
		ItemDefinition: itemDefinitionRef,
		Category:       categoryRef,
	}

	_, err = self.c.RemoveItemDefinitionCategory(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Save Item Localisation.
 */
func (self *Catalog_Client) SaveLocalisation(connectionId string, localisation *catalogpb.Localisation) error {

	rqst := &catalogpb.SaveLocalisationRequest{
		ConnectionId: connectionId,
		Localisation: localisation,
	}

	_, err := self.c.SaveLocalisation(globular.GetClientContext(self), rqst)
	return err
}
