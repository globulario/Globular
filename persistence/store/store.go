package store

import (
	base64 "encoding/base64"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"github.com/davecourtois/GoXapian"
	"github.com/davecourtois/Utility"
	"github.com/golang/protobuf/proto"
)

// operations
type persist_op struct {
	entity *persistencepb.Entity
	err    chan error
}

/**
 * The store where the value will be save.
 */
type Store struct {
	// channel
	path                   string
	persist_entity_channel chan *persist_op
}

var (
	s *Store
)

// Create a new data store.
func NewStore() *Store {
	if s == nil {
		s = new(Store)

		// Set the store path.
		s.path = os.Args[0]
		s.path, _ = filepath.Abs(filepath.Dir(os.Args[0]))

		// Set the channel.
		s.persist_entity_channel = make(chan *persist_op)

		go func() {
			s.run()
		}()
	}

	return s
}

// Take entity message and persist it inside the data store.
func (self *Store) PersistEntity(entity *persistencepb.Entity) error {

	op := new(persist_op)
	op.err = make(chan error)
	op.entity = entity

	// send the operation
	self.persist_entity_channel <- op

	// wait for it result.
	err := <-op.err

	// simple test.
	entities, _ := self.getEntitiesByTypeName(entity.Typename)
	log.Println("-----> entities: ", len(entities))
	log.Println("72 -----> entity: ", entity.UUID)
	// test single entity
	entity, err = self.getEntityByUuid(entity.Typename, entity.UUID)
	if err == nil {
		log.Println("-----> entity: ", entity.UUID)
	}
	return err
}

// Remove ambiquous query symbols % - . and replace it with _
func formalize(uuid string) string {
	return strings.TrimSpace(strings.ToLower(Utility.ToString(strings.Replace(strings.Replace(strings.Replace(uuid, "-", "_", -1), ".", "_", -1), "%", "_", -1))))
}

func (self *Store) run() {

	// so each data type will have it own store.
	wstores := make(map[string]xapian.WritableDatabase)

	// Keep store in memory.
	for {
		select {
		case op := <-s.persist_entity_channel:
			entity := op.entity

			path := s.path + "/" + entity.Typename + ".glass"
			if wstores[path] == nil {
				wstores[path] = xapian.NewWritableDatabase(path, xapian.DB_CREATE_OR_OPEN)
			}

			wstores[path].Begin_transaction()

			// So here I will index the property found in the entity.
			doc := xapian.NewDocument()

			// here I will marshal the entity and save it in the db.
			data, err := proto.Marshal(entity)

			if err == nil && len(data) > 0 {
				s.indexEntity(doc, entity)
				doc.Set_data(string(data))
				doc.Add_boolean_term("Q" + formalize(entity.UUID))
				wstores[path].Replace_document("Q"+formalize(entity.UUID), doc)
			}

			wstores[path].Commit_transaction()

			// Release the document memory.
			xapian.DeleteDocument(doc)

			// done return the error or nil
			op.err <- err
		}
	}

}

////////////////////////////////////////////////////////////////////////////////
// Indexation functionality
////////////////////////////////////////////////////////////////////////////////

// Generate prefix is use to create a indexation key for a given document.
// the field must be in the index or id's.
func generatePrefix(typeName string, field string) string {
	// remove the M_ part of the field name
	prefix := typeName
	if len(field) > 0 {
		prefix += "." + field
	}

	// replace unwanted character's
	prefix = strings.Replace(prefix, ".", "_", -1) + "%"
	prefix = "X" + strings.ToLower(prefix)

	return prefix
}

// Index entity string field.
func (this *Store) indexStringField(data string, field string, typeName string, termGenerator xapian.TermGenerator) {
	// I will index all string field to be able to found it back latter.
	termGenerator.Index_text(strings.ToLower(data), uint(1), strings.ToUpper(field))
	if Utility.IsUriBase64(data) {
		data_, err := base64.StdEncoding.DecodeString(data)
		if err == nil {
			if strings.Index(data, ":text/") > -1 || strings.Index(data, ":application/") > -1 {
				termGenerator.Index_text(strings.ToLower(string(data_)))
			}
		}
	} else if Utility.IsStdBase64(data) {
		data_, err := base64.StdEncoding.DecodeString(data)
		if err == nil {
			termGenerator.Index_text(strings.ToLower(string(data_)))
			termGenerator.Index_text(strings.ToLower(string(data)))
		}
	} else {
		termGenerator.Index_text(strings.ToLower(data))
	}
}

// Index entity field
func (this *Store) indexField(data interface{}, field string, fieldType string, typeName string, termGenerator xapian.TermGenerator, doc xapian.Document, index int) {
	// This will give possibility to search for given fields.
	if data != nil {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			s := reflect.ValueOf(data)
			for i := 0; i < s.Len(); i++ {
				// I will remove nil values.
				item := s.Index(i)
				if item.IsValid() {
					zeroValue := reflect.Zero(item.Type())
					if zeroValue != item {
						this.indexField(s.Index(i).Interface(), field, fieldType, typeName, termGenerator, doc, -1)
					} else {
						this.indexField(nil, field, fieldType, typeName, termGenerator, doc, -1)
					}
				}
			}
		} else {
			if index != -1 {
				doc.Add_value(uint(index), Utility.ToString(data))
			}
			if fieldType == "numeric" {
				value := Utility.ToNumeric(data)
				doc.Add_value(uint(index), xapian.Sortable_serialise(value))
			} else if reflect.TypeOf(data).Kind() == reflect.String {
				str := Utility.ToString(data)
				// If the the value is a valid entity reference i I will use boolean term.
				if Utility.IsValidEntityReferenceName(str) {
					term := generatePrefix(typeName, field) + formalize(str)
					doc.Add_boolean_term(term)
				} else {
					this.indexStringField(str, field, typeName, termGenerator)
				}
			}
		}
	} else {
		doc.Add_value(uint(index), "null")
	}
}

// index entity information.
func (self *Store) indexEntity(doc xapian.Document, entity *persistencepb.Entity) {

	// The term generator
	termGenerator := xapian.NewTermGenerator()

	// set english by default.
	stemmer := xapian.NewStem("en")

	termGenerator.Set_stemmer(stemmer)
	termGenerator.Set_document(doc)

	// Regular text indexation...
	termGenerator.Index_text(entity.Typename, uint(1), "TYPENAME")

	// Boolean term indexation exact match.
	typeNameIndex := generatePrefix(entity.Typename, "TYPENAME") + formalize(entity.Typename)
	doc.Add_boolean_term(typeNameIndex)

	// Indexation of attributes.
	for i := 0; i < len(entity.Attibutes); i++ {
		attr := entity.Attibutes[i]
		switch value := attr.Value.(type) {
		case *persistencepb.Attribute_BoolArr:
			self.indexField(value.BoolArr.Values, attr.Name, "[]bool", entity.Typename, termGenerator, doc, i)
		case *persistencepb.Attribute_BoolVal:
			self.indexField(value.BoolVal, attr.Name, "bool", entity.Typename, termGenerator, doc, i)
		case *persistencepb.Attribute_NumericArr:
			self.indexField(value.NumericArr.Values, attr.Name, "[]numeric", entity.Typename, termGenerator, doc, i)
		case *persistencepb.Attribute_NumericVal:
			self.indexField(value.NumericVal, attr.Name, "numeric", entity.Typename, termGenerator, doc, i)
		case *persistencepb.Attribute_StrArr:
			self.indexField(value.StrArr.Values, attr.Name, "[]string", entity.Typename, termGenerator, doc, i)
		case *persistencepb.Attribute_StrVal:
			self.indexField(value.StrVal, attr.Name, "string", entity.Typename, termGenerator, doc, i)
		}
	}

	xapian.DeleteStem(stemmer)
	xapian.DeleteTermGenerator(termGenerator)
}

//////////////////////////////////////////////////////////////////////////////////
// Search functionality
//////////////////////////////////////////////////////////////////////////////////

/**
 * Return the list of all entities for a given typename.
 */
func (self *Store) getEntitiesByTypeName(typeName string) ([]*persistencepb.Entity, error) {
	var err error
	path := self.path + "/" + typeName + ".glass"
	db := xapian.NewDatabase(path)

	// set the query.
	typeNameIndex := generatePrefix(typeName, "TYPENAME") + formalize(typeName)
	query := xapian.NewQuery(typeNameIndex)

	// Execute the search.
	enquire := xapian.NewEnquire(db)
	enquire.Set_query(query)
	mset := enquire.Get_mset(uint(0), uint(10000))

	// Return an array of entities
	entities := make([]*persistencepb.Entity, 0)

	// Now I will process the results.
	for i := 0; i < mset.Size(); i++ {
		it := mset.Get_hit(uint(i))
		doc := it.Get_document()
		entity := new(persistencepb.Entity)
		// In that case the data contain in the document are return.
		err = proto.Unmarshal([]byte(doc.Get_data()), entity)
		if err == nil {
			entities = append(entities, entity)
		}
		xapian.DeleteDocument(doc)
		xapian.DeleteMSetIterator(it)
	}

	db.Close()

	// release memory.
	xapian.DeleteMSet(mset)
	xapian.DeleteEnquire(enquire)
	xapian.DeleteQuery(query)
	xapian.DeleteDatabase(db)

	return entities, err
}

/**
 * Return an entity with a given uuid.
 */
func (self *Store) getEntityByUuid(typeName string, uuid string) (*persistencepb.Entity, error) {

	var err error
	path := self.path + "/" + typeName + ".glass"
	db := xapian.NewDatabase(path)

	// set the query.
	query := xapian.NewQuery("Q" + formalize(uuid))

	// Execute the search.
	enquire := xapian.NewEnquire(db)
	enquire.Set_query(query)
	mset := enquire.Get_mset(uint(0), uint(10000))

	// Return an array of entities
	var entity *persistencepb.Entity

	// Must be no more than one result...
	for i := 0; i < mset.Size(); i++ {
		it := mset.Get_hit(uint(i))
		doc := it.Get_document()
		// In that case the data contain in the document are return.
		entity = new(persistencepb.Entity)
		err = proto.Unmarshal([]byte(doc.Get_data()), entity)
		xapian.DeleteDocument(doc)
		xapian.DeleteMSetIterator(it)
	}

	db.Close()

	// release memory.
	xapian.DeleteMSet(mset)
	xapian.DeleteEnquire(enquire)
	xapian.DeleteQuery(query)
	xapian.DeleteDatabase(db)

	return entity, err
}
